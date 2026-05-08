package repository

import (
	"database/sql"
	"log"
	"time"
	"url-shortener/internal/model"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Save(url *model.URL) error {
	query := "INSERT INTO urls(original_url, short_code, created_at) VALUES (?,?,?)"
	_, err := r.db.Exec(query, url.OriginalURL, url.ShortCode, time.Now())
	return err
}

func (r *URLRepository) FindByShortCode(shortCode string) (*model.URL, error) {
	query := "SELECT id, original_url, short_code, clicks, created_at FROM urls WHERE short_code = ?"
	row := r.db.QueryRow(query, shortCode)

	url := &model.URL{}
	err := row.Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.Clicks, &url.CreatedAt)

	if err != nil {
		log.Println("FindByShortCode error:", err)        // debug
		log.Println("ShortCode yang dicari:", shortCode)
		return nil, err
	}

	return url, nil 
}

func (r *URLRepository) IncrementClick(shortCode string) error {
	query := "UPDATE urls SET clicks = clicks + 1 WHERE short_code = ?"
	_, err := r.db.Exec(query, shortCode)
	return err
}