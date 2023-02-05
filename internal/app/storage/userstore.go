package storage

import (
	"sync"

	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
)

type UserStore struct {
	sync.Mutex

	Users map[uint64]models.User
}

func NewUserStore() *UserStore {
	us := &UserStore{}
	us.Users = make(map[uint64]models.User)

	return us
}

func (us *UserStore) AddUser(user models.User) {
	us.Users[user.ID] = user
}

func (us *UserStore) GetUser(id uint64) (models.User, bool) {
	user, ok := us.Users[id]
	return user, ok
}
