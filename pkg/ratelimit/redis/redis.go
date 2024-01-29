package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/GeovaneCavalcante/rate-limit-api/pkg/ratelimit"
	"github.com/go-redis/redis/v8"
)

type RedisEventStorage struct {
	RedisClient *redis.Client
}

func NewRedisEventStorage(rc *redis.Client) *RedisEventStorage {
	return &RedisEventStorage{
		RedisClient: rc,
	}
}

func parserMinScore(minScore string) string {
	if minScore == "min" {
		return "-inf"
	}
	return minScore
}

func (res *RedisEventStorage) CountRange(ctx context.Context, key, min, max string) (int64, error) {
	min = parserMinScore(min)
	count, err := res.RedisClient.ZCount(ctx, key, min, max).Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (res *RedisEventStorage) FindRangeWithScores(ctx context.Context, key string, start, stop int64) ([]*ratelimit.Event, error) {
	eventsRd, err := res.RedisClient.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	var events []*ratelimit.Event
	for i, eventRd := range eventsRd {
		events = append(events, &ratelimit.Event{
			ID:    fmt.Sprint(i),
			Score: eventRd.Score,
			Value: eventRd.Member.(string),
		})
	}
	return events, nil
}

func (res *RedisEventStorage) RemoveRangeByScore(ctx context.Context, key, min, max string) error {
	min = parserMinScore(min)
	err := res.RedisClient.ZRemRangeByScore(ctx, key, min, max).Err()
	if err != nil {
		return err
	}
	return nil
}

func (res *RedisEventStorage) Add(ctx context.Context, key string, events ...*ratelimit.Event) ([]*ratelimit.Event, error) {
	var zEvents []*redis.Z
	for _, event := range events {
		zEvents = append(zEvents, &redis.Z{
			Score:  event.Score,
			Member: event.Value,
		})
	}
	err := res.RedisClient.ZAdd(ctx, key, zEvents...).Err()
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (res *RedisEventStorage) SetEventTLL(ctx context.Context, key string, ttl time.Duration) error {
	err := res.RedisClient.Expire(ctx, key, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}
