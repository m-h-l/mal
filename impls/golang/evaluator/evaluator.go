package evaluator

import (
	"fmt"
	"mal/core"
	"mal/env"
	"mal/types"
)

func isSpecialForm(name string) bool {
	m := map[string]interface{}{"def!": nil, "let*": nil}
	_, ok := m[name]
	return ok
}

func Eval(ast types.MalType, e *env.Env) (types.MalType, bool) {
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

	if isSpecialForm(fnExpr.GetStr()) {
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
		panic("boom")
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

	switch form.GetStr() {
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
	default:
		return list, true
	}
}
