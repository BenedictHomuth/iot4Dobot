const socket = new WebSocket("wss://localhost:8080/websocket")
sendBtn = document.querySelector("#test")

sendBtn.addEventListener("click", () =>{
    socket.send("Hello World")
})

socket.addEventListener("open", () => {
    console.log("Client -> Server WS connection established!")
})

socket.onmessage = ({data}) => {
    msg = JSON.parse(data)
    console.log(msg.data);
}


