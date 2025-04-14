package types

type MalString struct {
	str string
}

func NewMalString(str string) *MalString {
	return &MalString{
		str: str,
	}
}

func (str MalString) GetAtomTypeId() MalTypeId {
	return String
}

func (str MalString) GetTypeId() MalTypeId {
	return str.GetAtomTypeId()
}

func (str MalString) GetStr() string {
	return "\"" + str.str + "\""
}
