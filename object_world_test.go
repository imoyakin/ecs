package ecs

import "testing"

func TestObjectWorld_AddEntity(t *testing.T) {
	type args struct {
		entities []Entity
	}
	tests := []struct {
		name string
		w    *ObjectWorld
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.AddEntity(tt.args.entities...)
		})
	}
}
