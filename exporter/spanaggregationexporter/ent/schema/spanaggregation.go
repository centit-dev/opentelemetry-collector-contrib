package schema

import (
	"database/sql/driver"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Attributes struct {
	value *map[string]string
}

func (attributes *Attributes) Add(key string, value string) {
	if attributes.value == nil {
		attributes.value = &map[string]string{}
	}
	(*attributes.value)[key] = value
}

func (attributes *Attributes) Get(key string) string {
	if attributes.value == nil {
		return ""
	}
	return (*attributes.value)[key]
}

func (attributes *Attributes) Value() (driver.Value, error) {
	return attributes.value, nil
}

func (attributes *Attributes) Scan(src interface{}) error {
	switch src := src.(type) {
	case map[string]string:
		attributes.value = &src
		return nil
	default:
		return nil
	}
}

// a span representation aggregated by SpanId with the root span, the parent span, the children span information
type SpanAggregation struct {
	ent.Schema
}

func (SpanAggregation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "otel_span_aggregations"},
	}
}

// Fields of the SpanAggregation.
func (SpanAggregation) Fields() []ent.Field {
	return []ent.Field{
		field.Time("Timestamp").StorageKey("Timestamp"),
		field.String("TraceId").StorageKey("TraceId"),
		field.String("id").StorageKey("SpanId"),
		field.String("ParentSpanId").StorageKey("ParentSpanId"),
		field.String("PlatformName").StorageKey("PlatformName"),
		field.String("RootServiceName").StorageKey("RootServiceName").Optional(),
		field.String("RootSpanName").StorageKey("RootSpanName").Optional(),
		field.String("ServiceName").StorageKey("ServiceName"),
		field.String("SpanName").StorageKey("SpanName"),
		field.Other("ResourceAttributes", &Attributes{}).
			SchemaType(map[string]string{
				"clickhouse": "Map(String, String)",
			}).
			StorageKey("ResourceAttributes"),
		field.Other("SpanAttributes", &Attributes{}).
			SchemaType(map[string]string{
				"clickhouse": "Map(String, String)",
			}).
			StorageKey("SpanAttributes"),
		field.Int64("Duration").StorageKey("Duration"),
		field.Int64("Gap").StorageKey("Gap"),
		field.Int64("SelfDuration").StorageKey("SelfDuration"),
	}
}

// Edges of the SpanAggregation.
func (SpanAggregation) Edges() []ent.Edge {
	return nil
}
