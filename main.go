package main

import (
	"log"
	"net/http"
	"sync"

	"golang-battle/game"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for demo
	},
}

type GameManager struct {
	waitingPlayers []*game.Player
	activeGames    map[string]*game.GameRoom
	mu             sync.Mutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		waitingPlayers: make([]*game.Player, 0),
		activeGames:    make(map[string]*game.GameRoom),
	}
}

func (gm *GameManager) handleNewPlayer(player *game.Player) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Send waiting message
	player.SendMessage(game.Message{
		Type:    "status",
		Message: "Waiting for opponent...",
	})

	// Add to waiting players
	gm.waitingPlayers = append(gm.waitingPlayers, player)

	// Check if we can start a game
	if len(gm.waitingPlayers) >= 2 {
		player1 := gm.waitingPlayers[0]
		player2 := gm.waitingPlayers[1]
		gm.waitingPlayers = gm.waitingPlayers[2:]

		// Create new game room
		room := game.NewGameRoom(player1, player2)
		gm.activeGames[room.ID] = room

		// Notify players
		startMsg := game.Message{
			Type:    "gameStart",
			Message: "Game starting! Click/tap to attack!",
			HP:      100,
		}
		player1.SendMessage(startMsg)
		player2.SendMessage(startMsg)
	}
}

func main() {
	gameManager := NewGameManager()

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("static")))

	// WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error upgrading to WebSocket: %v", err)
			return
		}

		player := game.NewPlayer(conn)
		gameManager.handleNewPlayer(player)
	})

	log.Println("Server starting on 0.0.0.0:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
