package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// ApplicationStructure holds the schema definition for the ApplicationStructure entity.
type ApplicationStructure struct {
	ent.Schema
}

func (ApplicationStructure) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tb_application_structure"},
	}
}

// Fields of the ApplicationStructure.
func (ApplicationStructure) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().StorageKey("code").MaxLen(100),
		field.String("parentCode").MaxLen(100),
		field.Int("level"),
		field.Time("valid_date"),
		field.Time("create_time"),
		field.Time("update_time"),
	}
}

// Edges of the ApplicationStructure.
func (ApplicationStructure) Edges() []ent.Edge {
	return nil
}
