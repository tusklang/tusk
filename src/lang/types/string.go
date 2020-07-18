package types

type OmmString struct {
  String *[]*OmmRune
  Length     uint64
}

func (str *OmmString) FromGoType(val string) {

  var rarray []*OmmRune

  for _, v := range val {
    var r OmmRune
    r.FromGoType(v)
    rarray = append(rarray, &r)
  }

  str.String = &rarray
  str.Length = uint64(len(val))
}

func (str OmmString) ToGoType() string {
  var gostr string

  for _, v := range *str.String {
    gostr+=string((*v).ToGoType())
  }

  return gostr
}

func (str OmmString) Exists(idx int64) bool {
  return str.Length != 0 && uint64(idx) < str.Length && idx >= 0
}

func (str OmmString) At(idx int64) *OmmRune {
  return (*str.String)[idx]
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
