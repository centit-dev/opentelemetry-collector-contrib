package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type SoftwareDefinitionCondition struct {
	Column string
	Op     string
	Value  interface{}
}

// SoftwareDefinition holds the schema definition for the SoftwareDefinition entity.
type SoftwareDefinition struct {
	ent.Schema
}

// Fields of the MiddlewareDefinition.
func (SoftwareDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Positive(),
		field.String("name").MaxLen(100).NotEmpty(),
		field.Int16("type"),
		field.JSON("span_conditions", []SoftwareDefinitionCondition{}).Optional(),
		field.Int("is_valid").Default(1),
		field.Time("create_time").Immutable().Default(time.Now),
		field.Time("update_time").Default(time.Now),
	}
}

// Edges of the MiddlewareDefinition.
func (SoftwareDefinition) Edges() []ent.Edge {
	return nil
}

func (SoftwareDefinition) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tb_software_define"},
	}
}
