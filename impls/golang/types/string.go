package types

import (
	"strings"
)

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

func Parse(str string) string {
	if len(str) == 0 {
		return ""
	}
	if len(str) == 1 {
		return str
	}

	if str[0] == '\\' && str[1] == '\\' {
		return "\\" + Parse(str[2:])
	}
	if str[0] == '\\' && str[1] == 'n' {
		return "\n" + Parse(str[2:])
	}
	if str[0] == '\\' && str[1] == '"' {
		return "\"" + Parse(str[2:])
	}
	return string(str[0]) + Parse(str[1:])
}
