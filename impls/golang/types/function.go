package types

type MalFunction struct {
	fn         func(*Env, ...MalType) MalType
	binds      []MalType
	body       MalType
	closureEnv *Env
	isBuiltin  bool
}

func NewFunction(fn func(*Env, ...MalType) MalType) *MalFunction {
	return &MalFunction{
		fn:        fn,
		isBuiltin: true,
	}
}

func NewTCOFunction(binds []MalType, body MalType, closureEnv *Env) *MalFunction {
	return &MalFunction{
		binds:      binds,
		body:       body,
		closureEnv: closureEnv,
		isBuiltin:  false,
	}
}

func (f *MalFunction) CanTCO() bool {
	return !f.isBuiltin
}

func (fn MalFunction) GetAtomTypeId() MalTypeId {
	return Function
}

func (fn MalFunction) GetTypeId() MalTypeId {
	return fn.GetAtomTypeId()
}

func (fn MalFunction) GetStr(readable bool) string {
	return "#<function>"
}

func (fn MalFunction) Prepare(args ...MalType) (MalType, *Env, bool) {
	if fn.isBuiltin {
		return nil, nil, false
	}

	newEnv := NewEnv(fn.closureEnv)

	for i, param := range fn.binds {
		if param.(*MalSymbol).GetAsString() == "&" {
			bind := fn.binds[i+1].(*MalSymbol)
			newEnv.Add(*bind, NewMalList(List, args[i:]))
			break
		}
		if i < len(args) {
			newEnv.Add(*param.(*MalSymbol), args[i])
		}
	}
	return fn.body, newEnv, true
}

func (fn MalFunction) Apply(e *Env, args ...MalType) (MalType, bool) {
	if fn.isBuiltin {
		return fn.fn(e, args...), true
	}
	panic("TCO functions cannot be applied directly")
}
