let app = new Vue({
  el: '#app',
  http: { root: '/' },
  data() {
    return {
      ws: null,
      serverUrl: "ws://localhost:8080/ws",
      roomInput: null,
      rooms: [],
      activeRooms: [],
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
      this.connect()
    },
    fetchRooms() {
      this.$http.get("/rooms").then(res => {
        this.rooms = res.data
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