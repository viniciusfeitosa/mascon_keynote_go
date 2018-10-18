package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis"
	"github.com/nats-io/go-nats"
)

var (
	redisClient  *redis.Client
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
	redisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URI"),
		DB:   0,
	})
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

	brokerClient.Subscribe("ex3_domain3", func(u *User) {
		time.Sleep(50 * time.Millisecond)
		log.Println("ex3_domain3", u)
		if u.ID == "1" {
			u.Twitter = "@viniciuspach"
		} else {
			u.Twitter = "@nose"
		}
		data, _ := json.Marshal(u)
		redisClient.Set(fmt.Sprintf("ex3_%s", u.ID), string(data), 0)
	})

	<-signalCh
}
