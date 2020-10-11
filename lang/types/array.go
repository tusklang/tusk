package types

type TuskArray struct {
	Array  []*TuskType
	Length uint64
}

func (arr TuskArray) At(idx int64) *TuskType {

	length := arr.Length

	if uint64(idx) >= length || idx < 0 {
		var undef TuskType = TuskUndef{}
		return &undef
	}

	return arr.Array[idx]
}

func (arr TuskArray) Exists(idx int64) bool {
	return arr.Length != 0 && uint64(idx) < arr.Length && idx >= 0
}

func (arr *TuskArray) PushBack(val TuskType) {
	arr.Length++
	arr.Array = append(arr.Array, &val)
}

func (arr *TuskArray) PushFront(val TuskType) {
	arr.Length++
	arr.Array = append([]*TuskType{&val}, arr.Array...)
}

func (arr *TuskArray) PopBack(val TuskType) {
	arr.Length--
	arr.Array = arr.Array[:arr.Length]
}

func (arr *TuskArray) PopFront(val TuskType) {
	arr.Length--
	arr.Array = arr.Array[1:]
}

func (arr TuskArray) Format() string {
	var formatted = "["
	for _, v := range arr.Array {
		formatted += (*v).Format() + ", "
	}

	if len(formatted) > 1 {
		formatted = formatted[:len(formatted)-2] //remove the trailing ", "
	}

	formatted += "]"
	return formatted
}

func (arr TuskArray) Type() string {
	return "array"
}

func (arr TuskArray) TypeOf() string {
	return arr.Type()
}

func (arr TuskArray) Deallocate() {}

//Range ranges over an array
func (arr TuskArray) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {

	for k, v := range arr.Array {
		var key TuskNumber
		key.FromGoType(float64(k))
		var tusktypekey TuskType = key
		ret, e := fn(&tusktypekey, v)

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
