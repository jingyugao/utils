package main

import (
	"reflect"
	"testing"
)

func TestCheckCircle(t *testing.T) {
	tests := []struct {
		m           map[string][]string
		wantAllPath [][]string
	}{
		// a->b->c->a
		{map[string][]string{"a": []string{"b"}, "b": []string{"c"}, "c": []string{"a"}}, [][]string{[]string{"a", "b", "c"}}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run("tt.name", func(t *testing.T) {
			if gotAllPath := CheckCircle(tt.m); !reflect.DeepEqual(gotAllPath, tt.wantAllPath) {
				t.Errorf("CheckCircle() = %v, want %v", gotAllPath, tt.wantAllPath)
			}
		})
	}
}
