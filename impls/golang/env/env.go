package env

import "mal/types"

type Fn interface {
	HasName(name string) bool
	Apply(args ...types.MalType) []types.MalType
}

type BuiltInFn struct {
	name string
	fn   func(...types.MalType) []types.MalType
}

func NewBuiltinFn(name string, fn func(...types.MalType) []types.MalType) BuiltInFn {
	return BuiltInFn{
		name: name,
		fn:   fn,
	}
}

func (fn BuiltInFn) HasName(name string) bool {
	return fn.name == name
}

func (fn BuiltInFn) Apply(args ...types.MalType) []types.MalType {
	return fn.fn(args...)
}

type Env struct {
	Fns []Fn
}

func (env *Env) FindSymbol(symbol string) (*Fn, bool) {

	for _, fn := range env.Fns {
		if fn.HasName(symbol) {
			return &fn, true
		}
	}
	return nil, false
}
