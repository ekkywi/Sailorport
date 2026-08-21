package docker

import (
	"fmt"
	"os/exec"
	"strings"
)

func Build(workspace, imageTag, dockerfilePath string) error {
	dockerfilePath = strings.TrimSpace(dockerfilePath)
	if dockerfilePath == "" {
		dockerfilePath = "Dockerfile"
	}

	cmd := exec.Command("docker", "build", "-t", imageTag, "-f", dockerfilePath, ".")
	cmd.Dir = workspace
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build: %w\n%s", err, out)
	}
	return nil
}

func Run(containerName, imageTag string, hostPort int) (containerID string, err error) {
	if err := Remove(containerName); err != nil {
		return "", err
	}

	out, err := runDocker(
		"run",
		"-d",
		"--name", containerName,
		"-p", fmt.Sprintf("%d:8080", hostPort),
		imageTag,
	)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
