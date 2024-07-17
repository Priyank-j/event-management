package events

import (
	"bytes"
	"fmt"
	"time"

	"code.cloudfoundry.org/go-batching"
	kafka "github.com/segmentio/kafka-go"
)

func InitEvents() {
	go eventWorker()

	writer := batching.ByteWriterFunc(func(batch [][]byte) {
		result := bytes.Join(batch, nil)

		fmt.Printf("Inside writer %s\n", result)
		messages := byteTokafkaMessage(batch)
		go WriteMessageToKafka(messages)
	})

	//TODO:  move size and interval to constant
	Batcher = batching.NewByteBatcher(100, time.Second*5, writer)
	EventChan = make(chan []byte, 100)
	Done = make(chan struct{})

	topic := "quickstart-events"
	KafkaConn = &kafka.Writer{
		// TODO: Add server urls in Config/constants

		Addr:     kafka.TCP("localhost:9092"),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

}
