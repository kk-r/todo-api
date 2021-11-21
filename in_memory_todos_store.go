package main

import "sync"

func NewInMemoryTodoStore() *InMemoryTodoStore {
	return &InMemoryTodoStore{
		map[int]Todo{},
		sync.RWMutex{},
	}
}

type InMemoryTodoStore struct {
	store map[int]Todo
	lock  sync.RWMutex
}

func (i *InMemoryTodoStore) GetTodos() map[int]Todo {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.store
}

func (i *InMemoryTodoStore) CreateTodo(todo Todo) {
	i.lock.Lock()
	defer i.lock.Unlock()
	idx := len(i.store)
	i.store[idx] = todo
}
