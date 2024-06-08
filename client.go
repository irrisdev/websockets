package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	room *Room

	receiver chan *Message

	conn *websocket.Conn

	username string
}

const (
	response = `<div hx-swap-oob="beforeend" id="chat_room"><p class="text-gray-100">%s: %s</p></div>`
)

func (c *Client) receive() {

	defer c.disconnect()

	for {

		select {
		// Replying only on msg and ignore ok means unable to
		// determine whether the channel has been closed or not,
		// and might mistakenly process a zero value received from a closed channel as a valid message.
		case msg, ok := <-c.receiver:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(msg.MessageType)
			if err != nil {
				return
			}

			w.Write([]byte(fmt.Sprintf(response, msg.From, msg.Message)))

			n := len(c.receiver)
			for i := 0; i < n; i++ {
				queued := <-c.receiver
				w.Write([]byte(fmt.Sprintf(response, queued.From, queued.Message)))
			}

			if err := w.Close(); err != nil {
				return
			}

		}

	}

}

func (c *Client) broadcast() {
	defer c.disconnect()

	var m *Message

	for {

		t, buf, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			return
		}

		if err := json.Unmarshal(buf, &m); err != nil {
			log.Println(err)
			return
		}

		m.MessageType = t
		m.From = c.username

		c.room.Receiver <- m
	}

}

func (c *Client) disconnect() {
	c.room.RemoveClient <- c
	c.conn.Close()
	log.Println("Client disconnected: ", c.username)
}

func NewClient(room *Room, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error: ", err)
		return
	}

	client := &Client{
		room:     room,
		receiver: make(chan *Message),
		conn:     conn,
		username: r.URL.Query().Get("username"),
	}
	room.AddClient <- client
	log.Println("New client connected: ", client.username)

	go client.receive()
	go client.broadcast()

}
