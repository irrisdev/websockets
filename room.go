package main

type Room struct {
	Clients      map[*Client]struct{}
	Receiver     chan *Message
	AddClient    chan *Client
	RemoveClient chan *Client
	History      *History
}

func NewRoom(max int) *Room {
	return &Room{
		Clients:      make(map[*Client]struct{}),
		Receiver:     make(chan *Message),
		AddClient:    make(chan *Client),
		RemoveClient: make(chan *Client),
		History:      NewHistory(max),
	}
}

func (r *Room) start() {

	for {
		select {
		case buf := <-r.Receiver:
			r.History.AddMessage(buf)
			for client := range r.Clients {
				select {
				case client.receiver <- buf:
				default:
					close(client.receiver)
					delete(r.Clients, client)
				}

			}
		case add := <-r.AddClient:
			r.Clients[add] = struct{}{}
		case remove := <-r.RemoveClient:
			if _, exist := r.Clients[remove]; exist {
				delete(r.Clients, remove)
				close(remove.receiver)
			}
		}
	}

}
