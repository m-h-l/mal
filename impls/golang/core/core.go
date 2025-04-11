package core

import (
	"mal/env"
	"mal/types"
)

type BuiltInFn struct {
	name string
	fn   func(*env.Env, ...types.MalType) types.MalType
}

func NewBuiltinFn(name string, fn func(*env.Env, ...types.MalType) types.MalType) BuiltInFn {
	return BuiltInFn{
		name: name,
		fn:   fn,
	}
}

func (fn BuiltInFn) GetName() string {
	return fn.name
}

func (fn BuiltInFn) Apply(e *env.Env, args ...types.MalType) (types.MalType, bool) {
	return fn.fn(e, args...), true
}

func (fn BuiltInFn) GetStr() string {
	return "builtin<" + fn.name + ">"
}

func (fn BuiltInFn) GetTypeId() types.MalTypeId {
	return types.BuiltInFunction
}

func Plus() BuiltInFn {
	return NewBuiltinFn("+", func(e *env.Env, args ...types.MalType) types.MalType {
		acc := types.NewMalNumber(0)
		for _, arg := range args {
			if arg.GetTypeId() != types.Number {
				panic("boom!")
			}
			n := acc.Add(*arg.(*types.MalNumber))
			acc = &n
		}
		return acc
	})
}

func Multiply() BuiltInFn {
	return NewBuiltinFn("*", func(e *env.Env, args ...types.MalType) types.MalType {
		acc := types.NewMalNumber(1)
		for _, arg := range args {
			if arg.GetTypeId() != types.Number {
				panic("boom!")
			}
			n := acc.Multiply(*arg.(*types.MalNumber))
			acc = &n
		}
		return acc
	})
}

func Subtract() BuiltInFn {
	return NewBuiltinFn("-", func(e *env.Env, args ...types.MalType) types.MalType {
		for _, arg := range args {
			if arg.GetTypeId() != types.Number {
				panic("boom!")
			}
		}
		if len(args) == 0 {
			panic("boom!")
		}
		if len(args) == 1 {
			n := types.NewMalNumber(-1).Multiply(*args[0].(*types.MalNumber))
			return &n
		}

		acc := *args[0].(*types.MalNumber)

		for _, arg := range args[1:] {
			acc = acc.Minus(*arg.(*types.MalNumber))
		}
		return &acc
	})
}

func Divide() BuiltInFn {
	return NewBuiltinFn("/", func(e *env.Env, args ...types.MalType) types.MalType {
		for _, arg := range args {
			if arg.GetTypeId() != types.Number {
				panic("boom!")
			}
		}
		if len(args) < 2 {
			panic("boom!")
		}

		acc := *args[0].(*types.MalNumber)

		for _, arg := range args[1:] {
			acc = acc.Divide(*arg.(*types.MalNumber))
		}
		return &acc
	})
}

func AddCoreToEnv(e *env.Env) {
	builtins := []BuiltInFn{
		Plus(),
		Multiply(),
		Subtract(),
		Divide(),
	}

	for _, builtin := range builtins {
		e.Add(*types.NewMalGenericAtom(builtin.name), builtin)
	}
}
