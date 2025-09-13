package types

type MalBool struct {
	state bool
}

func NewMalBool(state bool) *MalBool {
	return &MalBool{
		state: state,
	}
}

func (bool MalBool) GetElementTypeId() MalTypeId {
	return Boolean
}
func (bool MalBool) GetTypeId() MalTypeId {
	return bool.GetElementTypeId()
}

func (bool MalBool) GetStr(readable bool) string {
	if bool.state {
		return "true"
	}
	return "false"
}

func (bool MalBool) GetState() bool {
	return bool.state
}

func (bool MalBool) And(other MalBool) MalBool {
	return *NewMalBool(bool.state && other.state)
}

func (bool MalBool) Or(other MalBool) MalBool {
	return *NewMalBool(bool.state || other.state)
}

func (bool MalBool) Xor(other MalBool) MalBool {
	return *NewMalBool(bool.state || other.state && !(bool.state && other.state))
}

func (bool MalBool) Not() MalBool {
	return *NewMalBool(!bool.state)
}
