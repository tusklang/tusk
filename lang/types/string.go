package types

type KaString struct {
	runel  []rune
	Length uint64
}

func (str *KaString) FromGoType(val string) {

	var arr = make([]rune, len(val))

	for k, v := range val {
		arr[k] = rune(v)
	}

	str.runel = arr
	str.Length = uint64(len(val))
}

func (str *KaString) FromRuneList(val []rune) {
	str.runel = val
	str.Length = uint64(len(val))
}

func (str KaString) ToGoType() string {

	if str.runel == nil {
		return ""
	}

	var gostr string

	for _, v := range str.runel {
		gostr += string(v)
	}

	return gostr
}

func (str KaString) ToRuneList() []rune {
	return str.runel
}

func (str KaString) Exists(idx int64) bool {
	return str.Length != 0 && uint64(idx) < str.Length && idx >= 0
}

func (str KaString) At(idx int64) *KaRune {

	if idx < 0 || uint64(idx) >= str.Length {
		return nil
	}

	var gotype = str.runel[idx]
	var karune KaRune
	karune.FromGoType(gotype)

	return &karune
}

func (str KaString) Format() string {
	return str.ToGoType()
}

func (str KaString) Type() string {
	return "string"
}

func (str KaString) TypeOf() string {
	return str.Type()
}

func (str KaString) Deallocate() {}

//Range ranges over a string
func (str KaString) Range(fn func(val1, val2 *KaType) Returner) *Returner {

	for k, v := range str.runel {
		var key KaNumber
		key.FromGoType(float64(k))
		var val KaRune
		val.FromGoType(v)

		var keykatype KaType = key
		var valkatype KaType = val
		ret := fn(&keykatype, &valkatype)

		if ret.Type == "break" {
			break
		} else if ret.Type == "return" {
			return &ret
		}
	}

	return nil
}
