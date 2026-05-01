package main

import (
	"math/rand/v2"
	"net/http"
	"time"
)

func injectDelay(cfg Config) {
	if cfg.DelayRate <= 0 || cfg.MaxDelay <= 0 {
		return
	}

	if rand.Float64() < cfg.DelayRate {
		delay := time.Duration(rand.Int64N(int64(cfg.MaxDelay)))
		cfg.Logger.Info("injecting delay", "duration", delay)
		time.Sleep(delay)
	}
}

func shouldDrop(cfg Config) bool {
	return cfg.DropRate > 0 && rand.Float64() < cfg.DropRate
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
