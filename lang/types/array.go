package types

//TuskArray represents an array in tusk
type TuskArray struct {
	array  []*TuskType
	length uint64
}

//At fetches the value at an index of an array
func (arr TuskArray) At(idx int64) *TuskType {

	length := arr.length

	if uint64(idx) >= length || idx < 0 {
		var undef TuskType = TuskUndef{}
		return &undef
	}

	return arr.array[idx]
}

//Length gets the length of the given array
func (arr TuskArray) Length() uint64 {
	return arr.length
}

//Exists determines if a given index exists in the array
func (arr TuskArray) Exists(idx int64) bool {
	return arr.length != 0 && uint64(idx) < arr.length && idx >= 0
}

//PushBack appends an item to the array
func (arr *TuskArray) PushBack(val TuskType) {
	arr.length++
	arr.array = append(arr.array, &val)
}

//PushFront prepends an item to the array
func (arr *TuskArray) PushFront(val TuskType) {
	arr.length++
	arr.array = append([]*TuskType{&val}, arr.array...)
}

//PopBack pops an item from the end of an array
func (arr *TuskArray) PopBack(val TuskType) {
	arr.length--
	arr.array = arr.array[:arr.length]
}

//PopFront pops an item from the start of an array
func (arr *TuskArray) PopFront(val TuskType) {
	arr.length--
	arr.array = arr.array[1:]
}

//Format formats the array as a string
func (arr TuskArray) Format() string {
	var formatted = "["
	for _, v := range arr.array {
		formatted += (*v).Format() + ", "
	}

	if len(formatted) > 1 {
		formatted = formatted[:len(formatted)-2] //remove the trailing ", "
	}

	formatted += "]"
	return formatted
}

//Type returns the type of a value
func (arr TuskArray) Type() string {
	return "array"
}

//TypeOf returns the type of a value's object
func (arr TuskArray) TypeOf() string {
	return arr.Type()
}

//Deallocate deallocates any hanging values associated with the array
func (arr TuskArray) Deallocate() {
	for _, v := range arr.array {
		(*v).Deallocate()
	}
}

//Clone clones the value into a new pointer
func (arr TuskArray) Clone() *TuskType {
	var a = arr.array

	var cloned = make([]*TuskType, arr.length)

	for k, v := range a {
		//clone each value in the array
		cloned[k] = (*v).Clone()
	}

	var tusktype TuskType = TuskArray{
		array:  cloned,
		length: arr.length,
	}

	return &tusktype
}

//Range ranges over an array
func (arr TuskArray) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {

	for k, v := range arr.array {
		var key TuskInt
		key.FromGoType(int64(k))
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
