package types

type OmmString struct {
  string *[]*OmmRune
  Length     uint64
}

func (str *OmmString) FromGoType(val string) {

  var rarray []*OmmRune

  for _, v := range val {
    var r OmmRune
    r.FromGoType(v)
    rarray = append(rarray, &r)
  }

  str.string = &rarray
  str.Length = uint64(len(val))
}

func (str OmmString) ToGoType() string {
  var gostr string

  for _, v := range *str.string {
    gostr+=string((*v).ToGoType())
  }

  return gostr
}

func (str OmmString) Exists(idx int64) bool {
  return str.Length != 0 && uint64(idx) < str.Length && idx >= 0
}

func (str OmmString) At(idx int64) *OmmRune {
  return (*str.string)[idx]
}

func (str OmmString) Format() string {
  return str.ToGoType()
}

func (arr OmmString) Type() string {
  return "string"
}