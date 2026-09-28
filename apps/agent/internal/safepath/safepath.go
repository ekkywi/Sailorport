package safepath

import (
	"fmt"
	"path/filepath"
	"strings"
)

// JoinUnder joins root and a single path segment (e.g. service name) and
// ensures the result stays under root after Clean/Abs.
func JoinUnder(root, name string) (string, error) {
	root = strings.TrimSpace(root)
	name = strings.TrimSpace(name)
	if root == "" {
		return "", fmt.Errorf("workspace root is empty")
	}
	if name == "" {
		return "", fmt.Errorf("name is empty")
	}
	if name == "." || name == ".." {
		return "", fmt.Errorf("invalid name %q", name)
	}
	if strings.ContainsAny(name, `/\`) {
		return "", fmt.Errorf("name must not contain path separators")
	}
	if filepath.IsAbs(name) {
		return "", fmt.Errorf("name must be relative")
	}

	cleanRoot, err := filepath.Abs(filepath.Clean(root))
	if err != nil {
		return "", fmt.Errorf("resolve root: %w", err)
	}
	joined := filepath.Join(cleanRoot, name)
	cleanJoined, err := filepath.Abs(filepath.Clean(joined))
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}
	if err := underRoot(cleanRoot, cleanJoined); err != nil {
		return "", err
	}
	return cleanJoined, nil
}

// UnderRoot resolves path and ensures it is root or a descendant of root.
func UnderRoot(root, path string) (string, error) {
	root = strings.TrimSpace(root)
	path = strings.TrimSpace(path)
	if root == "" {
		return "", fmt.Errorf("workspace root is empty")
	}
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}

	cleanRoot, err := filepath.Abs(filepath.Clean(root))
	if err != nil {
		return "", fmt.Errorf("resolve root: %w", err)
	}
	cleanPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}
	if err := underRoot(cleanRoot, cleanPath); err != nil {
		return "", err
	}
	return cleanPath, nil
}

// Dockerfile returns a path suitable for docker build -f, relative to workDir,
// after ensuring it does not escape the build context.
func Dockerfile(workDir, dockerfilePath string) (string, error) {
	dockerfilePath = strings.TrimSpace(dockerfilePath)
	if dockerfilePath == "" {
		dockerfilePath = "Dockerfile"
	}
	if filepath.IsAbs(dockerfilePath) {
		return "", fmt.Errorf("dockerfile_path must be relative")
	}
	cleaned := filepath.Clean(dockerfilePath)
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("dockerfile_path escapes workspace")
	}

	workDir, err := filepath.Abs(filepath.Clean(strings.TrimSpace(workDir)))
	if err != nil {
		return "", fmt.Errorf("resolve work dir: %w", err)
	}
	abs := filepath.Join(workDir, cleaned)
	abs, err = filepath.Abs(filepath.Clean(abs))
	if err != nil {
		return "", fmt.Errorf("resolve dockerfile: %w", err)
	}
	if err := underRoot(workDir, abs); err != nil {
		return "", fmt.Errorf("dockerfile_path: %w", err)
	}
	return cleaned, nil
}

func underRoot(cleanRoot, cleanPath string) error {
	sep := string(filepath.Separator)
	if cleanPath == cleanRoot {
		return nil
	}
	if !strings.HasPrefix(cleanPath+sep, cleanRoot+sep) {
		return fmt.Errorf("path %q escapes workspace root %q", cleanPath, cleanRoot)
	}
	return nil
}
