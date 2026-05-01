package main

import (
	"flag"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Listen    string
	DelayRate float64
	MaxDelay  time.Duration
	DropRate  float64
	Logger    *slog.Logger
}

func NewConfig() Config {
	handler := slog.NewJSONHandler(os.Stdout, nil)

	logger := slog.New(handler)
	cfg := Config{}
	cfg.Logger = logger

	flag.StringVar(&cfg.Listen, "listen", ":8080", "address to listen on")
	flag.Float64Var(&cfg.DelayRate, "delay-rate", 0.2, "probability of delaying a request")
	flag.DurationVar(&cfg.MaxDelay, "max-delay", time.Second, "maximum random delay")
	flag.Float64Var(&cfg.DropRate, "drop-rate", 0.02, "probability of dropping a request")
	flag.Parse()

	return cfg
}
