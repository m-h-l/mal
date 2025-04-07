package evaluator

import (
	"fmt"
	"mal/env"
	"mal/types"
)

func Eval(ast types.MalType, e env.Env) types.MalType {
	switch v := ast.(type) {
	case *types.MalSymbol:
		res, ok := e.FindSymbol(*v)
		if !ok {
			fmt.Println("Symbol '%s' not found", v)
			return v

		}
		return *res
	case *types.MalList:
		return evalList(v, e)
	default:
		// literals like numbers, strings, etc.
		return ast
	}
}

func evalList(list *types.MalList, e env.Env) types.MalType {
	if list.IsEmpty() {
		return list
	}

	fnExpr := Eval(list.First(), e)
	fn, ok := fnExpr.(types.Fn)
	if !ok {
		args := []types.MalType{}
		for _, item := range list.Children() {
			args = append(args, Eval(item, e))
		}
		return types.NewMalList(list.GetTypeId(), args)
	}

	args := []types.MalType{}
	for _, item := range list.Tail() {
		args = append(args, Eval(item, e))
	}
	return fn.Apply(args...)
}
