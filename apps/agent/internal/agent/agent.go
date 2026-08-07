package agent

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ekkywi/sailorport/apps/agent/internal/client"
	"github.com/ekkywi/sailorport/apps/agent/internal/config"
)

type Agent struct {
	cfg config.Config
	client *client.APIClient
}

func New(cfg config.Config, c *client.APIClient) *Agent {
	return &Agent{cfg: cfg, client: c}
}

func (a *Agent) Run(ctx context.Context) error {
	hostname, _ := os.Hostname()
	labels := map[string]any{
		"role": "agent",
	}

	w, err := a.client.Register(a.cfg.WorkerName, hostname, labels)
	if err != nil {
		return err
	}
	log.Printf("Registered worker id=%s name=%s", w.ID, w.Name)

	ticker := time.NewTicker(a.cfg.HeartbeatInterval)
	defer ticker.Stop()

	if _, err := a.client.Heartbeat(w.ID, "online"); err != nil {
		log.Printf("Heartbeat error: %v", err)
	} else {
		log.Printf("Heartbeat ok status=online")
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down")
			return nil
		case <-ticker.C:
			if _, err := a.client.Heartbeat(w.ID, "online"); err != nil {
				log.Printf("Heartbeat error: %v", err)
				continue
			}
			log.Printf("Heartbeat ok status=online")
		}
	}
}