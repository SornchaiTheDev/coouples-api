package main

import (
	"couuple/constants/messages"
	status "couuple/constants/room_status"
	"couuple/models"
	"encoding/json"
	"log"
	"math/rand"
	"strconv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
	}))

	rooms := make([]models.Room, 0)
	roomChan := make(chan *models.Room)

	app.Post("/create", func(c *fiber.Ctx) error {
		room := strconv.Itoa(rand.Intn(900000) + 100000)

		for {
			_room := findRoom(rooms, room)
			if _room == nil {
				break
			}

			room = strconv.Itoa(rand.Intn(900000) + 100000)
		}

		rooms = append(rooms, models.Room{Number: room, Status: status.Waiting})

		roomChan <- &rooms[len(rooms)-1]

		return c.JSON(fiber.Map{
			"code": room,
		})
	})

	app.Get("/game/:roomNumber", websocket.New(func(c *websocket.Conn) {
		defer c.Close()

		roomNumber := c.Params("roomNumber")
		room := findRoom(rooms, roomNumber)
		if room == nil {
			c.WriteMessage(websocket.TextMessage, []byte("Room not found"))
			return
		}

		if len(room.Players) == 2 {
			c.WriteMessage(websocket.TextMessage, []byte("Room is full"))
			return
		}

		playerID := uint(len(room.Players) + 1)

		room.AddPlayer(c, playerID)

		log.Println("Player ID : ", playerID, "has joined the game")

		room.NotifyAll(map[string]any{
			"type": "player_joined",
		})

		for {
			var msg models.Message
			if err := c.ReadJSON(&msg); err != nil {
				if err != nil {
					if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
						room.RemovePlayer(playerID)
						room.NotifyOther(playerID, map[string]any{
							"type": "player_left",
						})
					} else {
						log.Println("WebSocket error:", err)
					}
					break
				}

				log.Println("❌ Error reading message:", err)
			}

			if msg.Type == messages.START {
				if len(room.Players) < 2 {
					bytes, err := json.Marshal(map[string]string{
						"type": "not_enough_players",
						"data": "Not enough players",
					})
					if err != nil {
						log.Println("❌ Error marshaling message:", err)
					}

					c.WriteMessage(websocket.TextMessage, bytes)
					continue
				}

				if room.Status != status.Waiting {
					bytes, err := json.Marshal(map[string]string{
						"type": "game_already_started",
					})
					if err != nil {
						log.Println("❌ Error marshaling message:", err)
					}

					c.WriteMessage(websocket.TextMessage, bytes)
					continue
				}

				room.Start()
			}

			if msg.Type == messages.GET_DETAIL {
				detail := map[string]any{
					"status": room.Status,
					"players": []string{
						room.Players[0].Avatar,
						room.Players[1].Avatar,
					},
					"score": room.Score,
				}

				bytes, err := json.Marshal(map[string]any{
					"type": "game_detail",
					"data": detail,
				})
				if err != nil {
					log.Println("❌ Error marshaling message:", err)
				}

				c.WriteMessage(websocket.TextMessage, bytes)

			}

			if msg.Type == messages.CREATE_AVATAR {
				room.SetAvatar(playerID, msg.Data)
				log.Println("Player ", playerID, "had setup the avatar")

				bytes, err := json.Marshal(map[string]any{
					"type": "avatar_created",
					"data": "Avatar created",
				})
				if err != nil {
					log.Println("❌ Error marshaling message:", err)
				}

				c.WriteMessage(websocket.TextMessage, bytes)
			}

			if msg.Type == messages.PICK {
				cardID, err := strconv.Atoi(msg.Data)
				if err != nil {
					log.Println("❌ Error converting cardID:", err)
				}

				room.Pick(playerID, uint(cardID))

				bytes, err := json.Marshal(map[string]any{
					"type": "card_picked",
				})
				if err != nil {
					log.Println("❌ Error marshaling message:", err)
				}

				c.WriteMessage(websocket.TextMessage, bytes)
			}
		}

	}))

	go func() {
		for {
			select {
			case room := <-roomChan:
				go room.GameLoop()
			}
		}
	}()

	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server started on :8080")

}
