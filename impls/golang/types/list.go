package types

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

func (list MalList) First() MalType {
	return list.children[0]
}

func (list MalList) Tail() []MalType {
	return list.children[1:]
}

func (list MalList) Children() []MalType {
	return list.children
}

func (list MalList) IsEmpty() bool {
	return len(list.children) == 0
}

func (list MalList) Size() int {
	return len(list.children)
}

func (list MalList) GetTypeId() MalTypeId {
	return list.kind
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
