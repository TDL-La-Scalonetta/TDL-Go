package main

type Message struct {
  Action  string    `json:"action"`
  Value   string    `json:"value"`
  Room    string    `json:"room"`
  User    *User
}
