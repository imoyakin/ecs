package ecs_test

import (
	"ecs"
	"fmt"
	"testing"
)

type Soldier struct {
	id         string
	name       string
	mask       uint64
	components map[string]ecs.Component
}

func (s *Soldier) ID() string {
	return s.id
}

func (s *Soldier) Name() string {
	return s.name
}

func (s *Soldier) Mask() uint64 {
	return s.mask
}

func (s *Soldier) Type() string {
	return "soldier"
}

func (s *Soldier) Add(cs ...ecs.Component) {
	for _, c := range cs {
		s.components[c.ID()] = c
		s.mask |= 1 << c.Mask()
	}
}

func (s *Soldier) Remove(cs ...ecs.Component) {
	for _, c := range cs {
		delete(s.components, c.ID())
		s.mask &= ^(1 << c.Mask())
	}
}

func (s *Soldier) GetComponent(id string) ecs.Component {
	return s.components[id]
}

type Team struct {
	num int
}

func (t *Team) ID() string {
	return "team"
}

func (t *Team) Mask() uint64 {
	return 1
}

type Health struct {
	Current int
	Max     int
}

func (h *Health) ID() string {
	return "health"
}

func (h *Health) Mask() uint64 {
	return 1
}

type Position struct {
	X int
	Y int
}

func (p *Position) ID() string {
	return "position"
}

func (p *Position) Mask() uint64 {
	return 1
}

type Velocity struct {
	X int
	Y int
}

func (v *Velocity) ID() string {
	return "velocity"
}

func (v *Velocity) Mask() uint64 {
	return 1
}

type Weapon struct {
	Damage int
	Range  int
	Name   string
}

func (w *Weapon) ID() string {
	return "weapon"
}

func (w *Weapon) Mask() uint64 {
	return 1
}

type ObjectWorld struct {
	entities []ecs.Entity
	systems  []ecs.System
}

func (w *ObjectWorld) AddEntity(es ...ecs.Entity) {
	w.entities = append(w.entities, es...)
}

func (w *ObjectWorld) Entities() []ecs.Entity {
	return w.entities
}

func (w *ObjectWorld) FilterByType(types ...string) []ecs.Entity {
	var result []ecs.Entity
	for _, e := range w.entities {
		for _, t := range types {
			if (e).Type() == t {
				result = append(result, e)
			}
		}
	}
	return result
}

func (w *ObjectWorld) FilterByMask(mask uint64) []ecs.Entity {
	var result []ecs.Entity
	for _, e := range w.entities {
		if (e).Mask()&mask == mask {
			result = append(result, e)
		}
	}
	return result
}

func (w *ObjectWorld) FilterByNames(names ...string) []ecs.Entity {
	var result []ecs.Entity
	for _, e := range w.entities {
		//entity 中包含names
		if (e).Name() == names[0] {
			result = append(result, e)
		}

	}
	return result
}

func (w *ObjectWorld) Get(id string) ecs.Entity {
	for _, e := range w.entities {
		if (e).ID() == id {
			return e
		}
	}
	return nil
}

func (w *ObjectWorld) RemoveEntity(e ecs.Entity) {
	for i, e := range w.entities {
		if (e).ID() == e.ID() {
			w.entities = append(w.entities[:i], w.entities[i+1:]...)
			break
		}
	}
}

func (w *ObjectWorld) AddSystem(ss ...ecs.System) {
	w.systems = append(w.systems, ss...)
}

type CombatSystem struct {
}

func (s *CombatSystem) Setup() {}

/*
Simple implementation, there are problems, and it cannot cope with the team game mode, only support 1v1
*/
func (s *CombatSystem) Process(em ecs.EntityManager) (state ecs.State) {
	// every soldier attack enemy
	// 
	for _, e := range em.FilterByType((&Soldier{}).Type()) {
		p := e.GetComponent("position").(*Position)
		w := e.GetComponent("weapon").(*Weapon)
		t := e.GetComponent("team").(*Team)
		for _, e2 := range em.FilterByType((&Soldier{}).Type()) {
			p2 := e2.GetComponent("position").(*Position)
			t2 := e2.GetComponent("team").(*Team)
			if t.num == t2.num {
				continue
			}
			// check if in range
			if (p.X-p2.X)*(p.X-p2.X)+(p.Y-p2.Y)*(p.Y-p2.Y) <= w.Range*w.Range {
				h := e2.GetComponent("health").(*Health)
				h.Current -= w.Damage
				if h.Current <= 0 {
					em.RemoveEntity(e2)
				}
			} else {
				v := e.GetComponent("velocity").(*Velocity)
				// move the coordinates of p to p2, limited to the range of v
				if p.X < p2.X {
					p.X += v.X
				}
				if p.X > p2.X {
					p.X -= v.X
				}
				if p.Y < p2.Y {
					p.Y += v.Y
				}
				if p.Y > p2.Y {
					p.Y -= v.Y
				}
			}
		}
	}

	return ecs.StateContinue
}
func (s *CombatSystem) Teardown() {}

type HealthCheckSystem struct{}

func (h *HealthCheckSystem) Setup() {}
func (h *HealthCheckSystem) Process(em ecs.EntityManager) (state ecs.State) {
	// check health of all soldiers
	for _, e := range em.FilterByNames("health") {
		h := e.GetComponent("health").(*Health)
		if h.Current <= 0 {
			fmt.Printf("%s is dead\n", e.Name())
			em.RemoveEntity(e)
		}
	}
	// check if there is only one soldier left
	if len(em.Entities()) == 1 {
		return ecs.StateEnd
	}
	return ecs.StateContinue
}

func (h *HealthCheckSystem) Teardown() {}

// TestECS is a test case that creates an object world with two soldiers and runs the combat and movement systems until one soldier dies
func TestECS(t *testing.T) {
	world := &ObjectWorld{} // create an object world

	soldier1 := Soldier{
		id:         "soldier1",
		name:       "Alice",
		mask:       1,
		components: make(map[string]ecs.Component),
	}

	soldier2 := Soldier{
		id:         "soldier2",
		name:       "Bob",
		mask:       1,
		components: make(map[string]ecs.Component),
	}

	soldier1.Add(
		&Weapon{
			Damage: 3,
			Range:  5,
			Name:   "AK-47",
		},
		&Health{
			Current: 10,
			Max:     10,
		},
		&Team{
			num: 1,
		},
	)
	soldier2.Add(
		&Weapon{
			Damage: 4,
			Range:  3,
			Name:   "M4A1",
		},
		&Health{
			Current: 10,
			Max:     10,
		},
		&Team{
			num: 2,
		},
	)
	soldier1.Add(&Position{X: 0, Y: 0}, &Velocity{X: 2, Y: 1}) // add position and velocity components to soldier1
	soldier2.Add(&Position{X: 6, Y: 0}, &Velocity{X: 2, Y: 1}) // add position and velocity components to soldier2

	world.AddEntity(&soldier1, &soldier2) // add the soldiers to the world

	combat := &CombatSystem{} // create a combat system // create a movement system
	health := &HealthCheckSystem{}

	world.AddSystem(combat, health) // add the systems to the world

	var state ecs.State
	for {
		for _, system := range world.systems {
			// process each system
			state = system.Process(world)
		}
		var print string
		for _, e := range world.entities {
			h := e.GetComponent("health")
			if h != nil {
				print += fmt.Sprintf("%s has %d health! ", e.Name(), h.(*Health).Current)
			}
		}
		t.Logf("%s \n", print)
		if state == ecs.StateEnd {
			break
		}
	}

	t.Log("Game over")
}
