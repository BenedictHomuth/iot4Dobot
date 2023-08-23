package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

type RoboEvent struct {
	ID          uuid.UUID `json:"id,omitempty"`
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
	useTLS := getEnv("USE_TLS", "false")
	dbUser := getEnv("DB_USER", "root")
	// dbPassword := getEnv("DB_PW", "123456")
	dbHostname := getEnv("DB_HOSTNAME", "localhost")
	dbPort := getEnv("DB_PORT", "26257")
	// dbDatabaseName := getEnv("DB_Name", "robo_events")

	// Connect to NATS message queue
	nc, err := nats.Connect(natsHost)
	if err != nil {
		fmt.Println("Error while connecting to nats!", err)
	} else {
		fmt.Println("Successfully found nats server")
	}
	defer nc.Close()

	// Database stuff -v1
	// Connect to the CockroachDB database
	// conString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=insecure", dbUser, dbPassword, dbHostname, dbPort, dbDatabaseName)
	conString := fmt.Sprintf("postgresql://%s@%s:%s?sslmode=disable", dbUser, dbHostname, dbPort)
	// db, err := sql.Open("postgres", "postgresql://<username>:<password>@<hostname>:<port>/<database>?sslmode=insecure")
	db, err := sql.Open("postgres", conString)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	} else {
		fmt.Println("Successfully connected to cockroachdb!")
	}
	defer db.Close()

	// Create the table if it doesn't exist
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS robo_events (
		id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		x FLOAT,
		y FLOAT,
		z FLOAT,
		r FLOAT,
		j1 FLOAT,
		j2 FLOAT,
		j3 FLOAT,
		j4 FLOAT
	)`); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Define the HTTP endpoint

	// DB Endpoints
	http.HandleFunc("/createEvent", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var event RoboEvent
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		// Insert the RoboEvent into the database
		_, err = db.Exec("INSERT INTO robo_events (x, y, z, r, j1,j2,j3,j4) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			event.X, event.Y, event.Z, event.R, event.JointAngles[0], event.JointAngles[1], event.JointAngles[2], event.JointAngles[3])
		if err != nil {
			log.Fatal("Failed to insert data into the database:", err)
		}

		// Send a success response
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Successfully inserted event into the database")
	})

	// Define the HTTP endpoint
	http.HandleFunc("/getEvents", func(w http.ResponseWriter, r *http.Request) {
		// Query all events from the database
		rows, err := db.Query("SELECT * FROM robo_events")
		if err != nil {
			log.Fatal("Failed to query data from the database:", err)
		}
		defer rows.Close()

		// Create a slice to store the events
		var events []RoboEvent

		// Iterate over the rows and populate the events slice
		for rows.Next() {
			var event RoboEvent
			var j1, j2, j3, j4 float64
			var id uuid.UUID

			err := rows.Scan(&id, &event.X, &event.Y, &event.Z, &event.R, &j1, &j2, &j3, &j4)
			event.JointAngles = []float64{j1, j2, j3, j4}
			event.ID = id
			if err != nil {
				log.Fatal("Failed to scan row data:", err)
			}
			events = append(events, event)
		}
		if err = rows.Err(); err != nil {
			log.Fatal("Error occurred while iterating over rows:", err)
		}

		// Convert the events slice to JSON
		jsonData, err := json.Marshal(events)
		if err != nil {
			log.Fatal("Failed to marshal events to JSON:", err)
		}

		// Send the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

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
	if useTLS == "false" {
		err = http.ListenAndServe(":8080", nil)
	} else {
		err = http.ListenAndServeTLS(":8080", "domain.crt", "domain.key", nil)
	}
	if err != nil {
		fmt.Println("Error while serving the api!", err)
	}
}
