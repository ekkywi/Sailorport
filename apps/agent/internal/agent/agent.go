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

	w, err := a.client.Register(a.cfg.WorkerName, hostname, a.cfg.Labels)
	if err != nil {
		return err
	}
	log.Printf("Registered worker id=%s name=%s labels=%v", w.ID, w.Name, a.cfg.Labels)

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
			if _, err := a.client.Heartbeat(workerID, "offline"); err != nil {
				log.Printf("offline heartbeat error: %v", err)
			} else {
				log.Printf("Heartbeat ok status=offline")
			}
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
			if err := a.handleRuntime(workerID); err != nil {
				log.Printf("runtime error: %v", err)
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

	log.Printf("claimed job id=%s service=%s env=%s path=%s", job.ID, job.ServiceName, job.EnvironmentSlug, job.WorkspacePath)

	_ = a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{Status: "building"})

	imageTag := fmt.Sprintf("sailorport/%s:%s", job.ServiceName, job.ID[:8])
	containerName := docker.ContainerName(job.ServiceName, job.EnvironmentSlug)
	port, err := docker.AllocateHostPort(containerName, a.cfg.PortBase, a.cfg.PortCount)
	if err != nil {
		_ = a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{
			Status: "failed", ErrorMessage: err.Error(),
		})
		return err
	}
	log.Printf("host port=%d container=%s", port, containerName)

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

func (a *Agent) handleRuntime(workerID string) error {
	job, err := a.client.ClaimRuntimeNext(workerID)
	if err != nil {
		return err
	}
	if job == nil {
		return nil
	}

	containerName := docker.ContainerName(job.ServiceName, job.EnvironmentSlug)
	log.Printf("runtime job id=%s action=%s env %s container=%s", job.ID, job.Action, job.EnvironmentSlug, containerName)

	var runErr error
	switch job.Action {
	case "stop":
		runErr = docker.Stop(containerName)
	case "start":
		runErr = docker.Start(containerName)
	case "remove":
		runErr = docker.Remove(containerName)
	case "logs":
		text, err := docker.Logs(containerName, 200)
		if err != nil {
			runErr = err
			break
		}
		const maxLogBytes = 64 * 1024
		if len(text) > maxLogBytes {
			text = text[len(text)-maxLogBytes:]
		}
		return a.client.UpdateRuntime(job.ID, client.UpdateRuntimeRequest{
			Status: "done",
			Output: text,
		})
	default:
		runErr = fmt.Errorf("unknown action %q", job.Action)
	}

	if runErr != nil {
		_ = a.client.UpdateRuntime(job.ID, client.UpdateRuntimeRequest{
			Status: "failed", ErrorMessage: runErr.Error(),
		})
		return runErr
	}
	return a.client.UpdateRuntime(job.ID, client.UpdateRuntimeRequest{Status: "done"})
}
