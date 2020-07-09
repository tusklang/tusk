package types

type OmmUndef struct {
  value struct{}
}

func (b *OmmUndef) FromGoType(val struct{}) {
  b.value = struct{}{}
}

func (b OmmUndef) ToGoType() struct{} {
  var nilv struct{}
  return nilv
}

func (_ OmmUndef) ValueFunc() {}
