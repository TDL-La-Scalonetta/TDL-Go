package main

import (
	"fmt"
	"strconv"
	"time"
	"sort"
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

func (room *Room) closeSubasta() {
	fmt.Println("Fin de la subasta!")
	room.Finished = true
	newRoomMessage := room.createMessage("Ended")
	room.notify(newRoomMessage)
	fmt.Println("El ganador es:", room.Winner.Name, room.BaseValue)
}

// Process messages
func (room *Room) RegisterUser(user *User) {
	room.users[user] = true
	newRoomMessage := room.createMessage("Enter")
	room.notify(newRoomMessage)
}

func (room *Room) UnregisterUser(user *User) {
	delete(room.users, user)
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
	fmt.Println("Iniciamos la subasta si los usuarios no-owner son mas de tres: ", len(list))
	if len(list) >= 3 {
		if (room.Timer == nil) {
			go room.Clock()
		}
		if (room.Started == false) {
			room.Started = true
			fmt.Println("Creamos la notificacion que se inicia.")
			newRoomMessage := room.createMessage("Started")
			room.notify(newRoomMessage)
			fmt.Println("Enviamos notificacion de inicio a ", len(room.users))
		}
	}
}

func (room *Room) EndSubasta() {
	var list = make([]*User, 0)
	for user, _ := range room.users {
		if (user.Name != room.Owner.Name) {
			list = append(list, user)
		}
	}

	if (len(list) == 1) {
		// Terminamos si queda uno solo
		room.Winner = *list[0]
		room.BaseValue = list[0].LastOffer
		room.finish <-true
	} else {
		// Si se fue el que era el winner, ponemos al que haya hecho la oferta mayor.
		sort.Slice(list, func(i, j int) bool { return list[i].LastOffer > list[j].LastOffer })
		room.Winner = *list[0]
		room.BaseValue = list[0].LastOffer
	}
}

// Routines
func (room *Room) Clock() {
	room.Timer = time.NewTicker(time.Second)

	defer func() {
		fmt.Println("Se detiene el timer")
		room.Timer.Stop()
	}()

	for {
		select {
			case <-room.stopClock:
				return
			case <-room.Timer.C:
				room.TimeLeft -= time.Second
				newRoomMessage := room.createMessage("Time")
				room.notify(newRoomMessage)

				if (room.TimeLeft == 0) {
					fmt.Println("Se termino el tiempo, cerramos la subasta.")
					room.finish<-true
				}
		}
	}
}

func (room *Room) Run() {
	defer func() {
		room.stopClock<-true
		room.closeSubasta()
	}()

	for {
		select {
			case <-room.finish:
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
