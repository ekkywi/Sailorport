package config

import (
	"os"
	"time"
)

type Config struct {
	APIURL            string
	WorkerName        string
	HeartbeatInterval time.Duration
}

func Load() Config {
	interval := 15 * time.Second
	if v := os.Getenv("SAILORPORT_HEARTBEAT_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			interval = d
		}
	}

	name := os.Getenv("SAILORPORT_WORKER_NAME")
	if name == "" {
		name, _ = os.Hostname()
		if name == "" {
			name = "sailorport-agent"
		}
	}

	apiURL := os.Getenv("SAILORPORT_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	return Config{
		APIURL:            apiURL,
		WorkerName:        name,
		HeartbeatInterval: interval,
	}
}
