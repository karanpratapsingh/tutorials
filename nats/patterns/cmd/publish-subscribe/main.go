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
		log.Println("Subscriber 1:", string(msg.Data))
	})

	nc.Subscribe("foo", func(msg *nats.Msg) {
		log.Println("Subscriber 2:", string(msg.Data))
	})

	nc.Subscribe("foo", func(msg *nats.Msg) {
		log.Println("Subscriber 3:", string(msg.Data))
	})

	if err := nc.Publish("foo", []byte("Message")); err != nil {
		log.Fatalln(err)
	}

	time.Sleep(2 * time.Second)
}
