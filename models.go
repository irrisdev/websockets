package main

type Message struct {
	From        string
	Message     string `json:"message"`
	MessageType int
}

type PageData struct {
	Username string
	ChatData []Message
}
