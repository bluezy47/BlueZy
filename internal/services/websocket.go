// W E B S O C K E T   S E R V I C E
package services
import (
	"fmt"
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
)
//
// User struct
type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Conn *websocket.Conn // websocket connection
}
//
// map to store all the users
var ConnectedUsers = map[string]*User{};
//
// struct for the Message
type Message struct {
	Content string `json:"content"`
	To      string `json:"to"`
};

func WebsocketHandler(c *websocket.Conn) {
	// get the user email from the query params
	email := c.Query("email")
	// todo: for now create a dummy user
	user := &User{
		ID: "1",
		Username: "bluezy",
		Email: email,
		Conn: c,
	}
	//
	// add the user to the connected users
	ConnectedUsers[email] = user;
	fmt.Println("User connected: ", email);
	//
	// listen for the messages
	for {
		messageType, message, err := c.ReadMessage();
		if err != nil {
			fmt.Println("Error reading message: ", err);
			return;
		}
		//
		// Handle the message
		if messageType == websocket.TextMessage {
			var msg Message;
			if err := json.Unmarshal(message, &msg); err != nil {
				fmt.Println("Error unmarshalling message: ", err);
				return;
			}
			//
			// Check if the user is connected
			if user, ok := ConnectedUsers[msg.To]; ok {
				err := user.Conn.WriteMessage(websocket.TextMessage, message);
				if err != nil {
					fmt.Println("Error writing message: ", err);
					continue;
				};
				fmt.Println("Message sent to: ", msg.To);
			} else {
				fmt.Println("User not connected: ", msg.To);
				// todo: add the message to the offline messages...
			}
		};
	}
}
