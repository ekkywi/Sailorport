package docker

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const defaultPortCount = 32

// AllocateHostPort picks a host port for a workload container.
// If that container already has a published port, reuse it (redeploy).
// Otherwise take the first free port in [base, base+count).
func AllocateHostPort(containerName string, base, count int) (int, error) {
	if base <= 0 {
		base = 18080
	}
	if count <= 0 {
		count = defaultPortCount
	}

	if existing, err := publishedHostPort(containerName); err != nil {
		return 0, err
	} else if existing > 0 {
		return existing, nil
	}

	used, err := usedHostPorts()
	if err != nil {
		return 0, err
	}
	return firstFreePort(base, count, used, portAvailable)
}

func firstFreePort(base, count int, used map[int]struct{}, available func(int) bool) (int, error) {
	for i := 0; i < count; i++ {
		p := base + i
		if _, taken := used[p]; taken {
			continue
		}
		if available != nil && !available(p) {
			continue
		}
		return p, nil
	}
	return 0, fmt.Errorf("no free host port in %d–%d", base, base+count-1)
}

func portAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

func publishedHostPort(containerName string) (int, error) {
	out, err := runDocker(
		"inspect",
		"-f", "{{range $p, $b := .HostConfig.PortBindings}}{{range $b}}{{.HostPort}} {{end}}{{end}}",
		containerName,
	)
	if err != nil {
		if isNotFound(out) || isNoSuchObject(out) {
			return 0, nil
		}
		return 0, err
	}
	return firstPort(string(out)), nil
}

func usedHostPorts() (map[int]struct{}, error) {
	idsOut, err := runDocker("ps", "-aq", "--filter", "name=sailorport-")
	if err != nil {
		return nil, err
	}
	ids := strings.Fields(string(idsOut))
	if len(ids) == 0 {
		return map[int]struct{}{}, nil
	}

	args := []string{
		"inspect",
		"-f", "{{range $p, $b := .HostConfig.PortBindings}}{{range $b}}{{.HostPort}} {{end}}{{end}}",
	}
	args = append(args, ids...)
	out, err := runDocker(args...)
	if err != nil {
		return nil, err
	}

	used := map[int]struct{}{}
	for _, p := range parseHostPorts(string(out)) {
		used[p] = struct{}{}
	}
	return used, nil
}

func parseHostPorts(s string) []int {
	var out []int
	seen := map[int]struct{}{}
	for _, tok := range strings.Fields(s) {
		p, err := strconv.Atoi(strings.TrimSpace(tok))
		if err != nil || p <= 0 {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}

func firstPort(s string) int {
	ports := parseHostPorts(s)
	if len(ports) == 0 {
		return 0
	}
	return ports[0]
}

func isNoSuchObject(out []byte) bool {
	s := strings.ToLower(string(out))
	return strings.Contains(s, "no such object") || strings.Contains(s, "no such container")
}
