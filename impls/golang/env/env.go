package env

import (
	"mal/types"
)

type Env struct {
	outer *Env
	data  map[string]types.MalType
}

func NewEnv(outter *Env) *Env {
	return &Env{
		outer: outter,
		data:  map[string]types.MalType{},
	}
}

func (env *Env) Add(symbol types.MalSymbol, value types.MalType) {
	env.data[symbol.GetAsString()] = value
}

func (env Env) Get(symbol types.MalSymbol) (*types.MalType, bool) {
	v, ok := env.data[symbol.GetAsString()]
	if ok {
		return &v, true
	}

	if env.outer != nil {
		return env.outer.Get(symbol)
	}

	return nil, false
}
