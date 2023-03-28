package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

type RoboEvent struct {
	X      string       `json:"x"`
	Y      string       `json:"y"`
	Z      string       `json:"z"`
	R      string       `json:"r"`
	Angles JointsAngles `json:"jointAngle"`
}

type JointsAngles struct {
	Rotations [4]string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func main() {
	fmt.Println("Starting server on port :8080")
	// Connect to NATS message queue
	// nc, err := nats.Connect(nats.DefaultURL)
	// nc, err := nats.Connect("91.107.199.56:4222")
	nc, err := nats.Connect("65.109.172.100:4222")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully found nats server")
	}
	defer nc.Close()

	// Health endpoint
	http.HandleFunc("/health", healthHandler)

	// Set up websocket handler
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		ch := make(chan *nats.Msg, 64)
		sub, err := nc.ChanSubscribe("roboPos", ch)
		if err != nil {
			panic(err)
		}

		for event := range ch {
			var posEvent RoboEvent
			err = json.Unmarshal(event.Data, &posEvent)
			if err != nil {
				panic(err)
			}
			fmt.Println(posEvent)
			time.Sleep(500 * time.Millisecond)

			// Send message data to websocket client
			err = conn.WriteJSON(posEvent)
			if err != nil {
				fmt.Println("Error while writing msg to websocket.")
			}
		}
		// Unsubscribe if needed
		sub.Unsubscribe()
		close(ch)
	})

	// Host static files
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	// Start HTTP server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
