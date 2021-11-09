



package main

type Room struct {
	Owner         User
	Name 					string
	Product 			string
	BaseValue 		float32
	Participants  int
	users 				map[*User]bool
	register 			chan *User
	unregister 		chan *User
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
	}
}
