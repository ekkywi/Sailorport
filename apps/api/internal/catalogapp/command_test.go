package catalogapp

import "testing"

func TestResolveCommand_Empty(t *testing.T) {
	got, err := ResolveCommand(nil, nil)
	if err != nil || got != nil {
		t.Fatalf("got %#v err %v", got, err)
	}
}

func TestResolveCommand_OK(t *testing.T) {
	got, err := ResolveCommand(
		[]string{"redis-server", "--requirepass", "${REDIS_PASSWORD}"},
		map[string]string{"REDIS_PASSWORD": "s3cret"},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"redis-server", "--requirepass", "s3cret"}
	if len(got) != 3 || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Fatalf("got %#v want %#v", got, want)
	}
}

func TestResolveCommand_MissingEnv(t *testing.T) {
	_, err := ResolveCommand([]string{"${REDIS_PASSWORD}"}, map[string]string{})
	if err == nil {
		t.Fatalf("expected error")
	}
}
