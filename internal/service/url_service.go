package service

import (
	"errors"
	"math/rand"
	"net/url"
	"time"
	"url-shortener/internal/model"
	"url-shortener/internal/repository"
)

type URLService struct {
	repo *repository.URLRepository
	baseURL string
}

func NewURLService(repo *repository.URLRepository, baseURL string) *URLService {
	return &URLService{repo: repo, baseURL: baseURL}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(length int) string {
	seeded := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)

	for i := range code {
		code[i] = charset[seeded.Intn(len(charset))]
	}

	return string(code)
}

func validateURL(rawURL string) error {
	parsed, err := url.ParseRequestURI(rawURL)

	if err != nil {
		return errors.New("Invalid URL Format")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("URL must start with http or https")
	}
	return nil
}

func (s *URLService) Shorten(originalURL string) (*model.ShortenResponse, error) {
	if err := validateURL(originalURL); err != nil {
		return nil, err
	}

	shortCode := generateShortCode(6)

	url := &model.URL{
		OriginalURL: originalURL,
		ShortCode: shortCode,
	}

	if err := s.repo.Save(url); err != nil{
		return nil, errors.New("Failed to save URL")
	}

	return &model.ShortenResponse{
		ShortCode: shortCode,
		ShortURL: s.baseURL + "/" + shortCode,
		OriginalURL: originalURL,
	}, nil
}

func (s *URLService) Redirect(shortCode string) (string, error) {
	url, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return "", errors.New("short code not found")
	}

	if err := s.repo.IncrementClick(shortCode); err != nil {
		return "", errors.New("failed to update click count")
	}

	return url.OriginalURL, nil
}

func (s *URLService) GetStats(shortCode string) (*model.StatsResponse, error){
	url, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return nil, errors.New("short code not found")
	}

	return &model.StatsResponse{
		ShortCode: url.ShortCode,
		OriginalURL: url.OriginalURL,
		Clicks: url.Clicks,
		CreatedAt: url.CreatedAt,
	}, nil
}