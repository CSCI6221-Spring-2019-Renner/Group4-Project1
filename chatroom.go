// WebChat project main.go
package main

import (
	"fmt"
	"net/http"
	"time"

	"encoding/json"

	"strings"

	"golang.org/x/net/websocket"
)

type UserMsg struct {
	UserName string
	Msg      string
	DataType string
}

type UserData struct {
	UserName string
}

type Datas struct {
	UserMsgs  []UserMsg
	UserDatas []UserData
}

var datas Datas
var users map[*websocket.Conn]string

func main() {
	fmt.Println("start:")
	fmt.Println(time.Now())

	fmt.Println("init")
	datas = Datas{}
	users = make(map[*websocket.Conn]string)

	fmt.Println("bind root")
	http.HandleFunc("/", h_index)
	fmt.Println("bind socket")
	http.Handle("/webSocket", websocket.Handler(h_webSocket))
	fmt.Println("listening")
	http.ListenAndServe(":8888", nil)
	fmt.Println("Done")
}

func h_index(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "index.html")
}

func h_webSocket(ws *websocket.Conn) {

	var userMsg UserMsg
	var data string
	for {

		if _, ok := users[ws]; !ok {
			users[ws] = "Annymous"
		}
		userMsgsLen := len(datas.UserMsgs)
		fmt.Println("UserMsgs", userMsgsLen, "number of users:", len(users))

		if userMsgsLen > 0 {
			b, errMarshl := json.Marshal(datas)
			if errMarshl != nil {
				break
			}
			for key, _ := range users {
				errMarshl = websocket.Message.Send(key, string(b))
				if errMarshl != nil {
					delete(users, key)
					break
				}
			}
			datas.UserMsgs = make([]UserMsg, 0)
		}

		err := websocket.Message.Receive(ws, &data)
		fmt.Println("data：", data)
		if err != nil {
			delete(users, ws)
			break
		}

		data = strings.Replace(data, "\n", "", 0)
		err = json.Unmarshal([]byte(data), &userMsg)
		if err != nil {
			break
		}

		switch userMsg.DataType {
		case "send":
			if _, ok := users[ws]; ok {
				users[ws] = userMsg.UserName
				datas.UserDatas = make([]UserData, 0)
				for _, item := range users {
					userData := UserData{UserName: item}
					datas.UserDatas = append(datas.UserDatas, userData)
				}
			}
			datas.UserMsgs = append(datas.UserMsgs, userMsg)
		}
	}

}
