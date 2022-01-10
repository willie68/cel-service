package celproc

import (
	"sort"
	"sync"
	"time"

	"github.com/google/cel-go/cel"
)

type LRUEntry struct {
	LastAccess time.Time   `json:"lastAccess"`
	ID         string      `json:"id"`
	Expression string      `json:"expression"`
	Program    cel.Program `json:"program"`
}

type LRUList struct {
	MaxCount int
	entries  []LRUEntry
	dmu      sync.Mutex
	cstLock  sync.Mutex
}

func (l *LRUList) Init() {
	l.entries = make([]LRUEntry, 0)
}

func (l *LRUList) Size() int {
	return len(l.entries)
}

func (l *LRUList) Add(e LRUEntry) bool {
	l.dmu.Lock()
	defer l.dmu.Unlock()
	e.LastAccess = time.Now()
	l.entries = l.insertSorted(l.entries, e)
	return true
}

func (l *LRUList) Update(e LRUEntry) {
	id := e.ID
	l.dmu.Lock()
	defer l.dmu.Unlock()
	i := sort.Search(len(l.entries), func(i int) bool { return l.entries[i].ID >= id })
	if i < len(l.entries) && l.entries[i].ID == id {
		e.LastAccess = time.Now()
		l.entries[i] = e
	}
}

func (l *LRUList) UpdateAccess(id string) {
	l.dmu.Lock()
	defer l.dmu.Unlock()
	i := sort.Search(len(l.entries), func(i int) bool { return l.entries[i].ID >= id })
	if i < len(l.entries) && l.entries[i].ID == id {
		l.entries[i].LastAccess = time.Now()
	}
}

func (l *LRUList) GetFullIDList() []string {
	l.dmu.Lock()
	defer l.dmu.Unlock()
	ids := make([]string, len(l.entries))
	for x, e := range l.entries {
		ids[x] = e.ID
	}
	return ids
}

func (l *LRUList) HandleContrains() {
	l.cstLock.Lock()
	defer l.cstLock.Unlock()
	for {
		id := l.getSingleContrained()
		if id != "" {
			l.Delete(id)
		} else {
			break
		}
	}
}

func (l *LRUList) getSingleContrained() string {
	var id string
	l.dmu.Lock()
	defer l.dmu.Unlock()
	if len(l.entries) > int(l.MaxCount) {
		// remove oldest entry from cache
		oldest := l.getOldest()
		id = l.entries[oldest].ID
	}
	return id
}

func (l *LRUList) Has(id string) bool {
	l.dmu.Lock()
	defer l.dmu.Unlock()
	i := sort.Search(len(l.entries), func(i int) bool { return l.entries[i].ID >= id })
	if i < len(l.entries) && l.entries[i].ID == id {
		return true
	}
	return false
}

func (l *LRUList) Get(id string) (LRUEntry, bool) {
	l.dmu.Lock()
	defer l.dmu.Unlock()
	i := sort.Search(len(l.entries), func(i int) bool { return l.entries[i].ID >= id })
	if i < len(l.entries) && l.entries[i].ID == id {
		l.entries[i].LastAccess = time.Now()
		return l.entries[i], true
	}
	return LRUEntry{}, false
}

func (l *LRUList) Delete(id string) string {
	l.dmu.Lock()
	defer l.dmu.Unlock()
	i := sort.Search(len(l.entries), func(i int) bool { return l.entries[i].ID >= id })
	if i < len(l.entries) && l.entries[i].ID == id {
		ret := make([]LRUEntry, 0)
		ret = append(ret, l.entries[:i]...)
		l.entries = append(ret, l.entries[i+1:]...)
		return id
	}
	return ""
}

func (l *LRUList) getOldest() int {
	oldest := 0
	for x, e := range l.entries {
		if e.LastAccess.Before(l.entries[oldest].LastAccess) {
			oldest = x
		}
	}
	return oldest
}

func (l *LRUList) insertSorted(data []LRUEntry, v LRUEntry) []LRUEntry {
	i := sort.Search(len(data), func(i int) bool { return data[i].ID >= v.ID })
	return l.insertEntryAt(data, i, v)
}

func (l *LRUList) insertEntryAt(data []LRUEntry, i int, v LRUEntry) []LRUEntry {
	if i == len(data) {
		// Insert at end is the easy case.
		return append(data, v)
	}

	// Make space for the inserted element by shifting
	// values at the insertion index up one index. The call
	// to append does not allocate memory when cap(data) is
	// greater ​than len(data).
	data = append(data[:i+1], data[i:]...)

	// Insert the new element.
	data[i] = v

	// Return the updated slice.
	return data
}
