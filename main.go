package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	nats "github.com/nats-io/nats.go"
)

type PubMessage struct {
	Data string `json:"data"`
}

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

func main() {
	// nc, err := nats.Connect(nats.DefaultURL)
	nc, err := nats.Connect("65.109.172.100:4222")
	if err != nil {
		fmt.Println("Oh no. NATS server not found")
	}
	defer nc.Close()

	index := 0
	var msg PubMessage
	var armPos RoboEvent
	for {
		armPos.X = strconv.Itoa(index)
		armPos.Y = strconv.Itoa(index + 1)
		armPos.Z = strconv.Itoa(index + 2)
		armPos.R = strconv.Itoa(index + 3)
		armPos.Angles.Rotations[0] = "Test #1"
		armPos.Angles.Rotations[1] = "Test #2"
		armPos.Angles.Rotations[2] = "Test #3"
		armPos.Angles.Rotations[3] = "Test #4"

		msg.Data = strconv.Itoa(index)
		jmsg, err := json.Marshal(armPos)
		if err != nil {
			fmt.Println("Could not marshal message to publish!")
		}
		nc.Publish("roboPos", []byte(jmsg))
		fmt.Println("published: " + msg.Data)
		time.Sleep(500 * time.Millisecond)
		index++
	}

}
