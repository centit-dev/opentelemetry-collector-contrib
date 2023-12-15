// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent/predicate"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent/spanaggregation"
)

// SpanAggregationQuery is the builder for querying SpanAggregation entities.
type SpanAggregationQuery struct {
	config
	ctx        *QueryContext
	order      []spanaggregation.OrderOption
	inters     []Interceptor
	predicates []predicate.SpanAggregation
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the SpanAggregationQuery builder.
func (saq *SpanAggregationQuery) Where(ps ...predicate.SpanAggregation) *SpanAggregationQuery {
	saq.predicates = append(saq.predicates, ps...)
	return saq
}

// Limit the number of records to be returned by this query.
func (saq *SpanAggregationQuery) Limit(limit int) *SpanAggregationQuery {
	saq.ctx.Limit = &limit
	return saq
}

// Offset to start from.
func (saq *SpanAggregationQuery) Offset(offset int) *SpanAggregationQuery {
	saq.ctx.Offset = &offset
	return saq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (saq *SpanAggregationQuery) Unique(unique bool) *SpanAggregationQuery {
	saq.ctx.Unique = &unique
	return saq
}

// Order specifies how the records should be ordered.
func (saq *SpanAggregationQuery) Order(o ...spanaggregation.OrderOption) *SpanAggregationQuery {
	saq.order = append(saq.order, o...)
	return saq
}

// First returns the first SpanAggregation entity from the query.
// Returns a *NotFoundError when no SpanAggregation was found.
func (saq *SpanAggregationQuery) First(ctx context.Context) (*SpanAggregation, error) {
	nodes, err := saq.Limit(1).All(setContextOp(ctx, saq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{spanaggregation.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (saq *SpanAggregationQuery) FirstX(ctx context.Context) *SpanAggregation {
	node, err := saq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first SpanAggregation ID from the query.
// Returns a *NotFoundError when no SpanAggregation ID was found.
func (saq *SpanAggregationQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = saq.Limit(1).IDs(setContextOp(ctx, saq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{spanaggregation.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (saq *SpanAggregationQuery) FirstIDX(ctx context.Context) string {
	id, err := saq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single SpanAggregation entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one SpanAggregation entity is found.
// Returns a *NotFoundError when no SpanAggregation entities are found.
func (saq *SpanAggregationQuery) Only(ctx context.Context) (*SpanAggregation, error) {
	nodes, err := saq.Limit(2).All(setContextOp(ctx, saq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{spanaggregation.Label}
	default:
		return nil, &NotSingularError{spanaggregation.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (saq *SpanAggregationQuery) OnlyX(ctx context.Context) *SpanAggregation {
	node, err := saq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only SpanAggregation ID in the query.
// Returns a *NotSingularError when more than one SpanAggregation ID is found.
// Returns a *NotFoundError when no entities are found.
func (saq *SpanAggregationQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = saq.Limit(2).IDs(setContextOp(ctx, saq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{spanaggregation.Label}
	default:
		err = &NotSingularError{spanaggregation.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (saq *SpanAggregationQuery) OnlyIDX(ctx context.Context) string {
	id, err := saq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of SpanAggregations.
func (saq *SpanAggregationQuery) All(ctx context.Context) ([]*SpanAggregation, error) {
	ctx = setContextOp(ctx, saq.ctx, "All")
	if err := saq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*SpanAggregation, *SpanAggregationQuery]()
	return withInterceptors[[]*SpanAggregation](ctx, saq, qr, saq.inters)
}

// AllX is like All, but panics if an error occurs.
func (saq *SpanAggregationQuery) AllX(ctx context.Context) []*SpanAggregation {
	nodes, err := saq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of SpanAggregation IDs.
func (saq *SpanAggregationQuery) IDs(ctx context.Context) (ids []string, err error) {
	if saq.ctx.Unique == nil && saq.path != nil {
		saq.Unique(true)
	}
	ctx = setContextOp(ctx, saq.ctx, "IDs")
	if err = saq.Select(spanaggregation.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (saq *SpanAggregationQuery) IDsX(ctx context.Context) []string {
	ids, err := saq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (saq *SpanAggregationQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, saq.ctx, "Count")
	if err := saq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, saq, querierCount[*SpanAggregationQuery](), saq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (saq *SpanAggregationQuery) CountX(ctx context.Context) int {
	count, err := saq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (saq *SpanAggregationQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, saq.ctx, "Exist")
	switch _, err := saq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (saq *SpanAggregationQuery) ExistX(ctx context.Context) bool {
	exist, err := saq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the SpanAggregationQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (saq *SpanAggregationQuery) Clone() *SpanAggregationQuery {
	if saq == nil {
		return nil
	}
	return &SpanAggregationQuery{
		config:     saq.config,
		ctx:        saq.ctx.Clone(),
		order:      append([]spanaggregation.OrderOption{}, saq.order...),
		inters:     append([]Interceptor{}, saq.inters...),
		predicates: append([]predicate.SpanAggregation{}, saq.predicates...),
		// clone intermediate query.
		sql:  saq.sql.Clone(),
		path: saq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Timestamp time.Time `json:"Timestamp,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.SpanAggregation.Query().
//		GroupBy(spanaggregation.FieldTimestamp).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (saq *SpanAggregationQuery) GroupBy(field string, fields ...string) *SpanAggregationGroupBy {
	saq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &SpanAggregationGroupBy{build: saq}
	grbuild.flds = &saq.ctx.Fields
	grbuild.label = spanaggregation.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Timestamp time.Time `json:"Timestamp,omitempty"`
//	}
//
//	client.SpanAggregation.Query().
//		Select(spanaggregation.FieldTimestamp).
//		Scan(ctx, &v)
func (saq *SpanAggregationQuery) Select(fields ...string) *SpanAggregationSelect {
	saq.ctx.Fields = append(saq.ctx.Fields, fields...)
	sbuild := &SpanAggregationSelect{SpanAggregationQuery: saq}
	sbuild.label = spanaggregation.Label
	sbuild.flds, sbuild.scan = &saq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a SpanAggregationSelect configured with the given aggregations.
func (saq *SpanAggregationQuery) Aggregate(fns ...AggregateFunc) *SpanAggregationSelect {
	return saq.Select().Aggregate(fns...)
}

func (saq *SpanAggregationQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range saq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, saq); err != nil {
				return err
			}
		}
	}
	for _, f := range saq.ctx.Fields {
		if !spanaggregation.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if saq.path != nil {
		prev, err := saq.path(ctx)
		if err != nil {
			return err
		}
		saq.sql = prev
	}
	return nil
}

func (saq *SpanAggregationQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*SpanAggregation, error) {
	var (
		nodes = []*SpanAggregation{}
		_spec = saq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*SpanAggregation).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &SpanAggregation{config: saq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, saq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (saq *SpanAggregationQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := saq.querySpec()
	_spec.Node.Columns = saq.ctx.Fields
	if len(saq.ctx.Fields) > 0 {
		_spec.Unique = saq.ctx.Unique != nil && *saq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, saq.driver, _spec)
}

func (saq *SpanAggregationQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(spanaggregation.Table, spanaggregation.Columns, sqlgraph.NewFieldSpec(spanaggregation.FieldID, field.TypeString))
	_spec.From = saq.sql
	if unique := saq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if saq.path != nil {
		_spec.Unique = true
	}
	if fields := saq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, spanaggregation.FieldID)
		for i := range fields {
			if fields[i] != spanaggregation.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := saq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := saq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := saq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := saq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (saq *SpanAggregationQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(saq.driver.Dialect())
	t1 := builder.Table(spanaggregation.Table)
	columns := saq.ctx.Fields
	if len(columns) == 0 {
		columns = spanaggregation.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if saq.sql != nil {
		selector = saq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if saq.ctx.Unique != nil && *saq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range saq.predicates {
		p(selector)
	}
	for _, p := range saq.order {
		p(selector)
	}
	if offset := saq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := saq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// SpanAggregationGroupBy is the group-by builder for SpanAggregation entities.
type SpanAggregationGroupBy struct {
	selector
	build *SpanAggregationQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (sagb *SpanAggregationGroupBy) Aggregate(fns ...AggregateFunc) *SpanAggregationGroupBy {
	sagb.fns = append(sagb.fns, fns...)
	return sagb
}

// Scan applies the selector query and scans the result into the given value.
func (sagb *SpanAggregationGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, sagb.build.ctx, "GroupBy")
	if err := sagb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*SpanAggregationQuery, *SpanAggregationGroupBy](ctx, sagb.build, sagb, sagb.build.inters, v)
}

func (sagb *SpanAggregationGroupBy) sqlScan(ctx context.Context, root *SpanAggregationQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(sagb.fns))
	for _, fn := range sagb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*sagb.flds)+len(sagb.fns))
		for _, f := range *sagb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*sagb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := sagb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// SpanAggregationSelect is the builder for selecting fields of SpanAggregation entities.
type SpanAggregationSelect struct {
	*SpanAggregationQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (sas *SpanAggregationSelect) Aggregate(fns ...AggregateFunc) *SpanAggregationSelect {
	sas.fns = append(sas.fns, fns...)
	return sas
}

// Scan applies the selector query and scans the result into the given value.
func (sas *SpanAggregationSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, sas.ctx, "Select")
	if err := sas.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*SpanAggregationQuery, *SpanAggregationSelect](ctx, sas.SpanAggregationQuery, sas, sas.inters, v)
}

func (sas *SpanAggregationSelect) sqlScan(ctx context.Context, root *SpanAggregationQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(sas.fns))
	for _, fn := range sas.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*sas.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := sas.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}