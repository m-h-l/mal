package evaluator

import (
	"fmt"
	"mal/core"
	"mal/env"
	"mal/types"
)

func Eval(ast types.MalType, e *env.Env) types.MalType {
	switch v := ast.(type) {
	case *types.MalSymbol:
		res, ok := e.Get(*v)
		if !ok {
			fmt.Printf("Symbol '%s' not found\n", v.GetAsString())
			return v
		}
		return *res
	case *types.MalList:
		return evalList(v, e)
	default:
		return ast
	}
}

func evalList(list *types.MalList, e *env.Env) types.MalType {
	if list.IsEmpty() {
		return list
	}

	fnExpr := Eval(list.First(), e)
	switch fnExpr.GetTypeId() {
	case types.BuiltInFunction:
		evalFn := func(i types.MalType) types.MalType {
			return Eval(i, e)
		}
		return fnExpr.(core.BuiltInFn).Apply(e, evalFn, list.Tail()...)
	case types.DefinedFunction:
		panic("boom")
	default:
		args := []types.MalType{}
		for _, item := range list.Children() {
			args = append(args, Eval(item, e))
		}
		return types.NewMalList(list.GetTypeId(), args)
	}
}
