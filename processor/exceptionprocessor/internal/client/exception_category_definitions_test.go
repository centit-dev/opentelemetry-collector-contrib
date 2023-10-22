package client

import (
	"reflect"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/exceptionprocessor/ent/schema"
)

func Test_getDefinitionValue(t *testing.T) {
	type args struct {
		item *schema.ExceptionDefinitionCondition
	}
	tests := []struct {
		name string
		args args
		want ExceptionDefinitionValue
	}{
		{
			name: "simple-string-value",
			args: args{
				item: &schema.ExceptionDefinitionCondition{
					Column: "column",
					Op:     "simple-op",
					Value:  "simple-value",
				},
			},
			want: ExceptionDefinitionValue{
				StringValues: []string{"simple-value"},
			},
		},
		{
			name: "simple-number-value",
			args: args{
				item: &schema.ExceptionDefinitionCondition{
					Column: "column",
					Op:     "simple-op",
					Value:  1.,
				},
			},
			want: ExceptionDefinitionValue{
				numberValues: []float64{1.},
			},
		},
		{
			name: "simple-boolean-value",
			args: args{
				item: &schema.ExceptionDefinitionCondition{
					Column: "column",
					Op:     "simple-op",
					Value:  true,
				},
			},
			want: ExceptionDefinitionValue{
				boolValue: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDefinitionValue(tt.args.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDefinitionValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExceptionCategoryDefinition_Match(t *testing.T) {
	type fields struct {
		column string
		op     string
		value  ExceptionDefinitionValue
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
				value: ExceptionDefinitionValue{
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
				value: ExceptionDefinitionValue{
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
			definition := &ExceptionCategoryDefinition{
				column: tt.fields.column,
				op:     tt.fields.op,
				value:  tt.fields.value,
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
