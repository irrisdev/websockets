package main

import "sync"

type History struct {
	messages []Message
	mu       sync.Mutex
	maxSize  int
}

func NewHistory(max int) *History {
	return &History{
		messages: make([]Message, 0),
		mu:       sync.Mutex{},
		maxSize:  max,
	}
}

func (h *History) AddMessage(m *Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.messages) >= h.maxSize {
		h.messages = h.messages[1:]
	}
	h.messages = append(h.messages, *m)
}

func (h *History) RetrieveMessages() []Message {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]Message(nil), h.messages...)
}
