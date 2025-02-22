package game

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type GameRoom struct {
	ID      string
	Player1 *Player
	Player2 *Player
	mu      sync.Mutex
}

func NewGameRoom(p1, p2 *Player) *GameRoom {
	room := &GameRoom{
		ID:      generateRoomID(),
		Player1: p1,
		Player2: p2,
	}

	p1.Room = room
	p2.Room = room

	// Start game handlers for both players
	go room.handlePlayerMessages(p1)
	go room.handlePlayerMessages(p2)

	return room
}

func (r *GameRoom) handlePlayerMessages(player *Player) {
	for {
		var msg Message
		err := player.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			r.handlePlayerDisconnect(player)
			return
		}

		switch msg.Type {
		case "attack":
			r.handleAttack(player)
		}
	}
}

func (r *GameRoom) handleAttack(attacker *Player) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Determine the target player
	var target *Player
	if attacker == r.Player1 {
		target = r.Player2
	} else {
		target = r.Player1
	}

	// Calculate random damage between 1 and 2
	damage := 1 + rand.Intn(2)
	target.HP -= damage

	// Calculate healing (half of damage)
	healing := damage / 2
	attacker.HP = min(100, attacker.HP+healing) // Cap HP at 100

	// Send updated HP to both players
	attackMsg := Message{
		Type:    "damage",
		Damage:  damage,
		Message: "You dealt damage and healed!",
		HP:      attacker.HP,
		EnemyHP: target.HP,
	}
	attacker.SendMessage(attackMsg)

	receivedMsg := Message{
		Type:    "damage",
		Damage:  damage,
		Message: "You received damage!",
		HP:      target.HP,
		EnemyHP: attacker.HP,
	}
	target.SendMessage(receivedMsg)

	// Check for game over
	if target.HP <= 0 {
		r.handleGameOver(attacker, target)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (r *GameRoom) handleGameOver(winner, loser *Player) {
	winnerMsg := Message{
		Type:    "gameOver",
		Message: "You Won!",
	}
	loserMsg := Message{
		Type:    "gameOver",
		Message: "You Lost!",
	}

	// Send messages before closing connections
	winner.SendMessage(winnerMsg)
	loser.SendMessage(loserMsg)

	// Add a small delay before closing connections to ensure messages are sent
	go func() {
		time.Sleep(100 * time.Millisecond)
		winner.Close()
		loser.Close()
	}()
}

func (r *GameRoom) handlePlayerDisconnect(player *Player) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Send game over message to remaining player
	var remainingPlayer *Player
	if player == r.Player1 {
		remainingPlayer = r.Player2
	} else {
		remainingPlayer = r.Player1
	}

	if remainingPlayer != nil && !remainingPlayer.IsClosed {
		// Send message before closing connection
		remainingPlayer.SendMessage(Message{
			Type:    "gameOver",
			Message: "You Won! Opponent disconnected.",
		})

		// Add a small delay before closing connection
		go func() {
			time.Sleep(100 * time.Millisecond)
			remainingPlayer.Close()
		}()
	}

	player.Close()
}

func generateRoomID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
