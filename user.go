package main

import (
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  4096,
    WriteBufferSize: 4096,
}

type User struct {
	Name 	 	 string
	Amount 	 float32
  conn     *websocket.Conn
	rooms    map[*Room]bool
}

func newUser(name string) *User {
  return &User{
    Name:     name,
    Amount: 	10000,
    rooms:    make(map[*Room]bool),
  }
}
