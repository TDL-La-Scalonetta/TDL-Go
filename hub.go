package main

type Hub struct {
  users      map[*User]bool
  rooms      map[*Room]bool
  register   chan *User
  unregister chan *User
  broadcast  chan []byte
}

func NewHub() *Hub {
  return &Hub{
    users:    	make(map[*User]bool),
    rooms:      make(map[*Room]bool),
    register:   make(chan *User),
    unregister: make(chan *User),
    broadcast:  make(chan []byte),
  }
}

func (hub *Hub) Start() {
  for {
    select {
	    case user := <-hub.register:
	      hub.registerUser(user)
	    case user := <-hub.unregister:
	      hub.unregisterUser(user)
	  }
	}
}

func (hub *Hub) registerUser(user *User) {
  hub.users[user] = true
}

func (hub *Hub) unregisterUser(user *User) {
  if _, ok := hub.users[user]; ok {
    delete(hub.users, user)
  }
}
