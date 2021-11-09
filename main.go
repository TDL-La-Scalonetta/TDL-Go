package main

import (
    "flag"
    "log"
    "net/http"
    "encoding/json"
)

var addr = flag.String("addr", ":8080", "http server address")

func main() {
    flag.Parse()
    fs := http.FileServer(http.Dir("./public"))
    hub := NewHub()
    go hub.Start()

    http.Handle("/", fs)
    http.HandleFunc("/rooms", func (w http.ResponseWriter, r *http.Request) {
      fetchRooms(hub, w, r)
    })
    http.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request) {
      UserConn(hub, w, r)
    })

    log.Fatal(http.ListenAndServe(*addr, nil))
}

func fetchRooms(hub *Hub, w http.ResponseWriter, r *http.Request) {
    list := make([]Room, 0)
    if (len(hub.rooms) == 0) { 
      var room1 = newRoom("Tecnologia", "Iphone 7", 13000.99)
      var room2 = newRoom("Ropa", "Campera Harley", 9000.99)

      hub.rooms[room1] = true
      hub.rooms[room2] = true
    }

    w.Header().Set("Content-Type", "application/json")

    for  room, _ := range hub.rooms {
       list = append(list, *room)
    }

    var data, err = json.Marshal(list)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    w.Write(data)
}