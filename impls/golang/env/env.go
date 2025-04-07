package env

import "mal/types"

type Env struct {
	Fns []types.Fn
}

func (env *Env) FindSymbol(symbol types.MalSymbol) (*types.Fn, bool) {

	for _, fn := range env.Fns {
		if fn.HasName(symbol.GetAsString()) {
			return &fn, true
		}
	}
	return nil, false
}
