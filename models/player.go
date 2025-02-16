package models

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Player struct {
	ID     uint
	Avatar string
	Conn   *websocket.Conn
	Mu     sync.Mutex
}
