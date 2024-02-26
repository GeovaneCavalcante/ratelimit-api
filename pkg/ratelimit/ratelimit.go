package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const minScore = "min"

type RateLimiterInterface interface {
	CountEventsBeforeCurrent(ctx context.Context, key string, currentTimestamp int64) (int64, error)
	RemoveExpiredEvents(ctx context.Context, key string, recentTimestamp, intervalSecund float64) error
	AddEvent(ctx context.Context, key string, timestamp int64) error
	Limiter(ctx context.Context, key string, opt *Options) (bool, error)
}

type Options struct {
	NameSpace      string
	MaxInInterval  int64
	IntervalSecund int64
}
type RateLimiter struct {
	EventStorage EventStorageInterface
	Options
}

func New(es EventStorageInterface, ns string, max int64, inter time.Duration) (*RateLimiter, error) {

	return &RateLimiter{
		EventStorage: es,
		Options:      Options{NameSpace: ns, MaxInInterval: max, IntervalSecund: int64(inter.Seconds())},
	}, nil
}

func (rl *RateLimiter) CountEventsBeforeCurrent(ctx context.Context, key string, currentTimestamp int64) (int64, error) {
	count, err := rl.EventStorage.CountRange(ctx, key, minScore, strconv.FormatInt(currentTimestamp+1, 10))
	if err != nil {
		return 0, fmt.Errorf("error when counting the number of events: %w", err)
	}
	return count, nil
}

func (rl *RateLimiter) RemoveExpiredEvents(ctx context.Context, key string, recentTimestamp, intervalSecund float64) error {
	oldestEvent, err := rl.EventStorage.FindRangeWithScores(ctx, key, 0, 0)
	if err != nil {
		return fmt.Errorf("error when finding the oldest event: %w", err)
	}
	if len(oldestEvent) == 0 {
		return nil
	}

	oldestTimestamp := oldestEvent[0].Score

	if recentTimestamp-oldestTimestamp > intervalSecund {
		err := rl.EventStorage.RemoveRangeByScore(ctx, key, minScore, strconv.FormatFloat(recentTimestamp, 'f', -1, 64))
		if err != nil {
			return fmt.Errorf("error when removing expired events remove range by score: %w", err)
		}
	}

	return nil
}

func (rl *RateLimiter) AddEvent(ctx context.Context, key string, timestamp int64) error {
	id := uuid.New().String()
	value := fmt.Sprintf("event:%s:%d", id, timestamp)
	score := float64(timestamp)

	_, err := rl.EventStorage.Add(ctx, key, &Event{
		Score: score,
		Value: value,
	})

	if err != nil {
		return err
	}

	return nil
}

func (rl *RateLimiter) Limiter(ctx context.Context, key string, opt *Options) (bool, error) {

	timestamp := time.Now().Unix()

	nameSpace := chooseString(rl.NameSpace, opt, func(o *Options) string { return o.NameSpace })
	maxInInterval := chooseInt64(rl.MaxInInterval, opt, func(o *Options) int64 { return o.MaxInInterval })
	IntervalSecund := chooseInt64(rl.IntervalSecund, opt, func(o *Options) int64 { return o.IntervalSecund })

	bucketName := fmt.Sprintf("%s:%s", nameSpace, key)

	c, err := rl.CountEventsBeforeCurrent(ctx, bucketName, timestamp)

	if err != nil {
		return false, fmt.Errorf("error when counting the number of events: %w", err)
	}

	if c < maxInInterval {
		err := rl.AddEvent(ctx, bucketName, timestamp)
		if err != nil {
			return false, fmt.Errorf("error when adding event: %w", err)
		}
		return false, nil
	}

	err = rl.RemoveExpiredEvents(ctx, bucketName, float64(timestamp), float64(IntervalSecund))

	if err != nil {
		return false, fmt.Errorf("error when removing expired events: %w", err)
	}

	return true, nil
}

func chooseString(defaultVal string, opt *Options, optSelector func(*Options) string) string {
	if opt != nil {
		optVal := optSelector(opt)
		if optVal != "" {
			return optVal
		}
	}
	return defaultVal
}

func chooseInt64(defaultVal int64, opt *Options, optSelector func(*Options) int64) int64 {
	if opt != nil {
		optVal := optSelector(opt)
		if optVal != 0 {
			return optVal
		}
	}
	return defaultVal
}
