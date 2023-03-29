package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func wsConncetionString(w http.ResponseWriter, r *http.Request) {

}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	fmt.Println("Starting server on port :8080")

	// Getting environment vars
	natsHost := getEnv("NATS_CONN_STRING", nats.DefaultURL)

	// Connect to NATS message queue
	nc, err := nats.Connect(natsHost)
	if err != nil {
		fmt.Println("Error while connecting to nats!", err)
	} else {
		fmt.Println("Successfully found nats server")
	}
	defer nc.Close()

	// Health endpoint
	http.HandleFunc("/health", healthHandler)

	// Set up websocket handler
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error while connecting the ws conn!", err)
		}
		defer conn.Close()

		ch := make(chan *nats.Msg, 64)
		sub, err := nc.ChanSubscribe("roboPos", ch)
		if err != nil {
			fmt.Println("Error while creating nats channel subscription!", err)
		}

		for event := range ch {
			var posEvent RoboEvent
			err = json.Unmarshal(event.Data, &posEvent)
			if err != nil {
				fmt.Println("Error while unmarshaling a event message", err)
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
	err = http.ListenAndServeTLS(":8080", "domain.crt", "domain.key", nil)
	if err != nil {
		fmt.Println("Error while serving the api!", err)
	}
}
