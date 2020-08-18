package types

type OmmArray struct {
  Array  []*OmmType
  Length    uint64
}

func (arr OmmArray) At(idx int64) *OmmType {

  length := arr.Length

  if uint64(idx) >= length || idx < 0 {
    var undef OmmType = OmmUndef{}
    return &undef
  }

  return arr.Array[idx]
}

func (arr OmmArray) Exists(idx int64) bool {
  return arr.Length != 0 && uint64(idx) < arr.Length && idx >= 0
}

func (arr *OmmArray) PushBack(val OmmType) {
  arr.Length++
  arr.Array = append(arr.Array, &val)
}

func (arr *OmmArray) PushFront(val OmmType) {
  arr.Length++
  arr.Array = append([]*OmmType{ &val }, arr.Array...)
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
    formatted+=(*v).Format() + ", "
  }

  if len(formatted) > 1 {
    formatted = formatted[:len(formatted) - 2] //remove the trailing ", "
  }

  formatted+="]"
  return formatted
}

func (arr OmmArray) Type() string {
  return "array"
}

func (arr OmmArray) TypeOf() string {
  return arr.Type()
}

func (_ OmmArray) Deallocate() {}
