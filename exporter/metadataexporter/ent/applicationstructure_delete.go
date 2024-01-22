// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent/applicationstructure"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent/predicate"
)

// ApplicationStructureDelete is the builder for deleting a ApplicationStructure entity.
type ApplicationStructureDelete struct {
	config
	hooks    []Hook
	mutation *ApplicationStructureMutation
}

// Where appends a list predicates to the ApplicationStructureDelete builder.
func (asd *ApplicationStructureDelete) Where(ps ...predicate.ApplicationStructure) *ApplicationStructureDelete {
	asd.mutation.Where(ps...)
	return asd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (asd *ApplicationStructureDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, asd.sqlExec, asd.mutation, asd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (asd *ApplicationStructureDelete) ExecX(ctx context.Context) int {
	n, err := asd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (asd *ApplicationStructureDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(applicationstructure.Table, sqlgraph.NewFieldSpec(applicationstructure.FieldID, field.TypeString))
	if ps := asd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, asd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	asd.mutation.done = true
	return affected, err
}

// ApplicationStructureDeleteOne is the builder for deleting a single ApplicationStructure entity.
type ApplicationStructureDeleteOne struct {
	asd *ApplicationStructureDelete
}

// Where appends a list predicates to the ApplicationStructureDelete builder.
func (asdo *ApplicationStructureDeleteOne) Where(ps ...predicate.ApplicationStructure) *ApplicationStructureDeleteOne {
	asdo.asd.mutation.Where(ps...)
	return asdo
}

// Exec executes the deletion query.
func (asdo *ApplicationStructureDeleteOne) Exec(ctx context.Context) error {
	n, err := asdo.asd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{applicationstructure.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (asdo *ApplicationStructureDeleteOne) ExecX(ctx context.Context) {
	if err := asdo.Exec(ctx); err != nil {
		panic(err)
	}
}
