package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GeovaneCavalcante/rate-limit-api/configs"
	mock_ratelimit "github.com/GeovaneCavalcante/rate-limit-api/pkg/ratelimit/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RateLimiterTestSuite struct {
	suite.Suite
	RateLimitToken *mock_ratelimit.MockRateLimiterInterface
	RateLimitIp    *mock_ratelimit.MockRateLimiterInterface
	TokensConfig   []configs.TokenConfigLimit
}

func (suite *RateLimiterTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.RateLimitToken = mock_ratelimit.NewMockRateLimiterInterface(ctrl)
	suite.RateLimitIp = mock_ratelimit.NewMockRateLimiterInterface(ctrl)
	suite.TokensConfig = []configs.TokenConfigLimit{
		{
			Token:           "123",
			MaxRequests:     1,
			BlockTimeSecond: 1,
		},
	}
}

func (suite *RateLimiterTestSuite) TestRateLimiter() {
	suite.Run("should return request successfully", func() {
		suite.RateLimitToken.EXPECT().Limiter(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		req, err := http.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", "123")

		rr := httptest.NewRecorder()

		m := NewLimiter(suite.RateLimitToken, suite.RateLimitIp, suite.TokensConfig)

		handler := m.RateLimiter(testHandler)

		handler.ServeHTTP(rr, req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusOK, rr.Code)
		assert.Equal(suite.T(), "", rr.Body.String())
	})

	suite.Run("should return an error when the limiter options are empty", func() {
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		req, err := http.NewRequest("GET", "/", nil)

		rr := httptest.NewRecorder()
		req.Header.Set("API_KEY", "123")

		m := NewLimiter(suite.RateLimitToken, suite.RateLimitIp, []configs.TokenConfigLimit{})

		handler := m.RateLimiter(testHandler)

		handler.ServeHTTP(rr, req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusUnauthorized, rr.Code)
		assert.Equal(suite.T(), "token not found", rr.Body.String())
	})

	suite.Run("should return an error when the limiter by token returns an error", func() {
		suite.RateLimitToken.EXPECT().Limiter(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, assert.AnError)
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		req, err := http.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", "123")

		rr := httptest.NewRecorder()

		m := NewLimiter(suite.RateLimitToken, suite.RateLimitIp, suite.TokensConfig)

		handler := m.RateLimiter(testHandler)

		handler.ServeHTTP(rr, req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusInternalServerError, rr.Code)
		assert.Equal(suite.T(), "error when executing the RateLimiter", rr.Body.String())
	})

	suite.Run("should return an error when the limiter by ip returns an error", func() {
		suite.RateLimitIp.EXPECT().Limiter(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, assert.AnError)
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		req, err := http.NewRequest("GET", "/", nil)

		rr := httptest.NewRecorder()

		m := NewLimiter(suite.RateLimitToken, suite.RateLimitIp, suite.TokensConfig)

		handler := m.RateLimiter(testHandler)

		handler.ServeHTTP(rr, req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusInternalServerError, rr.Code)
		assert.Equal(suite.T(), "error when executing the RateLimiter", rr.Body.String())
	})

	suite.Run("should return an error when the limiter by token returns true", func() {
		suite.RateLimitToken.EXPECT().Limiter(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		req, err := http.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", "123")

		rr := httptest.NewRecorder()

		m := NewLimiter(suite.RateLimitToken, suite.RateLimitIp, suite.TokensConfig)

		handler := m.RateLimiter(testHandler)

		handler.ServeHTTP(rr, req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusTooManyRequests, rr.Code)
		assert.Equal(suite.T(), "you have reached the maximum number of requests or actions allowed within a certain time frame", rr.Body.String())
	})

	suite.Run("should return an error when the limiter by ip returns true", func() {
		suite.RateLimitToken.EXPECT().Limiter(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		req, err := http.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", "123")
		req.Header.Set("X-Forwarded-For", "127.0.0.1")

		rr := httptest.NewRecorder()

		m := NewLimiter(suite.RateLimitToken, suite.RateLimitIp, suite.TokensConfig)

		handler := m.RateLimiter(testHandler)

		handler.ServeHTTP(rr, req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusTooManyRequests, rr.Code)
		assert.Equal(suite.T(), "you have reached the maximum number of requests or actions allowed within a certain time frame", rr.Body.String())
	})

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterTestSuite))
}
