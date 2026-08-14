package docker

import "testing"

func TestParseHostPorts(t *testing.T) {
	got := parseHostPorts("18080 18080 18081\n")
	if len(got) != 2 || got[0] != 18080 || got[1] != 18081 {
		t.Fatalf("got %v", got)
	}
}

func TestFirstFreePort_skipsUsedAndUnavailable(t *testing.T) {
	used := map[int]struct{}{18080: {}}
	blocked := map[int]struct{}{18081: {}}
	p, err := firstFreePort(18080, 4, used, func(port int) bool {
		_, no := blocked[port]
		return !no
	})
	if err != nil {
		t.Fatal(err)
	}
	if p != 18082 {
		t.Fatalf("got %d, want 18082", p)
	}
}

func TestFirstFreePort_exhausted(t *testing.T) {
	used := map[int]struct{}{10: {}, 11: {}}
	_, err := firstFreePort(10, 2, used, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
