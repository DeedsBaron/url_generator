package url_generator_service

import (
	"context"

	"github.com/DeedsBaron/url_generator/pkg/errz"
)

func (s *urlGeneratorService) GetStringByUrl(ctx context.Context, id string) (string, error) {
	resString, err := s.repo.GetStringByUrl(ctx, id)
	if err != nil {
		return "", errz.Wrap(err, "urlGenerator: can't get string by url")
	}

	err = s.repo.MarkUrlAsUsed(ctx, id)
	if err != nil {
		return "", errz.Wrap(err, "urlGenerator: can't mark url as used")
	}

	return resString, nil
}
