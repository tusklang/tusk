package types

type TuskString struct {
	runel  []rune
	length uint64
}

func (str *TuskString) FromGoType(val string) {

	var arr = make([]rune, len(val))

	for k, v := range val {
		arr[k] = rune(v)
	}

	str.runel = arr
	str.length = uint64(len(val))
}

func (str *TuskString) FromRuneList(val []rune) {
	str.runel = val
	str.length = uint64(len(val))
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
	return str.length != 0 && uint64(idx) < str.length && idx >= 0
}

func (str TuskString) At(idx int64) *TuskRune {

	if idx < 0 || uint64(idx) >= str.length {
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

func (str TuskString) Length() uint64 {
	return str.length
}

func (str TuskString) Type() string {
	return "string"
}

func (str TuskString) TypeOf() string {
	return str.Type()
}

func (str TuskString) Deallocate() {}

func (str TuskString) Clone() *TuskType {
	var tmp = str.ToRuneList() //convert it to a go type
	var tuskstr TuskString
	tuskstr.FromRuneList(append(tmp, []rune{}...)) //clone tmp
	var returner TuskType = tuskstr
	return &returner
}

//Range ranges over a string
func (str TuskString) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {

	for k, v := range str.runel {
		var key TuskInt
		key.FromGoType(int64(k))
		var val TuskRune
		val.FromGoType(v)

		var keykatype TuskType = key
		var valkatype TuskType = val
		ret, e := fn(&keykatype, &valkatype)

		if e != nil {
			return nil, e
		}

		if ret.Type == "break" {
			break
		} else if ret.Type == "return" {
			return &ret, nil
		}
	}

	return nil, nil
}
