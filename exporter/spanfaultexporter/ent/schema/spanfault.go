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

// SpanFault holds the schema definition for the SpanFault entity.
type SpanFault struct {
	ent.Schema
}

func (SpanFault) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "otel_trace_faults"},
	}
}

// Fields of the SpanFault.
func (SpanFault) Fields() []ent.Field {
	return []ent.Field{
		field.Time("Timestamp").StorageKey("Timestamp").Optional(),
		field.String("id").StorageKey("TraceId"),
		field.String("PlatformName").StorageKey("PlatformName"),
		field.String("AppCluster").StorageKey("AppCluster"),
		field.String("InstanceName").StorageKey("InstanceName"),
		field.String("RootServiceName").StorageKey("RootServiceName"),
		field.String("RootSpanName").StorageKey("RootSpanName"),
		field.Int64("RootDuration").StorageKey("RootDuration"),
		field.String("ParentSpanId").StorageKey("ParentSpanId"),
		field.String("SpanId").StorageKey("SpanId"),
		field.String("ServiceName").StorageKey("ServiceName"),
		field.String("SpanName").StorageKey("SpanName"),
		field.String("SpanKind").StorageKey("SpanKind"),
		field.String("FaultKind").StorageKey("FaultKind"),
		field.Int64("Gap").StorageKey("Gap"),
		field.Int64("SelfDuration").StorageKey("SelfDuration"),
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
	}
}

// Edges of the SpanFault.
func (SpanFault) Edges() []ent.Edge {
	return nil
}
