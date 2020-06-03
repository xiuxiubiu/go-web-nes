package game

import (
	"sync"
)

// 等待游戏排队队列
type Queue struct {
	head  *People
	tail  *People
	mutex sync.Mutex
}

// 玩家加入排队
func (q *Queue) Push(people *People) {

	q.mutex.Lock()

	if q.head == nil {
		q.head = people
	}

	if q.tail != nil {
		people.prev = q.tail
		q.tail.next = people
	}

	q.tail = people

	q.mutex.Unlock()
}

// 获取队列头部玩家
func (q *Queue) Pop() *People {

	if q.head == nil {
		return nil
	}

	q.mutex.Lock()

	head := q.head

	if head != nil {
		q.head = head.next
	}

	if q.head != nil {
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

	if prev != nil && prev.next != nil {
		prev.next = next
	}

	if next != nil && next.prev != nil {
		next.prev = prev
	}

	q.mutex.Unlock()
}
