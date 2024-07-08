package websocket

import "github.com/gofiber/contrib/websocket"

type miniClient map[string]*websocket.Conn // Modified type

type ClientObject struct {
	USER string
	conn *websocket.Conn
}

type BroadcastObject struct {
	MSG  string
	FROM ClientObject
}

var clients = make(miniClient) // Initialized as a nested map
var register = make(chan ClientObject)
var broadcast = make(chan BroadcastObject)
var unregister = make(chan ClientObject)
