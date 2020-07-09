package types

type OmmArray struct {
  array []OmmType
  Length  uint64
}

func (arr OmmArray) At(idx int64) OmmType {

  length := arr.Length

  if uint64(idx) >= length || idx < 0 {
    var undef OmmUndef
    return undef
  }

  return arr.array[idx]
}

func (arr *OmmArray) Set(idx int64, val OmmType) {
  arr.array[idx] = val
}

func (arr OmmArray) Exists(idx int64) bool {
  return uint64(idx) < arr.Length || idx >= 0
}

func (arr *OmmArray) PushBack(val OmmType) {
  arr.Length++
  arr.array = append(arr.array, val)
}

func (arr *OmmArray) PushFront(val OmmType) {
  arr.Length++
  arr.array = append([]OmmType{ val }, arr.array...)
}

func (arr *OmmArray) PopBack(val OmmType) {
  arr.Length--
  arr.array = arr.array[:arr.Length]
}

func (arr *OmmArray) PopFront(val OmmType) {
  arr.Length--
  arr.array = arr.array[1:]
}

func (_ OmmArray) ValueFunc() {}
