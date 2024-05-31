package main

import (
	"fmt"
	"sync"
	"time"
)

type Operation struct {
	Type  string
	Key   string
	Value string
	Wal   chan struct{}
}

type Database struct {
	data          map[string]string
	log           []Operation
	logSize       int
	logMaxSize    int
	logCh         chan struct{}
	mu            sync.Mutex
	commitTimeout time.Duration
}

func (db *Database) commitLog() {
	ticker := time.NewTicker(db.commitTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			db.commit()
		case <-db.logCh:
			if db.logSize >= db.logMaxSize {
				db.commit()
				ticker.Reset(db.commitTimeout)
			}
		}
	}
}

func (db *Database) commit() {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, v := range db.log {
		v.Wal <- struct{}{}
		fmt.Println("Committing log:", v)
	}

	db.log = make([]Operation, 0)
	db.logSize = 0
}

func (db *Database) Set(key string, value string) {
	db.mu.Lock()

	db.data[key] = value
	op := Operation{Type: "set", Key: key, Value: value, Wal: make(chan struct{})}
	db.log = append(db.log, op)
	db.logSize++

	db.mu.Unlock()

	db.logCh <- struct{}{}

	<-op.Wal
}

func (db *Database) Get(key string) (string, bool) {
	db.mu.Lock()
	value, exists := db.data[key]

	op := Operation{Type: "get", Key: key, Wal: make(chan struct{})}
	db.log = append(db.log, op)
	db.logSize++

	db.mu.Unlock()

	db.logCh <- struct{}{}

	<-op.Wal

	return value, exists
}

func NewDatabase(logMaxSize int) *Database {
	db := &Database{
		data:          make(map[string]string),
		log:           make([]Operation, 0),
		logMaxSize:    logMaxSize,
		logCh:         make(chan struct{}),
		commitTimeout: time.Second * 1,
	}
	go db.commitLog()
	return db
}

func main() {
	db := NewDatabase(5)
	db.Set("key1", "value1")
	db.Set("key2", "value2")
	db.Get("key1")
	db.Set("key3", "value3")
	db.Get("key2")
	db.Set("key4", "value4")
	db.Get("key3")
	db.Set("key5", "value5")
	db.Get("key4")
	db.Set("key6", "value6")
	db.Get("key5")
}
