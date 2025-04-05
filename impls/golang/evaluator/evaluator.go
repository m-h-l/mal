package evaluator

import (
	"mal/env"
	"mal/types"
)

func Eval(ast types.MalType, env env.Env) types.MalType {
	switch ast.GetTypeId() {
	case types.List:

	default:
		return ast
	}
}
