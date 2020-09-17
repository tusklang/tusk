package types

import "strings"

type TuskHash struct {
	Hash   map[string]*TuskType
	keys   []string
	Length uint64
}

func (hash TuskHash) At(idx string) *TuskType {

	if _, exists := hash.Hash[idx]; !exists {
		var undef TuskType = TuskUndef{}
		hash.Hash[idx] = &undef
	}

	return hash.Hash[idx]
}

func (hash *TuskHash) Set(idx string, val TuskType) {

	if hash.Hash == nil {
		hash.Hash = make(map[string]*TuskType)
	}

	if _, exists := hash.Hash[idx]; !exists {
		hash.keys = append(hash.keys, idx)
		hash.Length++
	}

	hash.Hash[idx] = &val
}

func (hash TuskHash) Exists(idx string) bool {
	_, exists := hash.Hash[idx]
	return exists
}

func (hash TuskHash) Format() string {

	return func() string {

		if len(hash.Hash) == 0 {
			return "[]"
		}

		var formatted = "["

		for k, v := range hash.Hash {

			vFormatted := (*v).Format()

			switch (*v).(type) {
			case TuskHash: //if it is another hash, add the indents
				if vFormatted != "[]" {
					newlineSplit := strings.Split(vFormatted, "\n")

					vFormatted = ""

					for _, val := range newlineSplit {
						vFormatted += strings.Repeat(" ", 2) + val + "\n"
					}

					vFormatted = strings.TrimSpace(vFormatted) //remove the trailing \n (because an extra was added) and the leading two spaces (because it will be on the same line)
				}
			}

			formatted += "\n" + strings.Repeat(" ", 2) + k + " = " + vFormatted + ","
		}

		return formatted + "\n]"
	}() //staring with 2
}

func (hash TuskHash) Type() string {
	return "hash"
}

func (hash TuskHash) TypeOf() string {
	return hash.Type()
}

func (hash TuskHash) Deallocate() {}

//Range ranges over a hash
func (hash TuskHash) Range(fn func(val1, val2 *TuskType) Returner) *Returner {

	for _, keyi := range hash.keys {

		k, v := keyi, hash.Hash[keyi]

		var key TuskString
		key.FromGoType(k)
		var katypekey TuskType = key
		ret := fn(&katypekey, v)

		if ret.Type == "break" {
			break
		} else if ret.Type == "return" {
			return &ret
		}
	}

	return nil
}
