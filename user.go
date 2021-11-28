package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"fmt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type User struct {
	Name   		string            `json:"Name"`
	LastOffer float32           `json:"LastOffer"`
	Conn   		*websocket.Conn   `json:"-"`
	Room  		*Room             `json:"-"`
	Sender 		chan *RoomMessage `json:"-"`
}

type UserMessage struct {
	Action string
	Offer  string
	Room   string
	User   *User
}

func newUser(name string) *User {
	return &User{
		Name:   name,
		LastOffer: 0,
		Sender: make(chan *RoomMessage),
	}
}

func (user *User) handleMessage(rawMessage []byte) {
	var message UserMessage
	if err := json.Unmarshal(rawMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}

	switch message.Action {
		case "offer":
			message.User = user
			user.Room.offers <- &message
		case "leave":
			user.Room.unregister <- user
		case "finish":
			user.Room.finish <- true
	}
}

func (user *User) disconnect() {
	fmt.Println("Cerramos la conexion del socket", user)
	user.Conn.Close()
}

func (user *User) updateStatus(rM *RoomMessage) {
	parsed, err1 := json.Marshal(rM)
	if err1 != nil {
		log.Println(err1)
	}

	err2 := user.Conn.WriteMessage(websocket.TextMessage, parsed)
	if err2 != nil {
		log.Println("Error escribiendo en el socket")
		log.Println(err2)
	}
}

func (user *User) ReadSocket() {
	defer func() {
		user.disconnect()
	}()
	for {
		_, rawMessage, err := user.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		user.handleMessage(rawMessage)
	}
}
