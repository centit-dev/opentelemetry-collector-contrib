package schema

import (
	"database/sql/driver"
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type FaultKindCondition struct {
	Column string      `json:"column"`
	Op     string      `json:"op"`
	Value  interface{} `json:"value"`
}

type FaultKindDefinitions struct {
	System   []*FaultKindCondition `json:"system"`
	Business []*FaultKindCondition `json:"business"`
}

func (d *FaultKindDefinitions) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, d)
	case string:
		return json.Unmarshal([]byte(v), d)
	default:
		return nil
	}
}

func (d FaultKindDefinitions) Value() (driver.Value, error) {
	return json.Marshal(d)
}

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
		field.String("id").StorageKey("code").Unique(),
		field.Other("value", &FaultKindDefinitions{}).
			SchemaType(map[string]string{
				dialect.Postgres: "text",
			}),
	}
}

// Edges of the SystemParameter.
func (SystemParameter) Edges() []ent.Edge {
	return nil
}
