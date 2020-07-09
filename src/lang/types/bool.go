package types

type OmmBool struct {
  boolean *bool
}

func (b *OmmBool) FromGoType(val bool) {
  b.boolean = &val
}

func (b OmmBool) ToGoType() bool {
  return *b.boolean
}

func (_ OmmBool) ValueFunc() {}
