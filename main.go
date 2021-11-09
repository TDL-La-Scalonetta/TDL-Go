package main

import (
    "flag"
    "log"
    "strconv"
    "net/http"
    "encoding/json"
)

var addr = flag.String("addr", ":8080", "http server address")
var rooms = make(map[*Room]bool)
var users = make(map[*User]bool)

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

    http.HandleFunc("/rooms", func (w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "application/json")
      list := make([]*Room, 0)
      for  room, _ := range rooms {
         list = append(list, room)
      }
      var data, err = json.Marshal(list)
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
      }
      w.Write(data)
    })

    http.HandleFunc("/rooms/new", func (w http.ResponseWriter, req *http.Request) {
      decoder := json.NewDecoder(req.Body)
      var roomParams *RoomParams

      err := decoder.Decode(&roomParams)
      if err != nil {
          panic(err)
      }

      user := findUser(roomParams.Owner)

      converted, e := strconv.ParseFloat(roomParams.Value, 16)
      if (e != nil) {
          panic(e)
      }
      newRoom := newRoom(*user, roomParams.Name, roomParams.Product, float32(converted))

      // start room server

      rooms[newRoom] = true
    })

    http.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request) {
      name, ok := r.URL.Query()["name"]

      if !ok || len(name[0]) < 1 {
          log.Println("Url Param 'name' is missing")
          return
      }

      conn, err := upgrader.Upgrade(w, r, nil)
      if err != nil {
        log.Println(err)
        return
      }

      user := findUser(name[0])
      user.conn = conn
      users[user] = true
    })

    log.Fatal(http.ListenAndServe(*addr, nil))
}

func findUser(name string) *User {
  for u, _ := range users {
    if u.Name == name {
        return u
    }
  }

  return newUser(name)
}

func fetchRooms(list *[]Room, w http.ResponseWriter, r *http.Request) {

}
