package types

func NewMalVector(children []MalType) *MalVector {
	return &MalVector{
		children: children,
	}
}

type MalVector struct {
	children []MalType
}

func (list MalVector) First() MalType {
	return list.children[0]
}

func (list MalVector) Tail() []MalType {
	return list.children[1:]
}

func (list MalVector) Children() []MalType {
	return list.children
}

func (list MalVector) IsEmpty() bool {
	return len(list.children) == 0
}

func (list MalVector) Size() int {
	return len(list.children)
}

func (list MalVector) GetTypeId() MalTypeId {
	return Vector
}

func (list MalVector) GetStr(readable bool) string {
	start, end := "[", "]"
	str := start
	for i, child := range list.children {
		str += child.GetStr(readable)
		if i < len(list.children)-1 {
			str += " "
		}
	}
	str += end
	return str
}
