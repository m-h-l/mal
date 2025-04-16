package types

import "fmt"

type MalNumber struct {
	num int64
}

func NewMalNumber(num int64) *MalNumber {
	return &MalNumber{
		num: num,
	}
}

func (num MalNumber) GetAtomTypeId() MalTypeId {
	return Number
}
func (num MalNumber) GetTypeId() MalTypeId {
	return num.GetAtomTypeId()
}

func (num MalNumber) GetAsInt() int64 {
	return num.num
}

func (num MalNumber) GetStr() string {
	return fmt.Sprintf("%d", num.num)
}

func (num MalNumber) Add(other MalNumber) MalNumber {
	return *NewMalNumber(num.num + other.num)
}

func (num MalNumber) Multiply(other MalNumber) MalNumber {
	return *NewMalNumber(num.num * other.num)
}

func (num MalNumber) Minus(other MalNumber) MalNumber {
	return *NewMalNumber(num.num - other.num)
}

func (num MalNumber) Divide(other MalNumber) MalNumber {
	return *NewMalNumber(num.num / other.num)
}
