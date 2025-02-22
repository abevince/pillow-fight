package game

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	Conn     *websocket.Conn
	HP       int
	Room     *GameRoom
	IsClosed bool
}

func NewPlayer(conn *websocket.Conn) *Player {
	return &Player{
		Conn: conn,
		HP:   100, // Starting HP
	}
}

type Message struct {
	Type    string `json:"type"`
	Damage  int    `json:"damage,omitempty"`
	HP      int    `json:"hp,omitempty"`
	EnemyHP int    `json:"enemyHp,omitempty"`
	Message string `json:"message,omitempty"`
}

func (p *Player) SendMessage(msg Message) error {
	if p.IsClosed {
		return nil
	}
	return p.Conn.WriteJSON(msg)
}

func (p *Player) Close() {
	if !p.IsClosed {
		p.IsClosed = true
		p.Conn.Close()
	}
}
