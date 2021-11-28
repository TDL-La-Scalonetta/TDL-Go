package main

import (
	"fmt"
	"strconv"
	"time"
)

type Room struct {
	Owner      User    `json:"Owner"`
	Winner     User    `json:"Winner"`
	Name       string  `json:"Name"`
	Product    string  `json:"Product"`
	BaseValue  float32 `json:"BaseValue"`
	Started    bool    `json:"Started"`
	Finished   bool    `json:"Finished"`
	TimeLeft   time.Duration     `json:"TimeLeft"`
	Timer      *time.Ticker			 `json:"-"`
	users      map[*User]bool    `json:"-"`
	register   chan *User        `json:"-"`
	finish     chan bool 				 `json:"-"`
	unregister chan *User        `json:"-"`
	stopClock  chan bool 				 `json:"-"`
	offers     chan *UserMessage `json:"-"`
}

type RoomMessage struct {
	Type      string
	Owner     string
	Winner    string
	Name      string
	Product   string
	BaseValue float32
	Users     []string
	Amount    float32
	TimeLeft	string
	Started   bool
	Finished  bool
}

func newRoom(Owner User, Name string, Product string, BaseValue float32) *Room {
	duration, _ := time.ParseDuration("1m")
	return &Room{
		Owner:      Owner,
		Name:       Name,
		Product:    Product,
		BaseValue:  BaseValue,
		TimeLeft:   duration,
		Started:    false,
		Finished:   false,
		users:      make(map[*User]bool),
		register:   make(chan *User),
		unregister: make(chan *User),
		finish:     make(chan bool),
		stopClock:  make(chan bool),
		offers:     make(chan *UserMessage),
	}
}

func (room *Room) mapUsers() []string {
	var list = make([]string, 0)
	for user, _ := range room.users {
		list = append(list, user.Name)
	}

	return list
}

func (room *Room) notify(message *RoomMessage) {
	for u, _ := range room.users {
		u.updateStatus(message)
	}
}

func (room *Room) createMessage(messageType string) *RoomMessage {
	return &RoomMessage{
		Type:      messageType,
		Owner:     room.Owner.Name,
		Winner:    room.Winner.Name,
		Name:      room.Name,
		Product:   room.Product,
		BaseValue: room.BaseValue,
		TimeLeft:  room.TimeLeft.String(),
		Users:     room.mapUsers(),
		Started:   room.Started,
		Finished:  room.Finished,
	}
}

func (room *Room) closeSubasta(finished bool) {
	room.stopClock<-true
	room.Finished = finished
	newRoomMessage := room.createMessage("Ended")
	room.notify(newRoomMessage)

	for user, _ := range room.users {
		delete(room.users, user)
		user.disconnect()
	}
}

// Process messages
func (room *Room) RegisterUser(user *User) {
	room.users[user] = true
	newRoomMessage := room.createMessage("Enter")
	room.notify(newRoomMessage)
}

func (room *Room) UnregisterUser(user *User) {
	delete(room.users, user)
	user.disconnect()

	newRoomMessage := room.createMessage("Leave")
	room.notify(newRoomMessage)
}

func (room *Room) ProcessOffer(message *UserMessage) {
	currentOffer, e := strconv.ParseFloat(message.Offer, 16)
	if e != nil {
		fmt.Println("No se pudo parsear el valor", message.Offer)
		return
	}
	validOffer := float32(currentOffer)

	if (room.BaseValue < validOffer) {
		room.Winner = *message.User
		room.BaseValue = validOffer
		message.User.LastOffer = validOffer
	}

	newRoomMessage := room.createMessage("Offer")
	room.notify(newRoomMessage)
}

func (room *Room) StartSubasta() {
	var list = make([]string, 0)
	for user, _ := range room.users {
		if user.Name != room.Owner.Name {
			list = append(list, user.Name)
		}
	}

	if len(list) >= 3 {
		room.Started = true
		if (room.Timer == nil) {
			go room.Clock()
		}

		newRoomMessage := room.createMessage("Started")
		room.notify(newRoomMessage)
	}
}

func (room *Room) EndSubasta() {
	if (len(room.users) == 1) {
		for u, _ := range room.users {
			room.Winner = *u
			room.BaseValue = u.LastOffer
		}
		room.finish <-true
	}
}

// Routines
func (room *Room) Clock() {
	room.Timer = time.NewTicker(time.Second)

	defer func() {
		fmt.Println("Stop timer")
		room.Timer.Stop()
	}()

	for {
		select {
			case <-room.stopClock:
				return
			case <-room.Timer.C:
				room.TimeLeft -= time.Second
				fmt.Println("Tiempo restante", room.TimeLeft)
				newRoomMessage := room.createMessage("Time")
				room.notify(newRoomMessage)

				if (room.TimeLeft == 0) {
					fmt.Println("Se termino el tiempo")
					room.finish<-true
					return
				}
		}
	}
}

func (room *Room) Run() {
	for {
		select {
			case finished := <-room.finish:
				room.closeSubasta(finished)
				return
			case user := <-room.register:
				room.RegisterUser(user)
				room.StartSubasta()
			case user := <-room.unregister:
				room.UnregisterUser(user)
				room.EndSubasta()
			case offer := <-room.offers:
				room.ProcessOffer(offer)
		}
	}
}
