// Code generated by ent, DO NOT EDIT.

package spanfault

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the spanfault type in the database.
	Label = "span_fault"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "TraceId"
	// FieldTimestamp holds the string denoting the timestamp field in the database.
	FieldTimestamp = "Timestamp"
	// FieldPlatformName holds the string denoting the platformname field in the database.
	FieldPlatformName = "PlatformName"
	// FieldAppCluster holds the string denoting the appcluster field in the database.
	FieldAppCluster = "AppCluster"
	// FieldInstanceName holds the string denoting the instancename field in the database.
	FieldInstanceName = "InstanceName"
	// FieldRootServiceName holds the string denoting the rootservicename field in the database.
	FieldRootServiceName = "RootServiceName"
	// FieldRootSpanName holds the string denoting the rootspanname field in the database.
	FieldRootSpanName = "RootSpanName"
	// FieldRootDuration holds the string denoting the rootduration field in the database.
	FieldRootDuration = "RootDuration"
	// FieldParentSpanId holds the string denoting the parentspanid field in the database.
	FieldParentSpanId = "ParentSpanId"
	// FieldSpanId holds the string denoting the spanid field in the database.
	FieldSpanId = "SpanId"
	// FieldServiceName holds the string denoting the servicename field in the database.
	FieldServiceName = "ServiceName"
	// FieldSpanName holds the string denoting the spanname field in the database.
	FieldSpanName = "SpanName"
	// FieldFaultKind holds the string denoting the faultkind field in the database.
	FieldFaultKind = "FaultKind"
	// Table holds the table name of the spanfault in the database.
	Table = "otel_trace_faults"
)

// Columns holds all SQL columns for spanfault fields.
var Columns = []string{
	FieldID,
	FieldTimestamp,
	FieldPlatformName,
	FieldAppCluster,
	FieldInstanceName,
	FieldRootServiceName,
	FieldRootSpanName,
	FieldRootDuration,
	FieldParentSpanId,
	FieldSpanId,
	FieldServiceName,
	FieldSpanName,
	FieldFaultKind,
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

// OrderOption defines the ordering options for the SpanFault queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByTimestamp orders the results by the Timestamp field.
func ByTimestamp(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTimestamp, opts...).ToFunc()
}

// ByPlatformName orders the results by the PlatformName field.
func ByPlatformName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPlatformName, opts...).ToFunc()
}

// ByAppCluster orders the results by the AppCluster field.
func ByAppCluster(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAppCluster, opts...).ToFunc()
}

// ByInstanceName orders the results by the InstanceName field.
func ByInstanceName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldInstanceName, opts...).ToFunc()
}

// ByRootServiceName orders the results by the RootServiceName field.
func ByRootServiceName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRootServiceName, opts...).ToFunc()
}

// ByRootSpanName orders the results by the RootSpanName field.
func ByRootSpanName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRootSpanName, opts...).ToFunc()
}

// ByRootDuration orders the results by the RootDuration field.
func ByRootDuration(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRootDuration, opts...).ToFunc()
}

// ByParentSpanId orders the results by the ParentSpanId field.
func ByParentSpanId(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldParentSpanId, opts...).ToFunc()
}

// BySpanId orders the results by the SpanId field.
func BySpanId(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSpanId, opts...).ToFunc()
}

// ByServiceName orders the results by the ServiceName field.
func ByServiceName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldServiceName, opts...).ToFunc()
}

// BySpanName orders the results by the SpanName field.
func BySpanName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSpanName, opts...).ToFunc()
}

// ByFaultKind orders the results by the FaultKind field.
func ByFaultKind(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFaultKind, opts...).ToFunc()
}
