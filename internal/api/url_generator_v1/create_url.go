package url_generator_v1

import (
	"context"

	desc "github.com/DeedsBaron/url_generator/pkg/url_generator"
)

func (i *Implementation) CreateUrl(ctx context.Context, req *desc.CreateUrlReq) (*desc.CreateUrlResponse, error) {
	createdUrl, err := i.urlGeneratorService.CreateUrl(ctx, req.GetInputString())
	if err != nil {
		return nil, err
	}

	return &desc.CreateUrlResponse{Url: createdUrl}, nil
}
