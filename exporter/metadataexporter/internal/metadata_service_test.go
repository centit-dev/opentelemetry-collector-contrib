package internal

import (
	"context"
	"testing"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent"
	"go.uber.org/zap"
)

// disable for integration test
func DisabledTestMetadataServiceImpl_upsertAll(t *testing.T) {
	client := createTestClient()
	keyRepo := CreateQueryKeyRepository(client)
	valueRepo := CreateQueryValueRepository(client)
	logger := zap.NewExample()
	service := CreateMetadataService(&CacheConfig{}, &BatchConfig{}, 90, logger, keyRepo, valueRepo)

	// create all new keys and values
	queryKeys := make([]*ent.QueryKey, 0, 15)
	for i := 0; i < 10; i++ {
		queryKey := createQueryKey()
		for j := 0; j < 3; j++ {
			queryValue := createQueryValue()
			queryKey.Edges.Values = append(queryKey.Edges.Values, queryValue)
		}
		queryKeys = append(queryKeys, queryKey)
	}
	created, updated, err := service.upsertAll(context.Background(), queryKeys)
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
	if updated != nil {
		t.Fatalf("expected no updated keys, got %d", len(updated))
	}

	// update all keys with 3 more values
	for _, queryKey := range queryKeys {
		for j := 0; j < 3; j++ {
			queryValue := createQueryValue()
			queryKey.Edges.Values = append(queryKey.Edges.Values, queryValue)
		}
	}
	// add 5 more keys
	for i := 0; i < 5; i++ {
		queryKey := createQueryKey()
		for j := 0; j < 3; j++ {
			queryValue := createQueryValue()
			queryKey.Edges.Values = append(queryKey.Edges.Values, queryValue)
		}
		queryKeys = append(queryKeys, queryKey)
	}
	created, updated, err = service.upsertAll(context.Background(), queryKeys)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(created) != 5 {
		t.Fatalf("expected 5 created keys, got %d", len(created))
	}
	for _, queryKey := range created {
		if queryKey.ID == 0 {
			t.Fatalf("expected created key id to be non-zero")
		}
		if len(queryKey.Edges.Values) != 3 {
			t.Fatalf("expected 3 created value, got %d", len(queryKey.Edges.Values))
		}
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
