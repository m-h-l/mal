package types

type MalAtom struct {
	value *MalType
}

func NewMalAtom(value *MalType) *MalAtom {
	return &MalAtom{
		value: value,
	}
}

func (atom MalAtom) GetElementTypeId() MalTypeId {
	return Atom
}
func (atom MalAtom) GetTypeId() MalTypeId {
	return atom.GetElementTypeId()
}

func (atom MalAtom) GetStr(readable bool) string {
	if readable {
		return "(atom " + (*atom.value).GetStr(readable) + ")"
	}
	return "#<atom>"
}

func (atom *MalAtom) Set(value *MalType) {
	atom.value = value
}

func (atom MalAtom) Deref() *MalType {
	return atom.value
}
