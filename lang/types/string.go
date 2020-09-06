package types

type OmmString struct {
	runel  []rune
	Length uint64
}

func (str *OmmString) FromGoType(val string) {

	var arr = make([]rune, len(val))

	for k, v := range val {
		arr[k] = rune(v)
	}

	str.runel = arr
	str.Length = uint64(len(val))
}

func (str *OmmString) FromRuneList(val []rune) {
	str.runel = val
	str.Length = uint64(len(val))
}

func (str OmmString) ToGoType() string {

	if str.runel == nil {
		return ""
	}

	var gostr string

	for _, v := range str.runel {
		gostr += string(v)
	}

	return gostr
}

func (str OmmString) ToRuneList() []rune {
	return str.runel
}

func (str OmmString) Exists(idx int64) bool {
	return str.Length != 0 && uint64(idx) < str.Length && idx >= 0
}

func (str OmmString) At(idx int64) *OmmRune {

	if idx < 0 || uint64(idx) >= str.Length {
		return nil
	}

	var gotype = str.runel[idx]
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

	for k, v := range str.runel {
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
