package repository

import "context"

func (r *repo) CheckIfExists(_ context.Context, url string) bool {
	r.Lock()
	defer r.Unlock()

	_, ok := r.links[url]
	return ok
}
