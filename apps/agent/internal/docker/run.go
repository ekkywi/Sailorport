package docker

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ekkywi/sailorport/apps/agent/internal/safepath"
)

type VolumeMount struct {
	Name          string
	ContainerPath string
}

func Build(workspace, imageTag, dockerfilePath string) error {
	rel, err := safepath.Dockerfile(workspace, dockerfilePath)
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}

	cmd := exec.Command("docker", "build", "-t", imageTag, "-f", rel, ".")
	cmd.Dir = workspace
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build: %w\n%s", err, out)
	}
	return nil
}

func Pull(image string) error {
	image = strings.TrimSpace(image)
	if image == "" {
		return fmt.Errorf("pull: image is empty")
	}
	out, err := runDocker("pull", image)
	if err != nil {
		return fmt.Errorf("pull: %w\n%s", err, out)
	}
	return nil
}

func Run(containerName, imageTag string, hostPort, containerPort int, env []string, cmdArgs []string, volumes []VolumeMount) (containerID string, err error) {
	if err := Remove(containerName); err != nil {
		return "", err
	}

	if containerPort <= 0 {
		containerPort = 8080
	}

	args := []string{
		"run",
		"-d",
		"--name", containerName,
		"-p", fmt.Sprintf("%d:%d", hostPort, containerPort),
	}
	for _, e := range env {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		args = append(args, "-e", e)
	}
	for _, v := range volumes {
		name := strings.TrimSpace(v.Name)
		path := strings.TrimSpace(v.ContainerPath)
		if name == "" || path == "" {
			continue
		}
		if !strings.HasPrefix(path, "/") || strings.Contains(path, "..") {
			return "", fmt.Errorf("invalid volume path: %q", path)
		}
		args = append(args, "-v", name+":"+path)
	}
	args = append(args, imageTag)

	for _, a := range cmdArgs {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		args = append(args, a)
	}
	out, err := runDocker(args...)
	if err != nil {

		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
