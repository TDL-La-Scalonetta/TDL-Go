
package main

import (
	"log"
)
type Room struct {
	Owner         User
	Name 					string
	Product 			string
	BaseValue 		float32
	Participants  int
	users 				map[*User]bool
	register 			chan *User
	unregister 		chan *User
	offers        chan *Message
}

func newRoom(Owner User, Name string, Product string, BaseValue float32) *Room {
	return &Room {
		Owner:      Owner,
		Name: 			Name,
		Product: 		Product,
		BaseValue: 	BaseValue,
		users:      make(map[*User]bool),
		register:   make(chan *User),
    unregister: make(chan *User),
		offers:     make(chan *Message),
	}
}
func (room *Room) registerUser(user *User) {
	log.Println("Entra usuario", user)
}

func (room *Room) unregisterUser(user *User) {
	log.Println("Sale usuario", user)
}

func (room *Room) processOffer(message *Message) {
	log.Println("Nueva oferta", message)
}

func (room *Room) Run() {
	for {
			select {

			case user := <-room.register:
					room.registerUser(user)

			case user := <-room.unregister:
					room.unregisterUser(user)

			case offer := <-room.offers:
					room.processOffer(offer)
			}

	}
}
