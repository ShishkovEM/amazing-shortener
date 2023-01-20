package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/speps/go-hashids"
)

// Link Структура записи информации о гиперссылках
type Link struct {
	ID       int    `json:"id"`       // Идентификатор гиперссылки
	Original string `json:"original"` // Исходная (длинная) ссылка
	Short    string `json:"short"`    // Короткая ссылка
}

// LinkStore Структура для хранения записей типа Link в оперативной памяти
type LinkStore struct {
	sync.Mutex

	Links  map[string]Link
	nextID int
}

// NewLinkStore Создаёт новый LinkStore
func NewLinkStore() *LinkStore {
	ls := &LinkStore{}
	ls.Links = make(map[string]Link)
	ls.nextID = 0

	return ls
}

func (ls *LinkStore) AddLinkToMemStorage(link Link) {
	ls.Lock()
	defer ls.Unlock()

	link.ID = ls.nextID
	ls.Links[link.Short] = link
	ls.nextID++
}

// CreateLink создаёт новую запись в LinkStore
func (ls *LinkStore) CreateLink(longURL string) string {
	ls.Lock()
	defer ls.Unlock()

	link := Link{
		ID:       ls.nextID,
		Original: longURL,
		Short:    shorten(),
	}
	ls.Links[link.Short] = link
	ls.nextID++
	return link.Short
}

// shorten() Создаёт короткий хэш
func shorten() string {
	hd := hashids.NewData()
	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	short, _ := h.Encode([]int{int(now.UnixMicro())})
	return short
}

// GetLink получает запись об одной ссылке по её id
func (ls *LinkStore) GetLink(short string) (Link, error) {
	ls.Lock()
	defer ls.Unlock()

	l, ok := ls.Links[short]

	if ok {
		return l, nil
	} else {
		return Link{}, fmt.Errorf("link with id=%s not found", short)
	}
}

func (ls *LinkStore) GetSize() int {
	return len(ls.Links)
}
