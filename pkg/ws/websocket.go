package ws

import (
	"encoding/json"
	"net/http"

	"anne-hub/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {

		return true
	},
}

type Client struct {
	Conn *websocket.Conn
}

func UpgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*Client, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &Client{Conn: conn}, nil
}

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

func (c *Client) WriteMessage(data []byte) error {
	return c.Conn.WriteMessage(websocket.TextMessage, data)
}
