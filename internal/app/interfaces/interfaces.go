package interfaces

import (
	"context"
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
	"github.com/ShishkovEM/amazing-shortener/internal/app/responses"
	"github.com/jackc/pgx/v4"
)

type InMemoryLinkStorage interface {
	AddLinkToMemStorage(link models.Link)
	CreateLink(longURL string, userID uint32) (string, error)
	GetLink(short string) (models.Link, error)
	GetAll() []*models.Link
	GetSize() int
	DeleteUserRecordsByShortURLs(userID uint32, shortURLs []string)
	DeleteOne(dt *models.DeletionTask)
}

type LinkRepository interface {
	InitLinkStoreFromRepository(store InMemoryLinkStorage)
	WriteLinkToRepository(link *models.Link) error
	Refresh(inMemory InMemoryLinkStorage)
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

type ProcessorTarget interface {
	GetConn(ctx context.Context) (*pgx.Conn, error)
}
