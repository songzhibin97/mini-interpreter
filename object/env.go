package object

import "sync"

type Env struct {
	sync.Mutex
	store  map[string]Object
	parent *Env
}

func (e *Env) Get(key string) (Object, bool) {
	e.Lock()
	obj, ok := e.store[key]
	e.Unlock()
	if ok {
		return obj, true
	}
	if e.parent != nil {
		return e.parent.Get(key)
	}
	return nil, false
}

func (e *Env) Set(key string, obj Object) Object {
	e.Lock()
	defer e.Unlock()
	if e.store == nil {
		e.store = make(map[string]Object)
	}
	old := e.store[key]
	e.store[key] = obj
	return old
}

func NewEnv(parent *Env) *Env {
	return &Env{parent: parent}
}
