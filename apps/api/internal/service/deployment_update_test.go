package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type fakeDeployUpdateStore struct {
	fakeDeploymentsStore
	existing model.Deployment
	updated  model.UpdateDeploymentRequest
	getErr   error
	updErr   error
}

func (f *fakeDeployUpdateStore) Get(ctx context.Context, id string) (model.Deployment, error) {
	if f.getErr != nil {
		return model.Deployment{}, f.getErr
	}
	d := f.existing
	d.ID = id
	return d, nil
}

func (f *fakeDeployUpdateStore) Update(ctx context.Context, id string, req model.UpdateDeploymentRequest) (model.Deployment, error) {
	if f.updErr != nil {
		return model.Deployment{}, f.updErr
	}
	f.updated = req
	out := f.existing
	out.ID = id
	if req.Status != "" {
		out.Status = req.Status
	}
	return out, nil
}

func TestDeploymentsUpdate_RequiresMatchingWorker(t *testing.T) {
	wid := "worker-a"
	fake := &fakeDeployUpdateStore{
		existing: model.Deployment{
			ID:       "dep-1",
			WorkerID: &wid,
			Status:   "claimed",
		},
	}
	svc := newDeploymentsForListTest(fake)

	out, err := svc.Update(context.Background(), "dep-1", model.UpdateDeploymentRequest{
		Status:   "building",
		WorkerID: "worker-a",
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if out.Status != "building" {
		t.Fatalf("status=%q", out.Status)
	}
	if fake.updated.WorkerID != "worker-a" {
		t.Fatalf("store saw worker_id=%q", fake.updated.WorkerID)
	}
}

func TestDeploymentsUpdate_RejectsMissingWorkerID(t *testing.T) {
	wid := "worker-a"
	fake := &fakeDeployUpdateStore{
		existing: model.Deployment{WorkerID: &wid},
	}
	svc := newDeploymentsForListTest(fake)

	_, err := svc.Update(context.Background(), "dep-1", model.UpdateDeploymentRequest{Status: "building"})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("want ErrInvalid, got %v", err)
	}
}

func TestDeploymentsUpdate_RejectsMismatch(t *testing.T) {
	wid := "worker-a"
	fake := &fakeDeployUpdateStore{
		existing: model.Deployment{WorkerID: &wid},
	}
	svc := newDeploymentsForListTest(fake)

	_, err := svc.Update(context.Background(), "dep-1", model.UpdateDeploymentRequest{
		Status:   "building",
		WorkerID: "worker-b",
	})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("want ErrForbidden, got %v", err)
	}
}

func TestDeploymentsUpdate_RejectsUnclaimed(t *testing.T) {
	fake := &fakeDeployUpdateStore{
		existing: model.Deployment{Status: "pending"},
	}
	svc := newDeploymentsForListTest(fake)

	_, err := svc.Update(context.Background(), "dep-1", model.UpdateDeploymentRequest{
		Status:   "building",
		WorkerID: "worker-a",
	})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("want ErrConflict, got %v", err)
	}
}
