// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/exceptioncategory"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/predicate"
)

// ExceptionCategoryDelete is the builder for deleting a ExceptionCategory entity.
type ExceptionCategoryDelete struct {
	config
	hooks    []Hook
	mutation *ExceptionCategoryMutation
}

// Where appends a list predicates to the ExceptionCategoryDelete builder.
func (ecd *ExceptionCategoryDelete) Where(ps ...predicate.ExceptionCategory) *ExceptionCategoryDelete {
	ecd.mutation.Where(ps...)
	return ecd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ecd *ExceptionCategoryDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, ecd.sqlExec, ecd.mutation, ecd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (ecd *ExceptionCategoryDelete) ExecX(ctx context.Context) int {
	n, err := ecd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ecd *ExceptionCategoryDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(exceptioncategory.Table, sqlgraph.NewFieldSpec(exceptioncategory.FieldID, field.TypeInt64))
	if ps := ecd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ecd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	ecd.mutation.done = true
	return affected, err
}

// ExceptionCategoryDeleteOne is the builder for deleting a single ExceptionCategory entity.
type ExceptionCategoryDeleteOne struct {
	ecd *ExceptionCategoryDelete
}

// Where appends a list predicates to the ExceptionCategoryDelete builder.
func (ecdo *ExceptionCategoryDeleteOne) Where(ps ...predicate.ExceptionCategory) *ExceptionCategoryDeleteOne {
	ecdo.ecd.mutation.Where(ps...)
	return ecdo
}

// Exec executes the deletion query.
func (ecdo *ExceptionCategoryDeleteOne) Exec(ctx context.Context) error {
	n, err := ecdo.ecd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{exceptioncategory.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ecdo *ExceptionCategoryDeleteOne) ExecX(ctx context.Context) {
	if err := ecdo.Exec(ctx); err != nil {
		panic(err)
	}
}
