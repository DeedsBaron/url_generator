package url_generator_service

import (
	"context"
	"strconv"

	"github.com/DeedsBaron/url_generator/internal/config"
	"github.com/DeedsBaron/url_generator/pkg/errz"
)

const (
	defaultHost = "http://localhost:"
	route       = "/v1/url/get"
)

func (s *urlGeneratorService) CreateUrl(ctx context.Context, inputStr string) (string, error) {
	var generatedSeq string

	// примитивный способ решения коллизий - по-хорошему нужен chained-list
	for {
		generatedSeq = s.generateUrl(ctx)
		if !s.repo.CheckIfExists(ctx, generatedSeq) {
			break
		}
	}

	err := s.repo.CreateUrl(ctx, inputStr, generatedSeq)
	if err != nil {
		return "", errz.Wrap(err, "can't create url in repo")
	}

	return makeUrl(generatedSeq), nil
}

func makeUrl(generatedHash string) string {
	return defaultHost + strconv.Itoa(config.Data.HttpPort) + route + "/" + generatedHash
}
