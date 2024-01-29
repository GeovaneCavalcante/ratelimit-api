package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GeovaneCavalcante/rate-limit-api/configs"
	"github.com/GeovaneCavalcante/rate-limit-api/internal/infra/web/handlers"
	"github.com/GeovaneCavalcante/rate-limit-api/internal/infra/web/middlewares"
	"github.com/GeovaneCavalcante/rate-limit-api/internal/infra/web/webserver"
	"github.com/GeovaneCavalcante/rate-limit-api/pkg/logger"
	"github.com/GeovaneCavalcante/rate-limit-api/pkg/ratelimit"
	redisEventStorage "github.com/GeovaneCavalcante/rate-limit-api/pkg/ratelimit/redis"
	"github.com/go-redis/redis/v8"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configs.RedisHost, configs.RedisPort),
		Password: configs.RedisPassword,
		DB:       configs.RedisDB,
	})
	fmt.Printf("%s:%s", configs.RedisHost, configs.RedisPort)

	redisEventStorage := redisEventStorage.NewRedisEventStorage(rdb)

	rlIp, err := ratelimit.New(redisEventStorage, "ip", configs.IPConfigLimit.MaxRequests, time.Duration(configs.IPConfigLimit.BlockTimeSecond)*time.Second)
	if err != nil {
		logger.Error("error when executing the RateLimiter by ip", err)
		return
	}

	rlToken, err := ratelimit.New(redisEventStorage, "token", 0, 0*time.Second)
	if err != nil {
		logger.Error("error when executing the RateLimiter by token", err)
		return
	}

	m := middlewares.NewLimiter(rlToken, rlIp, configs.TokensConfigLimit)

	ws := webserver.New(configs.WebServerPort)
	h := handlers.NewHealthHandler()

	ws.AddHandler("/health", m.RateLimiter(http.HandlerFunc(h.HealthHandler)))
	ws.Start()
}
