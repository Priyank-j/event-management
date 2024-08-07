package websocket

import (
	"encoding/json"
	"fmt"
	"go-event-management/pkg/events"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func removeClient(user string) {
	if conn, ok := clients[user]; ok { // Check if client exists
		delete(clients, user)
		conn.Close()
	}
}

func EventCont(c *websocket.Conn) {
	clientObj := ClientObject{
		user: c.Locals("user").(string),
		conn: c,
	}
	defer func() {
		unregister <- clientObj
		c.Close()
	}()

	// Register the client
	register <- clientObj

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}

			return // Calls the deferred function, i.e. closes the connection on error
		}
		var EventMessage events.EventMessage
		fmt.Println(string(message))
		err = json.Unmarshal(message, &EventMessage)
		if err != nil {
			log.Println("can not Unmarshal message")
		}
		if messageType == websocket.TextMessage {
			// Broadcast the received message
			broadcast <- BroadcastObject{
				EventMessage: EventMessage,
				from:         clientObj,
			}
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}
func EventRequestMiddleWare(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		fmt.Println("here")
		// Your authentication process goes here. Get the Token from header and validate it
		// Extract the claims from the token and set them to the Locals
		// This is because you cannot access headers in the websocket.Conn object below
		c.Locals("user", string(c.Request().Header.Peek("user")))
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
