package safepath

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestJoinUnder_OK(t *testing.T) {
	root := t.TempDir()
	got, err := JoinUnder(root, "payments-api")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "payments-api")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestJoinUnder_RejectsTraversal(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"..", ".", "a/b", `a\b`, "../x"} {
		if _, err := JoinUnder(root, name); err == nil {
			t.Fatalf("expected error for name %q", name)
		}
	}
}

func TestUnderRoot_OK(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "svc")
	got, err := UnderRoot(root, child)
	if err != nil {
		t.Fatal(err)
	}
	if got != child {
		t.Fatalf("got %q want %q", got, child)
	}
}

func TestUnderRoot_RejectsOutside(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(filepath.Dir(root), "other")
	if _, err := UnderRoot(root, outside); err == nil {
		t.Fatal("expected error for path outside root")
	}
}

func TestDockerfile_OK(t *testing.T) {
	root := t.TempDir()
	got, err := Dockerfile(root, "deploy/Dockerfile")
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.Clean("deploy/Dockerfile") {
		t.Fatalf("got %q", got)
	}
}

func TestDockerfile_RejectsEscape(t *testing.T) {
	root := t.TempDir()
	for _, p := range []string{"/etc/passwd", "../Dockerfile", ".."} {
		if _, err := Dockerfile(root, p); err == nil {
			t.Fatalf("expected error for %q", p)
		}
	}
}

func TestDockerfile_Default(t *testing.T) {
	root := t.TempDir()
	got, err := Dockerfile(root, "  ")
	if err != nil {
		t.Fatal(err)
	}
	if got != "Dockerfile" {
		t.Fatalf("got %q", got)
	}
	if strings.Contains(got, "..") {
		t.Fatal("unexpected ..")
	}
}
