package internal

import (
	"math/rand"
	"time"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent"
	conventions "go.opentelemetry.io/collector/semconv/v1.18.0"
	"go.uber.org/zap"
)

func createQueryKey() *ent.QueryKey {
	source := []string{queryKeySourceResource, queryKeySourceSpan}[rand.Int31n(2)]
	var name string
	if source == queryKeySourceResource {
		count := len(conventions.GetResourceSemanticConventionAttributeNames())
		name = conventions.GetResourceSemanticConventionAttributeNames()[rand.Intn(count)]
	} else {
		count := len(conventions.GetTraceSemanticConventionAttributeNames())
		name = conventions.GetTraceSemanticConventionAttributeNames()[rand.Intn(count)]
	}
	return &ent.QueryKey{
		Name:      name,
		Source:    source,
		Type:      []string{queryValueTypeString, queryValueTypeNumber}[rand.Int31n(2)],
		ValidDate: time.Now().AddDate(0, 0, 90),
	}
}

func createQueryValue() *ent.QueryValue {
	return &ent.QueryValue{
		Value: []string{"value1", "value2", "value3", "value4", "value5"}[rand.Int31n(5)],
	}
}

func createTestClient() *DatabaseClient {
	logger := zap.NewExample()
	client, _ := CreateClient(&PostgresConfig{
		Host:     "host.docker.internal",
		Port:     25432,
		User:     "postgres",
		Pass:     "password",
		Database: "postgres",
		Debug:    true,
	}, logger)

	return client
}
