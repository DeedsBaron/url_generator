package url_generator_v1

import (
	"context"

	desc "github.com/DeedsBaron/url_generator/pkg/url_generator"
)

func (i *Implementation) GetStringByUrl(ctx context.Context, req *desc.GetStringByUrlRequest) (*desc.GetStringByUrlResponse, error) {
	resString, err := i.urlGeneratorService.GetStringByUrl(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetStringByUrlResponse{ResultString: resString}, nil
}
