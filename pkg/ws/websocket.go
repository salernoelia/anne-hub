package ws

import (
	"encoding/json"
	"net/http"

	"anne-hub/models"

	"github.com/gorilla/websocket"
)

// Define a global upgrader with an appropriate CheckOrigin function
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Adjust this as needed (restrict domain, etc.)
        return true
    },
}

// Client represents a single WebSocket client
type Client struct {
    Conn *websocket.Conn
}

// UpgradeToWebSocket upgrades an HTTP connection to a WebSocket connection.
func UpgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*Client, error) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return nil, err
    }
    return &Client{Conn: conn}, nil
}

// ReadAnneWearRequest reads a single message from the WebSocket and unmarshals it into AnneWearConversationRequest.
func (c *Client) ReadAnneWearRequest() (*models.AnneWearConversationRequest, error) {
    _, message, err := c.Conn.ReadMessage()
    if err != nil {
        return nil, err
    }

    var req models.AnneWearConversationRequest
    if err := json.Unmarshal(message, &req); err != nil {
        return nil, err
    }

    return &req, nil
}

// WriteMessage sends a text message back to the client.
func (c *Client) WriteMessage(data []byte) error {
    return c.Conn.WriteMessage(websocket.TextMessage, data)
}
