package ecs

type SystemManager interface {
	AddSystem(systems ...System)
	RemoveSystem(system System)
	System() []System
}
