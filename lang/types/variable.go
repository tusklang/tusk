package types

type KaVar struct {
	Name  string
	Value *KaType
}

func (v KaVar) Format() string {
	return v.Name
}
