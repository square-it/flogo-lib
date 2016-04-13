package util

import (
	"container/list"
	"sync"
)

// SyncQueue is a List backed queue
type SyncQueue struct {
	List *list.List
	lock sync.Mutex
}

//NewSyncQueue creates a new SyncQueue
func NewSyncQueue() *SyncQueue {
	return &SyncQueue{List: list.New(), lock: sync.Mutex{}}
}

// Push push item on to queue
func (sq *SyncQueue) Push(item interface{}) {
	sq.lock.Lock()
	defer sq.lock.Unlock()

	sq.List.PushFront(item)
}

// Pop pop item off of queue
func (sq *SyncQueue) Pop() (interface{}, bool) {
	sq.lock.Lock()
	defer sq.lock.Unlock()

	if sq.List.Len() == 0 {
		return nil, false
	}

	item := sq.List.Front()
	sq.List.Remove(item)

	return item.Value, true
}

// Size get the size of the queue
func (sq *SyncQueue) Size() int {
	sq.lock.Lock()
	defer sq.lock.Unlock()

	return sq.List.Len()
}

// IsEmpty indicates if the queue is empty
func (sq *SyncQueue) IsEmpty() bool {
	sq.lock.Lock()
	defer sq.lock.Unlock()

	return (sq.List.Len() == 0)
}

/*

From: https://github.com/eapache/queue

Package queue provides a fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
Using this instead of other, simpler, queue implementations (slice+append or linked list) provides
substantial memory and time benefits, and fewer GC pauses.
The queue implemented here is as fast as it is for an additional reason: it is *not* thread-safe.
*/

const minQueueLen = 16

// Queue represents a single instance of the queue data structure.
type Queue struct {
	buf               []interface{}
	head, tail, count int
}

// NewQueue constructs and returns a new Queue.
func NewQueue() *Queue {
	return &Queue{
		buf: make([]interface{}, minQueueLen),
	}
}

// Length returns the number of items currently stored in the queue.
func (q *Queue) Length() int {
	return q.count
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Queue) resize() {
	newBuf := make([]interface{}, q.count*2)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

// Add puts an item on the end of the queue.
func (q *Queue) Add(item interface{}) {
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = item
	q.tail = (q.tail + 1) % len(q.buf)
	q.count++
}

// Peek returns the item at the head of the queue. This call panics
// if the queue is empty.
func (q *Queue) Peek() interface{} {
	if q.count <= 0 {
		panic("queue: Peek() called on empty queue")
	}
	return q.buf[q.head]
}

// Get returns the item at index i in the queue. If the index is
// invalid, the call will panic.
func (q *Queue) Get(i int) interface{} {
	if i < 0 || i >= q.count {
		panic("queue: Get() called with index out of range")
	}
	return q.buf[(q.head+i)%len(q.buf)]
}

// Remove removes the item from the front of the queue. If you actually
// want the item, call Peek first. This call panics if the queue is empty.
func (q *Queue) Remove() {
	if q.count <= 0 {
		panic("queue: Remove() called on empty queue")
	}
	q.buf[q.head] = nil
	q.head = (q.head + 1) % len(q.buf)
	q.count--
	if len(q.buf) > minQueueLen && q.count*4 == len(q.buf) {
		q.resize()
	}
}
