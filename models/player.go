package models

import "github.com/gofiber/contrib/websocket"

type Player struct {
	ID   uint
	Conn *websocket.Conn
}
