package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ekkywi/sailorport/apps/agent/internal/client"
	"github.com/ekkywi/sailorport/apps/agent/internal/config"
	"github.com/ekkywi/sailorport/apps/agent/internal/docker"
)

type Agent struct {
	cfg    config.Config
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

	workerID := w.ID

	ticker := time.NewTicker(a.cfg.HeartbeatInterval)
	defer ticker.Stop()

	pollTicker := time.NewTicker(a.cfg.PollInterval)
	defer pollTicker.Stop()

	if _, err := a.client.Heartbeat(workerID, "online"); err != nil {
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
			if _, err := a.client.Heartbeat(workerID, "online"); err != nil {
				log.Printf("Heartbeat error: %v", err)
				continue
			}
			log.Printf("Heartbeat ok status=online")
		case <-pollTicker.C:
			if err := a.handleJob(ctx, workerID); err != nil {
				log.Printf("job error: %v", err)
			}
		}
	}
}

func (a *Agent) handleJob(ctx context.Context, workerID string) error {
	job, err := a.client.ClaimNext(workerID)
	if err != nil {
		return err
	}
	if job == nil {
		return nil
	}

	log.Printf("claimed job id=%s service=%s path=%s", job.ID, job.ServiceName, job.WorkspacePath)

	_ = a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{Status: "building"})

	imageTag := fmt.Sprintf("sailorport/%s:%s", job.ServiceName, job.ID[:8])
	containerName := "sailorport-" + job.ServiceName
	port := a.cfg.PortBase

	if err := docker.Build(job.WorkspacePath, imageTag); err != nil {
		_ = a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{
			Status: "failed", ErrorMessage: err.Error(),
		})
		return err
	}

	cid, err := docker.Run(containerName, imageTag, port)
	if err != nil {
		_ = a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{
			Status: "failed", ErrorMessage: err.Error(),
		})
		return err
	}

	return a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{
		Status:      "running",
		ImageTag:    imageTag,
		ContainerID: cid,
		Port:        &port,
	})
}
