// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/apptypeprocessor/ent/middlewaredefinition"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/apptypeprocessor/ent/predicate"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/apptypeprocessor/ent/schema"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeMiddlewareDefinition = "MiddlewareDefinition"
)

// MiddlewareDefinitionMutation represents an operation that mutates the MiddlewareDefinition nodes in the graph.
type MiddlewareDefinitionMutation struct {
	config
	op                    Op
	typ                   string
	id                    *int64
	name                  *string
	span_conditions       *[]schema.MiddlewareDefinitionCondition
	appendspan_conditions []schema.MiddlewareDefinitionCondition
	is_valid              *int
	addis_valid           *int
	create_time           *time.Time
	update_time           *time.Time
	clearedFields         map[string]struct{}
	done                  bool
	oldValue              func(context.Context) (*MiddlewareDefinition, error)
	predicates            []predicate.MiddlewareDefinition
}

var _ ent.Mutation = (*MiddlewareDefinitionMutation)(nil)

// middlewaredefinitionOption allows management of the mutation configuration using functional options.
type middlewaredefinitionOption func(*MiddlewareDefinitionMutation)

// newMiddlewareDefinitionMutation creates new mutation for the MiddlewareDefinition entity.
func newMiddlewareDefinitionMutation(c config, op Op, opts ...middlewaredefinitionOption) *MiddlewareDefinitionMutation {
	m := &MiddlewareDefinitionMutation{
		config:        c,
		op:            op,
		typ:           TypeMiddlewareDefinition,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withMiddlewareDefinitionID sets the ID field of the mutation.
func withMiddlewareDefinitionID(id int64) middlewaredefinitionOption {
	return func(m *MiddlewareDefinitionMutation) {
		var (
			err   error
			once  sync.Once
			value *MiddlewareDefinition
		)
		m.oldValue = func(ctx context.Context) (*MiddlewareDefinition, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().MiddlewareDefinition.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withMiddlewareDefinition sets the old MiddlewareDefinition of the mutation.
func withMiddlewareDefinition(node *MiddlewareDefinition) middlewaredefinitionOption {
	return func(m *MiddlewareDefinitionMutation) {
		m.oldValue = func(context.Context) (*MiddlewareDefinition, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m MiddlewareDefinitionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m MiddlewareDefinitionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of MiddlewareDefinition entities.
func (m *MiddlewareDefinitionMutation) SetID(id int64) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *MiddlewareDefinitionMutation) ID() (id int64, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *MiddlewareDefinitionMutation) IDs(ctx context.Context) ([]int64, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int64{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().MiddlewareDefinition.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetName sets the "name" field.
func (m *MiddlewareDefinitionMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *MiddlewareDefinitionMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the MiddlewareDefinition entity.
// If the MiddlewareDefinition object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MiddlewareDefinitionMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *MiddlewareDefinitionMutation) ResetName() {
	m.name = nil
}

// SetSpanConditions sets the "span_conditions" field.
func (m *MiddlewareDefinitionMutation) SetSpanConditions(sdc []schema.MiddlewareDefinitionCondition) {
	m.span_conditions = &sdc
	m.appendspan_conditions = nil
}

// SpanConditions returns the value of the "span_conditions" field in the mutation.
func (m *MiddlewareDefinitionMutation) SpanConditions() (r []schema.MiddlewareDefinitionCondition, exists bool) {
	v := m.span_conditions
	if v == nil {
		return
	}
	return *v, true
}

// OldSpanConditions returns the old "span_conditions" field's value of the MiddlewareDefinition entity.
// If the MiddlewareDefinition object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MiddlewareDefinitionMutation) OldSpanConditions(ctx context.Context) (v []schema.MiddlewareDefinitionCondition, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSpanConditions is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSpanConditions requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSpanConditions: %w", err)
	}
	return oldValue.SpanConditions, nil
}

// AppendSpanConditions adds sdc to the "span_conditions" field.
func (m *MiddlewareDefinitionMutation) AppendSpanConditions(sdc []schema.MiddlewareDefinitionCondition) {
	m.appendspan_conditions = append(m.appendspan_conditions, sdc...)
}

// AppendedSpanConditions returns the list of values that were appended to the "span_conditions" field in this mutation.
func (m *MiddlewareDefinitionMutation) AppendedSpanConditions() ([]schema.MiddlewareDefinitionCondition, bool) {
	if len(m.appendspan_conditions) == 0 {
		return nil, false
	}
	return m.appendspan_conditions, true
}

// ClearSpanConditions clears the value of the "span_conditions" field.
func (m *MiddlewareDefinitionMutation) ClearSpanConditions() {
	m.span_conditions = nil
	m.appendspan_conditions = nil
	m.clearedFields[middlewaredefinition.FieldSpanConditions] = struct{}{}
}

// SpanConditionsCleared returns if the "span_conditions" field was cleared in this mutation.
func (m *MiddlewareDefinitionMutation) SpanConditionsCleared() bool {
	_, ok := m.clearedFields[middlewaredefinition.FieldSpanConditions]
	return ok
}

// ResetSpanConditions resets all changes to the "span_conditions" field.
func (m *MiddlewareDefinitionMutation) ResetSpanConditions() {
	m.span_conditions = nil
	m.appendspan_conditions = nil
	delete(m.clearedFields, middlewaredefinition.FieldSpanConditions)
}

// SetIsValid sets the "is_valid" field.
func (m *MiddlewareDefinitionMutation) SetIsValid(i int) {
	m.is_valid = &i
	m.addis_valid = nil
}

// IsValid returns the value of the "is_valid" field in the mutation.
func (m *MiddlewareDefinitionMutation) IsValid() (r int, exists bool) {
	v := m.is_valid
	if v == nil {
		return
	}
	return *v, true
}

// OldIsValid returns the old "is_valid" field's value of the MiddlewareDefinition entity.
// If the MiddlewareDefinition object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MiddlewareDefinitionMutation) OldIsValid(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldIsValid is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldIsValid requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldIsValid: %w", err)
	}
	return oldValue.IsValid, nil
}

// AddIsValid adds i to the "is_valid" field.
func (m *MiddlewareDefinitionMutation) AddIsValid(i int) {
	if m.addis_valid != nil {
		*m.addis_valid += i
	} else {
		m.addis_valid = &i
	}
}

// AddedIsValid returns the value that was added to the "is_valid" field in this mutation.
func (m *MiddlewareDefinitionMutation) AddedIsValid() (r int, exists bool) {
	v := m.addis_valid
	if v == nil {
		return
	}
	return *v, true
}

// ResetIsValid resets all changes to the "is_valid" field.
func (m *MiddlewareDefinitionMutation) ResetIsValid() {
	m.is_valid = nil
	m.addis_valid = nil
}

// SetCreateTime sets the "create_time" field.
func (m *MiddlewareDefinitionMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the value of the "create_time" field in the mutation.
func (m *MiddlewareDefinitionMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "create_time" field's value of the MiddlewareDefinition entity.
// If the MiddlewareDefinition object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MiddlewareDefinitionMutation) OldCreateTime(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "create_time" field.
func (m *MiddlewareDefinitionMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the "update_time" field.
func (m *MiddlewareDefinitionMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the value of the "update_time" field in the mutation.
func (m *MiddlewareDefinitionMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "update_time" field's value of the MiddlewareDefinition entity.
// If the MiddlewareDefinition object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MiddlewareDefinitionMutation) OldUpdateTime(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdateTime: %w", err)
	}
	return oldValue.UpdateTime, nil
}

// ResetUpdateTime resets all changes to the "update_time" field.
func (m *MiddlewareDefinitionMutation) ResetUpdateTime() {
	m.update_time = nil
}

// Where appends a list predicates to the MiddlewareDefinitionMutation builder.
func (m *MiddlewareDefinitionMutation) Where(ps ...predicate.MiddlewareDefinition) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the MiddlewareDefinitionMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *MiddlewareDefinitionMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.MiddlewareDefinition, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *MiddlewareDefinitionMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *MiddlewareDefinitionMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (MiddlewareDefinition).
func (m *MiddlewareDefinitionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *MiddlewareDefinitionMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m.name != nil {
		fields = append(fields, middlewaredefinition.FieldName)
	}
	if m.span_conditions != nil {
		fields = append(fields, middlewaredefinition.FieldSpanConditions)
	}
	if m.is_valid != nil {
		fields = append(fields, middlewaredefinition.FieldIsValid)
	}
	if m.create_time != nil {
		fields = append(fields, middlewaredefinition.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, middlewaredefinition.FieldUpdateTime)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *MiddlewareDefinitionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case middlewaredefinition.FieldName:
		return m.Name()
	case middlewaredefinition.FieldSpanConditions:
		return m.SpanConditions()
	case middlewaredefinition.FieldIsValid:
		return m.IsValid()
	case middlewaredefinition.FieldCreateTime:
		return m.CreateTime()
	case middlewaredefinition.FieldUpdateTime:
		return m.UpdateTime()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *MiddlewareDefinitionMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case middlewaredefinition.FieldName:
		return m.OldName(ctx)
	case middlewaredefinition.FieldSpanConditions:
		return m.OldSpanConditions(ctx)
	case middlewaredefinition.FieldIsValid:
		return m.OldIsValid(ctx)
	case middlewaredefinition.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case middlewaredefinition.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	}
	return nil, fmt.Errorf("unknown MiddlewareDefinition field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *MiddlewareDefinitionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case middlewaredefinition.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case middlewaredefinition.FieldSpanConditions:
		v, ok := value.([]schema.MiddlewareDefinitionCondition)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSpanConditions(v)
		return nil
	case middlewaredefinition.FieldIsValid:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIsValid(v)
		return nil
	case middlewaredefinition.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case middlewaredefinition.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	}
	return fmt.Errorf("unknown MiddlewareDefinition field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *MiddlewareDefinitionMutation) AddedFields() []string {
	var fields []string
	if m.addis_valid != nil {
		fields = append(fields, middlewaredefinition.FieldIsValid)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *MiddlewareDefinitionMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case middlewaredefinition.FieldIsValid:
		return m.AddedIsValid()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *MiddlewareDefinitionMutation) AddField(name string, value ent.Value) error {
	switch name {
	case middlewaredefinition.FieldIsValid:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIsValid(v)
		return nil
	}
	return fmt.Errorf("unknown MiddlewareDefinition numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *MiddlewareDefinitionMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(middlewaredefinition.FieldSpanConditions) {
		fields = append(fields, middlewaredefinition.FieldSpanConditions)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *MiddlewareDefinitionMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *MiddlewareDefinitionMutation) ClearField(name string) error {
	switch name {
	case middlewaredefinition.FieldSpanConditions:
		m.ClearSpanConditions()
		return nil
	}
	return fmt.Errorf("unknown MiddlewareDefinition nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *MiddlewareDefinitionMutation) ResetField(name string) error {
	switch name {
	case middlewaredefinition.FieldName:
		m.ResetName()
		return nil
	case middlewaredefinition.FieldSpanConditions:
		m.ResetSpanConditions()
		return nil
	case middlewaredefinition.FieldIsValid:
		m.ResetIsValid()
		return nil
	case middlewaredefinition.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case middlewaredefinition.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	}
	return fmt.Errorf("unknown MiddlewareDefinition field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *MiddlewareDefinitionMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *MiddlewareDefinitionMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *MiddlewareDefinitionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *MiddlewareDefinitionMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *MiddlewareDefinitionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *MiddlewareDefinitionMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *MiddlewareDefinitionMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown MiddlewareDefinition unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *MiddlewareDefinitionMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown MiddlewareDefinition edge %s", name)
}
