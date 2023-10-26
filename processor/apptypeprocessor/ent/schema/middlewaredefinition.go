package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type MiddlewareDefinitionCondition struct {
	Column string
	Op     string
	Value  interface{}
}

// MiddlewareDefinition holds the schema definition for the MiddlewareDefinition entity.
type MiddlewareDefinition struct {
	ent.Schema
}

// Fields of the MiddlewareDefinition.
func (MiddlewareDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Positive(),
		field.String("name").MaxLen(100).NotEmpty(),
		field.JSON("span_conditions", []MiddlewareDefinitionCondition{}).Optional(),
		field.Int("is_valid").Default(1),
		field.Time("create_time").Immutable().Default(time.Now),
		field.Time("update_time").Default(time.Now),
	}
}

// Edges of the MiddlewareDefinition.
func (MiddlewareDefinition) Edges() []ent.Edge {
	return nil
}

func (MiddlewareDefinition) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tb_middleware_define"},
	}
}
