package ecs

// Componet contains only the data.
type Component interface {
	ID() string //TODO: change to uuid
	Mask() uint64
}
