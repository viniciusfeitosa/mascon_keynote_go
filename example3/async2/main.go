package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/nats-io/go-nats"
)

var (
	brokerClient *nats.EncodedConn
)

// User is an the domain struct
type User struct {
	ID        string
	FirstName string
	LastName  string
	Twitter   string
}

func init() {
	brokerConn, err := nats.Connect(os.Getenv("BROKER_URI"))
	if err != nil {
		log.Fatal(err)
	}
	brokerClient, err = nats.NewEncodedConn(brokerConn, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer brokerClient.Close()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	brokerClient.Subscribe("ex3_domain2", func(u *User) {
		log.Println("ex3_domain2", u)
		if u.ID == "1" {
			u.LastName = "Pacheco"
		} else {
			u.LastName = "Nadie"
		}
		brokerClient.Publish("ex3_domain3", u)
	})

	<-signalCh
}
