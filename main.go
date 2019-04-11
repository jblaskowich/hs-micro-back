package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	nats "github.com/nats-io/go-nats"
)

var (
	natsURL  = "demo.nats.io"
	natsPORT = ":4222"
)

func getEvents() {
	natsChan := os.Getenv("NATSCHAN")
	if natsChan == "" {
		natsChan = "zjnO12CgNkHD0IsuGd89zA"
	}

	nc, err := nats.Connect(natsURL + natsPORT)
	if err != nil {
		log.Println(err.Error())
	}
	nc.Subscribe(natsChan, func(m *nats.Msg) {
		fmt.Printf("Received message: %s\n", string(m.Data))
	})
}

func main() {
	getEvents()
	port := os.Getenv("HS-MICRO-BACK")
	if port == "" {
		port = ":9090"
	}
	rtr := mux.NewRouter()
	http.Handle("/", rtr)
	http.ListenAndServe(port, nil)
}
