package types

type KaArray struct {
	Array  []*KaType
	Length uint64
}

func (arr KaArray) At(idx int64) *KaType {

	length := arr.Length

	if uint64(idx) >= length || idx < 0 {
		var undef KaType = KaUndef{}
		return &undef
	}

	return arr.Array[idx]
}

func (arr KaArray) Exists(idx int64) bool {
	return arr.Length != 0 && uint64(idx) < arr.Length && idx >= 0
}

func (arr *KaArray) PushBack(val KaType) {
	arr.Length++
	arr.Array = append(arr.Array, &val)
}

func (arr *KaArray) PushFront(val KaType) {
	arr.Length++
	arr.Array = append([]*KaType{&val}, arr.Array...)
}

func (arr *KaArray) PopBack(val KaType) {
	arr.Length--
	arr.Array = arr.Array[:arr.Length]
}

func (arr *KaArray) PopFront(val KaType) {
	arr.Length--
	arr.Array = arr.Array[1:]
}

func (arr KaArray) Format() string {
	var formatted = "("
	for _, v := range arr.Array {
		formatted += (*v).Format() + ", "
	}

	if len(formatted) > 1 {
		formatted = formatted[:len(formatted)-2] //remove the trailing ", "
	}

	formatted += ")"
	return formatted
}

func (arr KaArray) Type() string {
	return "array"
}

func (arr KaArray) TypeOf() string {
	return arr.Type()
}

func (arr KaArray) Deallocate() {}

//Range ranges over an array
func (arr KaArray) Range(fn func(val1, val2 *KaType) Returner) *Returner {

	for k, v := range arr.Array {
		var key KaNumber
		key.FromGoType(float64(k))
		var katypekey KaType = key
		ret := fn(&katypekey, v)

		if ret.Type == "break" {
			break
		} else if ret.Type == "return" {
			return &ret
		}
	}

	return nil
}
