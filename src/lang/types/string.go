package types

type OmmString struct {
  String    *string
  Length     uint64
}

func (str *OmmString) FromGoType(val string) {
  str.String = &val
  str.Length = uint64(len(val))
}

func (str OmmString) ToGoType() string {

  if str.String == nil {
    return ""
  }

  return *str.String
}

func (str OmmString) Exists(idx int64) bool {
  return str.Length != 0 && uint64(idx) < str.Length && idx >= 0
}

func (str OmmString) At(idx int64) *OmmRune {

  var gotype = rune((*str.String)[idx])
  var ommrune OmmRune
  ommrune.FromGoType(gotype)

  return &ommrune
}

func (str OmmString) Format() string {
  return str.ToGoType()
}

func (str OmmString) Type() string {
  return "string"
}

func (str OmmString) TypeOf() string {
  return str.Type()
}
