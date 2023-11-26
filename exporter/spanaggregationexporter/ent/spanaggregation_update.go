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
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent/predicate"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent/schema"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/spanaggregationexporter/ent/spanaggregation"
)

// SpanAggregationUpdate is the builder for updating SpanAggregation entities.
type SpanAggregationUpdate struct {
	config
	hooks    []Hook
	mutation *SpanAggregationMutation
}

// Where appends a list predicates to the SpanAggregationUpdate builder.
func (sau *SpanAggregationUpdate) Where(ps ...predicate.SpanAggregation) *SpanAggregationUpdate {
	sau.mutation.Where(ps...)
	return sau
}

// SetTimestamp sets the "Timestamp" field.
func (sau *SpanAggregationUpdate) SetTimestamp(t time.Time) *SpanAggregationUpdate {
	sau.mutation.SetTimestamp(t)
	return sau
}

// SetNillableTimestamp sets the "Timestamp" field if the given value is not nil.
func (sau *SpanAggregationUpdate) SetNillableTimestamp(t *time.Time) *SpanAggregationUpdate {
	if t != nil {
		sau.SetTimestamp(*t)
	}
	return sau
}

// SetTraceId sets the "TraceId" field.
func (sau *SpanAggregationUpdate) SetTraceId(s string) *SpanAggregationUpdate {
	sau.mutation.SetTraceId(s)
	return sau
}

// SetNillableTraceId sets the "TraceId" field if the given value is not nil.
func (sau *SpanAggregationUpdate) SetNillableTraceId(s *string) *SpanAggregationUpdate {
	if s != nil {
		sau.SetTraceId(*s)
	}
	return sau
}

// SetParentSpanId sets the "ParentSpanId" field.
func (sau *SpanAggregationUpdate) SetParentSpanId(s string) *SpanAggregationUpdate {
	sau.mutation.SetParentSpanId(s)
	return sau
}

// SetNillableParentSpanId sets the "ParentSpanId" field if the given value is not nil.
func (sau *SpanAggregationUpdate) SetNillableParentSpanId(s *string) *SpanAggregationUpdate {
	if s != nil {
		sau.SetParentSpanId(*s)
	}
	return sau
}

// SetPlatformName sets the "PlatformName" field.
func (sau *SpanAggregationUpdate) SetPlatformName(s string) *SpanAggregationUpdate {
	sau.mutation.SetPlatformName(s)
	return sau
}

// SetNillablePlatformName sets the "PlatformName" field if the given value is not nil.
func (sau *SpanAggregationUpdate) SetNillablePlatformName(s *string) *SpanAggregationUpdate {
	if s != nil {
		sau.SetPlatformName(*s)
	}
	return sau
}

// SetRootServiceName sets the "RootServiceName" field.
func (sau *SpanAggregationUpdate) SetRootServiceName(s string) *SpanAggregationUpdate {
	sau.mutation.SetRootServiceName(s)
	return sau
}

// SetNillableRootServiceName sets the "RootServiceName" field if the given value is not nil.
func (sau *SpanAggregationUpdate) SetNillableRootServiceName(s *string) *SpanAggregationUpdate {
	if s != nil {
		sau.SetRootServiceName(*s)
	}
	return sau
}

// ClearRootServiceName clears the value of the "RootServiceName" field.
func (sau *SpanAggregationUpdate) ClearRootServiceName() *SpanAggregationUpdate {
	sau.mutation.ClearRootServiceName()
	return sau
}

// SetRootSpanName sets the "RootSpanName" field.
func (sau *SpanAggregationUpdate) SetRootSpanName(s string) *SpanAggregationUpdate {
	sau.mutation.SetRootSpanName(s)
	return sau
}

// SetNillableRootSpanName sets the "RootSpanName" field if the given value is not nil.
func (sau *SpanAggregationUpdate) SetNillableRootSpanName(s *string) *SpanAggregationUpdate {
	if s != nil {
		sau.SetRootSpanName(*s)
	}
	return sau
}

// ClearRootSpanName clears the value of the "RootSpanName" field.
func (sau *SpanAggregationUpdate) ClearRootSpanName() *SpanAggregationUpdate {
	sau.mutation.ClearRootSpanName()
	return sau
}

// SetServiceName sets the "ServiceName" field.
func (sau *SpanAggregationUpdate) SetServiceName(s string) *SpanAggregationUpdate {
	sau.mutation.SetServiceName(s)
	return sau
}

// SetNillableServiceName sets the "ServiceName" field if the given value is not nil.
func (sau *SpanAggregationUpdate) SetNillableServiceName(s *string) *SpanAggregationUpdate {
	if s != nil {
		sau.SetServiceName(*s)
	}
	return sau
}

// SetSpanName sets the "SpanName" field.
func (sau *SpanAggregationUpdate) SetSpanName(s string) *SpanAggregationUpdate {
	sau.mutation.SetSpanName(s)
	return sau
}

// SetNillableSpanName sets the "SpanName" field if the given value is not nil.
func (sau *SpanAggregationUpdate) SetNillableSpanName(s *string) *SpanAggregationUpdate {
	if s != nil {
		sau.SetSpanName(*s)
	}
	return sau
}

// SetResourceAttributes sets the "ResourceAttributes" field.
func (sau *SpanAggregationUpdate) SetResourceAttributes(s *schema.Attributes) *SpanAggregationUpdate {
	sau.mutation.SetResourceAttributes(s)
	return sau
}

// SetSpanAttributes sets the "SpanAttributes" field.
func (sau *SpanAggregationUpdate) SetSpanAttributes(s *schema.Attributes) *SpanAggregationUpdate {
	sau.mutation.SetSpanAttributes(s)
	return sau
}

// SetDuration sets the "Duration" field.
func (sau *SpanAggregationUpdate) SetDuration(i int64) *SpanAggregationUpdate {
	sau.mutation.ResetDuration()
	sau.mutation.SetDuration(i)
	return sau
}

// SetNillableDuration sets the "Duration" field if the given value is not nil.
func (sau *SpanAggregationUpdate) SetNillableDuration(i *int64) *SpanAggregationUpdate {
	if i != nil {
		sau.SetDuration(*i)
	}
	return sau
}

// AddDuration adds i to the "Duration" field.
func (sau *SpanAggregationUpdate) AddDuration(i int64) *SpanAggregationUpdate {
	sau.mutation.AddDuration(i)
	return sau
}

// SetGap sets the "Gap" field.
func (sau *SpanAggregationUpdate) SetGap(i int64) *SpanAggregationUpdate {
	sau.mutation.ResetGap()
	sau.mutation.SetGap(i)
	return sau
}

// SetNillableGap sets the "Gap" field if the given value is not nil.
func (sau *SpanAggregationUpdate) SetNillableGap(i *int64) *SpanAggregationUpdate {
	if i != nil {
		sau.SetGap(*i)
	}
	return sau
}

// AddGap adds i to the "Gap" field.
func (sau *SpanAggregationUpdate) AddGap(i int64) *SpanAggregationUpdate {
	sau.mutation.AddGap(i)
	return sau
}

// SetSelfDuration sets the "SelfDuration" field.
func (sau *SpanAggregationUpdate) SetSelfDuration(i int64) *SpanAggregationUpdate {
	sau.mutation.ResetSelfDuration()
	sau.mutation.SetSelfDuration(i)
	return sau
}

// SetNillableSelfDuration sets the "SelfDuration" field if the given value is not nil.
func (sau *SpanAggregationUpdate) SetNillableSelfDuration(i *int64) *SpanAggregationUpdate {
	if i != nil {
		sau.SetSelfDuration(*i)
	}
	return sau
}

// AddSelfDuration adds i to the "SelfDuration" field.
func (sau *SpanAggregationUpdate) AddSelfDuration(i int64) *SpanAggregationUpdate {
	sau.mutation.AddSelfDuration(i)
	return sau
}

// Mutation returns the SpanAggregationMutation object of the builder.
func (sau *SpanAggregationUpdate) Mutation() *SpanAggregationMutation {
	return sau.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (sau *SpanAggregationUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, sau.sqlSave, sau.mutation, sau.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (sau *SpanAggregationUpdate) SaveX(ctx context.Context) int {
	affected, err := sau.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (sau *SpanAggregationUpdate) Exec(ctx context.Context) error {
	_, err := sau.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sau *SpanAggregationUpdate) ExecX(ctx context.Context) {
	if err := sau.Exec(ctx); err != nil {
		panic(err)
	}
}

func (sau *SpanAggregationUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(spanaggregation.Table, spanaggregation.Columns, sqlgraph.NewFieldSpec(spanaggregation.FieldID, field.TypeString))
	if ps := sau.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := sau.mutation.Timestamp(); ok {
		_spec.SetField(spanaggregation.FieldTimestamp, field.TypeTime, value)
	}
	if value, ok := sau.mutation.TraceId(); ok {
		_spec.SetField(spanaggregation.FieldTraceId, field.TypeString, value)
	}
	if value, ok := sau.mutation.ParentSpanId(); ok {
		_spec.SetField(spanaggregation.FieldParentSpanId, field.TypeString, value)
	}
	if value, ok := sau.mutation.PlatformName(); ok {
		_spec.SetField(spanaggregation.FieldPlatformName, field.TypeString, value)
	}
	if value, ok := sau.mutation.RootServiceName(); ok {
		_spec.SetField(spanaggregation.FieldRootServiceName, field.TypeString, value)
	}
	if sau.mutation.RootServiceNameCleared() {
		_spec.ClearField(spanaggregation.FieldRootServiceName, field.TypeString)
	}
	if value, ok := sau.mutation.RootSpanName(); ok {
		_spec.SetField(spanaggregation.FieldRootSpanName, field.TypeString, value)
	}
	if sau.mutation.RootSpanNameCleared() {
		_spec.ClearField(spanaggregation.FieldRootSpanName, field.TypeString)
	}
	if value, ok := sau.mutation.ServiceName(); ok {
		_spec.SetField(spanaggregation.FieldServiceName, field.TypeString, value)
	}
	if value, ok := sau.mutation.SpanName(); ok {
		_spec.SetField(spanaggregation.FieldSpanName, field.TypeString, value)
	}
	if value, ok := sau.mutation.ResourceAttributes(); ok {
		_spec.SetField(spanaggregation.FieldResourceAttributes, field.TypeOther, value)
	}
	if value, ok := sau.mutation.SpanAttributes(); ok {
		_spec.SetField(spanaggregation.FieldSpanAttributes, field.TypeOther, value)
	}
	if value, ok := sau.mutation.Duration(); ok {
		_spec.SetField(spanaggregation.FieldDuration, field.TypeInt64, value)
	}
	if value, ok := sau.mutation.AddedDuration(); ok {
		_spec.AddField(spanaggregation.FieldDuration, field.TypeInt64, value)
	}
	if value, ok := sau.mutation.Gap(); ok {
		_spec.SetField(spanaggregation.FieldGap, field.TypeInt64, value)
	}
	if value, ok := sau.mutation.AddedGap(); ok {
		_spec.AddField(spanaggregation.FieldGap, field.TypeInt64, value)
	}
	if value, ok := sau.mutation.SelfDuration(); ok {
		_spec.SetField(spanaggregation.FieldSelfDuration, field.TypeInt64, value)
	}
	if value, ok := sau.mutation.AddedSelfDuration(); ok {
		_spec.AddField(spanaggregation.FieldSelfDuration, field.TypeInt64, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, sau.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{spanaggregation.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	sau.mutation.done = true
	return n, nil
}

// SpanAggregationUpdateOne is the builder for updating a single SpanAggregation entity.
type SpanAggregationUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SpanAggregationMutation
}

// SetTimestamp sets the "Timestamp" field.
func (sauo *SpanAggregationUpdateOne) SetTimestamp(t time.Time) *SpanAggregationUpdateOne {
	sauo.mutation.SetTimestamp(t)
	return sauo
}

// SetNillableTimestamp sets the "Timestamp" field if the given value is not nil.
func (sauo *SpanAggregationUpdateOne) SetNillableTimestamp(t *time.Time) *SpanAggregationUpdateOne {
	if t != nil {
		sauo.SetTimestamp(*t)
	}
	return sauo
}

// SetTraceId sets the "TraceId" field.
func (sauo *SpanAggregationUpdateOne) SetTraceId(s string) *SpanAggregationUpdateOne {
	sauo.mutation.SetTraceId(s)
	return sauo
}

// SetNillableTraceId sets the "TraceId" field if the given value is not nil.
func (sauo *SpanAggregationUpdateOne) SetNillableTraceId(s *string) *SpanAggregationUpdateOne {
	if s != nil {
		sauo.SetTraceId(*s)
	}
	return sauo
}

// SetParentSpanId sets the "ParentSpanId" field.
func (sauo *SpanAggregationUpdateOne) SetParentSpanId(s string) *SpanAggregationUpdateOne {
	sauo.mutation.SetParentSpanId(s)
	return sauo
}

// SetNillableParentSpanId sets the "ParentSpanId" field if the given value is not nil.
func (sauo *SpanAggregationUpdateOne) SetNillableParentSpanId(s *string) *SpanAggregationUpdateOne {
	if s != nil {
		sauo.SetParentSpanId(*s)
	}
	return sauo
}

// SetPlatformName sets the "PlatformName" field.
func (sauo *SpanAggregationUpdateOne) SetPlatformName(s string) *SpanAggregationUpdateOne {
	sauo.mutation.SetPlatformName(s)
	return sauo
}

// SetNillablePlatformName sets the "PlatformName" field if the given value is not nil.
func (sauo *SpanAggregationUpdateOne) SetNillablePlatformName(s *string) *SpanAggregationUpdateOne {
	if s != nil {
		sauo.SetPlatformName(*s)
	}
	return sauo
}

// SetRootServiceName sets the "RootServiceName" field.
func (sauo *SpanAggregationUpdateOne) SetRootServiceName(s string) *SpanAggregationUpdateOne {
	sauo.mutation.SetRootServiceName(s)
	return sauo
}

// SetNillableRootServiceName sets the "RootServiceName" field if the given value is not nil.
func (sauo *SpanAggregationUpdateOne) SetNillableRootServiceName(s *string) *SpanAggregationUpdateOne {
	if s != nil {
		sauo.SetRootServiceName(*s)
	}
	return sauo
}

// ClearRootServiceName clears the value of the "RootServiceName" field.
func (sauo *SpanAggregationUpdateOne) ClearRootServiceName() *SpanAggregationUpdateOne {
	sauo.mutation.ClearRootServiceName()
	return sauo
}

// SetRootSpanName sets the "RootSpanName" field.
func (sauo *SpanAggregationUpdateOne) SetRootSpanName(s string) *SpanAggregationUpdateOne {
	sauo.mutation.SetRootSpanName(s)
	return sauo
}

// SetNillableRootSpanName sets the "RootSpanName" field if the given value is not nil.
func (sauo *SpanAggregationUpdateOne) SetNillableRootSpanName(s *string) *SpanAggregationUpdateOne {
	if s != nil {
		sauo.SetRootSpanName(*s)
	}
	return sauo
}

// ClearRootSpanName clears the value of the "RootSpanName" field.
func (sauo *SpanAggregationUpdateOne) ClearRootSpanName() *SpanAggregationUpdateOne {
	sauo.mutation.ClearRootSpanName()
	return sauo
}

// SetServiceName sets the "ServiceName" field.
func (sauo *SpanAggregationUpdateOne) SetServiceName(s string) *SpanAggregationUpdateOne {
	sauo.mutation.SetServiceName(s)
	return sauo
}

// SetNillableServiceName sets the "ServiceName" field if the given value is not nil.
func (sauo *SpanAggregationUpdateOne) SetNillableServiceName(s *string) *SpanAggregationUpdateOne {
	if s != nil {
		sauo.SetServiceName(*s)
	}
	return sauo
}

// SetSpanName sets the "SpanName" field.
func (sauo *SpanAggregationUpdateOne) SetSpanName(s string) *SpanAggregationUpdateOne {
	sauo.mutation.SetSpanName(s)
	return sauo
}

// SetNillableSpanName sets the "SpanName" field if the given value is not nil.
func (sauo *SpanAggregationUpdateOne) SetNillableSpanName(s *string) *SpanAggregationUpdateOne {
	if s != nil {
		sauo.SetSpanName(*s)
	}
	return sauo
}

// SetResourceAttributes sets the "ResourceAttributes" field.
func (sauo *SpanAggregationUpdateOne) SetResourceAttributes(s *schema.Attributes) *SpanAggregationUpdateOne {
	sauo.mutation.SetResourceAttributes(s)
	return sauo
}

// SetSpanAttributes sets the "SpanAttributes" field.
func (sauo *SpanAggregationUpdateOne) SetSpanAttributes(s *schema.Attributes) *SpanAggregationUpdateOne {
	sauo.mutation.SetSpanAttributes(s)
	return sauo
}

// SetDuration sets the "Duration" field.
func (sauo *SpanAggregationUpdateOne) SetDuration(i int64) *SpanAggregationUpdateOne {
	sauo.mutation.ResetDuration()
	sauo.mutation.SetDuration(i)
	return sauo
}

// SetNillableDuration sets the "Duration" field if the given value is not nil.
func (sauo *SpanAggregationUpdateOne) SetNillableDuration(i *int64) *SpanAggregationUpdateOne {
	if i != nil {
		sauo.SetDuration(*i)
	}
	return sauo
}

// AddDuration adds i to the "Duration" field.
func (sauo *SpanAggregationUpdateOne) AddDuration(i int64) *SpanAggregationUpdateOne {
	sauo.mutation.AddDuration(i)
	return sauo
}

// SetGap sets the "Gap" field.
func (sauo *SpanAggregationUpdateOne) SetGap(i int64) *SpanAggregationUpdateOne {
	sauo.mutation.ResetGap()
	sauo.mutation.SetGap(i)
	return sauo
}

// SetNillableGap sets the "Gap" field if the given value is not nil.
func (sauo *SpanAggregationUpdateOne) SetNillableGap(i *int64) *SpanAggregationUpdateOne {
	if i != nil {
		sauo.SetGap(*i)
	}
	return sauo
}

// AddGap adds i to the "Gap" field.
func (sauo *SpanAggregationUpdateOne) AddGap(i int64) *SpanAggregationUpdateOne {
	sauo.mutation.AddGap(i)
	return sauo
}

// SetSelfDuration sets the "SelfDuration" field.
func (sauo *SpanAggregationUpdateOne) SetSelfDuration(i int64) *SpanAggregationUpdateOne {
	sauo.mutation.ResetSelfDuration()
	sauo.mutation.SetSelfDuration(i)
	return sauo
}

// SetNillableSelfDuration sets the "SelfDuration" field if the given value is not nil.
func (sauo *SpanAggregationUpdateOne) SetNillableSelfDuration(i *int64) *SpanAggregationUpdateOne {
	if i != nil {
		sauo.SetSelfDuration(*i)
	}
	return sauo
}

// AddSelfDuration adds i to the "SelfDuration" field.
func (sauo *SpanAggregationUpdateOne) AddSelfDuration(i int64) *SpanAggregationUpdateOne {
	sauo.mutation.AddSelfDuration(i)
	return sauo
}

// Mutation returns the SpanAggregationMutation object of the builder.
func (sauo *SpanAggregationUpdateOne) Mutation() *SpanAggregationMutation {
	return sauo.mutation
}

// Where appends a list predicates to the SpanAggregationUpdate builder.
func (sauo *SpanAggregationUpdateOne) Where(ps ...predicate.SpanAggregation) *SpanAggregationUpdateOne {
	sauo.mutation.Where(ps...)
	return sauo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (sauo *SpanAggregationUpdateOne) Select(field string, fields ...string) *SpanAggregationUpdateOne {
	sauo.fields = append([]string{field}, fields...)
	return sauo
}

// Save executes the query and returns the updated SpanAggregation entity.
func (sauo *SpanAggregationUpdateOne) Save(ctx context.Context) (*SpanAggregation, error) {
	return withHooks(ctx, sauo.sqlSave, sauo.mutation, sauo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (sauo *SpanAggregationUpdateOne) SaveX(ctx context.Context) *SpanAggregation {
	node, err := sauo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (sauo *SpanAggregationUpdateOne) Exec(ctx context.Context) error {
	_, err := sauo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sauo *SpanAggregationUpdateOne) ExecX(ctx context.Context) {
	if err := sauo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (sauo *SpanAggregationUpdateOne) sqlSave(ctx context.Context) (_node *SpanAggregation, err error) {
	_spec := sqlgraph.NewUpdateSpec(spanaggregation.Table, spanaggregation.Columns, sqlgraph.NewFieldSpec(spanaggregation.FieldID, field.TypeString))
	id, ok := sauo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "SpanAggregation.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := sauo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, spanaggregation.FieldID)
		for _, f := range fields {
			if !spanaggregation.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != spanaggregation.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := sauo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := sauo.mutation.Timestamp(); ok {
		_spec.SetField(spanaggregation.FieldTimestamp, field.TypeTime, value)
	}
	if value, ok := sauo.mutation.TraceId(); ok {
		_spec.SetField(spanaggregation.FieldTraceId, field.TypeString, value)
	}
	if value, ok := sauo.mutation.ParentSpanId(); ok {
		_spec.SetField(spanaggregation.FieldParentSpanId, field.TypeString, value)
	}
	if value, ok := sauo.mutation.PlatformName(); ok {
		_spec.SetField(spanaggregation.FieldPlatformName, field.TypeString, value)
	}
	if value, ok := sauo.mutation.RootServiceName(); ok {
		_spec.SetField(spanaggregation.FieldRootServiceName, field.TypeString, value)
	}
	if sauo.mutation.RootServiceNameCleared() {
		_spec.ClearField(spanaggregation.FieldRootServiceName, field.TypeString)
	}
	if value, ok := sauo.mutation.RootSpanName(); ok {
		_spec.SetField(spanaggregation.FieldRootSpanName, field.TypeString, value)
	}
	if sauo.mutation.RootSpanNameCleared() {
		_spec.ClearField(spanaggregation.FieldRootSpanName, field.TypeString)
	}
	if value, ok := sauo.mutation.ServiceName(); ok {
		_spec.SetField(spanaggregation.FieldServiceName, field.TypeString, value)
	}
	if value, ok := sauo.mutation.SpanName(); ok {
		_spec.SetField(spanaggregation.FieldSpanName, field.TypeString, value)
	}
	if value, ok := sauo.mutation.ResourceAttributes(); ok {
		_spec.SetField(spanaggregation.FieldResourceAttributes, field.TypeOther, value)
	}
	if value, ok := sauo.mutation.SpanAttributes(); ok {
		_spec.SetField(spanaggregation.FieldSpanAttributes, field.TypeOther, value)
	}
	if value, ok := sauo.mutation.Duration(); ok {
		_spec.SetField(spanaggregation.FieldDuration, field.TypeInt64, value)
	}
	if value, ok := sauo.mutation.AddedDuration(); ok {
		_spec.AddField(spanaggregation.FieldDuration, field.TypeInt64, value)
	}
	if value, ok := sauo.mutation.Gap(); ok {
		_spec.SetField(spanaggregation.FieldGap, field.TypeInt64, value)
	}
	if value, ok := sauo.mutation.AddedGap(); ok {
		_spec.AddField(spanaggregation.FieldGap, field.TypeInt64, value)
	}
	if value, ok := sauo.mutation.SelfDuration(); ok {
		_spec.SetField(spanaggregation.FieldSelfDuration, field.TypeInt64, value)
	}
	if value, ok := sauo.mutation.AddedSelfDuration(); ok {
		_spec.AddField(spanaggregation.FieldSelfDuration, field.TypeInt64, value)
	}
	_node = &SpanAggregation{config: sauo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, sauo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{spanaggregation.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	sauo.mutation.done = true
	return _node, nil
}
