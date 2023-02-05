package storage

import (
	"sync"

	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
)

type SessionStore struct {
	sync.Mutex

	Sessions map[int]models.Session
}

func NewSessionStore() *SessionStore {
	ss := &SessionStore{}
	ss.Sessions = make(map[int]models.Session)

	return ss
}
