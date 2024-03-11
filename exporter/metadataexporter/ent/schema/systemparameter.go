package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// SystemParameter holds the schema definition for the SystemParameter entity.
type SystemParameter struct {
	ent.Schema
}

func (SystemParameter) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tb_sys_parameter"},
	}
}

// Fields of the SystemParameter.
func (SystemParameter) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("code"),
		// TODO - this is a JSON field, let's assume it's a string array for now
		field.JSON("value", []string{}),
		field.Time("create_time"),
		field.Time("update_time"),
	}
}
