package url_generator_service

import (
	"context"
	"math/rand"
	"time"

	"github.com/DeedsBaron/url_generator/internal/config"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func (s *urlGeneratorService) generateUrl(_ context.Context) string {
	rand.NewSource(time.Now().UnixNano())
	b := make([]byte, config.Data.UrlLen)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}
