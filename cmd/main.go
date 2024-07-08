package main

import (
	internalWebsocket "clickhouse/internal/websocket"
	"flag"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {

	startWebsocketServer()

}

func startWebsocketServer() {

	app := fiber.New()

	app.Use("/event", internalWebsocket.EventRequestMiddleWare)
	go internalWebsocket.SocketHandler()

	app.Get("/event", websocket.New(internalWebsocket.EventCont))

	addr := flag.String("addr", ":3000", "http service address")
	flag.Parse()
	app.Listen(*addr)
}
