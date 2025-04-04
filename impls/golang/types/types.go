package types

type MalTypeId int

const (
	Atom MalTypeId = iota
	List
	Vector
	Map
)

type MalType interface {
	GetTypeId() MalTypeId
	GetStr() string
}

func NewMalAtom(atom string) *MalAtom {
	return &MalAtom{
		atom: atom,
	}
}

type MalAtom struct {
	atom string
}

func (atom *MalAtom) GetTypeId() MalTypeId {
	return Atom
}

func (atom MalAtom) GetStr() string {
	return atom.atom
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
