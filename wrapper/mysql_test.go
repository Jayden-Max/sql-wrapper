package wrapper

import (
	"reflect"
	"testing"
)

func TestNewMySQLTracerWrapper(t *testing.T) {
	tests := []struct {
		name string
		want *TracerWrapper
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMySQLTracerWrapper(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMySQLTracerWrapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMySQLTracerWrapperWithOpts(t *testing.T) {
	type args struct {
		options []TracerOption
	}
	tests := []struct {
		name string
		args args
		want *TracerWrapper
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMySQLTracerWrapperWithOpts(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMySQLTracerWrapperWithOpts() = %v, want %v", got, tt.want)
			}
		})
	}
}
