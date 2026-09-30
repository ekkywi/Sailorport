package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	APIURL            string
	WorkerName        string
	Labels            map[string]any
	HeartbeatInterval time.Duration
	PollInterval      time.Duration
	PortBase          int
	PortCount         int
	AgentToken        string
	WorkspaceDir      string
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

	portCount := 32
	if v := os.Getenv("SAILORPORT_DEPLOY_PORT_COUNT"); v != "" {
		fmt.Sscanf(v, "%d", &portCount)
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

	agentToken := os.Getenv("SAILORPORT_AGENT_TOKEN")
	if agentToken == "" {
		agentToken = "dev-agent-token"
	}

	workspaceDir := strings.TrimSpace(os.Getenv("SAILORPORT_WORKSPACE"))
	if workspaceDir == "" {
		workspaceDir = defaultWorkspaceDir()
	}

	return Config{
		APIURL:            apiURL,
		WorkerName:        name,
		HeartbeatInterval: interval,
		PollInterval:      poll,
		PortBase:          portBase,
		PortCount:         portCount,
		AgentToken:        agentToken,
		Labels:            parseWorkerLabels(),
		WorkspaceDir:      workspaceDir,
	}
}

// defaultWorkspaceDir prefers the repo data/workspaces folder (same as API
// scaffold) when the agent is started from apps/agent. Falls back to
// ./workspaces for isolated/agent-only layouts.
func defaultWorkspaceDir() string {
	candidates := []string{
		filepath.Join("..", "..", "data", "workspaces"),
		filepath.Join("..", "data", "workspaces"),
		filepath.Join("data", "workspaces"),
		filepath.Join(".", "workspaces"),
	}
	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		// Prefer Sailorport monorepo layout: .../data/workspaces next to templates/.
		repoRoot := filepath.Dir(filepath.Dir(abs))
		if info, err := os.Stat(filepath.Join(repoRoot, "templates")); err == nil && info.IsDir() {
			return abs
		}
		if info, err := os.Stat(abs); err == nil && info.IsDir() {
			return abs
		}
	}
	if abs, err := filepath.Abs(filepath.Join(".", "workspaces")); err == nil {
		return abs
	}
	return filepath.Join(".", "workspaces")
}

func parseWorkerLabels() map[string]any {
	labels := map[string]any{
		"role": "agent",
	}

	if raw := strings.TrimSpace(os.Getenv("SAILORPORT_WORKER_LABELS")); raw != "" {
		var extra map[string]any
		if err := json.Unmarshal([]byte(raw), &extra); err != nil {
			fmt.Fprintf(os.Stderr, "warn: invalid SAILORPORT_WORKER_LABELS JSON: %v\n", err)
		} else {
			for k, v := range extra {
				labels[k] = v
			}
		}
	}

	if tier := strings.TrimSpace(os.Getenv("SAILORPORT_WORKER_TIER")); tier != "" {
		labels["tier"] = tier
	}
	if envs := strings.TrimSpace(os.Getenv("SAILORPORT_WORKER_ENVIRONMENTS")); envs != "" {
		labels["environments"] = envs
	}

	return labels
}
