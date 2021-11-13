
package main

import (
	"log"
	"strconv"
)

type Room struct {
	Owner         User							`json:"Owner"`
	Winner  			User							`json:"Winner"`
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
  Owner   		string
	Winner  		string
  Name    		string
  Product 		string
  BaseValue   float32
  Users   		[]string
	Amount      float32
  Started 		bool
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

	if (len(room.users) >= 3) {
		room.Started = true
	}
	newRoomMessage := &RoomMessage {
		Owner: room.Owner.Name,
		Winner: room.Winner.Name,
		Name: room.Name,
		Product: room.Product,
		BaseValue: room.BaseValue,
		Amount: user.Amount,
		Users: room.mapUsers(),
		Started: room.Started,
	}

	for u, _ := range room.users {
    u.Sender <-newRoomMessage
  }
}

func (room *Room) unregisterUser(user *User) {
	log.Println("Sale usuario", user)
}

func (room *Room) processOffer(message *UserMessage) {
	currentOffer, e := strconv.ParseFloat(message.Offer, 16)
	if (e != nil) {
		log.Println("No se pudo parsear el valor", message.Offer)
		return
	}
	validOffer := float32(currentOffer)

	if (room.BaseValue < validOffer) {
		message.User.Amount = message.User.Amount - validOffer
		room.Winner = *message.User
		room.BaseValue = validOffer
	}

	for u, _ := range room.users {
		newRoomMessage := &RoomMessage {
			Owner: room.Owner.Name,
			Name: room.Name,
			Winner: room.Winner.Name,
			Product: room.Product,
			BaseValue: room.BaseValue,
			Users: room.mapUsers(),
			Started: room.Started,
			Amount: u.Amount,
		}
    u.Sender <-newRoomMessage
  }
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
