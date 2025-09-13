package evaluator

import (
	"fmt"
	"mal/types"
	"strings"
)

type EvalResult struct {
	Value    types.MalType
	Success  bool
	Continue bool
	Ast      types.MalType
	Env      *types.Env
}

func NewValue(val types.MalType, success bool) EvalResult {
	return EvalResult{
		Value:    val,
		Success:  success,
		Continue: false,
	}
}

func NewContinue(ast types.MalType, env *types.Env) EvalResult {
	return EvalResult{
		Value:    nil,
		Success:  true,
		Continue: true,
		Ast:      ast,
		Env:      env,
	}
}

func Eval(ast types.MalType, env *types.Env) (types.MalType, bool) {
	result := evalStep(ast, env)
	for result.Continue {
		result = evalStep(result.Ast, result.Env)
	}
	return result.Value, result.Success
}

func evalStep(ast types.MalType, env *types.Env) EvalResult {
	if v, ok := env.Get(*types.NewMalGenericAtom("DEBUG-EVAL")); ok {
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

	switch typedAst := ast.(type) {
	case *types.MalSymbol:
		res, ok := env.Get(*typedAst)
		if !ok && strings.HasPrefix(typedAst.GetAsString(), ":") {
			return NewValue(typedAst, true)
		}
		if fOk, _ := getSpecialForm(typedAst.GetAsString()); !ok && !fOk {
			fmt.Printf("Symbol '%s' not found\n", typedAst.GetAsString())
			return NewValue(typedAst, false)
		}
		if ok, _ := getSpecialForm(typedAst.GetAsString()); ok {
			return NewValue(typedAst, true)
		}
		return NewValue(*res, true)
	case *types.MalList:
		if typedAst.IsEmpty() {
			return NewValue(typedAst, true)
		}

		fnExpr, _ := Eval(typedAst.First(), env)

		if ok, form := getSpecialForm(fnExpr.GetStr(false)); ok {
			return form(*typedAst, env)
		}

		switch fnExpr.GetTypeId() {
		case types.Function:
			args := []types.MalType{}
			for _, arg := range typedAst.Tail() {
				r, ok := Eval(arg, env)
				if !ok {
					return NewValue(typedAst, false)
				}
				args = append(args, r)
			}
			if malFn, ok := fnExpr.(*types.MalFunction); ok && malFn.CanTCO() {
				if newAst, newEnv, isTailCall := malFn.Prepare(args...); isTailCall {
					return NewContinue(newAst, newEnv)
				}
			}
			result, success := fnExpr.(*types.MalFunction).Apply(env, args...)
			return NewValue(result, success)
		default:
			args := []types.MalType{}
			for _, item := range typedAst.Children() {
				r, ok := Eval(item, env)
				if !ok {
					return NewValue(typedAst, false)
				}
				args = append(args, r)
			}
			return NewValue(types.NewMalList(typedAst.GetTypeId(), args), true)
		}
	case *types.MalString:
		return NewValue(typedAst, true)
	default:
		return NewValue(ast, true)
	}
}

func evalIf(list types.MalList, e *types.Env) EvalResult {
	rest := list.Tail()
	if len(rest) < 2 {
		panic("boom!")
	}

	second, _ := Eval(rest[0], e)

	if len(rest) == 2 {
		if second.GetTypeId() == types.Nil {
			return NewValue(types.NewMalNil(), true)
		}
		if second.GetTypeId() == types.Boolean && !second.(*types.MalBool).GetState() {
			return NewValue(types.NewMalNil(), true)
		}
		return NewContinue(rest[1], e)
	}

	third := rest[1]
	fourth := rest[2]

	if second.GetTypeId() == types.Nil {
		return NewContinue(fourth, e)
	}
	if second.GetTypeId() == types.Boolean && !second.(*types.MalBool).GetState() {
		return NewContinue(fourth, e)
	}
	return NewContinue(third, e)
}

func evalDef(list types.MalList, e *types.Env) EvalResult {
	rest := list.Tail()
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
		return NewValue(list, false)
	}
	e.Add(*sy, r)
	return NewValue(r, true)
}

func evalLetN(list types.MalList, e *types.Env) EvalResult {
	rest := list.Tail()
	if len(rest) < 2 {
		panic("boom!")
	}
	second, third := rest[0], rest[1]
	if second.GetTypeId() != types.List && second.GetTypeId() != types.Vector {
		panic("boom!")
	}
	ne := types.NewEnv(e)
	children := second.(*types.MalList).Children()

	if len(children)%2 != 0 {
		panic("boom!")
	}

	for n := 0; n < len(children); n += 2 {
		k := *children[n].(*types.MalSymbol)
		v := children[n+1]
		r, ok := Eval(v, ne)
		if !ok {
			return NewValue(list, false)
		}
		ne.Add(k, r)
	}
	return NewContinue(third, ne)
}

func evalFnN(list types.MalList, e *types.Env) EvalResult {
	rest := list.Tail()
	first, second := rest[0], rest[1]
	if first.GetTypeId() != types.List && first.GetTypeId() != types.Vector {
		panic("puff!")
	}
	binds := first.(*types.MalList).Children()
	body := second
	fn := types.NewTCOFunction(binds, body, e)
	return NewValue(fn, true)
}

func evalDo(list types.MalList, e *types.Env) EvalResult {
	rest := list.Tail()
	if len(rest) == 0 {
		return NewValue(types.NewMalNil(), true)
	}
	for i := 0; i < len(rest)-1; i++ {
		_, ok := Eval(rest[i], e)
		if !ok {
			return NewValue(list, false)
		}
	}
	return NewContinue(rest[len(rest)-1], e)

}

func getSpecialForm(name string) (bool, func(types.MalList, *types.Env) EvalResult) {
	m := map[string]func(types.MalList, *types.Env) EvalResult{
		"if":   evalIf,
		"def!": evalDef,
		"let*": evalLetN,
		"fn*":  evalFnN,
		"do":   evalDo,
	}
	f, ok := m[name]
	return ok, f
}
