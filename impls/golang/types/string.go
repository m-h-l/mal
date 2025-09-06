package types

import (
	"strings"
)

type MalString struct {
	str string
}

func NewMalString(str string) *MalString {
	return &MalString{
		str: parse(str),
	}
}

func (str MalString) GetAtomTypeId() MalTypeId {
	return String
}

func (str MalString) GetTypeId() MalTypeId {
	return str.GetAtomTypeId()
}

func (ms MalString) GetStr(readable bool) string {
	if readable {
		var sb strings.Builder
		sb.WriteByte('"')
		for _, ch := range ms.str {
			switch ch {
			case '\\':
				sb.WriteString("\\\\")
			case '"':
				sb.WriteString("\\\"")
			case '\n':
				sb.WriteString("\\n")
			default:
				sb.WriteRune(ch)
			}
		}
		sb.WriteByte('"')
		return sb.String()
	}
	return ms.str
}

func (a MalString) Append(b MalString) MalString {
	return *NewMalString(a.str + b.str)
}

func parse(str string) string {
	str = strings.ReplaceAll(str, "\\\"", "\"")
	str = strings.ReplaceAll(str, "\\\\", "\\")
	str = strings.ReplaceAll(str, "\\n", "\n")
	return str
}
