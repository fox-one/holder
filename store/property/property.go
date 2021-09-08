package property

import (
	"context"

	"github.com/fox-one/pkg/property"
	"github.com/patrickmn/go-cache"
)

func Memory() property.Store {
	return &memoryStore{
		c: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

type memoryStore struct {
	c *cache.Cache
}

func (s *memoryStore) Get(_ context.Context, key string) (property.Value, error) {
	var value property.Value
	if v, ok := s.c.Get(key); ok {
		value = property.Parse(v)
	}

	return value, nil
}

func (s *memoryStore) Save(_ context.Context, key string, value interface{}) error {
	s.c.SetDefault(key, value)
	return nil
}

func (s *memoryStore) Expire(_ context.Context, key string) error {
	s.c.Delete(key)
	return nil
}

func (s *memoryStore) List(_ context.Context) (map[string]property.Value, error) {
	values := make(map[string]property.Value)
	items := s.c.Items()
	for k, v := range items {
		values[k] = property.Parse(v)
	}

	return values, nil
}
