package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ekkywi/sailorport/apps/agent/internal/client"
	"github.com/ekkywi/sailorport/apps/agent/internal/config"
	"github.com/ekkywi/sailorport/apps/agent/internal/docker"
	"github.com/ekkywi/sailorport/apps/agent/internal/git"
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

	source := strings.TrimSpace(job.SourceType)
	if source == "" {
		source = "scaffold"
	}

	if source == "catalog_app" {
		image := strings.TrimSpace(job.Image)
		if image == "" {
			err := fmt.Errorf("catalog_app job has empty image")
			_ = a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{
				Status: "failed", ErrorMessage: err.Error(),
			})
			return err
		}

		log.Printf("catalog_app pull image=%s container_port=%d", image, job.ContainerPort)
		if err := docker.Pull(image); err != nil {
			_ = a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{
				Status: "failed", ErrorMessage: err.Error(),
			})
			return err
		}

		containerPort := job.ContainerPort
		if containerPort <= 0 {
			containerPort = 8080
		}

		env := []string{"POSTGRES_PASSWORD=changeme"}

		cid, err := docker.Run(containerName, image, port, containerPort, env)
		if err != nil {
			_ = a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{
				Status: "failed", ErrorMessage: err.Error(),
			})
			return err
		}

		return a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{
			Status:      "running",
			ImageTag:    image,
			GitSHA:      "",
			ContainerID: cid,
			Port:        &port,
		})
	}

	workDir, gitSHA, err := a.resolveWorkDir(job)
	if err != nil {
		_ = a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{
			Status: "failed", ErrorMessage: err.Error(),
		})
		return err
	}
	log.Printf("work dir=%s source=%s git_sha=%s", workDir, job.SourceType, gitSHA)
	if err := docker.Build(workDir, imageTag, job.DockerfilePath); err != nil {
		_ = a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{
			Status: "failed", ErrorMessage: err.Error(),
		})
		return err
	}

	cid, err := docker.Run(containerName, imageTag, port, 8080, nil)
	if err != nil {
		_ = a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{
			Status: "failed", ErrorMessage: err.Error(),
		})
		return err
	}

	return a.client.UpdateDeployment(job.ID, client.UpdateDeploymentRequest{
		Status:      "running",
		ImageTag:    imageTag,
		GitSHA:      gitSHA,
		ContainerID: cid,
		Port:        &port,
	})
}

func (a *Agent) resolveWorkDir(job *client.DeploymentJob) (workDir string, gitSHA string, err error) {
	source := strings.TrimSpace(job.SourceType)
	if source == "" {
		source = "scaffold"
	}
	switch source {
	case "scaffold":
		path := strings.TrimSpace(job.WorkspacePath)
		if path == "" {
			return "", "", fmt.Errorf("scaffold service has empty workspace_path")
		}
		return path, "", nil
	case "git":
		if strings.TrimSpace(job.RepoURL) == "" {
			return "", "", fmt.Errorf("git service has empty repo_url")
		}
		dir := filepath.Join(a.cfg.WorkspaceDir, job.ServiceName)
		branch := strings.TrimSpace(job.Branch)
		if branch == "" {
			branch = "main"
		}
		wantSHA := strings.TrimSpace(job.GitSHA)
		log.Printf("git sync repo=%s branch=%s dir=%s sha=%q", job.RepoURL, branch, dir, wantSHA)
		if err := git.Sync(job.RepoURL, branch, dir, wantSHA); err != nil {
			return "", "", err
		}
		gotSHA, err := git.HeadSHA(dir)
		if err != nil {
			return "", "", err
		}
		if wantSHA != "" && gotSHA != wantSHA {
			return "", "", fmt.Errorf("requested sha %s but HEAD is %s", wantSHA, gotSHA)
		}
		log.Printf("git sha=%s", gotSHA)
		return dir, gotSHA, nil
	default:
		return "", "", fmt.Errorf("unsupported source_type %q", source)
	}
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
