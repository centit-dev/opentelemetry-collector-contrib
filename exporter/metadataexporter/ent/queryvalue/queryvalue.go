// Code generated by ent, DO NOT EDIT.

package queryvalue

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the queryvalue type in the database.
	Label = "query_value"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldKeyID holds the string denoting the key_id field in the database.
	FieldKeyID = "key_id"
	// FieldValue holds the string denoting the value field in the database.
	FieldValue = "value"
	// FieldValidDate holds the string denoting the valid_date field in the database.
	FieldValidDate = "valid_date"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// EdgeKey holds the string denoting the key edge name in mutations.
	EdgeKey = "key"
	// Table holds the table name of the queryvalue in the database.
	Table = "tb_query_value"
	// KeyTable is the table that holds the key relation/edge.
	KeyTable = "tb_query_value"
	// KeyInverseTable is the table name for the QueryKey entity.
	// It exists in this package in order to avoid circular dependency with the "querykey" package.
	KeyInverseTable = "tb_query_key"
	// KeyColumn is the table column denoting the key relation/edge.
	KeyColumn = "key_id"
)

// Columns holds all SQL columns for queryvalue fields.
var Columns = []string{
	FieldID,
	FieldKeyID,
	FieldValue,
	FieldValidDate,
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

// OrderOption defines the ordering options for the QueryValue queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByKeyID orders the results by the key_id field.
func ByKeyID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldKeyID, opts...).ToFunc()
}

// ByValue orders the results by the value field.
func ByValue(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldValue, opts...).ToFunc()
}

// ByValidDate orders the results by the valid_date field.
func ByValidDate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldValidDate, opts...).ToFunc()
}

// ByCreateTime orders the results by the create_time field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the update_time field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByKeyField orders the results by key field.
func ByKeyField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newKeyStep(), sql.OrderByField(field, opts...))
	}
}
func newKeyStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(KeyInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, KeyTable, KeyColumn),
	)
}
