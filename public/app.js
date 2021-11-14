let app = new Vue({
  el: '#app',
  http: { root: '/' },
  data() {
    return {
      ws: null,
      serverUrl: "ws://localhost:8080/ws",
      newRoom: {
        name: '',
        product: '',
        value: ''
      },
      offer: 0.0,
      rooms: [],
      currentRoom: {
        Name: null,
        Owner: null,
        Winner: null,
        BaseValue: 0,
        Started: false,
        Product: null,
        Amount: 0,
        Users: []
      },
      user: { name: null },
      authenticated: false
    }
  },
  created() {
    this.fetchRooms()
  },
  methods: {
    setUser() {
      this.authenticated = true
    },
    fetchRooms() {
      this.$http.get("/rooms").then(res => {
        this.rooms = res.data
      })
    },
    createRoom() {
      this.$http.post("/rooms/new", {
        "Owner": this.user.name,
        "Name": this.newRoom.name,
        "Product": this.newRoom.product,
        "Value": this.newRoom.value.toString()
      }).then(res => {
        this.newRoom = {
          name: '',
          product: '',
          value: ''
        }
        this.fetchRooms()
      })
    },
    connect(room) {
      try {
        this.ws = new WebSocket(`${this.serverUrl}?room=${room.Name}&name=${this.user.name}`);

        this.ws.addEventListener('open', (event) => {
          this.onWebsocketOpen(event)
        });

        this.ws.addEventListener('message', (event) => {
          this.handleNewMessage(event)
        });

      } catch(error) {
        console.log('Error al conectarse')
      }
    },
    disconnect() {
      console.log("Salir")
    },
    onWebsocketOpen() {
      console.log("connected to WS!");
    },
    handleNewMessage(event) {
      let data = JSON.parse(event.data);

      console.log('Data received', data)

      this.currentRoom.Owner =  data.Owner
      this.currentRoom.Name = data.Name
      this.currentRoom.Winner = data.Winner
      this.currentRoom.Product = data.Product
      this.currentRoom.BaseValue = data.BaseValue
      this.currentRoom.Started = data.Started
      this.currentRoom.Amount = data.Amount
      this.currentRoom.Users = data.Users
    },
    sendMessageX() {
      this.offer = Number(this.currentRoom.BaseValue) + 10
      this.sendMessage()
    },
    sendMessageXX() {
      this.offer = Number(this.currentRoom.BaseValue) + 20
      this.sendMessage()
    },
    sendMessageXXX() {
      this.offer = Number(this.currentRoom.BaseValue) + 30
      this.sendMessage()
    },
    sendMessage() {
      if (this.offer !== "") {
        this.ws.send(JSON.stringify({
          action: 'offer',
          offer: this.offer.toString(),
          room: this.currentRoom.Name
        }));
        this.offer = "";
      }
    }
  }
})
