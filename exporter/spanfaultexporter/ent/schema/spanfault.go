package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

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
		field.String("FaultKind").StorageKey("FaultKind"),
	}
}

// Edges of the SpanFault.
func (SpanFault) Edges() []ent.Edge {
	return nil
}
