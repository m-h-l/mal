package core

import (
	"fmt"
	"mal/evaluator"
	"mal/parser"
	"mal/reader"
	"mal/types"
	"os"
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
		return types.NewMalList(args)
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
		list := args[0].(types.MalSeq)
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
			list := args[0].(types.MalSeq)
			return types.NewMalNumber(int64(len(list.Children())))
		default:
			return types.NewMalNumber(0)
		}
	})
}

func eq(a types.MalType, b types.MalType) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if (a.GetTypeId() == types.List || a.GetTypeId() == types.Vector) &&
		(b.GetTypeId() == types.List || b.GetTypeId() == types.Vector) {
		listA := a.(types.MalSeq)
		listB := b.(types.MalSeq)

		childrenA := listA.Children()
		childrenB := listB.Children()

		if len(childrenA) != len(childrenB) {
			return false
		}

		for i := 0; i < len(childrenA); i++ {
			if !eq(childrenA[i], childrenB[i]) {
				return false
			}
		}
		return true
	}

	if a.GetTypeId() != b.GetTypeId() {
		return false
	}

	aStr := a.GetStr(false)
	bStr := b.GetStr(false)
	result := aStr == bStr
	return result
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

func Eval() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			panic("boom!")
		}
		ast := args[0]
		rootEnv := e.GetRoot()
		result, ok := evaluator.Eval(ast, rootEnv)
		if !ok {
			panic("boom!")
		}
		return result
	})
}

func ReadString() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			panic("boom!")
		}
		if args[0].GetTypeId() != types.String {
			panic("boom!")
		}
		input := args[0].(*types.MalString).GetStr(false)
		r := reader.NewReader(input)
		ast, ok := parser.Parse(r)
		if ok {
			return ast
		}
		return types.NewMalNil()
	})
}

func Slurp() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			panic("boom!")
		}
		if args[0].GetTypeId() != types.String {
			panic("boom!")
		}
		filename := args[0].(*types.MalString).GetStr(false)
		data, err := os.ReadFile(filename)
		if err != nil {
			panic("boom!")
		}
		return types.NewMalString(string(data))
	})
}

func Atom() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			panic("boom!")
		}
		return types.NewMalAtom(&args[0])
	})
}

func IsAtom() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			return types.NewMalBool(false)
		}
		if args[0].GetTypeId() == types.Atom {
			return types.NewMalBool(true)
		}
		return types.NewMalBool(false)
	})
}

func Deref() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			panic("boom!")
		}
		if args[0].GetTypeId() != types.Atom {
			panic("boom!")
		}
		atom := args[0].(*types.MalAtom)
		return *(atom.Deref())
	})
}

func LoadFile() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) == 0 {
			panic("boom!")
		}
		if args[0].GetTypeId() != types.String {
			panic("boom!")
		}
		filename := args[0].(*types.MalString).GetStr(false)
		data, err := os.ReadFile(filename)
		if err != nil {
			panic("boom!")
		}
		input := string(data)
		r := reader.NewReader(input)

		for {
			ast, ok := parser.Parse(r)
			if !ok {
				// Reached end of input
				break
			}

			_, ok = evaluator.Eval(ast, e)
			if !ok {
				panic("boom!")
			}
		}

		return types.NewMalNil()
	})
}

func Reset() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) < 2 {
			panic("boom!")
		}
		if args[0].GetTypeId() != types.Atom {
			panic("boom!")
		}
		atom := args[0].(*types.MalAtom)
		newValue := args[1]
		atom.Set(&newValue)
		return newValue
	})
}

func Swap() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) < 2 {
			panic("boom!")
		}
		if args[0].GetTypeId() != types.Atom {
			panic("boom!")
		}
		atom := args[0].(*types.MalAtom)
		fun := args[1]
		if fun.GetTypeId() != types.Function {
			panic("boom!")
		}
		f := fun.(*types.MalFunction)
		sArgs := append([]types.MalType{*(atom.Deref())}, args[2:]...)

		newValue, ok := evaluator.Eval(
			types.NewMalList(append([]types.MalType{f}, sArgs...)),
			e,
		)

		if ok {
			atom.Set(&newValue)
			return newValue
		}
		panic("boom!")
	})
}

func Cons() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) != 2 {
			panic("boom!")
		}

		element := args[0]
		sequence := args[1]

		// Handle nil (empty list)
		if sequence == nil || sequence.GetTypeId() == types.Nil {
			return types.NewMalList([]types.MalType{element})
		}

		// Handle lists and vectors
		if sequence.GetTypeId() == types.List || sequence.GetTypeId() == types.Vector {
			seq := sequence.(types.MalSeq)
			newElements := make([]types.MalType, len(seq.Children())+1)
			newElements[0] = element
			copy(newElements[1:], seq.Children())
			return types.NewMalList(newElements)
		}

		panic("boom!")
	})
}

func Concat() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		result := []types.MalType{}

		for _, arg := range args {
			if arg == nil || arg.GetTypeId() == types.Nil {
				// Empty list/nil contributes nothing
				continue
			}

			if arg.GetTypeId() == types.List || arg.GetTypeId() == types.Vector {
				seq := arg.(types.MalSeq)
				result = append(result, seq.Children()...)
			} else {
				panic("boom!")
			}
		}

		return types.NewMalList(result)
	})
}

func Vec() *types.MalFunction {
	return types.NewFunction(func(e *types.Env, args ...types.MalType) types.MalType {
		if len(args) != 1 {
			panic("boom!")
		}

		arg := args[0]

		// If it's already a vector, return it as-is
		if arg.GetTypeId() == types.Vector {
			return arg
		}

		// Handle nil (empty list) -> empty vector
		if arg == nil || arg.GetTypeId() == types.Nil {
			return types.NewMalVector([]types.MalType{})
		}

		// Convert list to vector
		if arg.GetTypeId() == types.List {
			list := arg.(types.MalSeq)
			return types.NewMalVector(list.Children())
		}

		panic("boom!")
	})
}

func AddCoreToEnv(e *types.Env) {
	builtins := map[string]*types.MalFunction{
		"+":           Plus(),
		"*":           Multiply(),
		"-":           Subtract(),
		"/":           Divide(),
		"list":        List(),
		"list?":       IsList(),
		"empty?":      Empty(),
		"count":       Count(),
		"=":           Equals(),
		"<":           Smaller(),
		"<=":          SmallerOrEqual(),
		">":           Bigger(),
		">=":          BiggerOrEqual(),
		"not":         Not(),
		"str":         Str(),
		"pr-str":      PrStr(),
		"prn":         Prn(),
		"println":     Println(),
		"eval":        Eval(),
		"read-string": ReadString(),
		"slurp":       Slurp(),
		"atom":        Atom(),
		"atom?":       IsAtom(),
		"swap!":       Swap(),
		"reset!":      Reset(),
		"deref":       Deref(),
		"load-file":   LoadFile(),
		"cons":        Cons(),
		"concat":      Concat(),
		"vec":         Vec(),
	}

	for name, builtin := range builtins {
		e.Add(*types.NewMalSymbol(name), builtin)
	}
}
func SetupArgv(e *types.Env, args []string) {
	malArgs := make([]types.MalType, len(args))
	for i, arg := range args {
		malArgs[i] = types.NewMalString(arg)
	}

	argv := types.NewMalList(malArgs)
	e.Add(*types.NewMalSymbol("*ARGV*"), argv)
}
