package core

import (
	"fmt"
	"mal/types"
	"strings"
)

func Plus() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
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

func Multiply() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
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

func Subtract() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
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

func Divide() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
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

func List() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		return types.NewMalList(types.List, args)
	})
}

func IsList() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			return types.NewMalBool(false)
		}
		if args[0].GetTypeId() == types.List {
			return types.NewMalBool(true)
		}
		return types.NewMalBool(false)
	})
}

func Empty() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			panic("boom!")
		}
		if args[0].GetTypeId() != types.List && args[0].GetTypeId() != types.Vector {
			panic("boom!")
		}
		list := args[0].(*types.MalList)
		if len(list.Children()) == 0 {
			return types.NewMalBool(true)
		}
		return types.NewMalBool(false)
	})
}

func Count() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			panic("boom!")
		}
		switch args[0].GetTypeId() {
		case types.List, types.Vector:
			list := args[0].(*types.MalList)
			return types.NewMalNumber(int64(len(list.Children())))
		default:
			return types.NewMalNumber(0)
		}
	})
}

func eq(a types.MalType, b types.MalType) bool {
	if (a.GetTypeId() == types.List || a.GetTypeId() == types.Vector) && (b.GetTypeId() == types.List || b.GetTypeId() == types.Vector) {
		listA := a.(*types.MalList)
		listB := b.(*types.MalList)
		if len(listA.Children()) != len(listB.Children()) {
			return false
		}
		for i := 0; i < len(listA.Children()); i++ {
			if !eq(listA.Children()[i], listB.Children()[i]) {
				return false
			}
		}
		return true
	}
	if a.GetTypeId() != b.GetTypeId() {
		return false
	}
	if a.GetStr(false) != b.GetStr(false) {
		return false
	}
	return true
}

func Equals() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			return types.NewMalBool(true)
		}
		if len(args) == 1 {
			return types.NewMalBool(false)
		}
		first := args[0]
		second := args[1]
		return types.NewMalBool(eq(first, second))
	})
}

func Smaller() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		first := args[0]
		second := args[1]
		if first.GetTypeId() != types.Number || second.GetTypeId() != types.Number {
			panic("boom!")
		}

		a := first.(*types.MalNumber)
		b := second.(*types.MalNumber)

		if a.GetAsInt() < b.GetAsInt() {
			return types.NewMalBool(true)
		}
		return types.NewMalBool(false)
	})
}

func SmallerOrEqual() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		first := args[0]
		second := args[1]
		if first.GetTypeId() != types.Number || second.GetTypeId() != types.Number {
			panic("boom!")
		}

		a := first.(*types.MalNumber)
		b := second.(*types.MalNumber)

		if a.GetAsInt() <= b.GetAsInt() {
			return types.NewMalBool(true)
		}
		return types.NewMalBool(false)
	})
}

func Bigger() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		first := args[0]
		second := args[1]
		if first.GetTypeId() != types.Number || second.GetTypeId() != types.Number {
			panic("boom!")
		}

		a := first.(*types.MalNumber)
		b := second.(*types.MalNumber)

		if a.GetAsInt() > b.GetAsInt() {
			return types.NewMalBool(true)
		}
		return types.NewMalBool(false)
	})
}

func BiggerOrEqual() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		first := args[0]
		second := args[1]
		if first.GetTypeId() != types.Number || second.GetTypeId() != types.Number {
			panic("boom!")
		}

		a := first.(*types.MalNumber)
		b := second.(*types.MalNumber)

		if a.GetAsInt() >= b.GetAsInt() {
			return types.NewMalBool(true)
		}
		return types.NewMalBool(false)
	})
}

func concatArgs(args []types.MalType, separator string, readable bool) string {
	result := []string{}

	for _, arg := range args {
		result = append(result, arg.GetStr(readable))
	}
	return strings.Join(result, separator)

}

func Prn() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		result := concatArgs(args, " ", true)
		fmt.Println(result)
		return types.NewMalNil()
	})
}

func Println() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		result := concatArgs(args, " ", false)
		fmt.Println(result)
		return types.NewMalNil()
	})
}

func PrStr() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		result := concatArgs(args, " ", true)
		return types.NewMalString(result)
	})
}

func Str() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		result := concatArgs(args, "", false)
		return types.NewMalString(result)
	})
}

func Not() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			fmt.Println("")
		} else {
			args[0].GetTypeId() // Ensure the type is checked
			if args[0].GetTypeId() == types.Nil || (args[0].GetTypeId() == types.Boolean && !args[0].(*types.MalBool).GetState()) {
				return types.NewMalBool(true)
			}
			return types.NewMalBool(false)
		}
		return types.NewMalNil()
	})
}

func AddCoreToEnv(e *types.Env) {
	builtins := map[string]*types.MalFunction{
		"+":       Plus(),
		"*":       Multiply(),
		"-":       Subtract(),
		"/":       Divide(),
		"list":    List(),
		"list?":   IsList(),
		"empty?":  Empty(),
		"count":   Count(),
		"=":       Equals(),
		"<":       Smaller(),
		"<=":      SmallerOrEqual(),
		">":       Bigger(),
		">=":      BiggerOrEqual(),
		"not":     Not(),
		"str":     Str(),
		"pr-str":  PrStr(),
		"prn":     Prn(),
		"println": Println(),
	}

	for name, builtin := range builtins {
		e.Add(*types.NewMalGenericAtom(name), builtin)
	}
}
