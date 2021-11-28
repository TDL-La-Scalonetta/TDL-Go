
package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	defaultTime = 30 * time.Minute
)

type Room struct {
	Owner      User    `json:"Owner"`
	Winner     User    `json:"Winner"`
	Name       string  `json:"Name"`
	Product    string  `json:"Product"`
	BaseValue  float32 `json:"BaseValue"`
	Started    bool    `json:"Started"`
	Timer      *time.Ticker
	users      map[*User]bool    `json:"-"`
	register   chan *User        `json:"-"`
	unregister chan *User        `json:"-"`
	offers     chan *UserMessage `json:"-"`
	StartTime  time.Time         `json:"-"`
	EndTime    time.Time         `json:"-"`
}

type RoomMessage struct {
	Owner     string
	Winner    string
	Name      string
	Product   string
	BaseValue float32
	Users     []string
	Amount    float32
	Started   bool
}

func newRoom(Owner User, Name string, Product string, BaseValue float32) *Room {
	return &Room{
		Owner:      Owner,
		Name:       Name,
		Product:    Product,
		BaseValue:  BaseValue,
		Started:    false,
		users:      make(map[*User]bool),
		register:   make(chan *User),
		unregister: make(chan *User),
		offers:     make(chan *UserMessage),
	}
}

func (room *Room) StartTimer() {
	ticker := time.NewTicker(time.Second)
	room.StartTime = time.Now()
	room.EndTime = room.StartTime.Add(defaultTime)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(defaultTime)
		done <- true
	}()
	for {
		select {
		case <-done:
			fmt.Println("Termino la Subasta!") // Ver como terminar la subasta por timer
			return
		case <-ticker.C:
			timeLeft := room.EndTime.Sub(time.Now())
			fmt.Println("Tiempo restante:", timeLeft)
		}
	}
}

func (room *Room) createMessage() *RoomMessage {
	return &RoomMessage{
		Owner:     room.Owner.Name,
		Winner:    room.Winner.Name,
		Name:      room.Name,
		Product:   room.Product,
		BaseValue: room.BaseValue,
		Users:     room.mapUsers(),
		Started:   room.Started,
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

	if len(room.users) >= 3 {
		room.Started = true
		if room.Timer == nil {
			go room.StartTimer()
		}
	}
	newRoomMessage := room.createMessage()
	for u, _ := range room.users {
		u.Sender <- newRoomMessage
	}
}

func (room *Room) unregisterUser(user *User) {
	log.Println("Sale usuario", user)
	room.users[user] = false
	users[user] = false

	user.disconnect()

	newRoomMessage := room.createMessage()

	for u, _ := range room.users {
		u.Sender <- newRoomMessage
	}

	onlineUsers := 0
	for _, value := range room.users {
		if value == true {
			onlineUsers++
		}
	}
	if (onlineUsers == 1) && ((room.EndTime.Sub(time.Now())) > 0) {
		for index, value := range room.users {
			if value == true {
				fmt.Println("El ganador es ", index.Name)
				room.BaseValue = index.Amount
			}
		}
	}
	return
}

func (room *Room) processOffer(message *UserMessage) {
	currentOffer, e := strconv.ParseFloat(message.Offer, 16)
	if e != nil {
		log.Println("No se pudo parsear el valor", message.Offer)
		return
	}
	validOffer := float32(currentOffer)

	if room.BaseValue < validOffer {
		room.Winner = *message.User
		message.User.Amount = validOffer
		room.BaseValue = validOffer
	}
	newRoomMessage := room.createMessage()
	for u, _ := range room.users {
		u.Sender <- newRoomMessage
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
