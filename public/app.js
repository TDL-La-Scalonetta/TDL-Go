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
        Finished: false,
        TimeLeft: "0",
        Product: null,
        Amount: 0,
        Users: []
      },
      user: { name: null },
      authenticated: false
    }
  },
  computed: {
    canOffer() {
      return this.currentRoom.Started && !this.currentRoom.Finished
    },
    isOwner() {
      return this.user.name == this.currentRoom.Owner
    }
  },
  methods: {
    fetchRooms() {
      this.$http.get("/rooms").then(res => {
        this.rooms = res.data
      })
    },
    setUser() {
      this.authenticated = true
      this.fetchRooms()
    },
    createRoom() {
      this.$http.post("/rooms/new", {
        "Owner": this.user.name,
        "Name": this.newRoom.name,
        "Product": this.newRoom.product,
        "Value": this.newRoom.value.toString()
      }).then(res => {
        this.newRoom = { name: '', product: '', value: '' }
        this.fetchRooms()
      })
    },
    connect(room) {
      try {
        this.ws = new WebSocket(`${this.serverUrl}?room=${room.Name}&name=${this.user.name}`);

        this.ws.addEventListener('message', (event) => {
          this.handleNewMessage(event)
        });

        this.ws.addEventListener('open', (event) => {
          console.log("Connected to WS!", event);
        });

        this.ws.addEventListener('error', (event) => {
          console.log("Error on WS!", event);
        });

        this.ws.addEventListener('close', (event) => {
          console.log("Disconnected from WS!", event);
          this.ws.close(1000);
        });
      } catch(error) {
        console.log('Error al conectarse')
      }
    },
    finish() {
      let m = {
        action: 'finish',
        room: this.currentRoom.Name
      }
      this.push(m)
    },
    disconnect() {
      this.ws.close(1000)
      this.currentRoom = {
        Name: null,
        Owner: null,
        Winner: null,
        BaseValue: 0,
        Started: false,
        Finished: false,
        TimeLeft: "0",
        Product: null,
        Amount: 0,
        Users: []
      }
    },
    handleNewMessage(event) {
      let data = JSON.parse(event.data);

      if (data.Type != "Time") console.log(data)

      this.currentRoom.Owner =  data.Owner
      this.currentRoom.Name = data.Name
      this.currentRoom.Winner = data.Winner
      this.currentRoom.Product = data.Product
      this.currentRoom.BaseValue = data.BaseValue
      this.currentRoom.Started = data.Started
      this.currentRoom.Finished = data.Finished
      this.currentRoom.TimeLeft = data.TimeLeft
      this.currentRoom.Amount = data.Amount
      this.currentRoom.Users = data.Users

      if (data.Type == 'Ended') this.ws.close(1000)
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
        let m = {
          action: 'offer',
          offer: this.offer.toString(),
          room: this.currentRoom.Name
        }
        this.push(m)
        this.offer = "";
      }
    },
    push(message) {
      this.ws.send(JSON.stringify(message));
    }
  }
})
