package nats

import (
	"log"

	"github.com/nats-io/stan.go"
)

func ConnectStan(clientID string) stan.Conn {
	clusterID := "test-cluster"
	url := "nats://127.0.0.1:4222"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url),
		stan.Pings(1, 3),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, url)
	}

	log.Println("Connected Nats")

	return sc
}
