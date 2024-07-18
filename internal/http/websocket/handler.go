package websocket

import (
	"fmt"
	"go-event-management/pkg/events"
	"log"
)

func SocketHandler() {
	for {
		select {
		case client := <-register:

			clients[client.user] = client.conn
			log.Println("client registered:", client.user)

		case message := <-broadcast:
			fmt.Println(message.msg)
			events.TrigerEvent(message.EventMessage)

		case client := <-unregister:
			removeClient(client.user) // Update client removal
			log.Println("client unregistered:", client.user)
		}
	}
}
