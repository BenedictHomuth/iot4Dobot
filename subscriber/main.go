package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

type RoboEvent struct {
	X           float64   `json:"x"`
	Y           float64   `json:"y"`
	Z           float64   `json:"z"`
	R           float64   `json:"r"`
	JointAngles []float64 `json:"jointAngles"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
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

	// Incoming Event, publish to NATS
	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		var event RoboEvent
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jEvent, err := json.Marshal(event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		nc.Publish("roboPos", jEvent)

		w.WriteHeader(http.StatusOK)
	})

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
