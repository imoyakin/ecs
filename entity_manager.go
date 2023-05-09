package ecs

type EntityManager interface {
	AddEntity(entities ...Entity)
	Entities() (entities []Entity)
	FilterByType(types ...string) (entities []Entity)
	FilterByMask(mask uint64) (entities []Entity)
	FilterByNames(names ...string) (entities []Entity)
	Get(id string) (entity Entity)
	RemoveEntity(entity Entity)
}
