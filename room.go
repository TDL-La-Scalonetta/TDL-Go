
package main

import (
	"log"
)
type Room struct {
	Owner         User							`json:"Owner"`
	Name 					string						`json:"Name"`
	Product 			string						`json:"Product"`
	BaseValue 		float32						`json:"BaseValue"`
	Started				bool							`json:"Started"`
	users 				map[*User]bool		`json:"-"`
	register 			chan *User				`json:"-"`
	unregister 		chan *User				`json:"-"`
	offers        chan *UserMessage	`json:"-"`
}

type RoomMessage struct {
  Owner   string
  Name    string
  Product string
  Value   float32
  Users   []string
  Started bool
}

func newRoom(Owner User, Name string, Product string, BaseValue float32) *Room {
	return &Room {
		Owner:      	Owner,
		Name: 				Name,
		Product: 			Product,
		BaseValue: 		BaseValue,
		Started:			false,
		users:      	make(map[*User]bool),
		register:   	make(chan *User),
    unregister: 	make(chan *User),
		offers:     	make(chan *UserMessage),
	}
}

func (room *Room) mapUsers() []string {
  var list = make([]string, 0)
  for user, _ := range room.users {
    list = append(list, user.Name)
  }

  return list
}

func (room *Room) registerUser(user *User) {
	room.users[user] = true
	newRoomMessage := &RoomMessage {
		Owner: room.Owner.Name,
		Name: room.Name,
		Product: room.Product,
		Value: room.BaseValue,
		Users: room.mapUsers(),
		Started: false,
	}

	for u, _ := range room.users {
    u.Sender <-newRoomMessage
  }
}

func (room *Room) unregisterUser(user *User) {
	log.Println("Sale usuario", user)
}

func (room *Room) processOffer(message *UserMessage) {
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
