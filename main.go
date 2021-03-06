package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
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
	flag.Parse()
	fs := http.FileServer(http.Dir("./public"))

	http.Handle("/", fs)

	http.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		list := make([]*Room, 0)
		for room, _ := range rooms {
			list = append(list, room)
		}
		var data, err = json.Marshal(list)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})

	http.HandleFunc("/rooms/new", func(w http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var roomParams *RoomParams

		err := decoder.Decode(&roomParams)
		if err != nil {
			panic(err)
		}

		user := findUser(roomParams.Owner)

		converted, e := strconv.ParseFloat(roomParams.Value, 16)
		if e != nil {
			panic(e)
		}
		newRoom := newRoom(*user, roomParams.Name, roomParams.Product, float32(converted))

		rooms[newRoom] = true

		go newRoom.Run()
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		userName, ok1 := r.URL.Query()["name"]

		if !ok1 || len(userName[0]) < 1 {
			log.Println("Url Param 'name' is missing")
			return
		}

		roomName, ok2 := r.URL.Query()["room"]

		if !ok2 || len(roomName[0]) < 1 {
			log.Println("Url Param 'room' is missing")
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		user := findUser(userName[0])
		room := findRoom(roomName[0])

		if room == nil {
			log.Println("No hay room con ese nombre")
			return
		}

		user.Conn = conn
		user.Room = room
		room.register <- user

		go user.ReadSocket()
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func findUser(name string) *User {
	return newUser(name)
}

func findRoom(name string) *Room {
	for r, _ := range rooms {
		if r.Name == name {
			return r
		}
	}
	return nil
}
