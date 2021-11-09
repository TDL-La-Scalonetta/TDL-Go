package main

import (
    "log"
    "net/http"
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
	hub      *Hub
	rooms    map[*Room]bool
}

func newUser(conn *websocket.Conn, hub *Hub, name string) *User {
  return &User{
    Name:     name,
    Amount: 	10000,
    conn:			conn,
    hub:			hub,
    rooms:    make(map[*Room]bool),
  }
}

func UserConn(hub *Hub, w http.ResponseWriter, r *http.Request)  {
	name, ok := r.URL.Query()["name"]

  if !ok || len(name[0]) < 1 {
      log.Println("Url Param 'name' is missing")
      return
  }

  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return
  }

  user := newUser(conn, hub, name[0])
  hub.register <-user
}
