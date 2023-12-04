package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// QueryValue holds the schema definition for the QueryValue entity.
type QueryValue struct {
	ent.Schema
}

func (QueryValue) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tb_query_value"},
	}
}

// Fields of the QueryValue.
func (QueryValue) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.Int64("key_id"),
		field.String("value"),
		field.Time("valid_date"),
		field.Time("create_time"),
		field.Time("update_time"),
	}
}

// Edges of the QueryValue.
func (QueryValue) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("key", QueryKey.Type).Ref("values").Field("key_id").Required().Unique(),
	}
}
