package configs

import (
	"encoding/json"

	"github.com/spf13/viper"
)

type TokenConfigLimit struct {
	Token           string `json:"token"`
	MaxRequests     int64  `json:"max_requests"`
	BlockTimeSecond int64  `json:"block_time_seconds"`
}

type IPConfigLimit struct {
	MaxRequests     int64 `json:"max_requests"`
	BlockTimeSecond int64 `json:"block_time_seconds"`
}

var (
	envVars *Environments
)

type Environments struct {
	WebServerPort         string `mapstructure:"WEB_SERVER_PORT"`
	TokensConfigLimitJson string `mapstructure:"TOKENS_CONFIG_LIMIT"`
	TokensConfigLimit     []TokenConfigLimit
	IPConfigLimitJson     string `mapstructure:"IP_CONFIG_LIMIT"`
	IPConfigLimit         IPConfigLimit
	RedisHost             string `mapstructure:"REDIS_HOST"`
	RedisPort             string `mapstructure:"REDIS_PORT"`
	RedisPassword         string `mapstructure:"REDIS_PASSWORD"`
	RedisDB               int    `mapstructure:"REDIS_DB"`
}

func LoadConfig(path string) (*Environments, error) {

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&envVars)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(envVars.TokensConfigLimitJson), &envVars.TokensConfigLimit)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(envVars.IPConfigLimitJson), &envVars.IPConfigLimit)
	if err != nil {
		return nil, err
	}

	return envVars, err
}

func GetEnvVars() *Environments {
	return envVars
}
