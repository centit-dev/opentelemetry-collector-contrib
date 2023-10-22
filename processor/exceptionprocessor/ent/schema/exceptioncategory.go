package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ExceptionCategory holds the schema definition for the ExceptionCategory entity.
type ExceptionCategory struct {
	ent.Schema
}

// Fields of the ExceptionCategory.
func (ExceptionCategory) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.String("name").MaxLen(100).NotEmpty(),
		field.Int("is_valid").Default(1),
		field.Time("create_time").Default(time.Now).Immutable(),
		field.Time("update_time").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the ExceptionCategory.
func (ExceptionCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("exception_definitions", ExceptionDefinition.Type),
	}
}

func (ExceptionCategory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tb_exception_category"},
	}
}
