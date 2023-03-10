package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/ShishkovEM/amazing-shortener/internal/app/interfaces"
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
	"github.com/ShishkovEM/amazing-shortener/internal/app/repository"
	"github.com/speps/go-hashids"
)

// LinkStore Структура для хранения записей типа Link в оперативной памяти
type LinkStore struct {
	sync.Mutex

	Links             map[string]models.Link
	nextID            int
	Repository        interfaces.LinkRepository
	DeletionProcessor interfaces.DeletionProcessor
}

// NewLinkStore Создаёт новый LinkStore
func NewLinkStore(repo interfaces.LinkRepository) *LinkStore {
	ls := &LinkStore{}
	ls.Links = make(map[string]models.Link)
	ls.nextID = 0
	if repo != nil {
		ls.Repository = repo
		repository.InitLinkStoreFromRepository(ls.Repository, ls)
	}
	return ls
}

func NewLinkStoreInMemory() *LinkStore {
	ls := &LinkStore{}
	ls.Links = make(map[string]models.Link)
	ls.nextID = 0

	return ls
}

func (ls *LinkStore) AddLinkToMemStorage(link models.Link) {
	ls.Lock()
	defer ls.Unlock()

	link.ID = ls.nextID
	ls.Links[link.Short] = link
	ls.nextID++
}

// CreateLink создаёт новую запись в LinkStore
func (ls *LinkStore) CreateLink(longURL string, userID uint32) (string, error) {
	ls.Lock()
	defer ls.Unlock()

	link := models.Link{
		ID:        ls.nextID,
		Original:  longURL,
		Short:     shorten(),
		UserID:    userID,
		IsDeleted: false,
	}
	ls.Links[link.Short] = link
	ls.nextID++

	if ls.Repository != nil {
		err := repository.WriteLinkToRepository(ls.Repository, &link)
		if err != nil {
			return "", err
		}
	}
	return link.Short, nil
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
func (ls *LinkStore) GetLink(short string) (models.Link, error) {
	ls.Lock()
	defer ls.Unlock()

	l, ok := ls.Links[short]

	if ok {
		return l, nil
	} else {
		return models.Link{}, fmt.Errorf("link with id=%s not found", short)
	}
}

func (ls *LinkStore) GetAll() []*models.Link {
	ls.Lock()
	defer ls.Unlock()

	var allLinks []*models.Link

	for _, link := range ls.Links {
		allLinks = append(allLinks, &link)
	}

	return allLinks
}

func (ls *LinkStore) GetSize() int {
	ls.Lock()
	defer ls.Unlock()

	return len(ls.Links)
}

func (ls *LinkStore) GetLinksByUserID(userID uint32) []models.Link {
	ls.Lock()
	defer ls.Unlock()

	var userLinks []models.Link

	for _, value := range ls.Links {
		if value.UserID == userID {
			userLinks = append(userLinks, value)
		}
	}

	return userLinks
}

func (ls *LinkStore) DeleteUserRecordsByShortURLs(userID uint32, shortURLs []string) {
	ls.Lock()
	defer ls.Unlock()

	for _, shortURL := range shortURLs {
		newTask := models.NewDeletionTask(userID, shortURL)
		ls.DeletionProcessor.AddTask(newTask)
	}
}

func (ls *LinkStore) DeleteOne(dt *models.DeletionTask) {
	ls.Lock()
	defer ls.Unlock()

	newLink := ls.Links[dt.URLToDelete]
	if newLink.UserID == dt.UserID {
		newLink.IsDeleted = true
		ls.Links[dt.URLToDelete] = newLink
	}

	if ls.Repository != nil {
		Refresh(ls, ls.Repository)
	}
}

func Refresh(inMemory interfaces.InMemoryLinkStorage, fileRepository interfaces.LinkRepository) {
	fileRepository.Refresh(inMemory)
}
