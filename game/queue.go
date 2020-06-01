package game

import (
	"sync"
)

// 等待游戏排队队列
type Queue struct {
	head *People
	tail *People
	mutex sync.Mutex
}

// 玩家加入排队
func (q *Queue) Push (p *People) {

	q.mutex.Lock()

	if q.head == nil {
		q.head = p
	}

	if q.tail != nil {
		p.prev = q.tail
		q.tail.next = p
	}

	q.tail = p

	q.mutex.Unlock()
}

// 获取队列头部玩家
func (q *Queue) Pop() *People {

	q.mutex.Lock()

	head := q.head

	if head != nil {
		q.head = head.next
		q.head.prev = nil
	}

	if head == q.tail {
		q.tail = nil
	}

	q.mutex.Unlock()

	return head
}

// 从队列删除指定的玩家
func (q *Queue) Del(people *People) {

	if people == nil {
		return
	}

	q.mutex.Lock()

	prev := people.prev
	next := people.next

	prev.next = next
	next.prev = prev

	q.mutex.Unlock()
}
