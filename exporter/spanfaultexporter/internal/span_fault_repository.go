package internal

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent/predicate"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanfaultexporter/ent/spanfault"
	"go.uber.org/zap"
)

type SpanFaultRepository interface {
	GetSpanFaultsByTraceIds(ctx context.Context, traceIds []string) ([]*ent.SpanFault, error)
	SaveAll(ctx context.Context, creates []*ent.SpanFault, updates []*ent.SpanFault) error
	Shutdown(ctx context.Context) error
}

type SpanFaultRepositoryImpl struct {
	client *ClickHouseClient
	logger *zap.Logger
}

func CreateSpanFaultRepository(client *ClickHouseClient, logger *zap.Logger) SpanFaultRepository {
	return &SpanFaultRepositoryImpl{client, logger}
}

func (repo *SpanFaultRepositoryImpl) GetSpanFaultsByTraceIds(ctx context.Context, traceIds []string) ([]*ent.SpanFault, error) {
	return repo.client.delegate.SpanFault.
		Query().
		Where(spanfault.TraceIdIn(traceIds...)).
		Select(spanfault.FieldTraceId).
		All(ctx)
}

func (repo *SpanFaultRepositoryImpl) SaveAll(ctx context.Context, creates []*ent.SpanFault, updates []*ent.SpanFault) error {
	// delete updates
	if len(updates) > 0 {
		spanIds := make([]string, len(updates))
		for index, update := range updates {
			spanIds[index] = update.ID
		}
		_, err := repo.client.delegate.SpanFault.Delete().
			Where(predicate.SpanFault(fieldIn[string](spanfault.FieldID, spanIds...))).Exec(ctx)
		if err != nil {
			return err
		}
	}

	// insert all
	inserts := make([]*ent.SpanFault, 0, len(creates)+len(updates))
	inserts = append(inserts, creates...)
	inserts = append(inserts, updates...)
	values := repo.buildValues(inserts...)
	query := fmt.Sprintf(
		"INSERT INTO %s "+
			"(%s, %s, "+
			"%s, %s, %s, %s, %s, "+
			"%s, %s, %s, %s, %s, %s) "+
			"VALUES %s",
		spanfault.Table,
		spanfault.FieldTimestamp, spanfault.FieldTraceId,
		spanfault.FieldPlatformName, spanfault.FieldClusterName, spanfault.FieldInstanceName, spanfault.FieldRootServiceName, spanfault.FieldRootSpanName,
		spanfault.FieldParentSpanId, spanfault.FieldID, spanfault.FieldServiceName, spanfault.FieldSpanName, spanfault.FieldFaultKind, spanfault.FieldIsRoot,
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
				"'%s', '%s', '%s', '%s', '%s', "+
				"'%s', '%s', '%s', '%s', '%s', %v)",
			date, nanoseconds, entity.TraceId,
			entity.PlatformName, entity.ClusterName, entity.InstanceName, entity.RootServiceName, entity.RootSpanName,
			entity.ParentSpanId, entity.ID, entity.ServiceName, entity.SpanName, entity.FaultKind, entity.IsRoot,
		)
	}
	return values
}

func (repo *SpanFaultRepositoryImpl) Shutdown(ctx context.Context) error {
	return repo.client.delegate.Close()
}

// fix `otel_span_aggregations`.`SpanId` for not working
func fieldIn[T any](name string, vs ...T) func(*sql.Selector) {
	return func(s *sql.Selector) {
		v := make([]any, len(vs))
		for i := range v {
			v[i] = vs[i]
		}
		s.Where(sql.In(name, v...))
	}
}
