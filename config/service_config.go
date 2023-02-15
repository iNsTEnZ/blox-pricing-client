package config

import (
	"blox-client-service/util"
	"time"
)

type ServiceConfig struct {
	Address       string
	FetchInterval time.Duration
	Timeout       time.Duration
	MaxCacheSize  int
	Currency      string
	Tokens        []string
}

func GetServiceConfig() ServiceConfig {
	return ServiceConfig{
		Address:       util.GetEnv("PRICING_ADDRESS", "localhost:8080"),
		MaxCacheSize:  util.IntEnv("MAX_CACHE_SIZE", 2000),
		FetchInterval: time.Second * time.Duration(util.IntEnv("FETCH_INTERVAL_SECONDS", 60)),
		Timeout:       time.Second * time.Duration(util.IntEnv("REQUEST_TIMEOUT_SECONDS", 10)),
		Currency:      util.GetEnv("CURRENCY", "USD"),
		Tokens:        util.GetListEnv("TOKENS", "BTC,ETH,BNB"),
	}
}
