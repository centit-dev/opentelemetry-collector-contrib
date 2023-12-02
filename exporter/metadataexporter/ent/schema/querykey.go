package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// QueryKey holds the schema definition for the QueryKey entity.
type QueryKey struct {
	ent.Schema
}

func (QueryKey) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tb_query_key"},
	}
}

// Fields of the QueryKey.
func (QueryKey) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.String("name"),
		field.String("type"),
		field.String("source"),
		field.Time("valid_date"),
		field.Time("create_time"),
		field.Time("update_time"),
	}
}

// Edges of the QueryKey.
func (QueryKey) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("values", QueryValue.Type),
	}
}
