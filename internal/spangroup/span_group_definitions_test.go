package spangroup

import (
	"reflect"
	"testing"
)

type Definition struct {
	Column string
	Op     string
	Value  interface{}
}

func Test_getDefinitionValue(t *testing.T) {
	type args struct {
		item *Definition
	}
	tests := []struct {
		name string
		args args
		want GroupDefinitionValue
	}{
		{
			name: "simple-string-value",
			args: args{
				item: &Definition{
					Column: "column",
					Op:     "simple-op",
					Value:  "simple-value",
				},
			},
			want: GroupDefinitionValue{
				StringValues: []string{"simple-value"},
			},
		},
		{
			name: "simple-number-value",
			args: args{
				item: &Definition{
					Column: "column",
					Op:     "simple-op",
					Value:  1.,
				},
			},
			want: GroupDefinitionValue{
				numberValues: []float64{1.},
			},
		},
		{
			name: "simple-boolean-value",
			args: args{
				item: &Definition{
					Column: "column",
					Op:     "simple-op",
					Value:  true,
				},
			},
			want: GroupDefinitionValue{
				boolValue: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateDefinitionValue(tt.args.item.Value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDefinitionValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExceptionCategoryDefinition_Match(t *testing.T) {
	type fields struct {
		column string
		op     string
		value  GroupDefinitionValue
	}
	type args struct {
		attributes *map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "simple-match",
			fields: fields{
				column: "span.kind",
				op:     "=",
				value: GroupDefinitionValue{
					StringValues: []string{"value"},
				},
			},
			args: args{
				attributes: &map[string]interface{}{
					"span.kind": "value",
				},
			},
			want: true,
		},
		{
			name: "simple-match",
			fields: fields{
				column: "span.kind",
				op:     ">",
				value: GroupDefinitionValue{
					numberValues: []float64{1.},
				},
			},
			args: args{
				attributes: &map[string]interface{}{
					"span.kind": 2.,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			definition := &SpanGroupDefinition{
				Column: tt.fields.column,
				Op:     tt.fields.op,
				Value:  tt.fields.value,
			}
			if got := definition.Match(tt.args.attributes); got != tt.want {
				t.Errorf("ExceptionCategoryDefinition.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayContain(t *testing.T) {
	type args struct {
		op       string
		expected interface{}
		actual   interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "simple-contain",
			args: args{
				op:       "in",
				expected: []float64{1., 2.},
				actual:   2.,
			},
			want: true,
		},
		{
			name: "simple-not-contain",
			args: args{
				op:       "not-in",
				expected: []float64{1., 2.},
				actual:   3.,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arrayContain(tt.args.op, tt.args.expected, tt.args.actual); got != tt.want {
				t.Errorf("arrayContain() = %v, want %v", got, tt.want)
			}
		})
	}
}
