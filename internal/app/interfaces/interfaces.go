package interfaces

import (
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
	"github.com/ShishkovEM/amazing-shortener/internal/app/responses"
)

type InMemoryLinkStorage interface {
	AddLinkToMemStorage(link models.Link)
	CreateLink(longURL string, userID uint32) (string, error)
	GetLink(short string) (models.Link, error)
	GetSize() int
}

type LinkRepository interface {
	InitLinkStoreFromRepository(store InMemoryLinkStorage)
	WriteLinkToRepository(link *models.Link) error
}

type DBLinkRepository interface {
	GetLink(shortID string) (models.OriginalURL, error)
	GetShortURIByOriginalURL(originalURL string) (string, error)
	CreateLink(shortID string, originalURL string, userID uint32) error
	GetLinksByUserID(userID uint32) []responses.ResponseShortOriginalLink
	DeleteUserRecordsByShortURLs(userID uint32, shortURLs []string)
}

type DeletionProcessor interface {
	AddTask(task *models.DeletionTask)
}
