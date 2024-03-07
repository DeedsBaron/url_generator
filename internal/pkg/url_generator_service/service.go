package url_generator_service

import "context"

type repo interface {
	CheckIfExists(_ context.Context, url string) bool
	GetStringByUrl(_ context.Context, id string) (string, error)
	CreateUrl(_ context.Context, inputStr, url string) error
	MarkUrlAsUsed(_ context.Context, id string) error
}
type urlGeneratorService struct {
	repo repo
}

func New(repo repo) *urlGeneratorService {
	return &urlGeneratorService{
		repo: repo,
	}
}
