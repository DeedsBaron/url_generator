package repository

import (
	"context"

	"github.com/DeedsBaron/url_generator/internal/models"
)

func (r *repo) CreateUrl(_ context.Context, inputStr, url string) error {
	r.Lock()
	defer r.Unlock()
	r.links[url] = &models.Url{
		InputString: inputStr,
		Used:        false,
	}
	return nil
}
