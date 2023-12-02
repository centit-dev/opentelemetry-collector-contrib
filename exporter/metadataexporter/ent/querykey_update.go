// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent/predicate"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent/querykey"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent/queryvalue"
)

// QueryKeyUpdate is the builder for updating QueryKey entities.
type QueryKeyUpdate struct {
	config
	hooks    []Hook
	mutation *QueryKeyMutation
}

// Where appends a list predicates to the QueryKeyUpdate builder.
func (qku *QueryKeyUpdate) Where(ps ...predicate.QueryKey) *QueryKeyUpdate {
	qku.mutation.Where(ps...)
	return qku
}

// SetName sets the "name" field.
func (qku *QueryKeyUpdate) SetName(s string) *QueryKeyUpdate {
	qku.mutation.SetName(s)
	return qku
}

// SetNillableName sets the "name" field if the given value is not nil.
func (qku *QueryKeyUpdate) SetNillableName(s *string) *QueryKeyUpdate {
	if s != nil {
		qku.SetName(*s)
	}
	return qku
}

// SetType sets the "type" field.
func (qku *QueryKeyUpdate) SetType(s string) *QueryKeyUpdate {
	qku.mutation.SetType(s)
	return qku
}

// SetNillableType sets the "type" field if the given value is not nil.
func (qku *QueryKeyUpdate) SetNillableType(s *string) *QueryKeyUpdate {
	if s != nil {
		qku.SetType(*s)
	}
	return qku
}

// SetSource sets the "source" field.
func (qku *QueryKeyUpdate) SetSource(s string) *QueryKeyUpdate {
	qku.mutation.SetSource(s)
	return qku
}

// SetNillableSource sets the "source" field if the given value is not nil.
func (qku *QueryKeyUpdate) SetNillableSource(s *string) *QueryKeyUpdate {
	if s != nil {
		qku.SetSource(*s)
	}
	return qku
}

// SetValidDate sets the "valid_date" field.
func (qku *QueryKeyUpdate) SetValidDate(t time.Time) *QueryKeyUpdate {
	qku.mutation.SetValidDate(t)
	return qku
}

// SetNillableValidDate sets the "valid_date" field if the given value is not nil.
func (qku *QueryKeyUpdate) SetNillableValidDate(t *time.Time) *QueryKeyUpdate {
	if t != nil {
		qku.SetValidDate(*t)
	}
	return qku
}

// SetCreateTime sets the "create_time" field.
func (qku *QueryKeyUpdate) SetCreateTime(t time.Time) *QueryKeyUpdate {
	qku.mutation.SetCreateTime(t)
	return qku
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (qku *QueryKeyUpdate) SetNillableCreateTime(t *time.Time) *QueryKeyUpdate {
	if t != nil {
		qku.SetCreateTime(*t)
	}
	return qku
}

// SetUpdateTime sets the "update_time" field.
func (qku *QueryKeyUpdate) SetUpdateTime(t time.Time) *QueryKeyUpdate {
	qku.mutation.SetUpdateTime(t)
	return qku
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (qku *QueryKeyUpdate) SetNillableUpdateTime(t *time.Time) *QueryKeyUpdate {
	if t != nil {
		qku.SetUpdateTime(*t)
	}
	return qku
}

// AddValueIDs adds the "values" edge to the QueryValue entity by IDs.
func (qku *QueryKeyUpdate) AddValueIDs(ids ...int64) *QueryKeyUpdate {
	qku.mutation.AddValueIDs(ids...)
	return qku
}

// AddValues adds the "values" edges to the QueryValue entity.
func (qku *QueryKeyUpdate) AddValues(q ...*QueryValue) *QueryKeyUpdate {
	ids := make([]int64, len(q))
	for i := range q {
		ids[i] = q[i].ID
	}
	return qku.AddValueIDs(ids...)
}

// Mutation returns the QueryKeyMutation object of the builder.
func (qku *QueryKeyUpdate) Mutation() *QueryKeyMutation {
	return qku.mutation
}

// ClearValues clears all "values" edges to the QueryValue entity.
func (qku *QueryKeyUpdate) ClearValues() *QueryKeyUpdate {
	qku.mutation.ClearValues()
	return qku
}

// RemoveValueIDs removes the "values" edge to QueryValue entities by IDs.
func (qku *QueryKeyUpdate) RemoveValueIDs(ids ...int64) *QueryKeyUpdate {
	qku.mutation.RemoveValueIDs(ids...)
	return qku
}

// RemoveValues removes "values" edges to QueryValue entities.
func (qku *QueryKeyUpdate) RemoveValues(q ...*QueryValue) *QueryKeyUpdate {
	ids := make([]int64, len(q))
	for i := range q {
		ids[i] = q[i].ID
	}
	return qku.RemoveValueIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (qku *QueryKeyUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, qku.sqlSave, qku.mutation, qku.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (qku *QueryKeyUpdate) SaveX(ctx context.Context) int {
	affected, err := qku.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (qku *QueryKeyUpdate) Exec(ctx context.Context) error {
	_, err := qku.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (qku *QueryKeyUpdate) ExecX(ctx context.Context) {
	if err := qku.Exec(ctx); err != nil {
		panic(err)
	}
}

func (qku *QueryKeyUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(querykey.Table, querykey.Columns, sqlgraph.NewFieldSpec(querykey.FieldID, field.TypeInt64))
	if ps := qku.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := qku.mutation.Name(); ok {
		_spec.SetField(querykey.FieldName, field.TypeString, value)
	}
	if value, ok := qku.mutation.GetType(); ok {
		_spec.SetField(querykey.FieldType, field.TypeString, value)
	}
	if value, ok := qku.mutation.Source(); ok {
		_spec.SetField(querykey.FieldSource, field.TypeString, value)
	}
	if value, ok := qku.mutation.ValidDate(); ok {
		_spec.SetField(querykey.FieldValidDate, field.TypeTime, value)
	}
	if value, ok := qku.mutation.CreateTime(); ok {
		_spec.SetField(querykey.FieldCreateTime, field.TypeTime, value)
	}
	if value, ok := qku.mutation.UpdateTime(); ok {
		_spec.SetField(querykey.FieldUpdateTime, field.TypeTime, value)
	}
	if qku.mutation.ValuesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   querykey.ValuesTable,
			Columns: []string{querykey.ValuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(queryvalue.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := qku.mutation.RemovedValuesIDs(); len(nodes) > 0 && !qku.mutation.ValuesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   querykey.ValuesTable,
			Columns: []string{querykey.ValuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(queryvalue.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := qku.mutation.ValuesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   querykey.ValuesTable,
			Columns: []string{querykey.ValuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(queryvalue.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, qku.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{querykey.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	qku.mutation.done = true
	return n, nil
}

// QueryKeyUpdateOne is the builder for updating a single QueryKey entity.
type QueryKeyUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *QueryKeyMutation
}

// SetName sets the "name" field.
func (qkuo *QueryKeyUpdateOne) SetName(s string) *QueryKeyUpdateOne {
	qkuo.mutation.SetName(s)
	return qkuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (qkuo *QueryKeyUpdateOne) SetNillableName(s *string) *QueryKeyUpdateOne {
	if s != nil {
		qkuo.SetName(*s)
	}
	return qkuo
}

// SetType sets the "type" field.
func (qkuo *QueryKeyUpdateOne) SetType(s string) *QueryKeyUpdateOne {
	qkuo.mutation.SetType(s)
	return qkuo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (qkuo *QueryKeyUpdateOne) SetNillableType(s *string) *QueryKeyUpdateOne {
	if s != nil {
		qkuo.SetType(*s)
	}
	return qkuo
}

// SetSource sets the "source" field.
func (qkuo *QueryKeyUpdateOne) SetSource(s string) *QueryKeyUpdateOne {
	qkuo.mutation.SetSource(s)
	return qkuo
}

// SetNillableSource sets the "source" field if the given value is not nil.
func (qkuo *QueryKeyUpdateOne) SetNillableSource(s *string) *QueryKeyUpdateOne {
	if s != nil {
		qkuo.SetSource(*s)
	}
	return qkuo
}

// SetValidDate sets the "valid_date" field.
func (qkuo *QueryKeyUpdateOne) SetValidDate(t time.Time) *QueryKeyUpdateOne {
	qkuo.mutation.SetValidDate(t)
	return qkuo
}

// SetNillableValidDate sets the "valid_date" field if the given value is not nil.
func (qkuo *QueryKeyUpdateOne) SetNillableValidDate(t *time.Time) *QueryKeyUpdateOne {
	if t != nil {
		qkuo.SetValidDate(*t)
	}
	return qkuo
}

// SetCreateTime sets the "create_time" field.
func (qkuo *QueryKeyUpdateOne) SetCreateTime(t time.Time) *QueryKeyUpdateOne {
	qkuo.mutation.SetCreateTime(t)
	return qkuo
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (qkuo *QueryKeyUpdateOne) SetNillableCreateTime(t *time.Time) *QueryKeyUpdateOne {
	if t != nil {
		qkuo.SetCreateTime(*t)
	}
	return qkuo
}

// SetUpdateTime sets the "update_time" field.
func (qkuo *QueryKeyUpdateOne) SetUpdateTime(t time.Time) *QueryKeyUpdateOne {
	qkuo.mutation.SetUpdateTime(t)
	return qkuo
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (qkuo *QueryKeyUpdateOne) SetNillableUpdateTime(t *time.Time) *QueryKeyUpdateOne {
	if t != nil {
		qkuo.SetUpdateTime(*t)
	}
	return qkuo
}

// AddValueIDs adds the "values" edge to the QueryValue entity by IDs.
func (qkuo *QueryKeyUpdateOne) AddValueIDs(ids ...int64) *QueryKeyUpdateOne {
	qkuo.mutation.AddValueIDs(ids...)
	return qkuo
}

// AddValues adds the "values" edges to the QueryValue entity.
func (qkuo *QueryKeyUpdateOne) AddValues(q ...*QueryValue) *QueryKeyUpdateOne {
	ids := make([]int64, len(q))
	for i := range q {
		ids[i] = q[i].ID
	}
	return qkuo.AddValueIDs(ids...)
}

// Mutation returns the QueryKeyMutation object of the builder.
func (qkuo *QueryKeyUpdateOne) Mutation() *QueryKeyMutation {
	return qkuo.mutation
}

// ClearValues clears all "values" edges to the QueryValue entity.
func (qkuo *QueryKeyUpdateOne) ClearValues() *QueryKeyUpdateOne {
	qkuo.mutation.ClearValues()
	return qkuo
}

// RemoveValueIDs removes the "values" edge to QueryValue entities by IDs.
func (qkuo *QueryKeyUpdateOne) RemoveValueIDs(ids ...int64) *QueryKeyUpdateOne {
	qkuo.mutation.RemoveValueIDs(ids...)
	return qkuo
}

// RemoveValues removes "values" edges to QueryValue entities.
func (qkuo *QueryKeyUpdateOne) RemoveValues(q ...*QueryValue) *QueryKeyUpdateOne {
	ids := make([]int64, len(q))
	for i := range q {
		ids[i] = q[i].ID
	}
	return qkuo.RemoveValueIDs(ids...)
}

// Where appends a list predicates to the QueryKeyUpdate builder.
func (qkuo *QueryKeyUpdateOne) Where(ps ...predicate.QueryKey) *QueryKeyUpdateOne {
	qkuo.mutation.Where(ps...)
	return qkuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (qkuo *QueryKeyUpdateOne) Select(field string, fields ...string) *QueryKeyUpdateOne {
	qkuo.fields = append([]string{field}, fields...)
	return qkuo
}

// Save executes the query and returns the updated QueryKey entity.
func (qkuo *QueryKeyUpdateOne) Save(ctx context.Context) (*QueryKey, error) {
	return withHooks(ctx, qkuo.sqlSave, qkuo.mutation, qkuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (qkuo *QueryKeyUpdateOne) SaveX(ctx context.Context) *QueryKey {
	node, err := qkuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (qkuo *QueryKeyUpdateOne) Exec(ctx context.Context) error {
	_, err := qkuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (qkuo *QueryKeyUpdateOne) ExecX(ctx context.Context) {
	if err := qkuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (qkuo *QueryKeyUpdateOne) sqlSave(ctx context.Context) (_node *QueryKey, err error) {
	_spec := sqlgraph.NewUpdateSpec(querykey.Table, querykey.Columns, sqlgraph.NewFieldSpec(querykey.FieldID, field.TypeInt64))
	id, ok := qkuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "QueryKey.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := qkuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, querykey.FieldID)
		for _, f := range fields {
			if !querykey.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != querykey.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := qkuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := qkuo.mutation.Name(); ok {
		_spec.SetField(querykey.FieldName, field.TypeString, value)
	}
	if value, ok := qkuo.mutation.GetType(); ok {
		_spec.SetField(querykey.FieldType, field.TypeString, value)
	}
	if value, ok := qkuo.mutation.Source(); ok {
		_spec.SetField(querykey.FieldSource, field.TypeString, value)
	}
	if value, ok := qkuo.mutation.ValidDate(); ok {
		_spec.SetField(querykey.FieldValidDate, field.TypeTime, value)
	}
	if value, ok := qkuo.mutation.CreateTime(); ok {
		_spec.SetField(querykey.FieldCreateTime, field.TypeTime, value)
	}
	if value, ok := qkuo.mutation.UpdateTime(); ok {
		_spec.SetField(querykey.FieldUpdateTime, field.TypeTime, value)
	}
	if qkuo.mutation.ValuesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   querykey.ValuesTable,
			Columns: []string{querykey.ValuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(queryvalue.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := qkuo.mutation.RemovedValuesIDs(); len(nodes) > 0 && !qkuo.mutation.ValuesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   querykey.ValuesTable,
			Columns: []string{querykey.ValuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(queryvalue.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := qkuo.mutation.ValuesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   querykey.ValuesTable,
			Columns: []string{querykey.ValuesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(queryvalue.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &QueryKey{config: qkuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, qkuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{querykey.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	qkuo.mutation.done = true
	return _node, nil
}
