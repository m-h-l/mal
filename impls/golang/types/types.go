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
	DefinedFunction
	BuiltInFunction
)

type MalType interface {
	GetTypeId() MalTypeId
	GetStr() string
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

func (symbol MalSymbol) GetStr() string {
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

func (nil MalNil) GetStr() string {
	return "nil"
}

type MalFnCtx interface {
}

type MalFn interface {
	MalType
	Apply(ctx MalFnCtx, args ...MalType)
	GetName() string
}
