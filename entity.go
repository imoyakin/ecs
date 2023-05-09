package ecs

type Entity interface {
	ID() string

	Mask() uint64

	Type() string

	Add(cn ...Component)

	GetComponent(c string) Component

	Name() string

	// removes the component by id.
	Remove(cn ...Component)
}