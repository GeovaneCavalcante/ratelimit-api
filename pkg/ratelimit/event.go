package ratelimit

import (
	"context"
	"time"
)

type Event struct {
	ID    string
	Score float64
	Value string
}

type EventStorageInterface interface {
	CountRange(ctx context.Context, key, min, max string) (int64, error)
	FindRangeWithScores(ctx context.Context, key string, start, stop int64) ([]*Event, error)
	RemoveRangeByScore(ctx context.Context, key, min, max string) error
	Add(ctx context.Context, key string, events ...*Event) ([]*Event, error)
	SetEventTLL(ctx context.Context, key string, ttl time.Duration) error
}
