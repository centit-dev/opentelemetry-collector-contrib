package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent/predicate"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent/schema"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent/spanaggregation"
	"go.uber.org/zap"
)

type SpanAggregationRepository interface {
	FindAllByTraceId(ctx context.Context, traceId string) ([]*ent.SpanAggregation, error)
	// save all the span aggregations
	SaveAll(ctx context.Context, creates []*ent.SpanAggregation, updates []*ent.SpanAggregation) error
	Shutdown(ctx context.Context) error
}

type SpanAggregationRepositoryImpl struct {
	client *ClickHouseClient
	logger *zap.Logger
}

func CreateSpanAggregationRepositoryImpl(client *ClickHouseClient, logger *zap.Logger) *SpanAggregationRepositoryImpl {
	return &SpanAggregationRepositoryImpl{client, logger}
}

func (repository *SpanAggregationRepositoryImpl) FindAllByTraceId(ctx context.Context, traceId string) ([]*ent.SpanAggregation, error) {
	aggregations, err := repository.client.delegate.SpanAggregation.Query().
		Where(spanaggregation.TraceIdEQ(traceId)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return aggregations, nil
}

// TBD aggregations pointers are released?
func (repository *SpanAggregationRepositoryImpl) SaveAll(ctx context.Context, creates []*ent.SpanAggregation, updates []*ent.SpanAggregation) error {
	// delete the updatings
	if len(updates) > 0 {
		spanIds := make([]string, len(updates))
		for i, toUpdate := range updates {
			spanIds[i] = toUpdate.ID
		}
		_, err := repository.client.delegate.SpanAggregation.Delete().
			Where(predicate.SpanAggregation(fieldIn[string](spanaggregation.FieldID, spanIds...))).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	// insert all
	// go ent doesn't work well with clickhouse bulk insert
	// so we have to use raw sql for now
	inserts := make([]*ent.SpanAggregation, 0, len(creates)+len(updates))
	inserts = append(inserts, creates...)
	inserts = append(inserts, updates...)
	query := "INSERT INTO `otel_span_aggregations` " +
		"(" +
		"`Timestamp`, `TraceId`, `ParentSpanId`, `SpanId`, " +
		"`PlatformName`, `RootServiceName`, `RootSpanName`, " +
		"`ServiceName`, `SpanName`, " +
		"`ResourceAttributes`, `SpanAttributes`, " +
		"`Duration`, `Gap`, `SelfDuration`)" +
		" VALUES " +
		repository.buildBulkInsert(inserts...)
	_, err := repository.client.driver.ExecContext(ctx, query)
	return err
}

func (SpanAggregationRepositoryImpl) buildBulkInsert(entities ...*ent.SpanAggregation) string {
	sql := ""
	for _, entity := range entities {
		date := entity.Timestamp.Format(time.DateTime)
		nanoseconds := entity.Timestamp.Nanosecond()
		value := fmt.Sprintf(
			"('%s.%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %v, %v, %d, %d, %d)",
			date, nanoseconds,
			entity.TraceId,
			entity.ParentSpanId,
			entity.ID,
			entity.PlatformName,
			entity.RootServiceName,
			entity.RootSpanName,
			entity.ServiceName,
			entity.SpanName,
			formatMap(entity.ResourceAttributes),
			formatMap(entity.SpanAttributes),
			entity.Duration,
			entity.Gap,
			entity.SelfDuration,
		)
		if sql == "" {
			sql = value
		} else {
			sql = fmt.Sprintf("%s,\n%s", sql, value)
		}
	}
	return sql
}

func (repository *SpanAggregationRepositoryImpl) Shutdown(ctx context.Context) error {
	return repository.client.Shutdown()
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

func formatMap(attributes *schema.Attributes) string {
	if attributes == nil {
		return ""
	}
	data, _ := attributes.Value()
	bytes, _ := json.Marshal(data)
	value := string(bytes)
	value = strings.ReplaceAll(value, "\"", "'")
	return value
}
