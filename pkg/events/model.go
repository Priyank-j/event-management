package events

import (
	"code.cloudfoundry.org/go-batching"
	kafka "github.com/segmentio/kafka-go"
)

var (
	Batcher   *batching.ByteBatcher
	EventChan chan []byte
	Done      chan struct{}
	KafkaConn *kafka.Writer
)

type EventData struct {
	EventType      string `json:"event_type"`
	SourceEntityID string `json:"source_entity_id"`
	DashboardType  string `json:"dashboard_type"`
	CreatedAt      string `json:"created_at"`
	UserID         string `json:"user_id"`
}
