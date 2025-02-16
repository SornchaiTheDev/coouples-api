package models

import (
	status "couuple/constants/room_status"
	"couuple/services/emoji"
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Room struct {
	Number  string
	Players []*Player
	Picks   [2]Pick
	Score   int
	Status  status.RoomStatus
	Mu      sync.Mutex
}

func (r *Room) AddPlayer(conn *websocket.Conn, playerID uint) {
	r.Players = append(r.Players, &Player{
		Conn: conn,
		ID:   playerID,
	})
}

func (r *Room) RemovePlayer(playerID uint) {
	for i, player := range r.Players {
		if player.ID == playerID {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
		}
	}
}

func (r *Room) Start() {
	defer func() {
		r.NotifyAll(map[string]any{
			"type": "setup_phase",
		})
	}()

	// Go to create avatar phrase
	r.Status = status.Setup
	log.Println("🎮 Room", r.Number, "started")
}

func (r *Room) NotifyAll(msg map[string]any) {
	log.Println("📡 Sending message to all players:", msg)
	r.Mu.Lock()
	defer r.Mu.Unlock()

	jsonMarshal, err := json.Marshal(msg)
	if err != nil {
		log.Println("❌ Error marshaling message:", err)
	}

	for _, player := range r.Players {
		player.Mu.Lock()
		err := player.Conn.WriteMessage(websocket.TextMessage, jsonMarshal)
		if err != nil {
			log.Println("❌ Error writing message:", err)
		}
		player.Mu.Unlock()
	}
}

func (r *Room) NotifyOther(playerID uint, msg map[string]any) {
	jsonMarshal, err := json.Marshal(msg)
	if err != nil {
		log.Println("❌ Error marshaling message:", err)
	}

	for _, player := range r.Players {
		player.Mu.Lock()
		if player.ID == playerID {
			player.Mu.Unlock()
			continue
		}

		player.Conn.WriteMessage(websocket.TextMessage, jsonMarshal)
		player.Mu.Unlock()
	}
}

func (r *Room) Pick(playerID uint, cardID uint) {
	r.Picks[playerID-1] = Pick{
		PlayerID: playerID,
		CardID:   cardID,
	}
}

func (r *Room) ResetPicks() {
	r.Picks = [2]Pick{}
}

func (r *Room) SetAvatar(playerID uint, avatar string) {
	for _, player := range r.Players {
		if player.ID == playerID {
			player.Avatar = avatar
		}
	}
}

func (r *Room) GetGameDetail() map[string]any {
	detail := map[string]any{
		"status":  r.Status,
		"players": []string{r.Players[0].Avatar, r.Players[1].Avatar},
		"score":   r.Score,
	}

	return detail
}

func (r *Room) GameLoop() {
	for {
		if r.Status == status.Setup {
			readyPlayers := 0
			for _, player := range r.Players {
				if player.Avatar != "" {
					readyPlayers++
				}
			}

			if readyPlayers == 2 {
				r.Status = status.Shuffling
				r.NotifyAll(map[string]any{
					"type": "game_started",
				})

				r.NotifyAll(map[string]any{
					"type": "game_detail",
					"data": r.GetGameDetail(),
				})

			}
			continue
		}

		if r.Status == status.Shuffling {
			emojis := emoji.GetSet(4)
			r.NotifyAll(map[string]any{
				"type": "shuffling",
				"data": emojis,
			})
			r.Status = status.Picking
		}

		if r.Status == status.Picking {
			isPicked := true
			for _, pick := range r.Picks {
				if pick.CardID == 0 {
					isPicked = false
				}
			}

			if isPicked {
				if r.Picks[0].CardID == r.Picks[1].CardID {
					r.Score++
				}

				r.ResetPicks()
				r.NotifyAll(map[string]any{
					"type": "game_detail",
					"data": r.GetGameDetail(),
				})

				r.Status = status.Shuffling
			}
		}
	}
}
