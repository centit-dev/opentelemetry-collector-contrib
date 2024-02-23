// Code generated by ent, DO NOT EDIT.

package exceptiondefinition

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldLTE(FieldID, id))
}

// CategoryID applies equality check predicate on the "category_id" field. It's identical to CategoryIDEQ.
func CategoryID(v int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldCategoryID, v))
}

// ShortName applies equality check predicate on the "short_name" field. It's identical to ShortNameEQ.
func ShortName(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldShortName, v))
}

// LongName applies equality check predicate on the "long_name" field. It's identical to LongNameEQ.
func LongName(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldLongName, v))
}

// IsValid applies equality check predicate on the "is_valid" field. It's identical to IsValidEQ.
func IsValid(v int) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldIsValid, v))
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldUpdateTime, v))
}

// CategoryIDEQ applies the EQ predicate on the "category_id" field.
func CategoryIDEQ(v int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldCategoryID, v))
}

// CategoryIDNEQ applies the NEQ predicate on the "category_id" field.
func CategoryIDNEQ(v int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNEQ(FieldCategoryID, v))
}

// CategoryIDIn applies the In predicate on the "category_id" field.
func CategoryIDIn(vs ...int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldIn(FieldCategoryID, vs...))
}

// CategoryIDNotIn applies the NotIn predicate on the "category_id" field.
func CategoryIDNotIn(vs ...int64) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNotIn(FieldCategoryID, vs...))
}

// ShortNameEQ applies the EQ predicate on the "short_name" field.
func ShortNameEQ(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldShortName, v))
}

// ShortNameNEQ applies the NEQ predicate on the "short_name" field.
func ShortNameNEQ(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNEQ(FieldShortName, v))
}

// ShortNameIn applies the In predicate on the "short_name" field.
func ShortNameIn(vs ...string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldIn(FieldShortName, vs...))
}

// ShortNameNotIn applies the NotIn predicate on the "short_name" field.
func ShortNameNotIn(vs ...string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNotIn(FieldShortName, vs...))
}

// ShortNameGT applies the GT predicate on the "short_name" field.
func ShortNameGT(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldGT(FieldShortName, v))
}

// ShortNameGTE applies the GTE predicate on the "short_name" field.
func ShortNameGTE(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldGTE(FieldShortName, v))
}

// ShortNameLT applies the LT predicate on the "short_name" field.
func ShortNameLT(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldLT(FieldShortName, v))
}

// ShortNameLTE applies the LTE predicate on the "short_name" field.
func ShortNameLTE(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldLTE(FieldShortName, v))
}

// ShortNameContains applies the Contains predicate on the "short_name" field.
func ShortNameContains(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldContains(FieldShortName, v))
}

// ShortNameHasPrefix applies the HasPrefix predicate on the "short_name" field.
func ShortNameHasPrefix(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldHasPrefix(FieldShortName, v))
}

// ShortNameHasSuffix applies the HasSuffix predicate on the "short_name" field.
func ShortNameHasSuffix(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldHasSuffix(FieldShortName, v))
}

// ShortNameEqualFold applies the EqualFold predicate on the "short_name" field.
func ShortNameEqualFold(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEqualFold(FieldShortName, v))
}

// ShortNameContainsFold applies the ContainsFold predicate on the "short_name" field.
func ShortNameContainsFold(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldContainsFold(FieldShortName, v))
}

// LongNameEQ applies the EQ predicate on the "long_name" field.
func LongNameEQ(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldLongName, v))
}

// LongNameNEQ applies the NEQ predicate on the "long_name" field.
func LongNameNEQ(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNEQ(FieldLongName, v))
}

// LongNameIn applies the In predicate on the "long_name" field.
func LongNameIn(vs ...string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldIn(FieldLongName, vs...))
}

// LongNameNotIn applies the NotIn predicate on the "long_name" field.
func LongNameNotIn(vs ...string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNotIn(FieldLongName, vs...))
}

// LongNameGT applies the GT predicate on the "long_name" field.
func LongNameGT(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldGT(FieldLongName, v))
}

// LongNameGTE applies the GTE predicate on the "long_name" field.
func LongNameGTE(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldGTE(FieldLongName, v))
}

// LongNameLT applies the LT predicate on the "long_name" field.
func LongNameLT(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldLT(FieldLongName, v))
}

// LongNameLTE applies the LTE predicate on the "long_name" field.
func LongNameLTE(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldLTE(FieldLongName, v))
}

// LongNameContains applies the Contains predicate on the "long_name" field.
func LongNameContains(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldContains(FieldLongName, v))
}

// LongNameHasPrefix applies the HasPrefix predicate on the "long_name" field.
func LongNameHasPrefix(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldHasPrefix(FieldLongName, v))
}

// LongNameHasSuffix applies the HasSuffix predicate on the "long_name" field.
func LongNameHasSuffix(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldHasSuffix(FieldLongName, v))
}

// LongNameEqualFold applies the EqualFold predicate on the "long_name" field.
func LongNameEqualFold(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEqualFold(FieldLongName, v))
}

// LongNameContainsFold applies the ContainsFold predicate on the "long_name" field.
func LongNameContainsFold(v string) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldContainsFold(FieldLongName, v))
}

// IsValidEQ applies the EQ predicate on the "is_valid" field.
func IsValidEQ(v int) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldIsValid, v))
}

// IsValidNEQ applies the NEQ predicate on the "is_valid" field.
func IsValidNEQ(v int) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNEQ(FieldIsValid, v))
}

// IsValidIn applies the In predicate on the "is_valid" field.
func IsValidIn(vs ...int) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldIn(FieldIsValid, vs...))
}

// IsValidNotIn applies the NotIn predicate on the "is_valid" field.
func IsValidNotIn(vs ...int) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNotIn(FieldIsValid, vs...))
}

// IsValidGT applies the GT predicate on the "is_valid" field.
func IsValidGT(v int) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldGT(FieldIsValid, v))
}

// IsValidGTE applies the GTE predicate on the "is_valid" field.
func IsValidGTE(v int) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldGTE(FieldIsValid, v))
}

// IsValidLT applies the LT predicate on the "is_valid" field.
func IsValidLT(v int) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldLT(FieldIsValid, v))
}

// IsValidLTE applies the LTE predicate on the "is_valid" field.
func IsValidLTE(v int) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldLTE(FieldIsValid, v))
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.FieldLTE(FieldUpdateTime, v))
}

// HasExceptionCategory applies the HasEdge predicate on the "exception_category" edge.
func HasExceptionCategory() predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ExceptionCategoryTable, ExceptionCategoryColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasExceptionCategoryWith applies the HasEdge predicate on the "exception_category" edge with a given conditions (other predicates).
func HasExceptionCategoryWith(preds ...predicate.ExceptionCategory) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(func(s *sql.Selector) {
		step := newExceptionCategoryStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ExceptionDefinition) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ExceptionDefinition) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ExceptionDefinition) predicate.ExceptionDefinition {
	return predicate.ExceptionDefinition(sql.NotPredicates(p))
}
