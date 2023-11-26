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
		entsql.Annotation{Table: "otel_span_faults"},
	}
}

// Fields of the SpanFault.
func (SpanFault) Fields() []ent.Field {
	return []ent.Field{
		field.Time("Timestamp").StorageKey("Timestamp").Optional(),
		field.String("TraceId").StorageKey("TraceId"),
		field.String("PlatformName").StorageKey("PlatformName"),
		field.String("ClusterName").StorageKey("ClusterName"),
		field.String("InstanceName").StorageKey("InstanceName"),
		field.String("RootServiceName").StorageKey("RootServiceName").Optional(),
		field.String("RootSpanName").StorageKey("RootSpanName").Optional(),
		field.String("ParentSpanId").StorageKey("ParentSpanId"),
		field.String("id").StorageKey("SpanId"),
		field.String("ServiceName").StorageKey("ServiceName"),
		field.String("SpanName").StorageKey("SpanName"),
		field.String("FaultKind").StorageKey("FaultKind"),
		field.Bool("IsRoot").StorageKey("IsRoot"),
	}
}

// Edges of the SpanFault.
func (SpanFault) Edges() []ent.Edge {
	return nil
}
