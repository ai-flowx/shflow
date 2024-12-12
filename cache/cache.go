package cache

import (
	"context"

	"github.com/hashicorp/go-hclog"
)

type Cache interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type Config struct {
	Logger hclog.Logger
}

type cache struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Cache {
	return &cache{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (c *cache) Init(_ context.Context) error {
	return nil
}

func (c *cache) Deinit(_ context.Context) error {
	return nil
}

func (c *cache) Run(ctx context.Context) error {
	return nil
}
