package object

func NewEnvironment() *Environment {
	return &Environment{store: map[string]Object{}}
}

type Environment struct {
	store  map[string]Object
	parent *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.parent != nil {
		return e.parent.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, object Object) Object {
	e.store[name] = object
	return object
}

func (e *Environment) NewEnclosedEnvironment() *Environment {
	newEnv := NewEnvironment()
	newEnv.parent = e
	return newEnv
}
