package middlewares

import (
	"fmt"
	"net/http"

	"github.com/GeovaneCavalcante/rate-limit-api/configs"
	"github.com/GeovaneCavalcante/rate-limit-api/pkg/logger"
	"github.com/GeovaneCavalcante/rate-limit-api/pkg/ratelimit"
)

type Limiter struct {
	TokenLimiter      ratelimit.RateLimiterInterface
	IPLimiter         ratelimit.RateLimiterInterface
	TokensConfigLimit []configs.TokenConfigLimit
}

func NewLimiter(tokenLimiter, ipLimiter ratelimit.RateLimiterInterface, tk []configs.TokenConfigLimit) Limiter {
	return Limiter{
		TokenLimiter:      tokenLimiter,
		IPLimiter:         ipLimiter,
		TokensConfigLimit: tk,
	}
}

func (l *Limiter) RateLimiter(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("API_KEY")
		ip := getIP(r)

		if token != "" {

			tokenOptions := findOpionsByToken(l.TokensConfigLimit, token)
			if tokenOptions == nil {
				logger.Error(fmt.Sprintf("token %s not found", token), nil)
				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`token not found`))
				return
			}
			tokenLimiter, err := l.TokenLimiter.Limiter(r.Context(), token, tokenOptions)
			if err != nil {
				logger.Error("error when executing the RateLimiter by token", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`error when executing the RateLimiter`))
				return
			}

			if tokenLimiter {
				logger.Warn("TOKENLIMIT - you have reached the maximum number of requests or actions allowed within a certain time frame", nil)
				w.WriteHeader(http.StatusTooManyRequests)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`you have reached the maximum number of requests or actions allowed within a certain time frame`))
				return
			}
			if !tokenLimiter {
				next.ServeHTTP(w, r)
				return
			}
		}

		ipLimiter, err := l.IPLimiter.Limiter(r.Context(), ip, nil)

		if err != nil {
			logger.Error("error when executing the RateLimiter by ip", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`error when executing the RateLimiter`))
			return
		}

		if ipLimiter {
			logger.Warn("IPLIMIT - you have reached the maximum number of requests or actions allowed within a certain time frame", nil)
			w.WriteHeader(http.StatusTooManyRequests)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`you have reached the maximum number of requests or actions allowed within a certain time frame`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func findOpionsByToken(tk []configs.TokenConfigLimit, token string) *ratelimit.Options {
	for _, t := range tk {
		if t.Token == token {
			return &ratelimit.Options{
				NameSpace:      "token",
				MaxInInterval:  t.MaxRequests,
				IntervalSecund: t.BlockTimeSecond,
			}
		}
	}

	return nil
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	return r.RemoteAddr
}
