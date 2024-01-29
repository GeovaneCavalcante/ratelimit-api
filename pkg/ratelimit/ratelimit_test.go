package ratelimit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GeovaneCavalcante/rate-limit-api/pkg/ratelimit"
	mock_storage "github.com/GeovaneCavalcante/rate-limit-api/pkg/ratelimit/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RateLimiterTestSuite struct {
	suite.Suite
	EventStorageMock *mock_storage.MockEventStorageInterface
}

func (suite *RateLimiterTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.EventStorageMock = mock_storage.NewMockEventStorageInterface(ctrl)
}

func (suite *RateLimiterTestSuite) TestLimiter() {
	suite.Run("should return true when the number of events is greater than the maximum allowed", func() {
		suite.EventStorageMock.EXPECT().CountRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), errors.New("error"))

		rl, err := ratelimit.New(suite.EventStorageMock, "test", 1, 1)
		if err != nil {
			suite.FailNow(err.Error())
		}

		value, err := rl.Limiter(context.Background(), "test", nil)
		assert.Equal(suite.T(), err.Error(), "error when counting the number of events: error when counting the number of events: error")
		assert.False(suite.T(), value)
	})

	suite.Run("should return false when the number of events is greater than the maximum allowed", func() {
		suite.EventStorageMock.EXPECT().CountRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(2), nil)
		suite.EventStorageMock.EXPECT().FindRangeWithScores(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))

		rl, err := ratelimit.New(suite.EventStorageMock, "test", 1, 1)
		if err != nil {
			suite.FailNow(err.Error())
		}

		value, err := rl.Limiter(context.Background(), "test", nil)
		assert.Equal(suite.T(), err.Error(), "error when removing expired events: error when finding the oldest event: error")
		assert.False(suite.T(), value)
	})

	suite.Run("should return true when user escapes the limitation", func() {

		suite.EventStorageMock.EXPECT().CountRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(2), nil)
		suite.EventStorageMock.EXPECT().FindRangeWithScores(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*ratelimit.Event{
			{
				ID:    "test",
				Score: float64(time.Now().Unix() + 1),
				Value: "test",
			},
		}, nil)

		rl, err := ratelimit.New(suite.EventStorageMock, "test", 1, 1)
		if err != nil {
			suite.FailNow(err.Error())
		}

		value, err := rl.Limiter(context.Background(), "test", nil)
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), value)
	})

	suite.Run("should return an error when adding an event fails", func() {

		suite.EventStorageMock.EXPECT().CountRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(2), nil)
		suite.EventStorageMock.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))

		rl, err := ratelimit.New(suite.EventStorageMock, "test", 6, 1)
		if err != nil {
			suite.FailNow(err.Error())
		}

		value, err := rl.Limiter(context.Background(), "test", nil)

		assert.Equal(suite.T(), err.Error(), "error when adding event: error")
		assert.False(suite.T(), value)
	})

	suite.Run("should return an true when adding an event", func() {

		suite.EventStorageMock.EXPECT().CountRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(2), nil)
		suite.EventStorageMock.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

		rl, err := ratelimit.New(suite.EventStorageMock, "test", 6, 1)
		if err != nil {
			suite.FailNow(err.Error())
		}

		value, err := rl.Limiter(context.Background(), "test", &ratelimit.Options{
			NameSpace:      "test",
			MaxInInterval:  6,
			IntervalSecund: 1,
		})

		assert.NoError(suite.T(), err)
		assert.False(suite.T(), value)
	})

	suite.Run("should return true when user escapes the limitation", func() {

		suite.EventStorageMock.EXPECT().CountRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(2), nil)
		suite.EventStorageMock.EXPECT().FindRangeWithScores(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*ratelimit.Event{}, nil)

		rl, err := ratelimit.New(suite.EventStorageMock, "test", 1, 1)
		if err != nil {
			suite.FailNow(err.Error())
		}

		value, err := rl.Limiter(context.Background(), "test", nil)
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), value)
	})

	suite.Run("should return error when there is an error in removal", func() {

		suite.EventStorageMock.EXPECT().CountRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(2), nil)
		suite.EventStorageMock.EXPECT().FindRangeWithScores(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*ratelimit.Event{
			{
				ID:    "test",
				Score: 0,
				Value: "test",
			},
		}, nil)

		suite.EventStorageMock.EXPECT().RemoveRangeByScore(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))

		rl, err := ratelimit.New(suite.EventStorageMock, "test", 1, 1)
		if err != nil {
			suite.FailNow(err.Error())
		}

		value, err := rl.Limiter(context.Background(), "test", nil)
		assert.Equal(suite.T(), err.Error(), "error when removing expired events: error when removing expired events remove range by score: error")
		assert.False(suite.T(), value)
	})
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterTestSuite))
}
