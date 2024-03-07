package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent/schema"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent/spanfault"
	"go.uber.org/zap"
)

type SpanFaultRepository interface {
	SaveAll(ctx context.Context, creates []*ent.SpanFault) error
	Shutdown(ctx context.Context) error
}

type SpanFaultRepositoryImpl struct {
	client *ClickHouseClient
	logger *zap.Logger
}

func CreateSpanFaultRepository(client *ClickHouseClient, logger *zap.Logger) SpanFaultRepository {
	return &SpanFaultRepositoryImpl{client, logger}
}

func (repo *SpanFaultRepositoryImpl) SaveAll(ctx context.Context, creates []*ent.SpanFault) error {
	values := repo.buildValues(creates...)
	query := fmt.Sprintf(
		"INSERT INTO %s "+
			"(%s, %s, "+
			"%s, %s, %s, "+
			"%s, %s, %s, "+
			"%s, %s, %s, %s, %s, %s, "+
			"%s, %s) "+
			"VALUES %s",
		spanfault.Table,
		spanfault.FieldTimestamp, spanfault.FieldID,
		spanfault.FieldPlatformName, spanfault.FieldAppCluster, spanfault.FieldInstanceName,
		spanfault.FieldRootServiceName, spanfault.FieldRootSpanName, spanfault.FieldRootDuration,
		spanfault.FieldParentSpanId, spanfault.FieldSpanId, spanfault.FieldServiceName, spanfault.FieldSpanName, spanfault.FieldSpanKind, spanfault.FieldFaultKind,
		spanfault.FieldResourceAttributes, spanfault.FieldSpanAttributes,
		values,
	)
	_, err := repo.client.driver.ExecContext(ctx, query)
	return err
}

func (repo *SpanFaultRepositoryImpl) buildValues(entities ...*ent.SpanFault) string {
	values := ""
	for index, entity := range entities {
		if index > 0 {
			values += ","
		}
		date := entity.Timestamp.Format(time.DateTime)
		nanoseconds := entity.Timestamp.Nanosecond()
		values += fmt.Sprintf(
			"('%s.%d', '%s', "+
				"'%s', '%s', '%s', "+
				"'%s', '%s', %v, "+
				"'%s', '%s', '%s', '%s', '%s', '%s', "+
				"%v, %v)",
			date, nanoseconds, entity.ID,
			entity.PlatformName, entity.AppCluster, entity.InstanceName,
			entity.RootServiceName, entity.RootSpanName, entity.RootDuration,
			entity.ParentSpanId, entity.SpanId, entity.ServiceName, entity.SpanName, entity.SpanKind, entity.FaultKind,
			formatMap(entity.ResourceAttributes), formatMap(entity.SpanAttributes),
		)
	}
	return values
}

func (repo *SpanFaultRepositoryImpl) Shutdown(ctx context.Context) error {
	return repo.client.delegate.Close()
}

func formatMap(attributes *schema.Attributes) string {
	if attributes == nil {
		return ""
	}
	data, _ := attributes.Value()
	bytes, _ := json.Marshal(data)
	value := string(bytes)
	value = strings.ReplaceAll(value, "'", "\\'")
	value = strings.ReplaceAll(value, "\"", "'")
	return value
}
