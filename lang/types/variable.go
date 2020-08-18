package types

type OmmVar struct {
  Name      string
  Value    *OmmType
}

func (v OmmVar) Format() string {
  return v.Name
}
