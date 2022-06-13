package nats

import (
	"log"

	"github.com/nats-io/stan.go"
)

func GetData(subject, qgroup, durable string, sc stan.Conn, q chan []byte) {

	mcb := func(msg *stan.Msg) {
		if err := msg.Ack(); err != nil {
			log.Printf("failed to ACK msg:%v", err)
		}
		q <- []byte(msg.Data)
	}

	_, err := sc.QueueSubscribe(subject,
		qgroup, mcb,
		stan.DeliverAllAvailable(),
		stan.SetManualAckMode(),
		stan.DurableName(durable))
	if err != nil {
		log.Println(err)
	}

}
