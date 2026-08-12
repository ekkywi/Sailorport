package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ekkywi/sailorport/apps/agent/internal/agent"
	"github.com/ekkywi/sailorport/apps/agent/internal/client"
	"github.com/ekkywi/sailorport/apps/agent/internal/config"
)

func main() {
	cfg := config.Load()
	log.Printf("agent starting api=%s name=%s interval=%s",
		cfg.APIURL, cfg.WorkerName, cfg.HeartbeatInterval)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	a := agent.New(cfg, client.New(cfg.APIURL, cfg.AgentToken))
	if err := a.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
