// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/exceptioncategory"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/exceptiondefinition"
)

// ExceptionDefinition is the model entity for the ExceptionDefinition schema.
type ExceptionDefinition struct {
	config `json:"-"`
	// ID of the ent.
	ID int64 `json:"id,omitempty"`
	// CategoryID holds the value of the "category_id" field.
	CategoryID int64 `json:"category_id,omitempty"`
	// ShortName holds the value of the "short_name" field.
	ShortName string `json:"short_name,omitempty"`
	// LongName holds the value of the "long_name" field.
	LongName string `json:"long_name,omitempty"`
	// IsValid holds the value of the "is_valid" field.
	IsValid int `json:"is_valid,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ExceptionDefinitionQuery when eager-loading is set.
	Edges        ExceptionDefinitionEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ExceptionDefinitionEdges holds the relations/edges for other nodes in the graph.
type ExceptionDefinitionEdges struct {
	// ExceptionCategory holds the value of the exception_category edge.
	ExceptionCategory *ExceptionCategory `json:"exception_category,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ExceptionCategoryOrErr returns the ExceptionCategory value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExceptionDefinitionEdges) ExceptionCategoryOrErr() (*ExceptionCategory, error) {
	if e.loadedTypes[0] {
		if e.ExceptionCategory == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: exceptioncategory.Label}
		}
		return e.ExceptionCategory, nil
	}
	return nil, &NotLoadedError{edge: "exception_category"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ExceptionDefinition) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case exceptiondefinition.FieldID, exceptiondefinition.FieldCategoryID, exceptiondefinition.FieldIsValid:
			values[i] = new(sql.NullInt64)
		case exceptiondefinition.FieldShortName, exceptiondefinition.FieldLongName:
			values[i] = new(sql.NullString)
		case exceptiondefinition.FieldCreateTime, exceptiondefinition.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ExceptionDefinition fields.
func (ed *ExceptionDefinition) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case exceptiondefinition.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ed.ID = int64(value.Int64)
		case exceptiondefinition.FieldCategoryID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field category_id", values[i])
			} else if value.Valid {
				ed.CategoryID = value.Int64
			}
		case exceptiondefinition.FieldShortName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field short_name", values[i])
			} else if value.Valid {
				ed.ShortName = value.String
			}
		case exceptiondefinition.FieldLongName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field long_name", values[i])
			} else if value.Valid {
				ed.LongName = value.String
			}
		case exceptiondefinition.FieldIsValid:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field is_valid", values[i])
			} else if value.Valid {
				ed.IsValid = int(value.Int64)
			}
		case exceptiondefinition.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				ed.CreateTime = value.Time
			}
		case exceptiondefinition.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				ed.UpdateTime = value.Time
			}
		default:
			ed.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ExceptionDefinition.
// This includes values selected through modifiers, order, etc.
func (ed *ExceptionDefinition) Value(name string) (ent.Value, error) {
	return ed.selectValues.Get(name)
}

// QueryExceptionCategory queries the "exception_category" edge of the ExceptionDefinition entity.
func (ed *ExceptionDefinition) QueryExceptionCategory() *ExceptionCategoryQuery {
	return NewExceptionDefinitionClient(ed.config).QueryExceptionCategory(ed)
}

// Update returns a builder for updating this ExceptionDefinition.
// Note that you need to call ExceptionDefinition.Unwrap() before calling this method if this ExceptionDefinition
// was returned from a transaction, and the transaction was committed or rolled back.
func (ed *ExceptionDefinition) Update() *ExceptionDefinitionUpdateOne {
	return NewExceptionDefinitionClient(ed.config).UpdateOne(ed)
}

// Unwrap unwraps the ExceptionDefinition entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ed *ExceptionDefinition) Unwrap() *ExceptionDefinition {
	_tx, ok := ed.config.driver.(*txDriver)
	if !ok {
		panic("ent: ExceptionDefinition is not a transactional entity")
	}
	ed.config.driver = _tx.drv
	return ed
}

// String implements the fmt.Stringer.
func (ed *ExceptionDefinition) String() string {
	var builder strings.Builder
	builder.WriteString("ExceptionDefinition(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ed.ID))
	builder.WriteString("category_id=")
	builder.WriteString(fmt.Sprintf("%v", ed.CategoryID))
	builder.WriteString(", ")
	builder.WriteString("short_name=")
	builder.WriteString(ed.ShortName)
	builder.WriteString(", ")
	builder.WriteString("long_name=")
	builder.WriteString(ed.LongName)
	builder.WriteString(", ")
	builder.WriteString("is_valid=")
	builder.WriteString(fmt.Sprintf("%v", ed.IsValid))
	builder.WriteString(", ")
	builder.WriteString("create_time=")
	builder.WriteString(ed.CreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(ed.UpdateTime.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// ExceptionDefinitions is a parsable slice of ExceptionDefinition.
type ExceptionDefinitions []*ExceptionDefinition
