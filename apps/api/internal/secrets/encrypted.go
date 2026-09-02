package secrets

import (
	"context"
	"fmt"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type EncryptedStore struct {
	env catalogEnvStore
	key []byte
}

func NewEncrypted(env catalogEnvStore, key []byte) (*EncryptedStore, error) {
	if err := validateSecretsKey(key); err != nil {
		return nil, err
	}
	return &EncryptedStore{env: env, key: append([]byte(nil), key...)}, nil
}

func (e *EncryptedStore) ReplaceAll(ctx context.Context, serviceID string, entries []model.CatalogEnv) error {
	toStore := make([]model.CatalogEnv, len(entries))
	for i, entry := range entries {
		toStore[i] = entry
		if !entry.Secret {
			continue
		}
		enc, err := encryptSecretValue(e.key, entry.Value)
		if err != nil {
			return fmt.Errorf("encrypt catalog env %q: %w", entry.Key, err)
		}
		toStore[i].Value = enc
	}
	return e.env.ReplaceAll(ctx, serviceID, toStore)
}

func (e *EncryptedStore) DeleteByServiceID(ctx context.Context, serviceID string) error {
	return e.env.DeleteByServiceID(ctx, serviceID)
}

func (e *EncryptedStore) ResolveForDeploy(ctx context.Context, serviceID string) (map[string]string, error) {
	rows, err := e.env.ListByServiceID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	plain, err := e.decryptRows(rows)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(plain))
	for _, row := range plain {
		out[row.Key] = row.Value
	}
	return out, nil
}

func (e *EncryptedStore) PublicView(ctx context.Context, serviceID string) (model.CatalogEnvPublic, error) {
	rows, err := e.env.ListByServiceID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	plain, err := e.decryptRows(rows)
	if err != nil {
		return nil, err
	}
	return catalogEnvPublic(plain), nil
}

func (e *EncryptedStore) decryptRows(rows []model.CatalogEnv) ([]model.CatalogEnv, error) {
	out := make([]model.CatalogEnv, len(rows))
	for i, row := range rows {
		out[i] = row
		if !isEncryptedValue(row.Value) {
			continue
		}
		plain, err := decryptSecretValue(e.key, row.Value)
		if err != nil {
			return nil, fmt.Errorf("decrypt catalog env %q: %w", row.Key, err)
		}
		out[i].Value = plain
	}
	return out, nil
}

var _ Store = (*EncryptedStore)(nil)
