package main

import (
	"flag"
	"fmt"
	"sync"
)

var addr = flag.String("addr", ":8080", "http server address")
var rooms = make(map[*Room]bool)

type RoomParams struct {
	Owner   string
	Product string
	Name    string
	Value   string
}

func main() {
	var wg sync.WaitGroup
	user := newUser("Owner")
	newRoom := newRoom(*user, "TestRoom", "Camiseta Ponzio", float32(2500))
	close := make(chan bool)
	wg.Add(1)

	go func(c chan bool) {
		defer func() {
			wg.Done()
		}()

		go newRoom.Run(c)

		add_users(1000, newRoom)
		newRoom.start <- true
		for {
			select {
			case <-c:
				return
			}
		}
	}(close)

	wg.Wait()
	fmt.Println("LLegamos al final")
}

func add_users(total int, room *Room) {
	for i := 1; i <= total; i++ {
		value := i
		name := fmt.Sprintf("User%d", value)
		user := newUser(name)
		room.register <- user
	}
}
