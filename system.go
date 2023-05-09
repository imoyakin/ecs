package ecs

type System interface {
	Setup()
	Process(em EntityManager) (state State)
	Teardown()
}
