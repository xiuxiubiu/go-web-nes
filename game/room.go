package game

import (
	"sync"
	"time"
)

// 房间
type Room struct {

	// 玩家1
	owner *People

	// 玩家2
	player *People

	// 等待游戏队列
	queue *Queue

	// 观战玩家
	watchers map[*People]interface{}

	// 玩家加入游戏通道
	join chan interface{}

	// 进入房间通道
	Enter chan *People

	// 退出房间通道
	Exit chan *People

	// 广播
	broadcast chan []byte

	// 锁
	mutex sync.Mutex
}

// 设置玩家1
func (r *Room) SetOwner(p *People) (err error) {

	if p != nil {
		p.Send(NewMessage([]byte("join 1")))
		p.ReadHandle = r.ownerReadHandle()
	}
	r.owner = p

	return
}

// 玩家1消息读取处理函数
func (r *Room) ownerReadHandle() func(message []byte) {

	return func(message []byte) {

		msg := NewMessage(message)
		if msg == nil {
			return
		}

		// nes数据
		if msg.mt == mtData {

			// 优先发送给玩家2
			r.player.Send(msg)

			// 广播给观战者
			r.broadcast <- msg.data
		}

		// 玩家2的操控命令
		if msg.mt == mtDown || msg.mt == mtUp {

			// 发送给玩家1自己的客户端
			r.owner.Send(msg)
		}
	}
}

// 设置玩家2
func (r *Room) SetPlayer(p *People) (err error) {

	if p != nil {
		p.Send(NewMessage([]byte("join 2")))
		p.ReadHandle = r.playerReadHandle()
	}
	r.player = p

	return
}

// 玩家2消息读取处理函数
func (r *Room) playerReadHandle() func(message []byte) {

	return func(message []byte) {

		msg := NewMessage(message)
		if msg == nil {
			return
		}

		// nes数据
		if msg.mt == mtData {

			// 发送给玩家2自己的客户端
			r.player.Send(msg)
		}

		// 玩家2的操控命令
		if msg.mt == mtDown || msg.mt == mtUp {

			// 将控制命令发送给玩家1
			r.owner.Send(msg)
		}
	}

}

// 获取等待队列的头部玩家
// 若队列为空阻塞直到队列有值
func (r *Room) BlockGetHead() *People {

	people := r.queue.Pop()
	for people == nil {
		time.Sleep(1 * time.Second)
		people = r.queue.Pop()
	}

	return people
}

// 加入游戏监听
func (r *Room) JoinListener() {

	for {

		<- r.join

		// 如果玩家1和2都不为nil
		// 等待玩家监听协程将已退出玩家设置为nil
		for r.player != nil && r.owner != nil {
			time.Sleep(10 * time.Millisecond)
		}

		// 获取队列头部玩家
		people := r.BlockGetHead()

		// 加锁防止与玩家监听协程冲突
		r.mutex.Lock()

		// 玩家2不为nil则升级为玩家1
		// 若升级玩家1失败
		// 重试直至成功并将玩家2设置为nil
		// 或者
		// 等待玩家监听协程将玩家2设置为nil
		for r.player != nil {
			if  r.SetOwner(r.player) != nil {
				continue
			}
			_ = r.SetPlayer(nil)
		}

		// 如果玩家1为空设置玩家1
		// 否则设置玩家2
		if r.owner == nil {

			// 设置玩家1不成功
			if r.SetOwner(people) != nil {
				r.join <- nil
			}

		} else {

			// 设置玩家2不成功
			if r.SetPlayer(people) != nil {
				r.join <- nil
			}
		}

		// 解锁
		r.mutex.Unlock()
	}
}

// 房间管理逻辑
func (r *Room) Manager() {

	for {
		select {

		// 用户加入
		case people := <- r.Enter:

			// 加入排队
			r.queue.Push(people)

			// 加入观战
			r.watchers[people] = nil

			// 设置用户房间
			people.Room = r

			// 消息处理
			people.ReadHandle = func(message []byte) {
				people.Send(NewMessage(message))
			}

		// 用户退出
		case people := <- r.Exit:

			// 从队列退出
			r.queue.Del(people)

			// 从观战中退出
			if _, ok := r.watchers[people]; ok {
				delete(r.watchers, people)
			}

			// 用户退出
			people.Exit()

		// 广播信息
		case message := <- r.broadcast:
			for people := range r.watchers {
				select {
				case people.send <- NewMessage(message):
				default:
					continue
				}
			}
		}
	}
}

// 创建游戏房间
func NewRoom() *Room {

	// 房间实例
	room := &Room{
		queue: &Queue{},
		watchers: map[*People]interface{}{},
		join: make(chan interface{}, 2),
		Enter: make(chan *People, 256),
		Exit: make(chan *People, 256),
		broadcast: make(chan []byte, 256),
	}

	// 需要两个玩家加入
	room.join <- nil
	room.join <- nil

	// 房间管理协程
	go room.Manager()

	// 玩家监听协程
	go room.JoinListener()

	return room
}