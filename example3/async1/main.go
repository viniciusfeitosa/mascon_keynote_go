package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/nats-io/go-nats"
)

const locationString = "http://localhost:8082/async/%s/data"

var (
	redisClient  *redis.Client
	brokerConn   *nats.Conn
	brokerClient *nats.EncodedConn
)

// User is an the domain struct
type User struct {
	ID        string
	FirstName string
	LastName  string
	Twitter   string
}

func input(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	u := &User{ID: vars["id"]}
	if err := brokerClient.Publish("ex3_domain1", u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	location := fmt.Sprintf(locationString, u.ID)

	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusAccepted)
}

func output(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	data, err := redisClient.Get(fmt.Sprintf("ex3_%s", vars["id"])).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func domain(ctx context.Context) {
	brokerClient.Subscribe("ex3_domain1", func(u *User) {
		log.Println("ex3_domain1", u)
		if u.ID == "1" {
			u.FirstName = "Vinicius"
		} else {
			u.FirstName = "Juan"
		}
		brokerClient.Publish("ex3_domain2", u)
	})
	select {
	case <-ctx.Done():
		log.Println("Exiting ex3_domain1")
		os.Exit(0)
	}
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

	c1, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-signalCh:
			cancel()
			return
		}
	}()
	go domain(c1)

	r := mux.NewRouter()
	r.HandleFunc("/async/{id:[0-9]+}", input)
	r.HandleFunc("/async/{id:[0-9]+}/data", output)

	if err := http.ListenAndServe(":5000", r); err != nil {
		log.Fatal(err)
	}
}
