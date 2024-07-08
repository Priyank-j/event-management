package websocket

import (
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
		USER: c.Locals("USER").(string),
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

		if messageType == websocket.TextMessage {
			// Broadcast the received message
			broadcast <- BroadcastObject{
				MSG:  string(message),
				FROM: clientObj,
			}
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}
func EventRequestMiddleWare(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		// Your authentication process goes here. Get the Token from header and validate it
		// Extract the claims from the token and set them to the Locals
		// This is because you cannot access headers in the websocket.Conn object below
		c.Locals("USER", string(c.Request().Header.Peek("USER")))
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
