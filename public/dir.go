package public

import (
	"container/list"
	"sync"
)

//自定义队列
type Queue struct {
	sync.RWMutex
	dirs *list.List
}

func (q *Queue) Push(v interface{}) *list.Element {
	q.Lock()
	defer q.Unlock()

	return q.dirs.PushBack(v)
}

func (q *Queue) Pop() *list.Element {
	q.Lock()
	defer q.Unlock()

	if q.dirs.Len() == 0 {
		return nil
	}

	element := q.dirs.Front()
	q.dirs.Remove(element)

	return element
}
