package evaluator

import (
	"fmt"
	"mal/core"
	"mal/env"
	"mal/types"
)

func isSpecialForm(name string) bool {
	m := map[string]interface{}{"if": nil, "def!": nil, "let*": nil, "fn*": nil, "do": nil}
	_, ok := m[name]
	return ok
}

func Eval(ast types.MalType, e *env.Env) (types.MalType, bool) {

	if v, ok := e.Get(*types.NewMalGenericAtom("DEBUG-EVAL")); ok {
		switch o := (*v).(type) {
		case *types.MalNil:
		case *types.MalBool:
			if !o.GetState() {
				fmt.Printf("EVAL: %s\n", ast.GetStr(false))
			}
		default:
			fmt.Printf("EVAL: %s\n", ast.GetStr(false))
		}
	}

	switch v := ast.(type) {
	case *types.MalSymbol:
		res, ok := e.Get(*v)
		if !ok && !isSpecialForm(v.GetAsString()) {
			fmt.Printf("Symbol '%s' not found\n", v.GetAsString())
			return v, false
		}
		if isSpecialForm(v.GetAsString()) {
			return v, true
		}
		return *res, true
	case *types.MalList:
		return evalList(v, e)
	default:
		return ast, true
	}
}

func evalList(list *types.MalList, e *env.Env) (types.MalType, bool) {
	if list.IsEmpty() {
		return list, true
	}

	fnExpr, _ := Eval(list.First(), e)

	if isSpecialForm(fnExpr.GetStr(false)) {
		return apply(*list, e)
	}

	switch fnExpr.GetTypeId() {
	case types.BuiltInFunction:
		args := []types.MalType{}
		for _, arg := range list.Tail() {
			r, ok := Eval(arg, e)
			if !ok {
				return list, false
			}
			args = append(args, r)
		}
		return fnExpr.(core.BuiltInFn).Apply(e, args...)
	case types.DefinedFunction:
		args := []types.MalType{}
		for _, arg := range list.Tail() {
			r, ok := Eval(arg, e)
			if !ok {
				return list, false
			}
			args = append(args, r)
		}
		return fnExpr.(*core.DefinedFn).Apply(e, args...)
	default:
		args := []types.MalType{}
		for _, item := range list.Children() {
			r, ok := Eval(item, e)
			if !ok {
				return list, false
			}
			args = append(args, r)
		}
		return types.NewMalList(list.GetTypeId(), args), true
	}
}

func apply(list types.MalList, e *env.Env) (types.MalType, bool) {

	form := list.First()
	rest := list.Tail()

	switch form.GetStr(false) {
	case "if":
		if len(rest) < 2 {
			panic("boom!")
		}

		first, _ := Eval(rest[0], e)

		if len(rest) == 2 {
			if first.GetTypeId() == types.Nil {
				return types.NewMalNil(), true
			}
			if first.GetTypeId() == types.Boolean && !first.(*types.MalBool).GetState() {
				return types.NewMalNil(), true
			}
			return Eval(rest[1], e)
		}

		third := rest[2]

		if first.GetTypeId() == types.Nil {
			return Eval(third, e)
		}
		if first.GetTypeId() == types.Boolean && !first.(*types.MalBool).GetState() {
			return Eval(third, e)
		}
		return Eval(rest[1], e)
	case "def!":
		if len(rest) < 2 {
			panic("boom!")
		}
		first, second := rest[0], rest[1]
		if first.GetTypeId() != types.Symbol {
			panic("boom!")
		}
		sy := first.(*types.MalSymbol)
		r, ok := Eval(second, e)
		if !ok {
			return list, false
		}
		e.Add(*sy, r)
		return r, true
	case "let*":
		if len(rest) < 2 {
			panic("boom!")
		}
		first, second := rest[0], rest[1]
		if first.GetTypeId() != types.List && first.GetTypeId() != types.Vector {
			panic("boom!")
		}
		ne := env.NewEnv(e)
		children := first.(*types.MalList).Children()

		if len(children)%2 != 0 {
			panic("boom!")
		}

		for n := 0; n < len(children); n += 2 {
			k := *children[n].(*types.MalSymbol)
			v := children[n+1]
			r, ok := Eval(v, ne)
			if !ok {
				return list, false
			}
			ne.Add(k, r)
		}
		return Eval(second, ne)
	case "fn*":
		first, second := rest[0], rest[1]
		if first.GetTypeId() != types.List && first.GetTypeId() != types.Vector {
			panic("puff!")
		}

		// Capture the current environment in a closure
		binds := first.(*types.MalList).Children()
		body := second
		return core.NewDefinedFn(func(closureEnv *env.Env, args ...types.MalType) types.MalType {
			ne := env.NewEnv(e)
			for n, param := range binds {
				sy := *param.(*types.MalSymbol)
				if sy.GetAsString() == "&" {
					b := *binds[n+1].(*types.MalSymbol)
					ne.Add(b, types.NewMalList(types.List, args[n:]))
					break
				}
				ne.Add(sy, args[n])
			}
			result, _ := Eval(body, ne)
			return result
		}), true
	case "do":
		if len(rest) == 0 {
			return types.NewMalNil(), true
		}
		for i := 0; i < len(rest)-1; i++ {
			_, ok := Eval(rest[i], e)
			if !ok {
				return list, false
			}
		}
		last, ok := Eval(rest[len(rest)-1], e)
		if !ok {
			return list, false
		}
		return last, true
	default:
		return list, true
	}
}
