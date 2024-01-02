package gormfx

import (
	"context"

	"github.com/yankeguo/ufx"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type Params struct {
	Debug bool `json:"debug"`
}

func ParamsFromConf(conf ufx.Conf) (params Params, err error) {
	err = conf.Bind(&params, "gorm")
	return
}

func NewConfig() *gorm.Config {
	return &gorm.Config{}
}

func NewClient(d gorm.Dialector, c *gorm.Config, params Params) (db *gorm.DB, err error) {
	if db, err = gorm.Open(d, c); err != nil {
		return
	}
	if err = db.Use(tracing.NewPlugin()); err != nil {
		return
	}
	if params.Debug {
		db = db.Debug()
	}
	return
}

func AddCheckerForClient(db *gorm.DB, v ufx.Prober) {
	v.AddChecker("gorm", func(ctx context.Context) error {
		return db.WithContext(ctx).Select("SELECT 1").Error
	})
}
