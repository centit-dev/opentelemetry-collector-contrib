package client

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"go.opentelemetry.io/collector/pdata/pcommon"
	conventions "go.opentelemetry.io/collector/semconv/v1.8.0"
	"go.uber.org/zap"
)

func TestCreateClient(t *testing.T) {
	config := &PostgresConfig{
		Host:     "host.docker.internal",
		Port:     25432,
		User:     "postgres",
		Pass:     "password",
		Database: "postgres",
	}
	logger, _ := zap.NewDevelopment()
	client, err := CreateClient(config, logger)
	if err != nil {
		t.Errorf("Error when creating client: %s\n", err)
	}
	definitions, _ := client.FindAllDefinitions(context.Background())
	for _, definition := range definitions {
		t.Logf("Definition: %s\n", definition.LongName)
		t.Logf("Definition: %s\n", definition.RelatedMiddlewareConditions)
	}

	service := &ExceptionCategoryService{logger, client, nil, nil}
	service.buildCache(context.Background())
	resourceAttributes := pcommon.NewMap()
	resourceAttributes.PutStr(conventions.AttributeServiceNamespace, "simple-app-local")
	resourceAttributes.PutStr(conventions.AttributeTelemetrySDKVersion, "1.29.0")
	spanAttributes := pcommon.NewMap()
	categories := service.getCategoriesByAttributes(&resourceAttributes, &spanAttributes, "java.lang.RuntimeException", "")
	if len(categories) == 0 {
		t.Errorf("Expected categories to be not empty")
	}
	for _, category := range categories {
		t.Logf("Category: %s\n", category)
	}
}
