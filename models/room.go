package models

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Room struct {
	Number  string
	Clients []*Player
	Score   int
	Mu      sync.Mutex
}

func (r *Room) AddClient(conn *websocket.Conn, playerID uint) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	r.Clients = append(r.Clients, &Player{
		Conn: conn,
		ID:   playerID,
	})
}

func (r *Room) NotifyAll(msg []byte) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	for _, client := range r.Clients {
		client.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func (r *Room) NotifyOther(playerID uint, msg []byte) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	for _, client := range r.Clients {
		if client.ID == playerID {
			continue
		}

		client.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
