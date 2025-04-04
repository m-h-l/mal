package types

import "fmt"

type MalTypeId int

const (
	Atom MalTypeId = iota
	Number
	String
	List
	Vector
	Map
)

type MalType interface {
	GetTypeId() MalTypeId
	GetStr() string
}

type MalAtom interface {
	MalType
	GetAtomTypeId() MalTypeId
}

type MalGenericAtom struct {
	atom string
}

func NewMalGenericAtom(atom string) *MalGenericAtom {
	return &MalGenericAtom{
		atom: atom,
	}
}

func (atom *MalGenericAtom) GetAtomTypeId() MalTypeId {
	return Atom
}

func (atom *MalGenericAtom) GetTypeId() MalTypeId {
	return atom.GetAtomTypeId()
}

func (atom MalGenericAtom) GetStr() string {
	return atom.atom
}

type MalNumber struct {
	num int64
}

func NewMalNumber(num int64) *MalNumber {
	return &MalNumber{
		num: num,
	}
}

func (num *MalNumber) GetAtomTypeId() MalTypeId {
	return Number
}
func (num *MalNumber) GetTypeId() MalTypeId {
	return num.GetAtomTypeId()
}

func (num MalNumber) GetStr() string {
	return fmt.Sprintf("%d", num.num)
}

type MalString struct {
	str string
}

func NewMalString(str string) *MalString {
	return &MalString{
		str: str,
	}
}

func (str *MalString) GetAtomTypeId() MalTypeId {
	return String
}

func (str *MalString) GetTypeId() MalTypeId {
	return str.GetAtomTypeId()
}

func (str MalString) GetStr() string {
	return "\"" + str.str + "\""
}

func NewMalList(kind MalTypeId, children []MalType) *MalList {
	return &MalList{
		children: children,
		kind:     kind,
	}
}

type MalList struct {
	children []MalType
	kind     MalTypeId
}

func (list *MalList) GetTypeId() MalTypeId {
	return List
}

func (list MalList) Limiters() (string, string) {
	switch list.kind {
	case List:
		return "(", ")"
	case Vector:
		return "[", "]"
	case Map:
		return "{", "}"
	default:
		panic("Unknown list type")
	}
}

func (list MalList) GetStr() string {
	start, end := list.Limiters()
	str := start
	for i, child := range list.children {
		str += child.GetStr()
		if i < len(list.children)-1 {
			str += " "
		}
	}
	str += end
	return str
}
