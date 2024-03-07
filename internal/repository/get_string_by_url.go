package repository

import (
	"context"

	"github.com/DeedsBaron/url_generator/pkg/errz"
)

var (
	errUrlIsAlreadyUsed = errz.Error(errz.InvalidArgument, "url is already used")
	errUrlIsNotFound    = errz.Error(errz.NotFound, "url is not found")
)

func (r *repo) GetStringByUrl(_ context.Context, id string) (string, error) {
	r.Lock()
	defer r.Unlock()
	if url, ok := r.links[id]; ok {
		if url.Used {
			return "", errUrlIsAlreadyUsed
		}

		return url.InputString, nil
	}

	return "", errUrlIsNotFound
}
