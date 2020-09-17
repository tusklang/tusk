package types

type TuskVar struct {
	Name  string
	Value *TuskType
}

func (v TuskVar) Format() string {
	return v.Name
}
