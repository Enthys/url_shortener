package redis

import (
	"context"
	"errors"
	"fmt"
	"internal/itoa"
	"os"

	"github.com/Enthys/url_shortener/pkg"
	"github.com/Enthys/url_shortener/pkg/repository"
	redisPkg "github.com/go-redis/redis/v9"
)

type redisRepository struct {
	client *redisPkg.Client
}

// generateConfig creates a configuration with provides details on how to connect to the Redis cache.
// It uses 3 environment variables(`DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_PASS`) of which only 2 are required(
// `DATABASE_HOST`, `DATABASE_PASS`).
//
// If either `DATABASE_HOST` or `DATABASE_PASS` are not provided a generic `errors.Error` error will be returned.
//
// If no `DATABASE_PASS` is provided no password will be used for the connection to the Redis cache.
func generateConfig() (*redisPkg.Options, error) {
	var host, port string
	if host = os.Getenv("DATABASE_HOST"); host == "" {
		return nil, errors.New("environment variable 'DATABASE_HOST' is missing")
	}

	if port = os.Getenv("DATABASE_PORT"); host == "" {
		return nil, errors.New("environment variable 'DATABASE_PORT' is missing")
	}

	return &redisPkg.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: os.Getenv("DATABASE_PASS"),
	}, nil
}

// NewRedisRepository creates a new repository which works with Redis. It uses the `generateConfig` function to generate
// the configuration for the connection to the Redis server.
func NewRedisRepository() (repository.LinkRepository, error) {
	config, err := generateConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to Redis. Error: %w", err)
	}

	client := redisPkg.NewClient(config)
	return &redisRepository{
		client: client,
	}, nil
}

// generateIdKey generates the key underwhich links are stored.
func generateIdKey(id string) string {
	return fmt.Sprintf("id_%s", id)
}

// generateLinkKey generates the key underwhich the id of the given link is stored.
func generateLinkKey(link string) string {
	return fmt.Sprintf("link_%s", link)
}

func (r *redisRepository) GenerateId() (string, error) {
	val, err := r.client.Incr(context.Background(), "id_gen").Result()
	if err != nil {
		return "", repository.ErrorIDGenerationFailed{}
	}

	return itoa.Itoa(int(val)) + pkg.RandomString(16), nil
}

// GetById retrieves from the Redis storage the link which corresponds to the given ID if such exists.
//
// If no link corresponds to the provided ID then a `repository.ErrorLinkNotFound` will be returned.
//
// If the retrieval of the record fails then a generic `errors.Error` will be returned.
func (r *redisRepository) GetById(id string) (string, error) {
	link, err := r.client.Get(context.Background(), generateIdKey(id)).Result()
	if err != nil {
		if errors.Is(err, redisPkg.Nil) {
			return "", repository.ErrorLinkNotFound{}
		} else {
			return "", repository.ErrorFailedRetrieval{Err: err}
		}
	}

	return link, nil
}

// GetLinkId retrieves from the Redis storage the ID of the given link if such is found.
//
// If no record is found for the given link then a `repository.ErrorIDNotFound` will be returned as the second
// argument.
//
// If the retrieval of the record fails a generic `errors.Error` will be returned.
func (r *redisRepository) GetLinkId(link string) (string, error) {
	id, err := r.client.Get(context.Background(), generateLinkKey(link)).Result()
	if err != nil {
		if errors.Is(err, redisPkg.Nil) {
			return "", repository.ErrorIDNotFound{}
		} else {
			return "", repository.ErrorFailedRetrieval{Err: err}
		}
	}

	return id, nil
}

// StoreLink stores 2 records in Redis. One of the records is for retrieving the link by the given ID and the other is
// intended to be used for quick lookup if the link is already in the storage and reuse them.
// The order in which they are stored is the record by which to retrieve the link and then the record through which to
// reuse link ids. The keys stored in the Redis cache are in the form of:
//
//	id_{link_hash}: {link}
//	link_{link}: {link_hash}
//
// If the insertion of the link key fails then a `repository.ErrorLinkSaveFailure` will be returned.
//
// If the insertion of the link lookup key fails then a rollback will be attempted and a `repository.ErrorIDSaveFailure`
// error will be returned.
func (r *redisRepository) StoreLink(id, link string) error {
	// Storing the link record first due to it being with higher priority
	_, err := r.client.Set(context.Background(), generateIdKey(id), link, 0).Result()
	if err != nil {
		return repository.ErrorLinkSaveFailure{Err: err}
	}

	_, err = r.client.Set(context.Background(), generateLinkKey(link), id, 0).Result()
	if err != nil {
		saveErr := repository.ErrorIDSaveFailure{Err: err}

		r.client.Del(context.Background(), generateLinkKey(link))

		return saveErr
	}

	return nil
}
