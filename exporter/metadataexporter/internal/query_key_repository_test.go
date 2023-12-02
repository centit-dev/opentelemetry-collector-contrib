package internal

import (
	"context"
	"testing"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent"
)

// disable for integration test
func DisabledTestQueryKeyRepositoryImpl_CreateAllAndUpdateAll(t *testing.T) {
	client := createTestClient()
	repo := CreateQueryKeyRepository(client)

	// create all new keys and values
	queryKeys := make([]*ent.QueryKey, 0, 10)
	for i := 0; i < 10; i++ {
		queryKey := createQueryKey()
		for j := 0; j < 3; j++ {
			queryValue := createQueryValue()
			queryKey.Edges.Values = append(queryKey.Edges.Values, queryValue)
		}
		queryKeys = append(queryKeys, queryKey)
	}
	created, err := repo.CreateAll(context.Background(), queryKeys)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(created) != 10 {
		t.Fatalf("expected 10 created keys, got %d", len(created))
	}
	for _, queryKey := range created {
		if queryKey.ID == 0 {
			t.Fatalf("expected created key id to be non-zero")
		}
		if len(queryKey.Edges.Values) != 3 {
			t.Fatalf("expected 1 created value, got %d", len(queryKey.Edges.Values))
		}
	}

	// update all keys with 3 more values
	for _, queryKey := range created {
		for j := 0; j < 3; j++ {
			queryValue := createQueryValue()
			queryKey.Edges.Values = append(queryKey.Edges.Values, queryValue)
		}
	}
	updated, err := repo.UpdateAll(context.Background(), created)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(updated) != 10 {
		t.Fatalf("expected 10 updated keys, got %d", len(updated))
	}
	for _, queryKey := range updated {
		if queryKey.ID == 0 {
			t.Fatalf("expected updated key id to be non-zero")
		}
		if len(queryKey.Edges.Values) != 6 {
			t.Fatalf("expected 6 updated values, got %d", len(queryKey.Edges.Values))
		}
	}
}

// disable for integration test
func DisabledTestQueryKeyRepositoryImpl_FindAll(t *testing.T) {
	client := createTestClient()
	repo := CreateQueryKeyRepository(client)

	queryKeys, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(queryKeys) == 0 {
		t.Fatalf("expected some query keys, got %d", len(queryKeys))
	}
	for _, queryKey := range queryKeys {
		if queryKey.ID == 0 {
			t.Fatalf("expected query key id to be non-zero")
		}
		if len(queryKey.Edges.Values) == 0 {
			t.Fatalf("expected query key to have some values")
		}
	}
}
