package game

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const (

	// 消息发送超时时间
	writeWait = 3 * time.Second

	// ping发送间隔时间
	pingWait = (pongWait / 5) * 4

	// pong接收间隔时间
	pongWait = 5 * time.Second

	// 消息最大长度
	maxMessageSize = 24 * 1024
)

// 玩家实例
type People struct {

	// websocket链接实例
	Conn *websocket.Conn

	// 房间实例
	Room *Room

	// 排在玩家后的玩家
	next *People

	// 排在玩家前的玩家
	prev *People

	// 消息发送通道
	send chan *Message

	// 消息读取处理方法
	ReadHandle func(message []byte)

	// 玩家是否退出
	isExit bool

	// 锁
	mutex sync.Mutex
}

// 读取消息逻辑
func (p *People) ReadPump() {

	defer func() {

		// 从房间退出
		p.Room.Exit <- p

		// 关闭websocket链接
		if err := p.Conn.Close(); err != nil {

		}
	}()

	// 设置读取消息大小
	p.Conn.SetReadLimit(maxMessageSize)

	// 设置读取超时时间
	if err := p.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {

	}

	// 设置pong消息处理函数
	p.Conn.SetPongHandler(func(appData string) error {

		// 读取pong后重新设置读取超时时间
		if err := p.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {

		}
		return nil
	})

	for {

		// 读取消息
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			return
		}

		// 发送消息给客户端
		p.ReadHandle(message)
	}
}

// 发送消息
func (p *People) Send(message *Message) {

	p.mutex.Lock()

	if !p.isExit {
		p.send <- message
	}

	p.mutex.Unlock()
}

// 消息发送处理逻辑
func (p *People) SendPump() {

	// ping消息定时器
	ticker := time.NewTicker(pingWait)

	defer func() {

		// 关闭定时器
		ticker.Stop()

		// 关闭websocket链接
		if err := p.Conn.Close(); err != nil {

		}
	}()

	for {

		select {
		case message, ok := <-p.send:

			// 设置写超时
			if err := p.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {

			}

			// 玩家主动退出Room关闭发送通道
			// 给客户端发送关闭的消息
			if !ok {
				_ = p.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 发送消息
			err := p.Conn.WriteMessage(websocket.TextMessage, message.data)
			if err != nil {

			}

		case <-ticker.C:

			// 设置消息发送超时时间
			if err := p.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {

			}

			// 发送ping消息
			if err := p.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {

			}
		}
	}
}

// 用户退出
func (p *People) Exit() {

	p.mutex.Lock()
	p.isExit = true
	p.mutex.Unlock()

	close(p.send)
}

func NewPeople(conn *websocket.Conn) *People {

	p := &People{
		Conn: conn,
		send: make(chan *Message, 256),
		isExit: false,
	}

	go p.ReadPump()
	go p.SendPump()

	return p
}
