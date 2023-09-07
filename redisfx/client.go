package redisfx

import (
	"context"
	"github.com/guoyk93/ufx"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type Params struct {
	URL string `json:"url" default:"redis://localhost:6379/0" validate:"required,url"`
}

func DecodeParams(conf ufx.Conf) (params Params, err error) {
	err = conf.Bind(&params, "redis")
	return
}

func NewOptions(params Params) (*redis.Options, error) {
	return redis.ParseURL(params.URL)
}

func NewClient(opts *redis.Options) (client *redis.Client, err error) {
	client = redis.NewClient(opts)
	if err = redisotel.InstrumentTracing(client); err != nil {
		return
	}
	if err = redisotel.InstrumentMetrics(client); err != nil {
		return
	}
	return
}

func AddCheckerForClient(client *redis.Client, v ufx.Prober) {
	v.AddChecker("redis", func(ctx context.Context) error {
		return client.Ping(ctx).Err()
	})
}
