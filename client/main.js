let ws = new WebSocket("ws://localhost:3001/ws")

ws.onmessage = (evt) => {
    console.log("From server: ", evt.data)
}
