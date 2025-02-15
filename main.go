package main

import (
	"couuple/models"
	"log"
	"math/rand"
	"strconv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func findRoom(rooms []models.Room, roomNumber string) *models.Room {
	for i := range rooms {
		if rooms[i].Number == roomNumber {
			return &rooms[i]
		}
	}
	return nil
}

func main() {

	app := fiber.New()

	rooms := make([]models.Room, 0)

	app.Post("/create", func(c *fiber.Ctx) error {
		room := strconv.Itoa(rand.Intn(900000) + 100000)

		for {
			_room := findRoom(rooms, room)
			if _room == nil {
				break
			}

			room = strconv.Itoa(rand.Intn(900000) + 100000)
		}

		rooms = append(rooms, models.Room{Number: room})

		return c.SendString(room)
	})

	app.Get("/game/:roomNumber", websocket.New(func(c *websocket.Conn) {
		defer c.Close()

		roomNumber := c.Params("roomNumber")
		room := findRoom(rooms, roomNumber)
		if room == nil {
			c.WriteMessage(websocket.TextMessage, []byte("Room not found"))
			return
		}

		if len(room.Clients) == 2 {
			c.WriteMessage(websocket.TextMessage, []byte("Room is full"))
			return
		}

		playerID := uint(len(room.Clients) + 1)

		room.AddClient(c, playerID)
		room.NotifyAll([]byte("New player joined"))

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}

			room.NotifyOther(playerID, msg)
		}
	}))

	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server started on :8080")

}
