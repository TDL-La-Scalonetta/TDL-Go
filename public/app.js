let app = new Vue({
  el: '#app',
  http: { root: '/' },
  data() {
    return {
      ws: null,
      serverUrl: "ws://localhost:8080/ws",
      roomInput: null,
      newRoom: {
        name: '',
        product: '',
        value: ''
      },
      rooms: [],
      activeRooms: [],
      user: { name: null },
      authenticated: false
    }
  },
  methods: {
    setUser() {
      this.authenticated = true
      this.connect()
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
        console.log('Success')
      })
    },
    connect() {
      this.ws = new WebSocket(`${this.serverUrl}?name=${this.user.name}`);

      this.ws.addEventListener('open', (event) => {
        this.onWebsocketOpen(event)
      });

      this.ws.addEventListener('message', (event) => {
        this.handleNewMessage(event)
      });
    },
    onWebsocketOpen() {
      console.log("connected to WS!");
    },
    handleNewMessage(event) {
      let data = event.data;
      // check data type (accepted, refused, roomUpdate)
      console.log('Data received', data)
    },
    sendMessage(room) {
      // send message to correct room.
      if (room.newMessage !== "") {
        this.ws.send(JSON.stringify({
          action: 'send-offer',
          offer: room.offer,
          target: room.name
        }));
        room.offer = "";
      }
    },
    joinRoom(room) {
      console.log("Join to room", room)
      console.log(this.ws)
      this.ws.send(JSON.stringify({
        action: 'join-room',
        message: room.name
      }));
    },
    leaveRoom(room) {
      this.ws.send(JSON.stringify({
        action: 'leave-room',
        message: room.name
      }));
    }
  }
})
