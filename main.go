package main

import (
	"flag"
	"go-event-management/conf"
	internalWebsocket "go-event-management/internal/http/websocket"
	"go-event-management/internal/repository/redis"
	"go-event-management/pkg/events"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {

	startWebsocketServer()

}

func startWebsocketServer() {
	defer close(events.Done)

	app := fiber.New()
	app.Use(limiter.New(limiter.Config{
		Expiration: time.Second,
		Max:        1000,
	}))

	// init redis
	enableSSL, _ := conf.RedisConf["SSL"].(bool)
	endpoint, _ := conf.RedisConf["Addr"].(string)
	replicaEndpoint, _ := conf.RedisConf["ReplicaAddr"].(string)

	redis.Init(enableSSL, endpoint, replicaEndpoint)

	go events.InitEvents()

	app.Use("/event", internalWebsocket.EventRequestMiddleWare)
	go internalWebsocket.SocketHandler()

	app.Get("/event", websocket.New(internalWebsocket.EventCont))

	addr := flag.String("addr", ":3332", "http service address")
	flag.Parse()
	app.Listen(*addr)
}
