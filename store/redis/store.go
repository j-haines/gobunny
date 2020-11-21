package redis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"gobunny/store"
	"gobunny/store/model"
)

// Config contains parameters for constructing and configuring a Redis-backed Store
type Config struct {
	HostAddress string
	Password    string
	DatabaseID  int
}

type redisStore struct {
	config Config
	client *redis.Client
	ctx    context.Context
}

// NewStore returns a store.Store implementation backed by Redis
func NewStore(ctx context.Context, config Config) (store.Store, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.HostAddress,
		Password: config.Password,
		DB:       config.DatabaseID,
	})

	return &redisStore{
		config: config,
		client: client,
		ctx:    ctx,
	}, nil
}

func (s *redisStore) Create(key model.Key, m model.Model) error {
	success, err := s.client.SetNX(s.ctx, key.String(), m.Value().String(), 0).Result()
	if err != nil {
		return err
	}

	if !success {
		return store.ErrKeyAlreadyExists
	}

	return nil
}

func (s *redisStore) Delete(key model.Key) error {
	if err := s.client.Del(s.ctx, key.String()).Err(); err != nil {
		return err
	}

	return nil
}

func (s *redisStore) Exists(key model.Key) (bool, error) {
	exists, err := s.client.Exists(s.ctx, key.String()).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func (s *redisStore) Get(key model.Key, m model.Model) error {
	value, err := s.client.Get(s.ctx, key.String()).Result()
	if err != nil {
		if err == redis.Nil {
			return store.ErrNotFound
		}

		return err
	}

	if err := m.Unmarshal(value); err != nil {
		return err
	}

	return nil
}

func (s *redisStore) Update(key model.Key, m model.Model) error {
	if err := s.client.Set(s.ctx, key.String(), m.Value().String(), 0).Err(); err != nil {
		return err
	}

	return nil
}
