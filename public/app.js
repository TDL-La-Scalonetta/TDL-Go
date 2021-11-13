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
      currentRoom: {},
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
        "Value": this.newRoom.value
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
        this.currentRoom = room
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
      let data = event.data;
      // check data type (accepted, refused, roomUpdate)
      console.log('Data received', data)
    },
    sendMessage() {
      // send message to correct room.
      if (this.offer !== "") {
        this.ws.send(JSON.stringify({
          action: 'offer',
          offer: this.offer,
          room: this.currentRoom.Name
        }));
        this.offer = "";
      }
    }
  }
})
