package chat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"chat_app_golang_js/utils"

	"github.com/gorilla/websocket"
)

type Chat struct {
	users    map[string]*User
	messages chan *Message
	join     chan *User
	leave    chan *User
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("%s %s %s %v", r.Method, r.Host, r.RequestURI, r.Proto)
		return r.Method == http.MethodGet

	},
}

func (c *Chat) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error on websocket connection:", err.Error())
	}
	keys := r.URL.Query()
	log.Println(keys)
	username := keys.Get("username")
	log.Println(username)
	log.Println(strings.TrimSpace(username))

	if strings.TrimSpace(username) == "" {
		username = fmt.Sprintf("anon-%d", utils.GetRandomI64())

	}
	fmt.Println("username here is:", username)
	user := &User{
		Username: username,
		Conn:     conn,
		Global:   c,
	}
	c.join <- user
	user.Read()

}

func (c *Chat) Run() {
	for {
		log.Println("at chat.go in line59")
		select {
		case user := <-c.join:
			log.Println("i am at chat.go line 61")
			c.add(user)
		case message := <-c.messages:
			log.Println("I am at chat.go and line64")
			c.broadcast(message)
		case user := <-c.leave:
			c.disconnect(user)

		}
		log.Println("hii")
	}

}
func (c *Chat) disconnect(user *User) {
	if _, ok := c.users[user.Username]; !ok {
		log.Println("Chat.go on line 76", c.users[user.Username])
		defer user.Conn.Close()
		delete(c.users, user.Username)
		log.Printf("user left the chat: %v, Total: %d", user.Username, len(c.users))
	}
}

func (c *Chat) broadcast(message *Message) {
	// log.Printf("Broadcast message is: %v\n", message)
	// log.Println("I am at chat.go and line85====>", c.users)
	log.Println("I am at l88 of chat.go", message)
	for _, user := range c.users {
		user.Write(message)
		writedata(message)

	}

}

func writedata(message *Message) {
	filename := "messageData.json"
	err := isFileAlreadyExist(filename)
	if err != nil {
		log.Println(err)
	}

	fileData, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Println(err)
	}
	data := []Message{}

	err = json.Unmarshal(fileData, &data)

	newStruct := &Message{
		ID:     message.ID,
		Body:   message.Body,
		Sender: message.Sender,
	}
	log.Println("at line 120 in chat.go::", data)

	data = append(data, *newStruct)

	if err != nil {
		log.Println(err)
	}
	log.Println(data)

	dataInBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	log.Println("at line 134::", dataInBytes)

	err = ioutil.WriteFile(filename, dataInBytes, 0644)
	if err != nil {
		log.Println(err)
	}

}

func isFileAlreadyExist(filename string) error {
	_, err := os.Stat(filename) //STat is used to get metadata (size, creation time..) it doesnot open file
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Chat) add(user *User) {
	log.Println("I am at line94 inside chat.go:", *user)
	if _, ok := c.users[user.Username]; !ok {
		log.Println(c.users, "at chat.go and line 93")
		c.users[user.Username] = user
		log.Println(c.users, "at chat.go and line95")
		fmt.Printf("Added user: %s, Total:%d", user.Username, len(c.users))
	}
}

func Start(port string) {
	log.Printf("Chat listening on http://localhost%s\n", port)
	c := &Chat{
		users:    make(map[string]*User),
		messages: make(chan *Message),
		join:     make(chan *User),
		leave:    make(chan *User),
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to GoWebchat"))
	})
	http.HandleFunc("/chat", c.Handler)
	go c.Run()
	log.Fatal(http.ListenAndServe(port, nil))

}
