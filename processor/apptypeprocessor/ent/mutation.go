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
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/apptypeprocessor/ent/predicate"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/apptypeprocessor/ent/schema"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/apptypeprocessor/ent/softwaredefinition"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeSoftwareDefinition = "SoftwareDefinition"
)

// SoftwareDefinitionMutation represents an operation that mutates the SoftwareDefinition nodes in the graph.
type SoftwareDefinitionMutation struct {
	config
	op                    Op
	typ                   string
	id                    *int64
	name                  *string
	_type                 *int16
	add_type              *int16
	span_conditions       *[]schema.SoftwareDefinitionCondition
	appendspan_conditions []schema.SoftwareDefinitionCondition
	is_valid              *int
	addis_valid           *int
	create_time           *time.Time
	update_time           *time.Time
	clearedFields         map[string]struct{}
	done                  bool
	oldValue              func(context.Context) (*SoftwareDefinition, error)
	predicates            []predicate.SoftwareDefinition
}

var _ ent.Mutation = (*SoftwareDefinitionMutation)(nil)

// softwaredefinitionOption allows management of the mutation configuration using functional options.
type softwaredefinitionOption func(*SoftwareDefinitionMutation)

// newSoftwareDefinitionMutation creates new mutation for the SoftwareDefinition entity.
func newSoftwareDefinitionMutation(c config, op Op, opts ...softwaredefinitionOption) *SoftwareDefinitionMutation {
	m := &SoftwareDefinitionMutation{
		config:        c,
		op:            op,
		typ:           TypeSoftwareDefinition,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withSoftwareDefinitionID sets the ID field of the mutation.
func withSoftwareDefinitionID(id int64) softwaredefinitionOption {
	return func(m *SoftwareDefinitionMutation) {
		var (
			err   error
			once  sync.Once
			value *SoftwareDefinition
		)
		m.oldValue = func(ctx context.Context) (*SoftwareDefinition, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().SoftwareDefinition.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withSoftwareDefinition sets the old SoftwareDefinition of the mutation.
func withSoftwareDefinition(node *SoftwareDefinition) softwaredefinitionOption {
	return func(m *SoftwareDefinitionMutation) {
		m.oldValue = func(context.Context) (*SoftwareDefinition, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SoftwareDefinitionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SoftwareDefinitionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of SoftwareDefinition entities.
func (m *SoftwareDefinitionMutation) SetID(id int64) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *SoftwareDefinitionMutation) ID() (id int64, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *SoftwareDefinitionMutation) IDs(ctx context.Context) ([]int64, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int64{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().SoftwareDefinition.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetName sets the "name" field.
func (m *SoftwareDefinitionMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *SoftwareDefinitionMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the SoftwareDefinition entity.
// If the SoftwareDefinition object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SoftwareDefinitionMutation) OldName(ctx context.Context) (v string, err error) {
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
func (m *SoftwareDefinitionMutation) ResetName() {
	m.name = nil
}

// SetType sets the "type" field.
func (m *SoftwareDefinitionMutation) SetType(i int16) {
	m._type = &i
	m.add_type = nil
}

// GetType returns the value of the "type" field in the mutation.
func (m *SoftwareDefinitionMutation) GetType() (r int16, exists bool) {
	v := m._type
	if v == nil {
		return
	}
	return *v, true
}

// OldType returns the old "type" field's value of the SoftwareDefinition entity.
// If the SoftwareDefinition object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SoftwareDefinitionMutation) OldType(ctx context.Context) (v int16, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldType is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldType requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldType: %w", err)
	}
	return oldValue.Type, nil
}

// AddType adds i to the "type" field.
func (m *SoftwareDefinitionMutation) AddType(i int16) {
	if m.add_type != nil {
		*m.add_type += i
	} else {
		m.add_type = &i
	}
}

// AddedType returns the value that was added to the "type" field in this mutation.
func (m *SoftwareDefinitionMutation) AddedType() (r int16, exists bool) {
	v := m.add_type
	if v == nil {
		return
	}
	return *v, true
}

// ResetType resets all changes to the "type" field.
func (m *SoftwareDefinitionMutation) ResetType() {
	m._type = nil
	m.add_type = nil
}

// SetSpanConditions sets the "span_conditions" field.
func (m *SoftwareDefinitionMutation) SetSpanConditions(sdc []schema.SoftwareDefinitionCondition) {
	m.span_conditions = &sdc
	m.appendspan_conditions = nil
}

// SpanConditions returns the value of the "span_conditions" field in the mutation.
func (m *SoftwareDefinitionMutation) SpanConditions() (r []schema.SoftwareDefinitionCondition, exists bool) {
	v := m.span_conditions
	if v == nil {
		return
	}
	return *v, true
}

// OldSpanConditions returns the old "span_conditions" field's value of the SoftwareDefinition entity.
// If the SoftwareDefinition object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SoftwareDefinitionMutation) OldSpanConditions(ctx context.Context) (v []schema.SoftwareDefinitionCondition, err error) {
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
func (m *SoftwareDefinitionMutation) AppendSpanConditions(sdc []schema.SoftwareDefinitionCondition) {
	m.appendspan_conditions = append(m.appendspan_conditions, sdc...)
}

// AppendedSpanConditions returns the list of values that were appended to the "span_conditions" field in this mutation.
func (m *SoftwareDefinitionMutation) AppendedSpanConditions() ([]schema.SoftwareDefinitionCondition, bool) {
	if len(m.appendspan_conditions) == 0 {
		return nil, false
	}
	return m.appendspan_conditions, true
}

// ClearSpanConditions clears the value of the "span_conditions" field.
func (m *SoftwareDefinitionMutation) ClearSpanConditions() {
	m.span_conditions = nil
	m.appendspan_conditions = nil
	m.clearedFields[softwaredefinition.FieldSpanConditions] = struct{}{}
}

// SpanConditionsCleared returns if the "span_conditions" field was cleared in this mutation.
func (m *SoftwareDefinitionMutation) SpanConditionsCleared() bool {
	_, ok := m.clearedFields[softwaredefinition.FieldSpanConditions]
	return ok
}

// ResetSpanConditions resets all changes to the "span_conditions" field.
func (m *SoftwareDefinitionMutation) ResetSpanConditions() {
	m.span_conditions = nil
	m.appendspan_conditions = nil
	delete(m.clearedFields, softwaredefinition.FieldSpanConditions)
}

// SetIsValid sets the "is_valid" field.
func (m *SoftwareDefinitionMutation) SetIsValid(i int) {
	m.is_valid = &i
	m.addis_valid = nil
}

// IsValid returns the value of the "is_valid" field in the mutation.
func (m *SoftwareDefinitionMutation) IsValid() (r int, exists bool) {
	v := m.is_valid
	if v == nil {
		return
	}
	return *v, true
}

// OldIsValid returns the old "is_valid" field's value of the SoftwareDefinition entity.
// If the SoftwareDefinition object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SoftwareDefinitionMutation) OldIsValid(ctx context.Context) (v int, err error) {
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
func (m *SoftwareDefinitionMutation) AddIsValid(i int) {
	if m.addis_valid != nil {
		*m.addis_valid += i
	} else {
		m.addis_valid = &i
	}
}

// AddedIsValid returns the value that was added to the "is_valid" field in this mutation.
func (m *SoftwareDefinitionMutation) AddedIsValid() (r int, exists bool) {
	v := m.addis_valid
	if v == nil {
		return
	}
	return *v, true
}

// ResetIsValid resets all changes to the "is_valid" field.
func (m *SoftwareDefinitionMutation) ResetIsValid() {
	m.is_valid = nil
	m.addis_valid = nil
}

// SetCreateTime sets the "create_time" field.
func (m *SoftwareDefinitionMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the value of the "create_time" field in the mutation.
func (m *SoftwareDefinitionMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "create_time" field's value of the SoftwareDefinition entity.
// If the SoftwareDefinition object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SoftwareDefinitionMutation) OldCreateTime(ctx context.Context) (v time.Time, err error) {
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
func (m *SoftwareDefinitionMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the "update_time" field.
func (m *SoftwareDefinitionMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the value of the "update_time" field in the mutation.
func (m *SoftwareDefinitionMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "update_time" field's value of the SoftwareDefinition entity.
// If the SoftwareDefinition object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SoftwareDefinitionMutation) OldUpdateTime(ctx context.Context) (v time.Time, err error) {
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
func (m *SoftwareDefinitionMutation) ResetUpdateTime() {
	m.update_time = nil
}

// Where appends a list predicates to the SoftwareDefinitionMutation builder.
func (m *SoftwareDefinitionMutation) Where(ps ...predicate.SoftwareDefinition) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the SoftwareDefinitionMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *SoftwareDefinitionMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.SoftwareDefinition, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *SoftwareDefinitionMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *SoftwareDefinitionMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (SoftwareDefinition).
func (m *SoftwareDefinitionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *SoftwareDefinitionMutation) Fields() []string {
	fields := make([]string, 0, 6)
	if m.name != nil {
		fields = append(fields, softwaredefinition.FieldName)
	}
	if m._type != nil {
		fields = append(fields, softwaredefinition.FieldType)
	}
	if m.span_conditions != nil {
		fields = append(fields, softwaredefinition.FieldSpanConditions)
	}
	if m.is_valid != nil {
		fields = append(fields, softwaredefinition.FieldIsValid)
	}
	if m.create_time != nil {
		fields = append(fields, softwaredefinition.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, softwaredefinition.FieldUpdateTime)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *SoftwareDefinitionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case softwaredefinition.FieldName:
		return m.Name()
	case softwaredefinition.FieldType:
		return m.GetType()
	case softwaredefinition.FieldSpanConditions:
		return m.SpanConditions()
	case softwaredefinition.FieldIsValid:
		return m.IsValid()
	case softwaredefinition.FieldCreateTime:
		return m.CreateTime()
	case softwaredefinition.FieldUpdateTime:
		return m.UpdateTime()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *SoftwareDefinitionMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case softwaredefinition.FieldName:
		return m.OldName(ctx)
	case softwaredefinition.FieldType:
		return m.OldType(ctx)
	case softwaredefinition.FieldSpanConditions:
		return m.OldSpanConditions(ctx)
	case softwaredefinition.FieldIsValid:
		return m.OldIsValid(ctx)
	case softwaredefinition.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case softwaredefinition.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	}
	return nil, fmt.Errorf("unknown SoftwareDefinition field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SoftwareDefinitionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case softwaredefinition.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case softwaredefinition.FieldType:
		v, ok := value.(int16)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetType(v)
		return nil
	case softwaredefinition.FieldSpanConditions:
		v, ok := value.([]schema.SoftwareDefinitionCondition)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSpanConditions(v)
		return nil
	case softwaredefinition.FieldIsValid:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIsValid(v)
		return nil
	case softwaredefinition.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case softwaredefinition.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	}
	return fmt.Errorf("unknown SoftwareDefinition field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *SoftwareDefinitionMutation) AddedFields() []string {
	var fields []string
	if m.add_type != nil {
		fields = append(fields, softwaredefinition.FieldType)
	}
	if m.addis_valid != nil {
		fields = append(fields, softwaredefinition.FieldIsValid)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *SoftwareDefinitionMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case softwaredefinition.FieldType:
		return m.AddedType()
	case softwaredefinition.FieldIsValid:
		return m.AddedIsValid()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SoftwareDefinitionMutation) AddField(name string, value ent.Value) error {
	switch name {
	case softwaredefinition.FieldType:
		v, ok := value.(int16)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddType(v)
		return nil
	case softwaredefinition.FieldIsValid:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIsValid(v)
		return nil
	}
	return fmt.Errorf("unknown SoftwareDefinition numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *SoftwareDefinitionMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(softwaredefinition.FieldSpanConditions) {
		fields = append(fields, softwaredefinition.FieldSpanConditions)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *SoftwareDefinitionMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *SoftwareDefinitionMutation) ClearField(name string) error {
	switch name {
	case softwaredefinition.FieldSpanConditions:
		m.ClearSpanConditions()
		return nil
	}
	return fmt.Errorf("unknown SoftwareDefinition nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *SoftwareDefinitionMutation) ResetField(name string) error {
	switch name {
	case softwaredefinition.FieldName:
		m.ResetName()
		return nil
	case softwaredefinition.FieldType:
		m.ResetType()
		return nil
	case softwaredefinition.FieldSpanConditions:
		m.ResetSpanConditions()
		return nil
	case softwaredefinition.FieldIsValid:
		m.ResetIsValid()
		return nil
	case softwaredefinition.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case softwaredefinition.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	}
	return fmt.Errorf("unknown SoftwareDefinition field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *SoftwareDefinitionMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *SoftwareDefinitionMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *SoftwareDefinitionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *SoftwareDefinitionMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *SoftwareDefinitionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *SoftwareDefinitionMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *SoftwareDefinitionMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown SoftwareDefinition unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *SoftwareDefinitionMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown SoftwareDefinition edge %s", name)
}