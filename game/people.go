package game

import "github.com/gorilla/websocket"

// 玩家实例
type People struct {
	Conn *websocket.Conn
	Partner *People
	next *People
	prev *People
}

// 向玩家发送信息
func (p *People) Send(m string) error {
	return p.Conn.WriteMessage(1, []byte(m))
}

type Player struct {
	People *People
	Partner *Player
}


