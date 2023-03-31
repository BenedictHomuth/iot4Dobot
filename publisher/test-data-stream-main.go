package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	nats "github.com/nats-io/nats.go"
)

type PubMessage struct {
	Data string `json:"data"`
}

type RoboEvent struct {
	X      float64      `json:"x"`
	Y      float64      `json:"y"`
	Z      float64      `json:"z"`
	R      float64      `json:"r"`
	Angles JointsAngles `json:"jointAngle"`
}

type JointsAngles struct {
	Rotations [4]float64
}

func randFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func main() {
	// nc, err := nats.Connect(nats.DefaultURL)
	nc, err := nats.Connect("159.69.125.114:4222")
	if err != nil {
		fmt.Println("Oh no. NATS server not found")
	}
	defer nc.Close()

	index := 0
	var msg PubMessage
	var armPos RoboEvent
	for {
		armPos.X = 0.0
		armPos.Y = 0.0
		armPos.Z = randFloats(0.0, 400.0)
		armPos.R = 0.0
		armPos.Angles.Rotations[0] = 0.0
		armPos.Angles.Rotations[1] = randFloats(0.0, 90.0)
		armPos.Angles.Rotations[2] = 0.0
		armPos.Angles.Rotations[3] = randFloats(0, 45.0)

		msg.Data = strconv.Itoa(index)
		jmsg, err := json.Marshal(armPos)
		if err != nil {
			fmt.Println("Could not marshal message to publish!")
		}
		nc.Publish("roboPos", []byte(jmsg))
		fmt.Println("published: " + fmt.Sprintf("%f", armPos))
		time.Sleep(500 * time.Millisecond)
		index++
	}

}
