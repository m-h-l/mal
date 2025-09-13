package types

type MalTypeId int

const (
	Symbol MalTypeId = iota
	Nil
	Boolean
	Number
	String
	List
	Vector
	Map
	Function
)

type MalType interface {
	GetTypeId() MalTypeId
	GetStr(bool) string
}

type MalAtom interface {
	MalType
	GetAtomTypeId() MalTypeId
}

type MalSymbol struct {
	symbol string
}

func NewMalGenericAtom(symbol string) *MalSymbol {
	return &MalSymbol{
		symbol: symbol,
	}
}

func (symbol MalSymbol) GetAsString() string {
	return symbol.symbol
}

func (symbol MalSymbol) GetAtomTypeId() MalTypeId {
	return Symbol
}

func (symbol MalSymbol) GetTypeId() MalTypeId {
	return symbol.GetAtomTypeId()
}

func (symbol MalSymbol) GetStr(readable bool) string {
	return symbol.symbol
}

type MalNil struct {
}

func NewMalNil() *MalNil {
	return &MalNil{}
}

func (nil MalNil) GetAtomTypeId() MalTypeId {
	return Nil
}

func (nil MalNil) GetTypeId() MalTypeId {
	return nil.GetAtomTypeId()
}

func (nil MalNil) GetStr(readable bool) string {
	return "nil"
}

type Env struct {
	outer *Env
	data  map[string]MalType
}

func NewEnv(outter *Env) *Env {
	return &Env{
		outer: outter,
		data:  map[string]MalType{},
	}
}

func (env *Env) Add(symbol MalSymbol, value MalType) {
	env.data[symbol.GetAsString()] = value
}

func (env *Env) Get(symbol MalSymbol) (*MalType, bool) {
	v, ok := env.data[symbol.GetAsString()]
	if ok {
		return &v, true
	}

	if env.outer != nil {
		return env.outer.Get(symbol)
	}

	return nil, false
}
