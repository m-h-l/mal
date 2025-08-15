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

func (str MalString) GetStr(readable bool) string {
	if readable {
		escaped := ""
		for _, ch := range str.str {
			switch ch {
			case '\n':
				escaped += "\\n"
			case '\t':
				escaped += "\\t"
			case '\r':
				escaped += "\\r"
			case '\\':
				escaped += "\\\\"
			case '"':
				escaped += "\\\""
			default:
				escaped += string(ch)
			}
		}
		return escaped
	}
	return str.str
}

func (a MalString) Append(b MalString) MalString {
	return *NewMalString(a.str + b.str)
}
