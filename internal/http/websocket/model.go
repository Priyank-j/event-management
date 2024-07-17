package websocket

import (
	"go-event-management/pkg/events"

	"github.com/gofiber/contrib/websocket"
)

type miniClient map[string]*websocket.Conn // Modified type

type ClientObject struct {
	user string
	conn *websocket.Conn
}

type BroadcastObject struct {
	eventData events.EventData
	msg       string
	from      ClientObject
}

var clients = make(miniClient) // Initialized as a nested map
var register = make(chan ClientObject)
var broadcast = make(chan BroadcastObject)
var unregister = make(chan ClientObject)
