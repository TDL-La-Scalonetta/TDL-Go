package main

import (
  "log"
  "encoding/json"
  "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  4096,
    WriteBufferSize: 4096,
}

type User struct {
	Name 	 	 string
	Amount 	 float32
	Conn     *websocket.Conn
  Room     *Room
  channel  chan []byte
}

func newUser(name string) *User {
  return &User{
    Name:     name,
    Amount: 	10000,
    channel:  make(chan []byte, 256),
  }
}

func (user *User) handleMessage(rawMessage []byte) {
  var message Message
  if err := json.Unmarshal(rawMessage, &message); err != nil {
      log.Printf("Error on unmarshal JSON message %s", err)
  }
  message.User = user
  // se lo pasamos al room para que lo procese
  user.Room.offers <-&message
}

func (user *User) Listen() {
    for {
        _, jsonMessage, err := user.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("unexpected close error: %v", err)
            }
            break
        }
        user.handleMessage(jsonMessage)
    }
}
