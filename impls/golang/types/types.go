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
	Atom
)

type MalType interface {
	GetTypeId() MalTypeId
	GetStr(bool) string
}

type MalElement interface {
	MalType
	GetElementTypeId() MalTypeId
}

type MalSeq interface {
	MalType
	Children() []MalType
	First() MalType
	Tail() []MalType
}

type MalSymbol struct {
	symbol string
}

func NewMalSymbol(symbol string) *MalSymbol {
	return &MalSymbol{
		symbol: symbol,
	}
}

func (symbol MalSymbol) GetAsString() string {
	return symbol.symbol
}

func (symbol MalSymbol) GetElementTypeId() MalTypeId {
	return Symbol
}

func (symbol MalSymbol) GetTypeId() MalTypeId {
	return symbol.GetElementTypeId()
}

func (symbol MalSymbol) GetStr(readable bool) string {
	return symbol.symbol
}

type MalNil struct {
}

func NewMalNil() *MalNil {
	return &MalNil{}
}

func (nil MalNil) GetElementTypeId() MalTypeId {
	return Nil
}

func (nil MalNil) GetTypeId() MalTypeId {
	return nil.GetElementTypeId()
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

func (env *Env) GetRoot() *Env {
	root := env
	for root.outer != nil {
		root = root.outer
	}
	return root
}
