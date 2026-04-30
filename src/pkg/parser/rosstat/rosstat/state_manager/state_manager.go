package state_manager

type StateMapFuncType[T any] map[int]func([]string, *StateManager[T])

type StateManager[T any] struct {
	state int

	parsedSlice []*T

	stateFuncMap StateMapFuncType[T]
}

func NewStateManager[T any](stateFuncMap StateMapFuncType[T]) *StateManager[T] {
	return &StateManager[T]{
		stateFuncMap: stateFuncMap,
	}
}

func (manager *StateManager[T]) getState() int {
	return manager.state
}

func (manager *StateManager[T]) ChangeState(state int) {
	manager.state = state
}

func (manager *StateManager[T]) DoOnStateState(row []string) {
	if fn, exists := manager.stateFuncMap[manager.state]; exists {
		fn(row, manager)
	}
}

func (manager *StateManager[T]) PutParsed(parsed *T) {
	manager.parsedSlice = append(manager.parsedSlice, parsed)
}

func (manager *StateManager[T]) GetAllParsed() []*T {
	return manager.parsedSlice
}
