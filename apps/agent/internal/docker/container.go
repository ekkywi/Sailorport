package docker

import (
	"fmt"
	"os/exec"
	"strings"
)

func ContainerName(serviceName string) string {
	return "sailorport-" + serviceName
}

func runDocker(args ...string) ([]byte, error) {
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker %s: %w\n%s", strings.Join(args, " "), err, out)
	}
	return out, nil
}

func isNotFound(out []byte) bool {
	s := strings.ToLower(string(out))
	return strings.Contains(s, "no such container")
}

func Stop(containerName string) error {
	out, err := runDocker("stop", containerName)
	if err != nil && isNotFound(out) {
		return nil
	}
	return err
}

func Start(containerName string) error {
	_, err := runDocker("start", containerName)
	return err
}

func Remove(containerName string) error {
	out, err := runDocker("rm", "-f", containerName)
	if err != nil && isNotFound(out) {
		return nil
	}
	return err
}
