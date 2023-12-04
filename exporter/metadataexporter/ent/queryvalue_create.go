// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent/querykey"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent/queryvalue"
)

// QueryValueCreate is the builder for creating a QueryValue entity.
type QueryValueCreate struct {
	config
	mutation *QueryValueMutation
	hooks    []Hook
}

// SetKeyID sets the "key_id" field.
func (qvc *QueryValueCreate) SetKeyID(i int64) *QueryValueCreate {
	qvc.mutation.SetKeyID(i)
	return qvc
}

// SetValue sets the "value" field.
func (qvc *QueryValueCreate) SetValue(s string) *QueryValueCreate {
	qvc.mutation.SetValue(s)
	return qvc
}

// SetValidDate sets the "valid_date" field.
func (qvc *QueryValueCreate) SetValidDate(t time.Time) *QueryValueCreate {
	qvc.mutation.SetValidDate(t)
	return qvc
}

// SetCreateTime sets the "create_time" field.
func (qvc *QueryValueCreate) SetCreateTime(t time.Time) *QueryValueCreate {
	qvc.mutation.SetCreateTime(t)
	return qvc
}

// SetUpdateTime sets the "update_time" field.
func (qvc *QueryValueCreate) SetUpdateTime(t time.Time) *QueryValueCreate {
	qvc.mutation.SetUpdateTime(t)
	return qvc
}

// SetID sets the "id" field.
func (qvc *QueryValueCreate) SetID(i int64) *QueryValueCreate {
	qvc.mutation.SetID(i)
	return qvc
}

// SetKey sets the "key" edge to the QueryKey entity.
func (qvc *QueryValueCreate) SetKey(q *QueryKey) *QueryValueCreate {
	return qvc.SetKeyID(q.ID)
}

// Mutation returns the QueryValueMutation object of the builder.
func (qvc *QueryValueCreate) Mutation() *QueryValueMutation {
	return qvc.mutation
}

// Save creates the QueryValue in the database.
func (qvc *QueryValueCreate) Save(ctx context.Context) (*QueryValue, error) {
	return withHooks(ctx, qvc.sqlSave, qvc.mutation, qvc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (qvc *QueryValueCreate) SaveX(ctx context.Context) *QueryValue {
	v, err := qvc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (qvc *QueryValueCreate) Exec(ctx context.Context) error {
	_, err := qvc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (qvc *QueryValueCreate) ExecX(ctx context.Context) {
	if err := qvc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (qvc *QueryValueCreate) check() error {
	if _, ok := qvc.mutation.KeyID(); !ok {
		return &ValidationError{Name: "key_id", err: errors.New(`ent: missing required field "QueryValue.key_id"`)}
	}
	if _, ok := qvc.mutation.Value(); !ok {
		return &ValidationError{Name: "value", err: errors.New(`ent: missing required field "QueryValue.value"`)}
	}
	if _, ok := qvc.mutation.ValidDate(); !ok {
		return &ValidationError{Name: "valid_date", err: errors.New(`ent: missing required field "QueryValue.valid_date"`)}
	}
	if _, ok := qvc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "QueryValue.create_time"`)}
	}
	if _, ok := qvc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "QueryValue.update_time"`)}
	}
	if _, ok := qvc.mutation.KeyID(); !ok {
		return &ValidationError{Name: "key", err: errors.New(`ent: missing required edge "QueryValue.key"`)}
	}
	return nil
}

func (qvc *QueryValueCreate) sqlSave(ctx context.Context) (*QueryValue, error) {
	if err := qvc.check(); err != nil {
		return nil, err
	}
	_node, _spec := qvc.createSpec()
	if err := sqlgraph.CreateNode(ctx, qvc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int64(id)
	}
	qvc.mutation.id = &_node.ID
	qvc.mutation.done = true
	return _node, nil
}

func (qvc *QueryValueCreate) createSpec() (*QueryValue, *sqlgraph.CreateSpec) {
	var (
		_node = &QueryValue{config: qvc.config}
		_spec = sqlgraph.NewCreateSpec(queryvalue.Table, sqlgraph.NewFieldSpec(queryvalue.FieldID, field.TypeInt64))
	)
	if id, ok := qvc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := qvc.mutation.Value(); ok {
		_spec.SetField(queryvalue.FieldValue, field.TypeString, value)
		_node.Value = value
	}
	if value, ok := qvc.mutation.ValidDate(); ok {
		_spec.SetField(queryvalue.FieldValidDate, field.TypeTime, value)
		_node.ValidDate = value
	}
	if value, ok := qvc.mutation.CreateTime(); ok {
		_spec.SetField(queryvalue.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := qvc.mutation.UpdateTime(); ok {
		_spec.SetField(queryvalue.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if nodes := qvc.mutation.KeyIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   queryvalue.KeyTable,
			Columns: []string{queryvalue.KeyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(querykey.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.KeyID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// QueryValueCreateBulk is the builder for creating many QueryValue entities in bulk.
type QueryValueCreateBulk struct {
	config
	err      error
	builders []*QueryValueCreate
}

// Save creates the QueryValue entities in the database.
func (qvcb *QueryValueCreateBulk) Save(ctx context.Context) ([]*QueryValue, error) {
	if qvcb.err != nil {
		return nil, qvcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(qvcb.builders))
	nodes := make([]*QueryValue, len(qvcb.builders))
	mutators := make([]Mutator, len(qvcb.builders))
	for i := range qvcb.builders {
		func(i int, root context.Context) {
			builder := qvcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*QueryValueMutation)
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
					_, err = mutators[i+1].Mutate(root, qvcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, qvcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int64(id)
				}
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
		if _, err := mutators[0].Mutate(ctx, qvcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (qvcb *QueryValueCreateBulk) SaveX(ctx context.Context) []*QueryValue {
	v, err := qvcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (qvcb *QueryValueCreateBulk) Exec(ctx context.Context) error {
	_, err := qvcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (qvcb *QueryValueCreateBulk) ExecX(ctx context.Context) {
	if err := qvcb.Exec(ctx); err != nil {
		panic(err)
	}
}
