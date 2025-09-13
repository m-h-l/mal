package types

func NewMalList(children []MalType) *MalList {
	return &MalList{
		children: children,
	}
}

type MalList struct {
	children []MalType
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
	return List
}

func (list MalList) GetStr(readable bool) string {
	start, end := "(", ")"
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
