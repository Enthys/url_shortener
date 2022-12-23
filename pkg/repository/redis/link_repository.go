package redis

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Enthys/url_shortener/pkg/repository"
	redis_pkg "github.com/go-redis/redis/v9"
)

type redisRepository struct {
	client *redis_pkg.Client
}

// generateConfig creates a configuration with provides details on how to connect to the Redis cache.
// It uses 3 environment variables(`DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_PASS`) of which only 2 are required(
// `DATABASE_HOST`, `DATABASE_PASS`).
//
// If either `DATABASE_HOST` or `DATABASE_PASS` are not provided a generic `errors.Error` error will be returned.
//
// If no `DATABASE_PASS` is provided no password will be used for the connection to the Redis cache.
func generateConfig() (*redis_pkg.Options, error) {
	var host, port string
	if host = os.Getenv("DATABASE_HOST"); host == "" {
		return nil, errors.New("environment variable 'DATABASE_HOST' is missing")
	}

	if port = os.Getenv("DATABASE_PORT"); host == "" {
		return nil, errors.New("environment variable 'DATABASE_PORT' is missing")
	}

	return &redis_pkg.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: os.Getenv("DATABASE_PASS"),
	}, nil
}

func NewRedisRepository() (repository.LinkRepository, error) {
	config, err := generateConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to Redis. Error: %w", err)
	}

	client := redis_pkg.NewClient(config)
	return &redisRepository{
		client: client,
	}, nil
}

func generateIdKey(id string) string {
	return fmt.Sprintf("id_%s", id)
}

func generateLinkKey(link string) string {
	return fmt.Sprintf("link_%s", link)
}

// GetById retrieves from the Redis storage the link which corresponds to the given ID if such exists.
//
// If no link corresponds to the provided ID then a `repository.ErrorLinkNotFound` will be returned.
//
// If the retrieval of the record fails then a generic `errors.Error` will be returned.
func (r *redisRepository) GetById(id string) (string, error) {
	link, err := r.client.Get(context.Background(), generateIdKey(id)).Result()
	if err != nil {
		if errors.Is(err, redis_pkg.Nil) {
			return "", repository.ErrorLinkNotFound{}
		} else {
			return "", fmt.Errorf("retrieval of link by id failed. Error: %w", err)
		}
	}

	if link == "" {
	}

	return link, nil
}

// GetLinkId retrieves from the Redis storage the ID of the given link if such is found.
//
// If no record is found for the given link then a `repository.ErrorPathNotFound` will be returned as the second
// argument.
//
// If the retrieval of the record fails a generic `errors.Error` will be returned.
func (r *redisRepository) GetLinkId(link string) (string, error) {
	id, err := r.client.Get(context.Background(), generateLinkKey(link)).Result()
	if err != nil {
		if errors.Is(err, redis_pkg.Nil) {
			return "", repository.ErrorPathNotFound{}
		} else {
			return "", fmt.Errorf("retrieval of link path failed. Error: %w", err)
		}
	}

	return id, nil
}

// StoreLink stores 2 records in Redis. One of the records is for retrieving the link by the given ID and the other is
// intended to be used for quick lookup if the link is already in the storage and reuse them.
// The order in which they are stored is the record by which to retrieve the link and then the record through which to
// reuse link ids. The keys stored in the Redis cache are in the form of:
//
//	id_{link_hash}: https://example.com
//	link_{link}: {link_hash}
//
// If the insertion of the link key fails then a `repository.ErrorLinkSaveFailure` will be returned.
//
// If the insertion of the link lookup key fails then a rollback will be attempted.
func (r *redisRepository) StoreLink(id, link string) error {
	// Storing the link record first due to it being with higher priority
	_, err := r.client.Set(context.Background(), generateIdKey(id), link, 0).Result()
	if err != nil {
		return repository.ErrorLinkSaveFailure{Err: err}
	}

	_, err = r.client.Set(context.Background(), generateLinkKey(link), id, 0).Result()
	if err != nil {
		saveErr := repository.ErrorPathSaveFailure{Err: err}

		r.client.Del(context.Background(), generateLinkKey(link))

		return saveErr
	}

	return nil
}
