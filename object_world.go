package ecs

type ObjectWorld struct {
	Entities map[string]Entity
	Systems []System
}

func (w *ObjectWorld) AddEntity(entities ...Entity) {
	for _, entity := range entities {
		w.Entities[entity.ID()] = entity
	}
}

func (w *ObjectWorld) FilterByMask(mask uint64) (entities []Entity) {
	for _, entity := range w.Entities {
		if entity.Mask()&mask == mask {
			entities = append(entities, entity)
		}
	}
	return
}

func (w *ObjectWorld) FilterByNames(names ...string) (entities []Entity) {
	for _, entity := range w.Entities {
		for _, name := range names {
			if entity.Name() == name {
				entities = append(entities, entity)
			}
		}
	}
	return
}

func (w *ObjectWorld) Get(id string) (entity Entity) {
	return w.Entities[id]
}

func (w *ObjectWorld) RemoveEntity(entity Entity) {
	delete(w.Entities, entity.ID())
}
