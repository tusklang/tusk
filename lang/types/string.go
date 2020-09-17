package types

type TuskString struct {
	runel  []rune
	Length uint64
}

func (str *TuskString) FromGoType(val string) {

	var arr = make([]rune, len(val))

	for k, v := range val {
		arr[k] = rune(v)
	}

	str.runel = arr
	str.Length = uint64(len(val))
}

func (str *TuskString) FromRuneList(val []rune) {
	str.runel = val
	str.Length = uint64(len(val))
}

func (str TuskString) ToGoType() string {

	if str.runel == nil {
		return ""
	}

	var gostr string

	for _, v := range str.runel {
		gostr += string(v)
	}

	return gostr
}

func (str TuskString) ToRuneList() []rune {
	return str.runel
}

func (str TuskString) Exists(idx int64) bool {
	return str.Length != 0 && uint64(idx) < str.Length && idx >= 0
}

func (str TuskString) At(idx int64) *TuskRune {

	if idx < 0 || uint64(idx) >= str.Length {
		return nil
	}

	var gotype = str.runel[idx]
	var karune TuskRune
	karune.FromGoType(gotype)

	return &karune
}

func (str TuskString) Format() string {
	return str.ToGoType()
}

func (str TuskString) Type() string {
	return "string"
}

func (str TuskString) TypeOf() string {
	return str.Type()
}

func (str TuskString) Deallocate() {}

//Range ranges over a string
func (str TuskString) Range(fn func(val1, val2 *TuskType) Returner) *Returner {

	for k, v := range str.runel {
		var key TuskNumber
		key.FromGoType(float64(k))
		var val TuskRune
		val.FromGoType(v)

		var keykatype TuskType = key
		var valkatype TuskType = val
		ret := fn(&keykatype, &valkatype)

		if ret.Type == "break" {
			break
		} else if ret.Type == "return" {
			return &ret
		}
	}

	return nil
}
