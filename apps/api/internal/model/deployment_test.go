package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDeploymentJob_GitTokenJSON(t *testing.T) {
	job := DeploymentJob{
		ServiceName: "priv",
		SourceType:  "git",
		GitToken:    "ghp_test",
	}
	b, err := json.Marshal(job)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"git_token":"ghp_test"`) {
		t.Fatalf("missing git_token in JSON: %s", b)
	}
	empty := DeploymentJob{ServiceName: "pub", SourceType: "git"}
	b2, _ := json.Marshal(empty)
	if strings.Contains(string(b2), "git_token") {
		t.Fatalf("omitempty failed: %s", b2)
	}
}
