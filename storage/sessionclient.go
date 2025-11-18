package storage

import "sync"

type SessionStore struct {
	sync.RWMutex
	Data map[string]int
}

func GetSessionStore() *SessionStore {
	return &SessionStore{Data: make(map[string]int)}
}
