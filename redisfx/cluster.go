package redisfx

import (
	"context"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/yankeguo/ufx"
)

type ClusterParams struct {
	URL string `json:"url"`
}

func DecodeClusterParams(conf ufx.Conf) (params ClusterParams, err error) {
	err = conf.Bind(&params, "redis", "cluster")
	return
}

func NewClusterOptions(params ClusterParams) (*redis.ClusterOptions, error) {
	return redis.ParseClusterURL(params.URL)
}

func NewClusterClient(opts *redis.ClusterOptions) (client *redis.ClusterClient, err error) {
	client = redis.NewClusterClient(opts)
	if err = redisotel.InstrumentTracing(client); err != nil {
		return
	}
	if err = redisotel.InstrumentMetrics(client); err != nil {
		return
	}
	return
}

func AddCheckerForClusterClient(client *redis.ClusterClient, v ufx.Prober) {
	v.AddChecker("redis-cluster", func(ctx context.Context) error {
		return client.Ping(ctx).Err()
	})
}
