package chat

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type User struct {
	Username string
	Conn     *websocket.Conn
	Global   *Chat
}

func (u *User) Read() {
	for {
		if _, message, err := u.Conn.ReadMessage(); err != nil {
			log.Println("Error on Read Message:", err.Error())

			break

		} else {
			u.Global.messages <- NewMessage(string(message), u.Username)
			log.Println("I am here")

		}
	}
	log.Println("I am above part")
	u.Global.leave <- u
	log.Println("I am  below part")
}

func (u *User) Write(message *Message) {
	// log.Println(message, "at line 35 of users.go")
	b, _ := json.Marshal(message) //this comverts struct obj in json string bytes
	// log.Println(b, "at line 37 of users.go")
	if err := u.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
		log.Println("error on write message:", err.Error())
	}

}
