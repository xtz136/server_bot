package websocket

import (
	"fmt"
	"strings"
	"sync"
)

var hub *Hub
var onceHub sync.Once
var mutex sync.Mutex

type ClientFilter func(*Client) bool

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]struct{}

	// Inbound messages from the clients.
	observe map[chan *Client]struct{}

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	onceHub.Do(func() {
		hub = &Hub{
			observe:    make(map[chan *Client]struct{}),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]struct{}),
		}
	})
	return hub
}

func (h *Hub) CollectClientsIdentity() []string {
	result := []string{}
	dup := map[string]struct{}{}

	for client := range h.clients {
		ident := client.GetIdentity()
		if _, ok := dup[ident]; !ok {
			dup[ident] = struct{}{}
			result = append(result, ident)
		}
	}

	return result
}

func (h *Hub) CollectClientsName() map[string]string {
	result := map[string]string{}
	ca := NewClientAlias()
	// FIXME 需要优化
	for _, clientIdent := range h.CollectClientsIdentity() {
		match := false
		for name, values := range ca.Alias {
			for _, value := range values {
				if value == clientIdent {
					result[name] = strings.Join(values, ",")
					match = true
				}
			}
		}

		if !match {
			result[clientIdent] = clientIdent
		}
	}
	return result
}

func (h *Hub) HasClient(key string) bool {
	if !strings.Contains(key, ":") {
		key += ":default"
	}
	for client := range h.clients {
		ident := strings.Split(client.ip, ":")[0] + ":" + client.group
		if ident == key {
			return true
		}
	}

	return false
}

func (h *Hub) ListClients() map[*Client]struct{} {
	return h.clients
}

func (h *Hub) SendMessageToAll(command []byte) {
	h.SendMessage(command, func(c *Client) bool {
		return true
	})
}

// 给客户端发送消息，ClientFilter用于筛选哪些客户端能收到消息
func (h *Hub) SendMessage(command []byte, clientFilter ClientFilter) {
	for client := range h.clients {
		if !clientFilter(client) {
			continue
		}

		select {
		case client.send <- command:
		default:
			mutex.Lock()
			close(client.send)
			delete(h.clients, client)
			mutex.Unlock()
		}
	}
}

func (h *Hub) ObserveRegister(notifyer chan *Client, done chan int) {
	defer close(notifyer)
	h.observe[notifyer] = struct{}{}

	select {
	case <-done:
		delete(h.observe, notifyer)
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			fmt.Println("pre new client: ", client.GetName())
			mutex.Lock()
			h.clients[client] = struct{}{}
			for observeChan := range h.observe {
				observeChan <- client
			}
			mutex.Unlock()

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				mutex.Lock()
				fmt.Println("drop client: ", client.GetName())
				delete(h.clients, client)
				close(client.send)
				mutex.Unlock()
			}
		}
	}
}
