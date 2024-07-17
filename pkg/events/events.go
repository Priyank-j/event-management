package events

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	kafka "github.com/segmentio/kafka-go"
	"k8s.io/apimachinery/pkg/util/json"
)

func eventWorker() {

	for {
		// Do a non-blocking read from a data source.
		select {
		case data := <-EventChan:
			// If read succeeds write it out. This will flush if the batch
			// exceeds the batch size.
			Batcher.Write(data)
		case <-Done:
			return
		default:
			// If read fails make sure to call Flush to ensure data doesn't
			// get stuck in the batch for long periods of time.
			Batcher.Flush()
		}
	}
}

func byteTokafkaMessage(batch [][]byte) []kafka.Message {
	var kafkaMessages []kafka.Message
	for _, data := range batch {
		kafkaMessages = append(kafkaMessages, kafka.Message{
			Value: data,
		})

	}
	return kafkaMessages
}

func TrigerEvent(event EventData) {

	ctx := context.Background()
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.WithContext(ctx).Errorf("[TrigerEvent] failed Marshal event. err: %v", err)
	}
	EventChan <- eventBytes

}

func WriteMessageToKafka(messages []kafka.Message) {
	ctx := context.Background()
	err := KafkaConn.WriteMessages(context.Background(), messages...)

	if err != nil {
		log.WithContext(ctx).Errorf("[WriteMessageToKafka] failed to write messages. err: %v", err)
	}
}
