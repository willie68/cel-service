package lrucache

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Entry struct {
	ID      string
	Payload string
}

const MAX_COUNT = 100

func TestCRUD(t *testing.T) {
	ast := assert.New(t)

	lru := New(MAX_COUNT)

	ast.NotNil(lru)
	ast.Equal(0, lru.Size())

	for i := 0; i < MAX_COUNT; i++ {
		e := Entry{
			ID:      fmt.Sprintf("%04d", i),
			Payload: fmt.Sprintf("%04d", i),
		}
		lru.Put(e.ID, e)
	}

	ast.Equal(MAX_COUNT, lru.Size())
	n, ok := lru.GetOldest()
	ast.True(ok)
	ast.Equal("0000", n.(Entry).ID)
	has := lru.Has("muckeuck")
	ast.False(has)

	e, ok := lru.Get("muckeuck")
	ast.False(ok)
	ast.Nil(e)

	for i := 0; i < MAX_COUNT; i++ {
		id := fmt.Sprintf("%04d", MAX_COUNT-i-1)
		has := lru.Has(id)
		ast.True(has)

		e, ok := lru.Get(id)
		ast.True(ok)
		ast.Equal(id, e.(Entry).Payload)

		Payload := fmt.Sprintf("cc_%04d", i)
		entry := e.(Entry)
		entry.Payload = Payload

		lru.Put(id, entry)
	}

	ast.Equal(MAX_COUNT, lru.Size())
	n, ok = lru.GetOldest()
	ast.True(ok)
	ast.Equal(fmt.Sprintf("%04d", MAX_COUNT-1), n.(Entry).ID)

	for i := 0; i < MAX_COUNT; i++ {
		id := fmt.Sprintf("%04d", MAX_COUNT-i-1)
		e, ok := lru.Get(id)
		ast.True(ok)
		Payload := fmt.Sprintf("cc_%04d", i)
		ast.Equal(Payload, e.(Entry).Payload)
	}

	ast.Equal(MAX_COUNT, lru.Size())

	n, ok = lru.GetOldest()
	ast.True(ok)
	lru.UpdateAccess(n.(Entry).ID)
	m, ok := lru.GetOldest()
	ast.True(ok)
	ast.NotEqual(n.(Entry).ID, m.(Entry).ID)

	for i := 0; i < MAX_COUNT; i++ {
		id := fmt.Sprintf("%04d", MAX_COUNT-i-1)
		ok := lru.Remove(id)
		ast.True(ok)
	}

	lru.Clear()

	ast.Equal(0, lru.Size())
}

func TestWrongID(t *testing.T) {
	ast := assert.New(t)

	lru := New(MAX_COUNT)

	ast.NotNil(lru)
	ast.Equal(0, lru.Size())

	for i := 0; i < MAX_COUNT; i++ {
		e := Entry{
			ID:      fmt.Sprintf("%04d", i),
			Payload: fmt.Sprintf("%04d", i),
		}
		lru.Put(e.ID, e)
	}

	ast.Equal(MAX_COUNT, lru.Size())
	n, ok := lru.GetOldest()
	ast.True(ok)
	ast.Equal("0000", n.(Entry).ID)

	has := lru.Has("muckeuck")
	ast.False(has)

	e, ok := lru.Get("muckeuck")
	ast.False(ok)
	ast.Nil(e)

	ok = lru.Remove("muckeuck")
	ast.False(ok)
}

func TestList(t *testing.T) {
	ast := assert.New(t)
	idList := make([]string, 0)

	lru := New(MAX_COUNT)
	ast.NotNil(lru)
	ast.Equal(0, lru.Size())

	for i := 0; i < MAX_COUNT; i++ {
		id := fmt.Sprintf("%04d", i)
		e := Entry{
			ID:      id,
			Payload: fmt.Sprintf("%04d", i),
		}
		lru.Put(e.ID, e)
		idList = append(idList, id)
	}
	sort.Strings(idList)

	ast.Equal(MAX_COUNT, lru.Size())

	list := lru.GetFullIDList()
	sort.Strings(list)
	ast.Equal(MAX_COUNT, len(list))

	ast.True(reflect.DeepEqual(idList, list))
}

func TestClear(t *testing.T) {
	ast := assert.New(t)

	lru := New(MAX_COUNT)
	ast.NotNil(lru)
	ast.Equal(0, lru.Size())

	for i := 0; i < MAX_COUNT; i++ {
		e := Entry{
			ID:      fmt.Sprintf("%04d", i),
			Payload: fmt.Sprintf("%04d", i),
		}
		lru.Put(e.ID, e)
	}
	ast.Equal(MAX_COUNT, lru.Size())

	lru.Clear()
	ast.Equal(0, lru.Size())

}

func TestEviction(t *testing.T) {
	ast := assert.New(t)

	lru := New(MAX_COUNT)
	ast.NotNil(lru)
	ast.Equal(0, lru.Size())

	for i := 0; i < MAX_COUNT*2; i++ {
		e := Entry{
			ID:      fmt.Sprintf("%04d", i),
			Payload: fmt.Sprintf("%04d", i),
		}
		lru.Put(e.ID, e)
	}
	ast.Equal(MAX_COUNT, lru.Size())

	for i := 0; i < MAX_COUNT; i++ {
		id := fmt.Sprintf("%04d", i)
		ast.False(lru.Has(id))
	}

	for i := MAX_COUNT; i < MAX_COUNT*2; i++ {
		id := fmt.Sprintf("%04d", i)
		ast.True(lru.Has(id))
	}

	lru.Clear()
	ast.Equal(0, lru.Size())
}
func TestGetOldest(t *testing.T) {
	ast := assert.New(t)

	lru := New(MAX_COUNT)
	ast.NotNil(lru)
	ast.Equal(0, lru.Size())
	e, ok := lru.GetOldest()
	ast.False(ok)
	ast.Nil(e)
}
