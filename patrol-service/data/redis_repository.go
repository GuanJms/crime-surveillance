package data

import (
	"context"
	"fmt"
	"log"
	"patrolServiceApp/ptr"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	fastdbTimeout   = time.Second * 2
	redisPort       = "6379"
	nullLocation    = -999.0
	locationTTL     = time.Second * 30
	patrolDirtyKeys = "dirty_keys_patrol"
)

type RedisRepository struct {
	client      *redis.Client
	persistRepo Repository
}

func (r *RedisRepository) Close() {
	r.client.Close()
}

func NewRedisRepository(repo Repository) *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("redis:%s", redisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisRepository{
		client:      rdb,
		persistRepo: repo,
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

	repo.client.SAdd(ctx, patrolDirtyKeys, key) // add the key to the dirty keys set

	return nil
}

func (repo *RedisRepository) PersistUpdatePatrolLocation(patrolId string, location *Location) error {
	log.Printf("Persisting patrol location for patrolId: %s", patrolId)
	if location == nil {
		return fmt.Errorf("persist location update received empty nil location")
	}
	var req UpdatePatrolInfoRequest
	req.UserId = patrolId
	req.Street = location.Street
	req.City = location.City
	req.State = location.State
	req.Latitude = location.Latitude
	req.Longitude = location.Longitude

	err := repo.persistRepo.PatchPatrolInfo(&req)
	if err != nil {
		return err
	}
	return nil
}

func (*RedisRepository) getUserLocationKey(id string) string {
	return fmt.Sprintf("user:%s:location", id)
}

func (*RedisRepository) parseRedisUserLocationKeyForUserID(key string) string {
	parts := strings.Split(key, ":")
	log.Printf("Parsing redis user location key for userID: %s", parts)
	return parts[1] // TODO: hardcoded wiring for retriving userID
}

func (repo *RedisRepository) SyncPatrolLocation() error {
	log.Printf("Syncing patrol location")
	dirtyIDs, err := repo.ReadRemovePatrolDirtyKeys()
	if err != nil {
		return err
	}
	log.Printf("Dirty IDs: %v", dirtyIDs)
	for _, redisKey := range dirtyIDs {
		log.Printf("Processing dirty ID: %s", redisKey)
		ctx, cancel := context.WithTimeout(context.Background(), fastdbTimeout)
		defer cancel()
		var l Location
		if err := repo.client.HGetAll(ctx, redisKey).Scan(&l); err != nil {
			log.Printf("failed to retreive entity: %v", err)
		} else {
			err = repo.PersistUpdatePatrolLocation(
				repo.parseRedisUserLocationKeyForUserID(redisKey),
				&l,
			)
			if err != nil {
				log.Printf("failed to persist update: %v", err)
			}
		}
	}

	return nil
}

func (repo *RedisRepository) ReadRemovePatrolDirtyKeys() ([]string, error) {
	// Lua script to read and remove the dirty keys with atomicity

	flushScript := redis.NewScript(`
		local ids = redis.call("SMEMBERS", KEYS[1])
		redis.call("DEL", KEYS[1])
		return ids
		`)

	ctx, cancel := context.WithTimeout(context.Background(), fastdbTimeout)
	defer cancel()

	dirtyIDs, err := flushScript.Run(ctx, repo.client, []string{patrolDirtyKeys}).StringSlice()
	if err != nil {
		return nil, err
	}

	return dirtyIDs, err
}
