package docker

import (
	"fmt"
	"os/exec"
	"strings"
)

func Build(workspace, imageTag string) error {
	cmd := exec.Command("docker", "build", "-t", imageTag, ".")
	cmd.Dir = workspace
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build: %w\n%s", err, out)
	}
	return nil
}

func Run(containerName, imageTag string, hostPort int) (containerID string, err error) {
	_ = exec.Command("docker", "rm", "-f", containerName).Run()

	cmd := exec.Command(
		"docker", "run", "-d", "--name", containerName, "-p", fmt.Sprintf("%d:8080", hostPort), imageTag,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker run: %w\n%s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}
