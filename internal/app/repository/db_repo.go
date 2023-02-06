package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
	"github.com/ShishkovEM/amazing-shortener/internal/app/responses"
)

type DBLinkStorage struct {
	DB *models.DB
}

func NewDBURLStorage(db *models.DB) *DBLinkStorage {
	return &DBLinkStorage{
		DB: db,
	}
}

func (d *DBLinkStorage) GetLink(shortID string) (string, error) {
	var originalURL string

	conn, err := d.DB.GetConn(context.Background())
	if err != nil {
		return "", err
	}

	defer d.DB.Close()

	err = conn.QueryRow(context.Background(), "SELECT original_url FROM urls WHERE short_url = $1 LIMIT 1", shortID).Scan(&originalURL)
	if err != nil {
		panic(err)
	}

	if originalURL == "" {
		return "", errors.New("not found")
	}

	return originalURL, nil
}

func (d *DBLinkStorage) CreateLink(shortID string, originalURL string, userID uint32) {
	conn, err := d.DB.GetConn(context.Background())
	if err != nil {
		panic(err)
	}

	defer d.DB.Close()

	_, err = conn.Exec(context.Background(), "INSERT INTO urls (short_url, original_url, user_id, created_at) VALUES ($1,$2,$3,$4)", shortID, originalURL, userID, time.Now())
	if err != nil {
		panic(err)
	}
}

func (d *DBLinkStorage) GetLinksByUserID(userID uint32) []responses.ResponseShortOriginalLink {
	userURLs := make([]responses.ResponseShortOriginalLink, 0)

	conn, err := d.DB.GetConn(context.Background())
	if err != nil {
		return userURLs
	}

	defer d.DB.Close()

	rows, err := conn.Query(context.Background(), "SELECT short_url, original_url FROM urls WHERE user_id = $1", userID)
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
