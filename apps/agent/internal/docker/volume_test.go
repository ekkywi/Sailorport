package docker

import "testing"

func TestVolumeName(t *testing.T) {
	got := VolumeName("my-pg", "dev", "data")
	want := "sailorport-my-pg-dev-data"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if VolumeName("x", "", "data") != "sailorport-x-dev-data" {
		t.Fatalf("empty env should default to dev")
	}
}
