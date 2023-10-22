// Code generated by ent, DO NOT EDIT.

package exceptiondefinition

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the exceptiondefinition type in the database.
	Label = "exception_definition"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCategoryID holds the string denoting the category_id field in the database.
	FieldCategoryID = "category_id"
	// FieldShortName holds the string denoting the short_name field in the database.
	FieldShortName = "short_name"
	// FieldLongName holds the string denoting the long_name field in the database.
	FieldLongName = "long_name"
	// FieldRelatedMiddlewareID holds the string denoting the related_middleware_id field in the database.
	FieldRelatedMiddlewareID = "related_middleware_id"
	// FieldRelatedMiddlewareConditions holds the string denoting the related_middleware_conditions field in the database.
	FieldRelatedMiddlewareConditions = "related_middleware_conditions"
	// FieldIsValid holds the string denoting the is_valid field in the database.
	FieldIsValid = "is_valid"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// EdgeExceptionCategory holds the string denoting the exception_category edge name in mutations.
	EdgeExceptionCategory = "exception_category"
	// Table holds the table name of the exceptiondefinition in the database.
	Table = "tb_exception_define"
	// ExceptionCategoryTable is the table that holds the exception_category relation/edge.
	ExceptionCategoryTable = "tb_exception_define"
	// ExceptionCategoryInverseTable is the table name for the ExceptionCategory entity.
	// It exists in this package in order to avoid circular dependency with the "exceptioncategory" package.
	ExceptionCategoryInverseTable = "tb_exception_category"
	// ExceptionCategoryColumn is the table column denoting the exception_category relation/edge.
	ExceptionCategoryColumn = "category_id"
)

// Columns holds all SQL columns for exceptiondefinition fields.
var Columns = []string{
	FieldID,
	FieldCategoryID,
	FieldShortName,
	FieldLongName,
	FieldRelatedMiddlewareID,
	FieldRelatedMiddlewareConditions,
	FieldIsValid,
	FieldCreateTime,
	FieldUpdateTime,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// CategoryIDValidator is a validator for the "category_id" field. It is called by the builders before save.
	CategoryIDValidator func(int64) error
	// ShortNameValidator is a validator for the "short_name" field. It is called by the builders before save.
	ShortNameValidator func(string) error
	// LongNameValidator is a validator for the "long_name" field. It is called by the builders before save.
	LongNameValidator func(string) error
	// DefaultIsValid holds the default value on creation for the "is_valid" field.
	DefaultIsValid int
	// DefaultCreateTime holds the default value on creation for the "create_time" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "update_time" field.
	DefaultUpdateTime func() time.Time
	// IDValidator is a validator for the "id" field. It is called by the builders before save.
	IDValidator func(int64) error
)

// OrderOption defines the ordering options for the ExceptionDefinition queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCategoryID orders the results by the category_id field.
func ByCategoryID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCategoryID, opts...).ToFunc()
}

// ByShortName orders the results by the short_name field.
func ByShortName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldShortName, opts...).ToFunc()
}

// ByLongName orders the results by the long_name field.
func ByLongName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLongName, opts...).ToFunc()
}

// ByRelatedMiddlewareID orders the results by the related_middleware_id field.
func ByRelatedMiddlewareID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRelatedMiddlewareID, opts...).ToFunc()
}

// ByIsValid orders the results by the is_valid field.
func ByIsValid(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsValid, opts...).ToFunc()
}

// ByCreateTime orders the results by the create_time field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the update_time field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByExceptionCategoryField orders the results by exception_category field.
func ByExceptionCategoryField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newExceptionCategoryStep(), sql.OrderByField(field, opts...))
	}
}
func newExceptionCategoryStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ExceptionCategoryInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ExceptionCategoryTable, ExceptionCategoryColumn),
	)
}
