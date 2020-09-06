package types

type OmmString struct {
	String []rune
	Length uint64
}

func (str *OmmString) FromGoType(val string) {

	var arr = make([]rune, len(val))

	for k, v := range val {
		arr[k] = rune(v)
	}

	str.String = arr
	str.Length = uint64(len(val))
}

func (str *OmmString) FromRuneList(val []rune) {
	str.String = val
	str.Length = uint64(len(val))
}

func (str OmmString) ToGoType() string {

	if str.String == nil {
		return ""
	}

	var gostr string

	for _, v := range str.String {
		gostr += string(v)
	}

	return gostr
}

func (str OmmString) ToRuneList() []rune {
	return str.String
}

func (str OmmString) Exists(idx int64) bool {
	return str.Length != 0 && uint64(idx) < str.Length && idx >= 0
}

func (str OmmString) At(idx int64) *OmmRune {

	var gotype = str.String[idx]
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

func (str OmmString) Deallocate() {}

//Range ranges over a string
func (str OmmString) Range(fn func(val1, val2 *OmmType) Returner) *Returner {

	for k, v := range str.String {
		var key OmmNumber
		key.FromGoType(float64(k))
		var val OmmRune
		val.FromGoType(v)

		var keyommtype OmmType = key
		var valommtype OmmType = val
		ret := fn(&keyommtype, &valommtype)

		if ret.Type == "break" {
			break
		} else if ret.Type == "return" {
			return &ret
		}
	}

	return nil
}
