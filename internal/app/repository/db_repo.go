package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ShishkovEM/amazing-shortener/internal/app/exceptions"
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
	"github.com/ShishkovEM/amazing-shortener/internal/app/responses"

	"github.com/jackc/pgerrcode"
)

type DBLinkStorage struct {
	DB *models.DB
}

func NewDBURLStorage(db *models.DB) *DBLinkStorage {
	return &DBLinkStorage{
		DB: db,
	}
}

func (d *DBLinkStorage) GetLink(shortID string) (models.OriginalURL, error) {
	var originalURL models.OriginalURL

	conn, err := d.DB.GetConn(context.Background())
	if err != nil {
		return originalURL, err
	}

	defer d.DB.Close()

	err = conn.QueryRow(context.Background(), "SELECT original_url, is_deleted FROM urls WHERE short_uri = $1 LIMIT 1", shortID).Scan(&originalURL.OriginalURL, &originalURL.IsDeleted)
	if err != nil {
		panic(err)
	}

	if originalURL.OriginalURL == "" {
		return originalURL, errors.New("not found")
	}

	return originalURL, nil
}

func (d *DBLinkStorage) GetShortURIByOriginalURL(originalURL string) (string, error) {
	var shortURI string

	conn, err := d.DB.GetConn(context.Background())
	if err != nil {
		return "", err
	}

	defer d.DB.Close()

	err = conn.QueryRow(context.Background(), "SELECT short_uri FROM urls WHERE original_url = $1 LIMIT 1", originalURL).Scan(&shortURI)
	if err != nil {
		panic(err)
	}

	if shortURI == "" {
		return "", errors.New("not found")
	}

	return shortURI, nil
}

func (d *DBLinkStorage) CreateLink(shortID string, originalURL string, userID uint32) error {
	conn, err := d.DB.GetConn(context.Background())
	if err != nil {
		panic(err)
	}

	defer d.DB.Close()

	_, err = conn.Exec(context.Background(), "INSERT INTO urls (short_uri, original_url, user_id, created_at) VALUES ($1,$2,$3,$4)", shortID, originalURL, userID, time.Now())

	if err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
		return &exceptions.LinkAlreadyExistsError{Value: originalURL}
	}

	if err != nil {
		panic(err)
	}

	return nil
}

func (d *DBLinkStorage) GetLinksByUserID(userID uint32) []responses.ResponseShortOriginalLink {
	userURLs := make([]responses.ResponseShortOriginalLink, 0)

	conn, err := d.DB.GetConn(context.Background())
	if err != nil {
		return userURLs
	}

	defer d.DB.Close()

	rows, err := conn.Query(context.Background(), "SELECT short_uri, original_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var urlData responses.ResponseShortOriginalLink

		err := rows.Scan(&urlData.ShortURL, &urlData.OriginalURL)
		if err != nil {
			return nil
		}

		userURLs = append(userURLs, urlData)
	}

	return userURLs
}

func (d *DBLinkStorage) DeleteUserRecordsByShortURLs(userID uint32, shortURLs []string) error {
	conn, err := d.DB.GetConn(context.Background())
	if err != nil {
		return err
	}

	defer d.DB.Close()

	_, err = conn.Exec(context.Background(), "UPDATE urls SET is_deleted = true WHERE short_uri = ANY($1) AND user_id = $2", shortURLs, userID)
	if err != nil {
		return err
	}

	return nil
}
