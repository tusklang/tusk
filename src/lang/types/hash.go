package types

type OmmHash struct {
  hash map[string]OmmType
  Length  uint64
}

func (hash OmmHash) At(idx string) OmmType {
  return hash.hash[idx]
}

func (hash *OmmHash) Set(idx string, val OmmType) {

  if hash.hash == nil {
    hash.hash = map[string]OmmType{}
  }

  if _, exists := hash.hash[idx]; !exists {
    hash.Length++
  }

  hash.hash[idx] = val
}

func (hash OmmHash) Exists(idx string) bool {
  _, exists := hash.hash[idx]
  return exists
}

func (_ OmmHash) ValueFunc() {}
