package linkstore

import (
	"fmt"
	"sync"
	"time"

	"github.com/speps/go-hashids"
)

// Link Структура записи информации о гиперссылках
type Link struct {
	ID       int    // Идентификатор гиперссылки
	Original string // Исходная (длинная) ссылка
	Short    string // Короткая ссылка
}

// LinkStore Структура для хранения записей типа Link в оперативной памяти
type LinkStore struct {
	sync.Mutex

	links  map[string]Link
	nextID int
}

// NewLinkStore Создаёт новый LinkStore
func NewLinkStore() *LinkStore {
	ls := &LinkStore{}
	ls.links = make(map[string]Link)
	ls.nextID = 0
	return ls
}

// CreateLink создаёт новую запись в LinkStore
func (ls *LinkStore) CreateLink(longURL string) string {
	ls.Lock()
	defer ls.Unlock()

	link := Link{
		ID:       ls.nextID,
		Original: longURL,
		Short:    shorten()}

	ls.links[link.Short] = link
	ls.nextID++
	return link.Short
}

// shorten() Создаёт короткий хэш
func shorten() string {
	hd := hashids.NewData()
	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	short, _ := h.Encode([]int{int(now.Unix())})
	return short
}

// GetLink получает запись об одной ссылке по её id
func (ls *LinkStore) GetLink(short string) (Link, error) {
	ls.Lock()
	defer ls.Unlock()

	l, ok := ls.links[short]

	if ok {
		return l, nil
	} else {
		return Link{}, fmt.Errorf("link with id=%s not found", short)
	}
}
