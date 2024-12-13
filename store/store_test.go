package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	apiTestStore   = "http://127.0.0.1:8082/"
	tokenTestStore = "token"
	valueTestStore = "testStore"
)

func initStoreTest(_ context.Context) store {
	cfg := Config{}
	cfg.Provider = ProviderChroma
	cfg.Api = apiTestStore
	cfg.Token = tokenTestStore

	return store{
		cfg: &cfg,
		st:  storeList[cfg.Provider],
	}
}

func TestStoreInit(t *testing.T) {
	ctx := context.Background()

	s := initStoreTest(ctx)

	err := s.Init(ctx)
	assert.Equal(t, nil, err)
}

func TestStoreDeinit(t *testing.T) {
	ctx := context.Background()

	s := initStoreTest(ctx)

	err := s.Deinit(ctx)
	assert.Equal(t, nil, err)
}

func TestStoreReset(t *testing.T) {
	ctx := context.Background()

	s := initStoreTest(ctx)

	_ = s.Init(ctx)

	defer func(s *store, ctx context.Context) {
		_ = s.Deinit(ctx)
	}(&s, ctx)

	err := s.Reset(ctx)
	assert.Equal(t, nil, err)
}

func TestStoreSave(t *testing.T) {
	ctx := context.Background()

	s := initStoreTest(ctx)

	_ = s.Init(ctx)

	defer func(s *store, ctx context.Context) {
		_ = s.Deinit(ctx)
	}(&s, ctx)

	value := valueTestStore
	meta := map[string]interface{}{
		"task":    "testTask",
		"quality": 0.5,
	}
	agent := "testAgent"

	err := s.Save(ctx, value, meta, agent)
	assert.Equal(t, nil, err)
}

func TestStoreSearch(t *testing.T) {
	ctx := context.Background()

	s := initStoreTest(ctx)

	_ = s.Init(ctx)

	defer func(s *store, ctx context.Context) {
		_ = s.Deinit(ctx)
	}(&s, ctx)

	query := valueTestStore
	limit := 3
	threshold := 0.35

	_, err := s.Search(ctx, query, limit, threshold)
	assert.Equal(t, nil, err)
}
