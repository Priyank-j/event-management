package websocket

import (
	"fmt"
	"log"
)

func SocketHandler() {
	for {
		select {
		case client := <-register:

			clients[client.USER] = client.conn
			log.Println("client registered:", client.USER)

		case message := <-broadcast:
			//TODO put this message into loacal batch
			fmt.Println(message.MSG)

		case client := <-unregister:
			removeClient(client.USER) // Update client removal
			log.Println("client unregistered:", client.USER)
		}
	}
}
