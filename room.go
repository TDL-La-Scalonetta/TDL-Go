package main

import (
	"fmt"
	"time"
)

type Room struct {
	Owner      User              `json:"Owner"`
	Winner     User              `json:"Winner"`
	Name       string            `json:"Name"`
	Product    string            `json:"Product"`
	BaseValue  float32           `json:"BaseValue"`
	Started    bool              `json:"Started"`
	Finished   bool              `json:"Finished"`
	TimeLeft   time.Duration     `json:"TimeLeft"`
	Timer      *time.Ticker      `json:"-"`
	users      map[*User]bool    `json:"-"`
	register   chan *User        `json:"-"`
	finish     chan bool         `json:"-"`
	unregister chan *User        `json:"-"`
	stopClock  chan bool         `json:"-"`
	offers     chan *UserMessage `json:"-"`
	start      chan bool
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
	TimeLeft  string
	Started   bool
	Finished  bool
}

func newRoom(Owner User, Name string, Product string, BaseValue float32) *Room {
	duration, _ := time.ParseDuration("15s")
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
		finish:     make(chan bool, 1),
		start:      make(chan bool),
		stopClock:  make(chan bool),
		offers:     make(chan *UserMessage),
	}
}

func (room *Room) notify() {
	fmt.Println("Ganador:", room.Winner.Name, "Valor:", room.BaseValue)
	if room.Finished {
		fmt.Println("Subasta Finalizada")
	} else {
		fmt.Println("Subasta en curso:", room.TimeLeft)
	}
}

// Process messages
func (room *Room) RegisterUser(user *User) {
	room.users[user] = true
	user.Room = room
}

func (room *Room) ProcessOffer(message *UserMessage) {
	if room.Finished {
		return
	}

	validOffer := message.Offer

	if room.BaseValue < validOffer {
		room.Winner = *message.User
		room.BaseValue = validOffer
		message.User.LastOffer = validOffer
	}
}

func (room *Room) StartSubasta() {
	room.Started = true
	go room.startRoomClock()
	for user, _ := range room.users {
		go user.start()
	}
}

func (room *Room) EndSubasta() {
	room.Finished = true
	for user, _ := range room.users {
		user.finished <- true
	}
	room.notify()
}

// Routines
//Creates a new ticker and starts the clock for the room
func (room *Room) startRoomClock() {
	room.Timer = time.NewTicker(time.Second)

	defer func() {
		room.Timer.Stop()
		room.finish <- true
	}()

	for {
		select {
		case <-room.stopClock:
			return
		case <-room.Timer.C:
			room.TimeLeft -= time.Second
			room.notify()

			if room.TimeLeft == 0 {
				return
			}
		}
	}
}

func (room *Room) Run(c chan bool) {
	defer func() {
		fmt.Println("Terminado")
		c <- true
	}()
	for {
		select {
		case <-room.start:
			room.StartSubasta()
		case user := <-room.register:
			room.RegisterUser(user)
		case <-room.finish:
			room.EndSubasta()
			close(room.offers)
		case user := <-room.unregister:
			delete(room.users, user)
			fmt.Println("Quedan:", len(room.users))
			if len(room.users) == 0 {
				return
			}
		case offer := <-room.offers:
			room.ProcessOffer(offer)
		}
	}
}
