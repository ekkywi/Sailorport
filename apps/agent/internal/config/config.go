package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	APIURL            string
	WorkerName        string
	HeartbeatInterval time.Duration
	PollInterval      time.Duration
	PortBase          int
}

func Load() Config {
	interval := 15 * time.Second
	if v := os.Getenv("SAILORPORT_HEARTBEAT_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			interval = d
		}
	}

	poll := 5 * time.Second
	if v := os.Getenv("SAILORPORT_POLL_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			poll = d
		}
	}

	portBase := 18080
	if v := os.Getenv("SAILORPORT_DEPLOY_PORT_BASE"); v != "" {
		fmt.Sscanf(v, "%d", &portBase)
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
		PollInterval:      poll,
		PortBase:          portBase,
	}
}
