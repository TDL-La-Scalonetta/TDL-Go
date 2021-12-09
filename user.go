package main

import (
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type User struct {
	Name      string    `json:"Name"`
	LastOffer float32   `json:"LastOffer"`
	Room      *Room     `json:"-"`
	started   chan bool `json:"-"`
	finished  chan bool `json:"-"`
}

type UserMessage struct {
	Offer float32
	User  *User
}

func newUser(name string) *User {
	return &User{
		Name:      name,
		LastOffer: 0,
		started:   make(chan bool),
		finished:  make(chan bool, 1),
	}
}

func (user *User) start() {
	defer func() {
		recover()
	}()

	for {
		select {
		case <-user.finished:
			return
		default:
			waitTime := rand.Intn(10)
			amount := float32(rand.Intn(10)) + user.Room.BaseValue

			time.Sleep(time.Duration(waitTime) * time.Second)
			if user.Room.Winner.Name != user.Name {
				message := UserMessage{User: user, Offer: amount}
				user.Room.offers <- &message
			}
		}
	}
}
