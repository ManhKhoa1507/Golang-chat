var app = new Vue({
  el: '#app',
  data: {
    ws: null,
    serverUrl: "ws://localhost:8080/ws",
    messages: [],
    newMessage: ""
  },

  mounted: function () {
    this.connectToWebsocket();
  },

  methods: {
    connect() {
      this.connectToWebsocket();
    },

    // Connect to websocket
    connectToWebsocket() {
      this.ws = new WebSocket(this.serverUrl + "?name=" + this.user.name);
      this.ws.addEventListener('open', (event) => { this.onWebsocketOpen(event) });
      this.ws.addEventListener('message', (event) => { this.handleNewMessage(event) });
    },

    // Notify to console
    onWebsocketOpen() {
      console.log("connected to WS!");
    },

    // handle new message 
    handleNewMessage(event) {
      let data = event.data;
      data = data.split(/\r?\n/);

      for (let i = 0; i < data.length; i++) {
        let msg = JSON.parse(data[i]);
        // display message in correct room 
        const room = this.findRoom(msg.target);

      }
    },

    // send new message to correct room
    sendMessage(room) {
      if (room.newMessage !== "") {
        this.ws.send(JSON.stringify({
          action: 'send-message',
          message: room.newMessage,
          target: room.name
        }));
        room.newMessage = "";
      }
    },

    findRoom(roomName) {
      for (let i = 0; i < this.rooms.length; i++) {
        if (this.room[i].name == roomName) {
          return this.room[i];
        }
      }
    },

    // joinRoom action
    joinRoom() {

      // send JSON message
      this.ws.send(JSON.stringify({
        action: 'join-room',
        message: this.roomInput
      }));

      this.message = [];

      this.room.push({
        "name": this.roomInput,
        "message": []
      });

      this.roomInput = "";
    },

    leaveRoom(room){
      this.ws.send(JSON.stringify({
        action: 'leave-room',
        message: room.name
      }));

      for (let i = 0; i < this.rooms.length; i++){
        if(this.rooms[i].name == room.name){
          this.rooms.splice(i, 1);
          break;
        }
      }
    }
  }
})