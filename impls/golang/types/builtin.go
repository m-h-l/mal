package types

func Plus() BuiltInFn {
	return NewBuiltinFn("+", func(args ...MalType) MalType {
		acc := NewMalNumber(0)
		for _, arg := range args {
			if arg.GetTypeId() != Number {
				panic("boom!")
			}
			n := acc.Add(*arg.(*MalNumber))
			acc = &n
		}
		return acc
	})
}

func Multiply() BuiltInFn {
	return NewBuiltinFn("*", func(args ...MalType) MalType {
		acc := NewMalNumber(1)
		for _, arg := range args {
			if arg.GetTypeId() != Number {
				panic("boom!")
			}
			n := acc.Multiply(*arg.(*MalNumber))
			acc = &n
		}
		return acc
	})
}

func Subtract() BuiltInFn {
	return NewBuiltinFn("-", func(args ...MalType) MalType {
		for _, arg := range args {
			if arg.GetTypeId() != Number {
				panic("boom!")
			}
		}
		if len(args) == 0 {
			panic("boom!")
		}
		if len(args) == 1 {
			n := NewMalNumber(-1).Multiply(*args[0].(*MalNumber))
			return &n
		}

		acc := *args[0].(*MalNumber)

		for _, arg := range args[1:] {
			acc = acc.Minus(*arg.(*MalNumber))
		}
		return &acc
	})
}

func Divide() BuiltInFn {
	return NewBuiltinFn("/", func(args ...MalType) MalType {
		for _, arg := range args {
			if arg.GetTypeId() != Number {
				panic("boom!")
			}
		}
		if len(args) < 2 {
			panic("boom!")
		}

		acc := *args[0].(*MalNumber)

		for _, arg := range args[1:] {
			acc = acc.Divide(*arg.(*MalNumber))
		}
		return &acc
	})
}

type Fn interface {
	MalType
	HasName(name string) bool
	Apply(args ...MalType) MalType
}

type BuiltInFn struct {
	name string
	fn   func(...MalType) MalType
}

func NewBuiltinFn(name string, fn func(...MalType) MalType) BuiltInFn {
	return BuiltInFn{
		name: name,
		fn:   fn,
	}
}

func (fn BuiltInFn) HasName(name string) bool {
	return fn.name == name
}

func (fn BuiltInFn) Apply(args ...MalType) MalType {
	return fn.fn(args...)
}

func (fn BuiltInFn) GetStr() string {
	return "builtin<" + fn.name + ">"
}

func (fn BuiltInFn) GetTypeId() MalTypeId {
	return BuiltInFunction
}
