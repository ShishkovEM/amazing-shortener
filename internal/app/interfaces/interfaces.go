package interfaces

import "github.com/ShishkovEM/amazing-shortener/internal/app/models"

type InMemoryStorage interface {
	AddLinkToMemStorage(link models.Link)
	CreateLink(longURL string) (string, error)
	GetLink(short string) (models.Link, error)
	GetSize() int
}

type LinkRepository interface {
	InitLinkStoreFromRepository(store InMemoryStorage)
	WriteLinkToRepository(link *models.Link) error
}
