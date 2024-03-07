package repository

import (
	"context"
)

func (r *repo) MarkUrlAsUsed(_ context.Context, id string) error {
	r.Lock()
	defer r.Unlock()

	r.links[id].Used = true
	return nil
}
