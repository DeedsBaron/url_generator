package repository

import (
	"sync"

	"github.com/DeedsBaron/url_generator/internal/models"
)

type repo struct {
	sync.Mutex
	links map[string]*models.Url
}

func NewRepo() *repo {
	return &repo{
		links: make(map[string]*models.Url),
	}
}
