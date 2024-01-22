// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent/applicationstructure"
)

// ApplicationStructureCreate is the builder for creating a ApplicationStructure entity.
type ApplicationStructureCreate struct {
	config
	mutation *ApplicationStructureMutation
	hooks    []Hook
}

// SetParentCode sets the "parentCode" field.
func (asc *ApplicationStructureCreate) SetParentCode(s string) *ApplicationStructureCreate {
	asc.mutation.SetParentCode(s)
	return asc
}

// SetLevel sets the "level" field.
func (asc *ApplicationStructureCreate) SetLevel(i int) *ApplicationStructureCreate {
	asc.mutation.SetLevel(i)
	return asc
}

// SetValidDate sets the "valid_date" field.
func (asc *ApplicationStructureCreate) SetValidDate(t time.Time) *ApplicationStructureCreate {
	asc.mutation.SetValidDate(t)
	return asc
}

// SetCreateTime sets the "create_time" field.
func (asc *ApplicationStructureCreate) SetCreateTime(t time.Time) *ApplicationStructureCreate {
	asc.mutation.SetCreateTime(t)
	return asc
}

// SetUpdateTime sets the "update_time" field.
func (asc *ApplicationStructureCreate) SetUpdateTime(t time.Time) *ApplicationStructureCreate {
	asc.mutation.SetUpdateTime(t)
	return asc
}

// SetID sets the "id" field.
func (asc *ApplicationStructureCreate) SetID(s string) *ApplicationStructureCreate {
	asc.mutation.SetID(s)
	return asc
}

// Mutation returns the ApplicationStructureMutation object of the builder.
func (asc *ApplicationStructureCreate) Mutation() *ApplicationStructureMutation {
	return asc.mutation
}

// Save creates the ApplicationStructure in the database.
func (asc *ApplicationStructureCreate) Save(ctx context.Context) (*ApplicationStructure, error) {
	return withHooks(ctx, asc.sqlSave, asc.mutation, asc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (asc *ApplicationStructureCreate) SaveX(ctx context.Context) *ApplicationStructure {
	v, err := asc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (asc *ApplicationStructureCreate) Exec(ctx context.Context) error {
	_, err := asc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (asc *ApplicationStructureCreate) ExecX(ctx context.Context) {
	if err := asc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (asc *ApplicationStructureCreate) check() error {
	if _, ok := asc.mutation.ParentCode(); !ok {
		return &ValidationError{Name: "parentCode", err: errors.New(`ent: missing required field "ApplicationStructure.parentCode"`)}
	}
	if v, ok := asc.mutation.ParentCode(); ok {
		if err := applicationstructure.ParentCodeValidator(v); err != nil {
			return &ValidationError{Name: "parentCode", err: fmt.Errorf(`ent: validator failed for field "ApplicationStructure.parentCode": %w`, err)}
		}
	}
	if _, ok := asc.mutation.Level(); !ok {
		return &ValidationError{Name: "level", err: errors.New(`ent: missing required field "ApplicationStructure.level"`)}
	}
	if _, ok := asc.mutation.ValidDate(); !ok {
		return &ValidationError{Name: "valid_date", err: errors.New(`ent: missing required field "ApplicationStructure.valid_date"`)}
	}
	if _, ok := asc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "ApplicationStructure.create_time"`)}
	}
	if _, ok := asc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "ApplicationStructure.update_time"`)}
	}
	if v, ok := asc.mutation.ID(); ok {
		if err := applicationstructure.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "ApplicationStructure.id": %w`, err)}
		}
	}
	return nil
}

func (asc *ApplicationStructureCreate) sqlSave(ctx context.Context) (*ApplicationStructure, error) {
	if err := asc.check(); err != nil {
		return nil, err
	}
	_node, _spec := asc.createSpec()
	if err := sqlgraph.CreateNode(ctx, asc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(string); ok {
			_node.ID = id
		} else {
			return nil, fmt.Errorf("unexpected ApplicationStructure.ID type: %T", _spec.ID.Value)
		}
	}
	asc.mutation.id = &_node.ID
	asc.mutation.done = true
	return _node, nil
}

func (asc *ApplicationStructureCreate) createSpec() (*ApplicationStructure, *sqlgraph.CreateSpec) {
	var (
		_node = &ApplicationStructure{config: asc.config}
		_spec = sqlgraph.NewCreateSpec(applicationstructure.Table, sqlgraph.NewFieldSpec(applicationstructure.FieldID, field.TypeString))
	)
	if id, ok := asc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := asc.mutation.ParentCode(); ok {
		_spec.SetField(applicationstructure.FieldParentCode, field.TypeString, value)
		_node.ParentCode = value
	}
	if value, ok := asc.mutation.Level(); ok {
		_spec.SetField(applicationstructure.FieldLevel, field.TypeInt, value)
		_node.Level = value
	}
	if value, ok := asc.mutation.ValidDate(); ok {
		_spec.SetField(applicationstructure.FieldValidDate, field.TypeTime, value)
		_node.ValidDate = value
	}
	if value, ok := asc.mutation.CreateTime(); ok {
		_spec.SetField(applicationstructure.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := asc.mutation.UpdateTime(); ok {
		_spec.SetField(applicationstructure.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	return _node, _spec
}

// ApplicationStructureCreateBulk is the builder for creating many ApplicationStructure entities in bulk.
type ApplicationStructureCreateBulk struct {
	config
	err      error
	builders []*ApplicationStructureCreate
}

// Save creates the ApplicationStructure entities in the database.
func (ascb *ApplicationStructureCreateBulk) Save(ctx context.Context) ([]*ApplicationStructure, error) {
	if ascb.err != nil {
		return nil, ascb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ascb.builders))
	nodes := make([]*ApplicationStructure, len(ascb.builders))
	mutators := make([]Mutator, len(ascb.builders))
	for i := range ascb.builders {
		func(i int, root context.Context) {
			builder := ascb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ApplicationStructureMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ascb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ascb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ascb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ascb *ApplicationStructureCreateBulk) SaveX(ctx context.Context) []*ApplicationStructure {
	v, err := ascb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ascb *ApplicationStructureCreateBulk) Exec(ctx context.Context) error {
	_, err := ascb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ascb *ApplicationStructureCreateBulk) ExecX(ctx context.Context) {
	if err := ascb.Exec(ctx); err != nil {
		panic(err)
	}
}
