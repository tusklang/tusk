package types

type OmmArray struct {
  Array   []OmmType
  Length    uint64
}

func (arr OmmArray) At(idx int64) OmmType {

  length := arr.Length

  if uint64(idx) >= length || idx < 0 {
    var undef OmmUndef
    return undef
  }

  return arr.Array[idx]
}

func (arr *OmmArray) Set(idx int64, val OmmType) {
  arr.Array[idx] = val
}

func (arr OmmArray) Exists(idx int64) bool {
  return uint64(idx) < arr.Length || idx >= 0
}

func (arr *OmmArray) PushBack(val OmmType) {
  arr.Length++
  arr.Array = append(arr.Array, val)
}

func (arr *OmmArray) PushFront(val OmmType) {
  arr.Length++
  arr.Array = append([]OmmType{ val }, arr.Array...)
}

func (arr *OmmArray) PopBack(val OmmType) {
  arr.Length--
  arr.Array = arr.Array[:arr.Length]
}

func (arr *OmmArray) PopFront(val OmmType) {
  arr.Length--
  arr.Array = arr.Array[1:]
}

func (arr OmmArray) Format() string {
  var formatted = "["
  for _, v := range arr.Array {
    formatted+=v.Format()
  }
  formatted+="]"
  return formatted
}

func (arr OmmArray) Type() string {
  return "array"
}
