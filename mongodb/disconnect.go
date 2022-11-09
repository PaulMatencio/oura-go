package mongodb

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

func (r *MongoDB) DisConnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := r.Client.Disconnect(ctx); err != nil {
		log.Error().Stack().Err(err)
	}
}
