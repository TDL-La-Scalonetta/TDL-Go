package main

type Room struct {
	Name 					string
	Product 			string
	BaseValue 		float32
	Participants  int
	users 				map[*User]bool
	register 			chan *User
	unregister 		chan *User
}

func newRoom(name string, product string, value float32) *Room {
	return &Room {
		Name: 			name,
		Product: 		product,
		BaseValue: 	value,
		users:      make(map[*User]bool),
		register:   make(chan *User),
    unregister: make(chan *User),
	}
}
