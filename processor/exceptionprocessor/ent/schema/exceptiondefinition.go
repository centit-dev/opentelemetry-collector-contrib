package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ExceptionDefinitionCondition struct {
	Column string
	Op     string
	Value  interface{}
}

// ExceptionDefinition holds the schema definition for the ExceptionDefinition entity.
type ExceptionDefinition struct {
	ent.Schema
}

// Fields of the ExceptionDefinition.
func (ExceptionDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Positive(),
		field.Int64("category_id").Positive(),
		field.String("short_name").MaxLen(100).NotEmpty(),
		field.Text("long_name").NotEmpty(),
		field.Int64("related_middleware_id").Optional(),
		field.JSON("related_middleware_conditions", []ExceptionDefinitionCondition{}).Optional(),
		field.Int("is_valid").Default(1),
		field.Time("create_time").Immutable().Default(time.Now),
		field.Time("update_time").Default(time.Now),
	}
}

// Edges of the ExceptionDefinition.
func (ExceptionDefinition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("exception_category", ExceptionCategory.Type).
			Required().
			Ref("exception_definitions").
			Field("category_id").
			Unique(),
	}
}

func (ExceptionDefinition) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tb_exception_define"},
	}
}
