package types

func NewMalMap(children []MalType) *MalMap {
	return &MalMap{
		children: children,
	}
}

type MalMap struct {
	children []MalType
}

func (list MalMap) First() MalType {
	return list.children[0]
}

func (list MalMap) Tail() []MalType {
	return list.children[1:]
}

func (list MalMap) Children() []MalType {
	return list.children
}

func (list MalMap) IsEmpty() bool {
	return len(list.children) == 0
}

func (list MalMap) Size() int {
	return len(list.children)
}

func (list MalMap) GetTypeId() MalTypeId {
	return Map
}

func (list MalMap) GetStr(readable bool) string {
	start, end := "{", "}"
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
