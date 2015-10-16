package keylock

import (
	"hash/crc32"
	"sync"
)

type Keylock struct {
	lock_count uint32
	locks      []sync.Mutex
	table      *crc32.Table
}

func New(lock_count uint32) *Keylock {
	table := crc32.MakeTable(0xD5828281)
	keylock := Keylock{locks: make([]sync.Mutex, lock_count), table: table}
	keylock.lock_count = lock_count
	return &keylock
}

func (this *Keylock) Lock(key string) {
	this.locks[this.keyToIndex(key)].Lock()
}

func (this *Keylock) Unlock(key string) {
	this.locks[this.keyToIndex(key)].Unlock()
}

func (this *Keylock) keyToIndex(key string) uint32 {
	return crc32.Checksum([]byte(key), this.table) % this.lock_count
}
