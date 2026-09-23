package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type fakeDeploymentsStore struct {
	all       []model.Deployment
	byOwner   map[string][]model.Deployment
	listCalls int
	ownerID   string
	err       error
}

func (f *fakeDeploymentsStore) List(ctx context.Context) ([]model.Deployment, error) {
	f.listCalls++
	if f.err != nil {
		return nil, f.err
	}
	return f.all, nil
}

func (f *fakeDeploymentsStore) ListByOwner(ctx context.Context, ownerUserID string) ([]model.Deployment, error) {
	f.ownerID = ownerUserID
	if f.err != nil {
		return nil, f.err
	}
	return f.byOwner[ownerUserID], nil
}

func (f *fakeDeploymentsStore) Create(ctx context.Context, serviceID, environmentID string, targetWorkerID *string, gitSHA string) (model.Deployment, error) {
	return model.Deployment{}, errors.New("not implemented")
}
func (f *fakeDeploymentsStore) Get(ctx context.Context, id string) (model.Deployment, error) {
	return model.Deployment{}, errors.New("not implemented")
}
func (f *fakeDeploymentsStore) ListByService(ctx context.Context, serviceID string) ([]model.Deployment, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeDeploymentsStore) ClaimNext(ctx context.Context, workerID string) (model.DeploymentJob, error) {
	return model.DeploymentJob{}, errors.New("not implemented")
}
func (f *fakeDeploymentsStore) Update(ctx context.Context, id string, req model.UpdateDeploymentRequest) (model.Deployment, error) {
	return model.Deployment{}, errors.New("not implemented")
}

func newDeploymentsForListTest(store deploymentsStore) *Deployments {
	return &Deployments{store: store}
}

func TestDeploymentsList_AdminUsesListAll(t *testing.T) {
	fake := &fakeDeploymentsStore{
		all: []model.Deployment{{ID: "d1"}, {ID: "d2"}},
	}
	svc := newDeploymentsForListTest(fake)

	out, err := svc.List(context.Background(), "anyone", "admin")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("want 2 deployments, got %d", len(out))
	}
	if fake.listCalls != 1 {
		t.Fatalf("want List called once, got %d", fake.listCalls)
	}
	if fake.ownerID != "" {
		t.Fatalf("ListByOwner should not run for admin, got ownerID=%q", fake.ownerID)
	}
}

func TestDeploymentsList_NonAdminUsesListByOwner(t *testing.T) {
	fake := &fakeDeploymentsStore{
		byOwner: map[string][]model.Deployment{
			"user-1": {{ID: "mine"}},
			"other":  {{ID: "theirs"}},
		},
	}
	svc := newDeploymentsForListTest(fake)

	out, err := svc.List(context.Background(), "user-1", "developer")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(out) != 1 || out[0].ID != "mine" {
		t.Fatalf("unexpected out: %+v", out)
	}
	if fake.ownerID != "user-1" {
		t.Fatalf("want ListByOwner(user-1), got %q", fake.ownerID)
	}
	if fake.listCalls != 0 {
		t.Fatalf("admin List must not run, got listCalls=%d", fake.listCalls)
	}
}

func TestDeploymentsList_MissingActorForbidden(t *testing.T) {
	fake := &fakeDeploymentsStore{}
	svc := newDeploymentsForListTest(fake)

	_, err := svc.List(context.Background(), "  ", "developer")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("want ErrForbidden, got %v", err)
	}
	if fake.listCalls != 0 || fake.ownerID != "" {
		t.Fatalf("stone must not be called when actor missing")
	}
}

func TestDeploymentsList_AdminRoleTrimmed(t *testing.T) {
	fake := &fakeDeploymentsStore{all: []model.Deployment{{ID: "d1"}}}
	svc := newDeploymentsForListTest(fake)

	_, err := svc.List(context.Background(), "", "  admin  ")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if fake.listCalls != 1 {
		t.Fatalf("trimmed admin should use list all")
	}
}
