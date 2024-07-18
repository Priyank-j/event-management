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

// type EventMessage struct {
// 	EventType      string `json:"event_type"`
// 	SourceEntityID string `json:"source_entity_id"`
// 	DashboardType  string `json:"dashboard_type"`
// 	CreatedAt      string `json:"created_at"`
// 	UserID         string `json:"user_id"`
// }

type EventMessage struct {
	EventType      string       `json:"event_type"`
	User           string       `json:"user"`
	UserType       string       `json:"user_type"`
	Action         string       `json:"action"`
	Name           string       `json:"name"`
	ObjectType     string       `json:"object_type"`
	ActionBy       string       `json:"action_by"`
	Timestamp      string       `json:"timestamp"`
	LoanMetaData   LoanMetaData `json:"loan_meta_data"`
	Screen         string       `json:"screen"`
	Component      string       `json:"component"`
	ElementData    string       `json:"element_data"`
	ActionDetails  string       `json:"action_details"`
	SessionDetails string       `json:"session_details"`
	Source         string       `json:"source"`
}

type LoanMetaData struct {
	LoanApplicationId string `json:"loan_application_id"`
	CustomerID        string `json:"customer_id"`
	Program           string `json:"program"`
	Status            string `json:"status"`
}
