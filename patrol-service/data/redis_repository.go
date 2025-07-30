package data

import (
	"context"
	"fmt"
	"patrolServiceApp/ptr"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	fastdbTimeout = time.Second * 2
	redisPort     = "6379"
	nullLocation  = -999.0
	locationTTL   = time.Second * 30
)

type RedisRepository struct {
	client *redis.Client
}

func (r *RedisRepository) Close() {
	r.client.Close()
}

func NewRedisRepository() *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("redis:%s", redisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisRepository{
		client: rdb,
	}
}

func (repo *RedisRepository) UpdatePatrolLocation(patrolId string, location *Location) error {
	if repo == nil {
		return fmt.Errorf("redis repository is nil")
	}

	if location == nil {
		return fmt.Errorf("no location is provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), fastdbTimeout)
	defer cancel()

	key := repo.getUserLocationKey(patrolId)

	if location.Latitude == nil {
		location.Latitude = ptr.Of(nullLocation)
	}

	if location.Longitude == nil {
		location.Longitude = ptr.Of(nullLocation)
	}

	if location.Timestamp == nil {
		location.Timestamp = ptr.Of(time.Now())
	}

	if _, err := repo.client.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HMSet(ctx, key, *location)
		rdb.Expire(ctx, key, locationTTL)
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (*RedisRepository) getUserLocationKey(id string) string {
	return fmt.Sprintf("user:%s:location", id)
}
