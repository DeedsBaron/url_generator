package url_generator_v1

import (
	"context"

	desc "github.com/DeedsBaron/url_generator/pkg/url_generator"
)

type UrlGeneratorService interface {
	CreateUrl(ctx context.Context, str string) (string, error)
	GetStringByUrl(ctx context.Context, id string) (string, error)
}
type Implementation struct {
	desc.UnimplementedUrlGeneratorServer
	urlGeneratorService UrlGeneratorService
}

func NewUrlGeneratorV1(urlGeneratorService UrlGeneratorService) *Implementation {
	return &Implementation{
		desc.UnimplementedUrlGeneratorServer{},
		urlGeneratorService,
	}
}
