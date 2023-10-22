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
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/exceptioncategory"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/exceptiondefinition"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/predicate"
)

// ExceptionCategoryUpdate is the builder for updating ExceptionCategory entities.
type ExceptionCategoryUpdate struct {
	config
	hooks    []Hook
	mutation *ExceptionCategoryMutation
}

// Where appends a list predicates to the ExceptionCategoryUpdate builder.
func (ecu *ExceptionCategoryUpdate) Where(ps ...predicate.ExceptionCategory) *ExceptionCategoryUpdate {
	ecu.mutation.Where(ps...)
	return ecu
}

// SetName sets the "name" field.
func (ecu *ExceptionCategoryUpdate) SetName(s string) *ExceptionCategoryUpdate {
	ecu.mutation.SetName(s)
	return ecu
}

// SetIsValid sets the "is_valid" field.
func (ecu *ExceptionCategoryUpdate) SetIsValid(i int) *ExceptionCategoryUpdate {
	ecu.mutation.ResetIsValid()
	ecu.mutation.SetIsValid(i)
	return ecu
}

// SetNillableIsValid sets the "is_valid" field if the given value is not nil.
func (ecu *ExceptionCategoryUpdate) SetNillableIsValid(i *int) *ExceptionCategoryUpdate {
	if i != nil {
		ecu.SetIsValid(*i)
	}
	return ecu
}

// AddIsValid adds i to the "is_valid" field.
func (ecu *ExceptionCategoryUpdate) AddIsValid(i int) *ExceptionCategoryUpdate {
	ecu.mutation.AddIsValid(i)
	return ecu
}

// SetUpdateTime sets the "update_time" field.
func (ecu *ExceptionCategoryUpdate) SetUpdateTime(t time.Time) *ExceptionCategoryUpdate {
	ecu.mutation.SetUpdateTime(t)
	return ecu
}

// AddExceptionDefinitionIDs adds the "exception_definitions" edge to the ExceptionDefinition entity by IDs.
func (ecu *ExceptionCategoryUpdate) AddExceptionDefinitionIDs(ids ...int64) *ExceptionCategoryUpdate {
	ecu.mutation.AddExceptionDefinitionIDs(ids...)
	return ecu
}

// AddExceptionDefinitions adds the "exception_definitions" edges to the ExceptionDefinition entity.
func (ecu *ExceptionCategoryUpdate) AddExceptionDefinitions(e ...*ExceptionDefinition) *ExceptionCategoryUpdate {
	ids := make([]int64, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecu.AddExceptionDefinitionIDs(ids...)
}

// Mutation returns the ExceptionCategoryMutation object of the builder.
func (ecu *ExceptionCategoryUpdate) Mutation() *ExceptionCategoryMutation {
	return ecu.mutation
}

// ClearExceptionDefinitions clears all "exception_definitions" edges to the ExceptionDefinition entity.
func (ecu *ExceptionCategoryUpdate) ClearExceptionDefinitions() *ExceptionCategoryUpdate {
	ecu.mutation.ClearExceptionDefinitions()
	return ecu
}

// RemoveExceptionDefinitionIDs removes the "exception_definitions" edge to ExceptionDefinition entities by IDs.
func (ecu *ExceptionCategoryUpdate) RemoveExceptionDefinitionIDs(ids ...int64) *ExceptionCategoryUpdate {
	ecu.mutation.RemoveExceptionDefinitionIDs(ids...)
	return ecu
}

// RemoveExceptionDefinitions removes "exception_definitions" edges to ExceptionDefinition entities.
func (ecu *ExceptionCategoryUpdate) RemoveExceptionDefinitions(e ...*ExceptionDefinition) *ExceptionCategoryUpdate {
	ids := make([]int64, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecu.RemoveExceptionDefinitionIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ecu *ExceptionCategoryUpdate) Save(ctx context.Context) (int, error) {
	ecu.defaults()
	return withHooks(ctx, ecu.sqlSave, ecu.mutation, ecu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ecu *ExceptionCategoryUpdate) SaveX(ctx context.Context) int {
	affected, err := ecu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ecu *ExceptionCategoryUpdate) Exec(ctx context.Context) error {
	_, err := ecu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecu *ExceptionCategoryUpdate) ExecX(ctx context.Context) {
	if err := ecu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ecu *ExceptionCategoryUpdate) defaults() {
	if _, ok := ecu.mutation.UpdateTime(); !ok {
		v := exceptioncategory.UpdateDefaultUpdateTime()
		ecu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ecu *ExceptionCategoryUpdate) check() error {
	if v, ok := ecu.mutation.Name(); ok {
		if err := exceptioncategory.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "ExceptionCategory.name": %w`, err)}
		}
	}
	return nil
}

func (ecu *ExceptionCategoryUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ecu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(exceptioncategory.Table, exceptioncategory.Columns, sqlgraph.NewFieldSpec(exceptioncategory.FieldID, field.TypeInt64))
	if ps := ecu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ecu.mutation.Name(); ok {
		_spec.SetField(exceptioncategory.FieldName, field.TypeString, value)
	}
	if value, ok := ecu.mutation.IsValid(); ok {
		_spec.SetField(exceptioncategory.FieldIsValid, field.TypeInt, value)
	}
	if value, ok := ecu.mutation.AddedIsValid(); ok {
		_spec.AddField(exceptioncategory.FieldIsValid, field.TypeInt, value)
	}
	if value, ok := ecu.mutation.UpdateTime(); ok {
		_spec.SetField(exceptioncategory.FieldUpdateTime, field.TypeTime, value)
	}
	if ecu.mutation.ExceptionDefinitionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   exceptioncategory.ExceptionDefinitionsTable,
			Columns: []string{exceptioncategory.ExceptionDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exceptiondefinition.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ecu.mutation.RemovedExceptionDefinitionsIDs(); len(nodes) > 0 && !ecu.mutation.ExceptionDefinitionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   exceptioncategory.ExceptionDefinitionsTable,
			Columns: []string{exceptioncategory.ExceptionDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exceptiondefinition.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ecu.mutation.ExceptionDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   exceptioncategory.ExceptionDefinitionsTable,
			Columns: []string{exceptioncategory.ExceptionDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exceptiondefinition.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ecu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{exceptioncategory.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ecu.mutation.done = true
	return n, nil
}

// ExceptionCategoryUpdateOne is the builder for updating a single ExceptionCategory entity.
type ExceptionCategoryUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ExceptionCategoryMutation
}

// SetName sets the "name" field.
func (ecuo *ExceptionCategoryUpdateOne) SetName(s string) *ExceptionCategoryUpdateOne {
	ecuo.mutation.SetName(s)
	return ecuo
}

// SetIsValid sets the "is_valid" field.
func (ecuo *ExceptionCategoryUpdateOne) SetIsValid(i int) *ExceptionCategoryUpdateOne {
	ecuo.mutation.ResetIsValid()
	ecuo.mutation.SetIsValid(i)
	return ecuo
}

// SetNillableIsValid sets the "is_valid" field if the given value is not nil.
func (ecuo *ExceptionCategoryUpdateOne) SetNillableIsValid(i *int) *ExceptionCategoryUpdateOne {
	if i != nil {
		ecuo.SetIsValid(*i)
	}
	return ecuo
}

// AddIsValid adds i to the "is_valid" field.
func (ecuo *ExceptionCategoryUpdateOne) AddIsValid(i int) *ExceptionCategoryUpdateOne {
	ecuo.mutation.AddIsValid(i)
	return ecuo
}

// SetUpdateTime sets the "update_time" field.
func (ecuo *ExceptionCategoryUpdateOne) SetUpdateTime(t time.Time) *ExceptionCategoryUpdateOne {
	ecuo.mutation.SetUpdateTime(t)
	return ecuo
}

// AddExceptionDefinitionIDs adds the "exception_definitions" edge to the ExceptionDefinition entity by IDs.
func (ecuo *ExceptionCategoryUpdateOne) AddExceptionDefinitionIDs(ids ...int64) *ExceptionCategoryUpdateOne {
	ecuo.mutation.AddExceptionDefinitionIDs(ids...)
	return ecuo
}

// AddExceptionDefinitions adds the "exception_definitions" edges to the ExceptionDefinition entity.
func (ecuo *ExceptionCategoryUpdateOne) AddExceptionDefinitions(e ...*ExceptionDefinition) *ExceptionCategoryUpdateOne {
	ids := make([]int64, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecuo.AddExceptionDefinitionIDs(ids...)
}

// Mutation returns the ExceptionCategoryMutation object of the builder.
func (ecuo *ExceptionCategoryUpdateOne) Mutation() *ExceptionCategoryMutation {
	return ecuo.mutation
}

// ClearExceptionDefinitions clears all "exception_definitions" edges to the ExceptionDefinition entity.
func (ecuo *ExceptionCategoryUpdateOne) ClearExceptionDefinitions() *ExceptionCategoryUpdateOne {
	ecuo.mutation.ClearExceptionDefinitions()
	return ecuo
}

// RemoveExceptionDefinitionIDs removes the "exception_definitions" edge to ExceptionDefinition entities by IDs.
func (ecuo *ExceptionCategoryUpdateOne) RemoveExceptionDefinitionIDs(ids ...int64) *ExceptionCategoryUpdateOne {
	ecuo.mutation.RemoveExceptionDefinitionIDs(ids...)
	return ecuo
}

// RemoveExceptionDefinitions removes "exception_definitions" edges to ExceptionDefinition entities.
func (ecuo *ExceptionCategoryUpdateOne) RemoveExceptionDefinitions(e ...*ExceptionDefinition) *ExceptionCategoryUpdateOne {
	ids := make([]int64, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ecuo.RemoveExceptionDefinitionIDs(ids...)
}

// Where appends a list predicates to the ExceptionCategoryUpdate builder.
func (ecuo *ExceptionCategoryUpdateOne) Where(ps ...predicate.ExceptionCategory) *ExceptionCategoryUpdateOne {
	ecuo.mutation.Where(ps...)
	return ecuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ecuo *ExceptionCategoryUpdateOne) Select(field string, fields ...string) *ExceptionCategoryUpdateOne {
	ecuo.fields = append([]string{field}, fields...)
	return ecuo
}

// Save executes the query and returns the updated ExceptionCategory entity.
func (ecuo *ExceptionCategoryUpdateOne) Save(ctx context.Context) (*ExceptionCategory, error) {
	ecuo.defaults()
	return withHooks(ctx, ecuo.sqlSave, ecuo.mutation, ecuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ecuo *ExceptionCategoryUpdateOne) SaveX(ctx context.Context) *ExceptionCategory {
	node, err := ecuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ecuo *ExceptionCategoryUpdateOne) Exec(ctx context.Context) error {
	_, err := ecuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecuo *ExceptionCategoryUpdateOne) ExecX(ctx context.Context) {
	if err := ecuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ecuo *ExceptionCategoryUpdateOne) defaults() {
	if _, ok := ecuo.mutation.UpdateTime(); !ok {
		v := exceptioncategory.UpdateDefaultUpdateTime()
		ecuo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ecuo *ExceptionCategoryUpdateOne) check() error {
	if v, ok := ecuo.mutation.Name(); ok {
		if err := exceptioncategory.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "ExceptionCategory.name": %w`, err)}
		}
	}
	return nil
}

func (ecuo *ExceptionCategoryUpdateOne) sqlSave(ctx context.Context) (_node *ExceptionCategory, err error) {
	if err := ecuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(exceptioncategory.Table, exceptioncategory.Columns, sqlgraph.NewFieldSpec(exceptioncategory.FieldID, field.TypeInt64))
	id, ok := ecuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "ExceptionCategory.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ecuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, exceptioncategory.FieldID)
		for _, f := range fields {
			if !exceptioncategory.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != exceptioncategory.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ecuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ecuo.mutation.Name(); ok {
		_spec.SetField(exceptioncategory.FieldName, field.TypeString, value)
	}
	if value, ok := ecuo.mutation.IsValid(); ok {
		_spec.SetField(exceptioncategory.FieldIsValid, field.TypeInt, value)
	}
	if value, ok := ecuo.mutation.AddedIsValid(); ok {
		_spec.AddField(exceptioncategory.FieldIsValid, field.TypeInt, value)
	}
	if value, ok := ecuo.mutation.UpdateTime(); ok {
		_spec.SetField(exceptioncategory.FieldUpdateTime, field.TypeTime, value)
	}
	if ecuo.mutation.ExceptionDefinitionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   exceptioncategory.ExceptionDefinitionsTable,
			Columns: []string{exceptioncategory.ExceptionDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exceptiondefinition.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ecuo.mutation.RemovedExceptionDefinitionsIDs(); len(nodes) > 0 && !ecuo.mutation.ExceptionDefinitionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   exceptioncategory.ExceptionDefinitionsTable,
			Columns: []string{exceptioncategory.ExceptionDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exceptiondefinition.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ecuo.mutation.ExceptionDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   exceptioncategory.ExceptionDefinitionsTable,
			Columns: []string{exceptioncategory.ExceptionDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exceptiondefinition.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &ExceptionCategory{config: ecuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ecuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{exceptioncategory.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ecuo.mutation.done = true
	return _node, nil
}
