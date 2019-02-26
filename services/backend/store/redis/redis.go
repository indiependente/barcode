package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type Store struct {
	client *redis.Client
}

// InitStore returns a pointer to a new redis.Store
func InitStore(addr string) *Store {
	return &Store{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (rs *Store) Get(id string) ([]byte, error) {
	val, err := rs.client.Get(id).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Could not get data")
	}
	return []byte(val), nil
}

func (rs *Store) Put(id string, data []byte) error {
	err := rs.client.Set(id, data, 0).Err()
	if err != nil {
		return errors.Wrap(err, "Could not set data")
	}
	return nil
}
