package core

import (
	"fmt"
	"mal/env"
	"mal/types"
	"strings"
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

func (fn BuiltInFn) GetStr(readable bool) string {
	return "builtin<" + fn.name + ">"
}

func (fn BuiltInFn) GetTypeId() types.MalTypeId {
	return types.BuiltInFunction
}

type DefinedFn struct {
	fn func(*env.Env, ...types.MalType) types.MalType
}

func NewDefinedFn(fn func(*env.Env, ...types.MalType) types.MalType) *DefinedFn {
	return &DefinedFn{
		fn: fn,
	}
}

func (fn DefinedFn) GetAtomTypeId() types.MalTypeId {
	return types.DefinedFunction
}

func (fn DefinedFn) GetTypeId() types.MalTypeId {
	return fn.GetAtomTypeId()
}

func (fn DefinedFn) GetStr(readable bool) string {
	return "#<function>"
}

func (fn DefinedFn) Apply(e *env.Env, args ...types.MalType) (types.MalType, bool) {
	return fn.fn(e, args...), true
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

func PrintStr() BuiltInFn {
	return NewBuiltinFn("pr-str", func(e *env.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			fmt.Println("")
		}
		if len(args) == 1 {
			fmt.Println(args[0].GetStr(true))
		}
		return types.NewMalNil()
	})
}

func List() BuiltInFn {
	return NewBuiltinFn("list", func(e *env.Env, args ...types.MalType) types.MalType {
		return types.NewMalList(types.List, args)
	})
}

func IsList() BuiltInFn {
	return NewBuiltinFn("list?", func(e *env.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			return types.NewMalBool(false)
		}
		if args[0].GetTypeId() == types.List {
			return types.NewMalBool(true)
		}
		return types.NewMalBool(false)
	})
}

func Empty() BuiltInFn {
	return NewBuiltinFn("empty?", func(e *env.Env, args ...types.MalType) types.MalType {
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

func Count() BuiltInFn {
	return NewBuiltinFn("count", func(e *env.Env, args ...types.MalType) types.MalType {
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

func Equals() BuiltInFn {
	return NewBuiltinFn("=", func(e *env.Env, args ...types.MalType) types.MalType {
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

func Smaller() BuiltInFn {
	return NewBuiltinFn("<", func(e *env.Env, args ...types.MalType) types.MalType {
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

func SmallerOrEqual() BuiltInFn {
	return NewBuiltinFn("<=", func(e *env.Env, args ...types.MalType) types.MalType {
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

func Bigger() BuiltInFn {
	return NewBuiltinFn(">", func(e *env.Env, args ...types.MalType) types.MalType {
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

func BiggerOrEqual() BuiltInFn {
	return NewBuiltinFn(">=", func(e *env.Env, args ...types.MalType) types.MalType {
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
		result = append(result, arg.GetStr(readable)) // Use raw representation
	}
	return strings.Join(result, separator)
}

func Prn() BuiltInFn {
	return NewBuiltinFn("prn", func(e *env.Env, args ...types.MalType) types.MalType {
		result := concatArgs(args, " ", true)
		fmt.Println(result)
		return types.NewMalNil()
	})
}

func Println() BuiltInFn {
	return NewBuiltinFn("println", func(e *env.Env, args ...types.MalType) types.MalType {
		result := concatArgs(args, " ", false)
		fmt.Println(result)
		return types.NewMalNil()
	})
}

func PrStr() BuiltInFn {
	return NewBuiltinFn("pr-str", func(e *env.Env, args ...types.MalType) types.MalType {
		result := concatArgs(args, " ", true)
		return types.NewMalString(result) // Return joined string
	})
}

func Str() BuiltInFn {
	return NewBuiltinFn("str", func(e *env.Env, args ...types.MalType) types.MalType {
		result := concatArgs(args, "", false) // Use raw representation
		return types.NewMalString(result)
	})
}

func Not() BuiltInFn {
	return NewBuiltinFn("not", func(e *env.Env, args ...types.MalType) types.MalType {
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

func AddCoreToEnv(e *env.Env) {
	builtins := []BuiltInFn{
		Plus(),
		Multiply(),
		Subtract(),
		Divide(),
		PrintStr(),
		List(),
		IsList(),
		Empty(),
		Count(),
		Equals(),
		Smaller(),
		SmallerOrEqual(),
		Bigger(),
		BiggerOrEqual(),
		Not(),
		Str(),
		PrStr(),
		Prn(),
		Println(),
	}

	for _, builtin := range builtins {
		e.Add(*types.NewMalGenericAtom(builtin.name), builtin)
	}
}
