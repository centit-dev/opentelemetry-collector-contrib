// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/apptypeprocessor/ent/predicate"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/apptypeprocessor/ent/schema"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/apptypeprocessor/ent/softwaredefinition"
)

// SoftwareDefinitionUpdate is the builder for updating SoftwareDefinition entities.
type SoftwareDefinitionUpdate struct {
	config
	hooks    []Hook
	mutation *SoftwareDefinitionMutation
}

// Where appends a list predicates to the SoftwareDefinitionUpdate builder.
func (sdu *SoftwareDefinitionUpdate) Where(ps ...predicate.SoftwareDefinition) *SoftwareDefinitionUpdate {
	sdu.mutation.Where(ps...)
	return sdu
}

// SetName sets the "name" field.
func (sdu *SoftwareDefinitionUpdate) SetName(s string) *SoftwareDefinitionUpdate {
	sdu.mutation.SetName(s)
	return sdu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (sdu *SoftwareDefinitionUpdate) SetNillableName(s *string) *SoftwareDefinitionUpdate {
	if s != nil {
		sdu.SetName(*s)
	}
	return sdu
}

// SetType sets the "type" field.
func (sdu *SoftwareDefinitionUpdate) SetType(i int16) *SoftwareDefinitionUpdate {
	sdu.mutation.ResetType()
	sdu.mutation.SetType(i)
	return sdu
}

// SetNillableType sets the "type" field if the given value is not nil.
func (sdu *SoftwareDefinitionUpdate) SetNillableType(i *int16) *SoftwareDefinitionUpdate {
	if i != nil {
		sdu.SetType(*i)
	}
	return sdu
}

// AddType adds i to the "type" field.
func (sdu *SoftwareDefinitionUpdate) AddType(i int16) *SoftwareDefinitionUpdate {
	sdu.mutation.AddType(i)
	return sdu
}

// SetSpanConditions sets the "span_conditions" field.
func (sdu *SoftwareDefinitionUpdate) SetSpanConditions(sdc []schema.SoftwareDefinitionCondition) *SoftwareDefinitionUpdate {
	sdu.mutation.SetSpanConditions(sdc)
	return sdu
}

// AppendSpanConditions appends sdc to the "span_conditions" field.
func (sdu *SoftwareDefinitionUpdate) AppendSpanConditions(sdc []schema.SoftwareDefinitionCondition) *SoftwareDefinitionUpdate {
	sdu.mutation.AppendSpanConditions(sdc)
	return sdu
}

// ClearSpanConditions clears the value of the "span_conditions" field.
func (sdu *SoftwareDefinitionUpdate) ClearSpanConditions() *SoftwareDefinitionUpdate {
	sdu.mutation.ClearSpanConditions()
	return sdu
}

// SetIsValid sets the "is_valid" field.
func (sdu *SoftwareDefinitionUpdate) SetIsValid(i int) *SoftwareDefinitionUpdate {
	sdu.mutation.ResetIsValid()
	sdu.mutation.SetIsValid(i)
	return sdu
}

// SetNillableIsValid sets the "is_valid" field if the given value is not nil.
func (sdu *SoftwareDefinitionUpdate) SetNillableIsValid(i *int) *SoftwareDefinitionUpdate {
	if i != nil {
		sdu.SetIsValid(*i)
	}
	return sdu
}

// AddIsValid adds i to the "is_valid" field.
func (sdu *SoftwareDefinitionUpdate) AddIsValid(i int) *SoftwareDefinitionUpdate {
	sdu.mutation.AddIsValid(i)
	return sdu
}

// SetUpdateTime sets the "update_time" field.
func (sdu *SoftwareDefinitionUpdate) SetUpdateTime(t time.Time) *SoftwareDefinitionUpdate {
	sdu.mutation.SetUpdateTime(t)
	return sdu
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (sdu *SoftwareDefinitionUpdate) SetNillableUpdateTime(t *time.Time) *SoftwareDefinitionUpdate {
	if t != nil {
		sdu.SetUpdateTime(*t)
	}
	return sdu
}

// Mutation returns the SoftwareDefinitionMutation object of the builder.
func (sdu *SoftwareDefinitionUpdate) Mutation() *SoftwareDefinitionMutation {
	return sdu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (sdu *SoftwareDefinitionUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, sdu.sqlSave, sdu.mutation, sdu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (sdu *SoftwareDefinitionUpdate) SaveX(ctx context.Context) int {
	affected, err := sdu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (sdu *SoftwareDefinitionUpdate) Exec(ctx context.Context) error {
	_, err := sdu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sdu *SoftwareDefinitionUpdate) ExecX(ctx context.Context) {
	if err := sdu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sdu *SoftwareDefinitionUpdate) check() error {
	if v, ok := sdu.mutation.Name(); ok {
		if err := softwaredefinition.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "SoftwareDefinition.name": %w`, err)}
		}
	}
	return nil
}

func (sdu *SoftwareDefinitionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := sdu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(softwaredefinition.Table, softwaredefinition.Columns, sqlgraph.NewFieldSpec(softwaredefinition.FieldID, field.TypeInt64))
	if ps := sdu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := sdu.mutation.Name(); ok {
		_spec.SetField(softwaredefinition.FieldName, field.TypeString, value)
	}
	if value, ok := sdu.mutation.GetType(); ok {
		_spec.SetField(softwaredefinition.FieldType, field.TypeInt16, value)
	}
	if value, ok := sdu.mutation.AddedType(); ok {
		_spec.AddField(softwaredefinition.FieldType, field.TypeInt16, value)
	}
	if value, ok := sdu.mutation.SpanConditions(); ok {
		_spec.SetField(softwaredefinition.FieldSpanConditions, field.TypeJSON, value)
	}
	if value, ok := sdu.mutation.AppendedSpanConditions(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, softwaredefinition.FieldSpanConditions, value)
		})
	}
	if sdu.mutation.SpanConditionsCleared() {
		_spec.ClearField(softwaredefinition.FieldSpanConditions, field.TypeJSON)
	}
	if value, ok := sdu.mutation.IsValid(); ok {
		_spec.SetField(softwaredefinition.FieldIsValid, field.TypeInt, value)
	}
	if value, ok := sdu.mutation.AddedIsValid(); ok {
		_spec.AddField(softwaredefinition.FieldIsValid, field.TypeInt, value)
	}
	if value, ok := sdu.mutation.UpdateTime(); ok {
		_spec.SetField(softwaredefinition.FieldUpdateTime, field.TypeTime, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, sdu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{softwaredefinition.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	sdu.mutation.done = true
	return n, nil
}

// SoftwareDefinitionUpdateOne is the builder for updating a single SoftwareDefinition entity.
type SoftwareDefinitionUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SoftwareDefinitionMutation
}

// SetName sets the "name" field.
func (sduo *SoftwareDefinitionUpdateOne) SetName(s string) *SoftwareDefinitionUpdateOne {
	sduo.mutation.SetName(s)
	return sduo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (sduo *SoftwareDefinitionUpdateOne) SetNillableName(s *string) *SoftwareDefinitionUpdateOne {
	if s != nil {
		sduo.SetName(*s)
	}
	return sduo
}

// SetType sets the "type" field.
func (sduo *SoftwareDefinitionUpdateOne) SetType(i int16) *SoftwareDefinitionUpdateOne {
	sduo.mutation.ResetType()
	sduo.mutation.SetType(i)
	return sduo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (sduo *SoftwareDefinitionUpdateOne) SetNillableType(i *int16) *SoftwareDefinitionUpdateOne {
	if i != nil {
		sduo.SetType(*i)
	}
	return sduo
}

// AddType adds i to the "type" field.
func (sduo *SoftwareDefinitionUpdateOne) AddType(i int16) *SoftwareDefinitionUpdateOne {
	sduo.mutation.AddType(i)
	return sduo
}

// SetSpanConditions sets the "span_conditions" field.
func (sduo *SoftwareDefinitionUpdateOne) SetSpanConditions(sdc []schema.SoftwareDefinitionCondition) *SoftwareDefinitionUpdateOne {
	sduo.mutation.SetSpanConditions(sdc)
	return sduo
}

// AppendSpanConditions appends sdc to the "span_conditions" field.
func (sduo *SoftwareDefinitionUpdateOne) AppendSpanConditions(sdc []schema.SoftwareDefinitionCondition) *SoftwareDefinitionUpdateOne {
	sduo.mutation.AppendSpanConditions(sdc)
	return sduo
}

// ClearSpanConditions clears the value of the "span_conditions" field.
func (sduo *SoftwareDefinitionUpdateOne) ClearSpanConditions() *SoftwareDefinitionUpdateOne {
	sduo.mutation.ClearSpanConditions()
	return sduo
}

// SetIsValid sets the "is_valid" field.
func (sduo *SoftwareDefinitionUpdateOne) SetIsValid(i int) *SoftwareDefinitionUpdateOne {
	sduo.mutation.ResetIsValid()
	sduo.mutation.SetIsValid(i)
	return sduo
}

// SetNillableIsValid sets the "is_valid" field if the given value is not nil.
func (sduo *SoftwareDefinitionUpdateOne) SetNillableIsValid(i *int) *SoftwareDefinitionUpdateOne {
	if i != nil {
		sduo.SetIsValid(*i)
	}
	return sduo
}

// AddIsValid adds i to the "is_valid" field.
func (sduo *SoftwareDefinitionUpdateOne) AddIsValid(i int) *SoftwareDefinitionUpdateOne {
	sduo.mutation.AddIsValid(i)
	return sduo
}

// SetUpdateTime sets the "update_time" field.
func (sduo *SoftwareDefinitionUpdateOne) SetUpdateTime(t time.Time) *SoftwareDefinitionUpdateOne {
	sduo.mutation.SetUpdateTime(t)
	return sduo
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (sduo *SoftwareDefinitionUpdateOne) SetNillableUpdateTime(t *time.Time) *SoftwareDefinitionUpdateOne {
	if t != nil {
		sduo.SetUpdateTime(*t)
	}
	return sduo
}

// Mutation returns the SoftwareDefinitionMutation object of the builder.
func (sduo *SoftwareDefinitionUpdateOne) Mutation() *SoftwareDefinitionMutation {
	return sduo.mutation
}

// Where appends a list predicates to the SoftwareDefinitionUpdate builder.
func (sduo *SoftwareDefinitionUpdateOne) Where(ps ...predicate.SoftwareDefinition) *SoftwareDefinitionUpdateOne {
	sduo.mutation.Where(ps...)
	return sduo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (sduo *SoftwareDefinitionUpdateOne) Select(field string, fields ...string) *SoftwareDefinitionUpdateOne {
	sduo.fields = append([]string{field}, fields...)
	return sduo
}

// Save executes the query and returns the updated SoftwareDefinition entity.
func (sduo *SoftwareDefinitionUpdateOne) Save(ctx context.Context) (*SoftwareDefinition, error) {
	return withHooks(ctx, sduo.sqlSave, sduo.mutation, sduo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (sduo *SoftwareDefinitionUpdateOne) SaveX(ctx context.Context) *SoftwareDefinition {
	node, err := sduo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (sduo *SoftwareDefinitionUpdateOne) Exec(ctx context.Context) error {
	_, err := sduo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sduo *SoftwareDefinitionUpdateOne) ExecX(ctx context.Context) {
	if err := sduo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sduo *SoftwareDefinitionUpdateOne) check() error {
	if v, ok := sduo.mutation.Name(); ok {
		if err := softwaredefinition.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "SoftwareDefinition.name": %w`, err)}
		}
	}
	return nil
}

func (sduo *SoftwareDefinitionUpdateOne) sqlSave(ctx context.Context) (_node *SoftwareDefinition, err error) {
	if err := sduo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(softwaredefinition.Table, softwaredefinition.Columns, sqlgraph.NewFieldSpec(softwaredefinition.FieldID, field.TypeInt64))
	id, ok := sduo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "SoftwareDefinition.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := sduo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, softwaredefinition.FieldID)
		for _, f := range fields {
			if !softwaredefinition.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != softwaredefinition.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := sduo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := sduo.mutation.Name(); ok {
		_spec.SetField(softwaredefinition.FieldName, field.TypeString, value)
	}
	if value, ok := sduo.mutation.GetType(); ok {
		_spec.SetField(softwaredefinition.FieldType, field.TypeInt16, value)
	}
	if value, ok := sduo.mutation.AddedType(); ok {
		_spec.AddField(softwaredefinition.FieldType, field.TypeInt16, value)
	}
	if value, ok := sduo.mutation.SpanConditions(); ok {
		_spec.SetField(softwaredefinition.FieldSpanConditions, field.TypeJSON, value)
	}
	if value, ok := sduo.mutation.AppendedSpanConditions(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, softwaredefinition.FieldSpanConditions, value)
		})
	}
	if sduo.mutation.SpanConditionsCleared() {
		_spec.ClearField(softwaredefinition.FieldSpanConditions, field.TypeJSON)
	}
	if value, ok := sduo.mutation.IsValid(); ok {
		_spec.SetField(softwaredefinition.FieldIsValid, field.TypeInt, value)
	}
	if value, ok := sduo.mutation.AddedIsValid(); ok {
		_spec.AddField(softwaredefinition.FieldIsValid, field.TypeInt, value)
	}
	if value, ok := sduo.mutation.UpdateTime(); ok {
		_spec.SetField(softwaredefinition.FieldUpdateTime, field.TypeTime, value)
	}
	_node = &SoftwareDefinition{config: sduo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, sduo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{softwaredefinition.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	sduo.mutation.done = true
	return _node, nil
}