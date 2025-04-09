package core

import (
	"mal/env"
	"mal/types"

	"github.com/kr/pretty"
)

type EvalFn func(types.MalType) types.MalType

type BuiltInFn struct {
	name string
	fn   func(*env.Env, EvalFn, ...types.MalType) types.MalType
}

func eval(evalFn EvalFn, items []types.MalType) []types.MalType {
	args := []types.MalType{}
	for _, item := range items {
		args = append(args, evalFn(item))
	}
	return args
}

func NewBuiltinFn(name string, fn func(*env.Env, EvalFn, ...types.MalType) types.MalType) BuiltInFn {
	return BuiltInFn{
		name: name,
		fn:   fn,
	}
}

func (fn BuiltInFn) GetName() string {
	return fn.name
}

func (fn BuiltInFn) Apply(e *env.Env, eval EvalFn, args ...types.MalType) types.MalType {
	return fn.fn(e, eval, args...)
}

func (fn BuiltInFn) GetStr() string {
	return "builtin<" + fn.name + ">"
}

func (fn BuiltInFn) GetTypeId() types.MalTypeId {
	return types.BuiltInFunction
}

func Plus() BuiltInFn {
	return NewBuiltinFn("+", func(e *env.Env, evalFn EvalFn, args ...types.MalType) types.MalType {
		acc := types.NewMalNumber(0)
		for _, arg := range eval(evalFn, args) {
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
	return NewBuiltinFn("*", func(e *env.Env, evalFn EvalFn, args ...types.MalType) types.MalType {
		acc := types.NewMalNumber(1)
		for _, arg := range eval(evalFn, args) {
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
	return NewBuiltinFn("-", func(e *env.Env, evalFn EvalFn, args ...types.MalType) types.MalType {
		evaledArgs := eval(evalFn, args)
		for _, arg := range evaledArgs {
			if arg.GetTypeId() != types.Number {
				panic("boom!")
			}
		}
		if len(evaledArgs) == 0 {
			panic("boom!")
		}
		if len(evaledArgs) == 1 {
			n := types.NewMalNumber(-1).Multiply(*evaledArgs[0].(*types.MalNumber))
			return &n
		}

		acc := *evaledArgs[0].(*types.MalNumber)

		for _, arg := range evaledArgs[1:] {
			acc = acc.Minus(*arg.(*types.MalNumber))
		}
		return &acc
	})
}

func Divide() BuiltInFn {
	return NewBuiltinFn("/", func(e *env.Env, evalFn EvalFn, args ...types.MalType) types.MalType {
		evaledArgs := eval(evalFn, args)
		for _, arg := range evaledArgs {
			if arg.GetTypeId() != types.Number {
				panic("boom!")
			}
		}
		if len(evaledArgs) < 2 {
			panic("boom!")
		}

		acc := *evaledArgs[0].(*types.MalNumber)

		for _, arg := range evaledArgs[1:] {
			acc = acc.Divide(*arg.(*types.MalNumber))
		}
		return &acc
	})
}

func Def() BuiltInFn {
	return NewBuiltinFn("def!", func(e *env.Env, evalFn EvalFn, args ...types.MalType) types.MalType {
		first, second := args[0], eval(evalFn, args[1:])[0]

		if len(args) < 2 {
			panic("boom!")
		}

		if first.GetTypeId() != types.Symbol {
			panic("boom!")
		}

		sy := first.(*types.MalSymbol)
		e.Add(*sy, second)

		return second
	})
}

func Let() BuiltInFn {
	return NewBuiltinFn("let*", func(e *env.Env, evalFn EvalFn, args ...types.MalType) types.MalType {
		first, second := args[0], args[1]

		if len(args) < 2 {
			panic("boom!")
		}

		if first.GetTypeId() != types.List {
			panic("boom!")
		}

		ne := env.NewEnv(e)
		children := first.(*types.MalList).Children()

		if len(children)%2 != 0 {
			panic("boom!")
		}

		pretty.Println(first, second)

		for n := 0; n < len(children); n += 2 {
			k := *children[n].(*types.MalSymbol)
			v := children[n+1]
			ne.Add(k, evalFn(v))
		}
		return evalFn(second)
	})
}

func AddCoreToEnv(e *env.Env) {
	builtins := []BuiltInFn{
		Def(),
		Let(),
		Plus(),
		Multiply(),
		Subtract(),
		Divide(),
	}

	for _, builtin := range builtins {
		e.Add(*types.NewMalGenericAtom(builtin.name), builtin)
	}
}
