package ecs

import (
	"github.com/lithammer/shortuuid"
)

type ObjectEntity struct {
	id         string
	mask       uint64
	name       string
	Components map[string]Component

	parent   *Entity
	children []*Entity
}

func CreateObjectEntity() Entity {
	return &ObjectEntity{
		id: shortuuid.New(),
	}
}

func (e *ObjectEntity) ID() string {
	return e.id
}

func (e *ObjectEntity) Mask() uint64 {
	return e.mask
}

func (e *ObjectEntity) Name() string {
	return e.name
}

func (e *ObjectEntity) Type() string {
	return "object"
}

func (e *ObjectEntity) Add(cn ...Component) {
	for _, c := range cn {
		e.Components[c.ID()] = c
	}
}

func (e *ObjectEntity) GetComponent(id string) Component {
	return e.Components[id]
}

func (e *ObjectEntity) Remove(cn ...Component) {
	for _, c := range cn {
		delete(e.Components, c.ID())
	}
}
