package lrucache

import (
	"container/list"
	"sync"
)

type LRUEntry struct {
	Identifier string      `json:"identifier"`
	Payload    interface{} `json:"payload"`
}

type LRUCache struct {
	capacity int
	queue    *list.List
	entries  map[string]*list.Element
	dmu      sync.Mutex
}

func New(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		queue:    list.New(),
		entries:  make(map[string]*list.Element, capacity),
	}
}

func (l *LRUCache) Size() int {
	return len(l.entries)
}

func (l *LRUCache) Put(id string, payload interface{}) bool {
	e := LRUEntry{
		Identifier: id,
		Payload:    payload,
	}
	l.dmu.Lock()
	defer l.dmu.Unlock()

	if n, ok := l.entries[id]; ok {
		// entry already exists, so only update the payload
		l.queue.MoveToFront(n)
		n.Value.(*list.Element).Value = e
	} else {
		// before adding we have to check the capacity
		if l.queue.Len() == l.capacity {
			id := l.queue.Back().Value.(*list.Element).Value.(LRUEntry).Identifier
			delete(l.entries, id)
			l.queue.Remove(l.queue.Back())
		}

		n := &list.Element{
			Value: e,
		}

		p := l.queue.PushFront(n)
		l.entries[id] = p
	}

	return true
}

func (l *LRUCache) UpdateAccess(id string) {
	l.dmu.Lock()
	defer l.dmu.Unlock()
	if n, ok := l.entries[id]; ok {
		// entry already exists, so only put it in front of the queue
		l.queue.MoveToFront(n)
	}
}

func (l *LRUCache) GetFullIDList() []string {
	l.dmu.Lock()
	defer l.dmu.Unlock()
	ids := make([]string, len(l.entries))
	x := 0
	for k := range l.entries {
		ids[x] = k
		x++
	}
	return ids
}

func (l *LRUCache) Has(id string) bool {
	l.dmu.Lock()
	defer l.dmu.Unlock()
	_, ok := l.entries[id]
	return ok
}

func (l *LRUCache) Get(id string) (interface{}, bool) {
	l.dmu.Lock()
	defer l.dmu.Unlock()

	if n, ok := l.entries[id]; ok {
		// entry already exists, so only update the payload
		v := n.Value.(*list.Element).Value.(LRUEntry)
		l.queue.MoveToFront(n)
		return v.Payload, true
	}
	return nil, false
}

func (l *LRUCache) Remove(id string) bool {
	l.dmu.Lock()
	defer l.dmu.Unlock()

	if n, ok := l.entries[id]; ok {
		delete(l.entries, id)
		l.queue.Remove(n)
		return true
	}
	return false
}

func (l *LRUCache) GetOldest() (interface{}, bool) {
	if len(l.entries) > 0 {
		n := l.queue.Back()
		p := n.Value.(*list.Element).Value.(LRUEntry).Payload
		return p, true
	}
	return nil, false
}

func (l *LRUCache) Clear() {
	l.entries = make(map[string]*list.Element, l.capacity)
	l.queue = list.New()
}
