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

	// 锁
	mutex sync.Mutex
}

// 创建房间
func (r *Room) Create() {

	// 加入游戏监听
	r.JoinListener()

	// 玩家状态监听
	r.PlayerListener()
}

// 设置玩家1
func (r *Room) SetOwner(p *People) (err error) {
	if p != nil {
		err = p.Send("join 1")
		if err != nil {
			return
		}
	}
	r.owner = p
	return
}

// 设置玩家2
func (r *Room) SetPlayer(p *People) (err error) {
	if p != nil {
		err = p.Send("join 2")
		if err != nil {
			return
		}
	}
	r.player = p
	return
}

// 获取等待队列的头部玩家
// 若队列为空阻塞直到队列有值
func (r *Room) BlockGetHead() *People {

	people := r.queue.Pop()
	for people == nil {
		time.Sleep(10 * time.Millisecond)
		people = r.queue.Pop()
	}

	return people
}

// 加入游戏监听
func (r *Room) JoinListener() {

	go func() {

		for {

			<- r.join

			// 如果玩家1和2都不为nil
			// 等待玩家监听协程将已退出玩家设置为nil
			for r.player != nil && r.owner != nil {
				time.Sleep(10 * time.Millisecond)
			}

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

			// 获取队列头部玩家
			people := r.BlockGetHead()

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
	}()
}

// 玩家状态监听
func (r *Room) PlayerListener()  {

	for {

		r.mutex.Lock()

		if r.owner != nil && r.owner.Send("ping") != nil {
			r.owner = nil
		}

		if r.player != nil && r.player.Send("ping") != nil {
			r.player = nil
		}

		r.mutex.Unlock()

		time.Sleep(1 * time.Second)
	}
}

// 玩家进入房间
func (r *Room) Enter(people *People) {

	// 加入排队
	r.queue.Push(people)

	// 加入观战
	r.watchers[people] = nil
}

// 玩家退出房间
func (r *Room) Exit(people *People) {
	if people != r.owner && people != r.player {
		r.queue.Del(people)
		delete(r.watchers, people)
	}
}

// 创建游戏房间
func Step() *Room {

	room := &Room{
		queue: &Queue{},
		watchers: map[*People]interface{}{},
		join: make(chan interface{}, 2),
	}

	// 需要两个玩家加入
	room.join <- nil
	room.join <- nil

	room.Create()

	return room
}