package memory

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ai-flowx/flowx/store"
)

type ShortTerm struct {
	Store store.Store
}

func (s *ShortTerm) Init(ctx context.Context) error {
	if err := s.Store.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init\n")
	}

	return nil
}

func (s *ShortTerm) Deinit(ctx context.Context) error {
	if err := s.Store.Deinit(ctx); err != nil {
		return errors.Wrap(err, "failed to deinit\n")
	}

	return nil
}

func (s *ShortTerm) Reset(ctx context.Context) error {
	if err := s.Store.Reset(ctx); err != nil {
		return errors.Wrap(err, "failed to reset\n")
	}

	return nil
}

func (s *ShortTerm) Save(ctx context.Context, value interface{}, meta map[string]interface{}, agent string) error {
	if err := s.Store.Save(ctx, value, meta, agent); err != nil {
		return errors.Wrap(err, "failed to save\n")
	}

	return nil
}

func (s *ShortTerm) Search(ctx context.Context, query string, limit int, threshold float64) ([]interface{}, error) {
	buf, err := s.Store.Search(ctx, query, limit, threshold)
	if err != nil {
		return nil, errors.Wrap(err, "failed to search\n")
	}

	return buf, nil
}