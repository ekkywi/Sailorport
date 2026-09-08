package model

import "testing"

func TestWorkerAccpetsDeploy(t *testing.T) {
	if !WorkerAcceptsDeploy(WorkerStatusOnline) {
		t.Fatal("online should accept deploy")
	}
	if WorkerAcceptsDeploy(WorkerStatusOffline) {
		t.Fatal("offline should not accept deploy")
	}
	if WorkerAcceptsDeploy(WorkerStatusDraining) {
		t.Fatal("draining should not accpet deploy")
	}
}

func TestIsWorkerStatus(t *testing.T) {
	if !IsWorkerStatus(WorkerStatusDraining) {
		t.Fatal("draining should be valid")
	}
	if IsWorkerStatus("nope") {
		t.Fatal("nope should be invalid")
	}
}