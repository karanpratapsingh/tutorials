package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatalln(err)
	}

	defer nc.Close()

	nc.Subscribe("foo", func(msg *nats.Msg) {
		log.Println("Request received:", string(msg.Data))

		nc.Publish(msg.Reply, []byte("I'm here"))
		msg.Respond([]byte("I'm here"))
	})

	reply, err := nc.Request("foo", []byte("I need help"), 4*time.Second)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Got Reply:", string(reply.Data))
}
