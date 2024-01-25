package microredis

import (
	"sync"
	"time"
)

type MicroRedis[Key any, Value any] struct {
	data    map[string]Value
	key     func(item *Key) string
	sync    *sync.RWMutex
	timeout time.Duration
}

func New[Value any](timeout time.Duration) *MicroRedis[string, Value] {
	return Custom[string, Value](timeout, func(item *string) string {
		return *item
	})
}

func Custom[Key any, Value any](timeout time.Duration, key func(item *Key) string) *MicroRedis[Key, Value] {
	return &MicroRedis[Key, Value]{
		data:    map[string]Value{},
		key:     key,
		sync:    &sync.RWMutex{},
		timeout: timeout,
	}
}

func (db *MicroRedis[Key, Value]) Get(key Key) *Value {
	db.sync.RLock()
	defer db.sync.RUnlock()

	el, ok := db.data[db.key(&key)]
	if !ok {
		return nil
	}
	return &el
}

func (db *MicroRedis[Key, Value]) Set(key Key, val Value) {
	db.sync.Lock()
	defer db.sync.Unlock()

	k := db.key(&key)
	db.data[k] = val
	go time.AfterFunc(db.timeout, func() {
		db.Del(key)
	})
}

func (db *MicroRedis[Key, Value]) Del(key Key) bool {
	db.sync.Lock()
	defer db.sync.Unlock()

	k := db.key(&key)
	if _, ok := db.data[k]; !ok {
		return false
	}

	delete(db.data, k)
	return true
}

func (db *MicroRedis[Key, Value]) Size() int {
	return len(db.data)
}
