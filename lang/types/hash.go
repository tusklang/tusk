package types

import "strings"

type TuskHash struct {
	Hash   map[*TuskType]*TuskType
	keys   []*TuskType
	Length uint64
}

func (hash TuskHash) At(idx *TuskType) *TuskType {

	v, exists := hash.Hash[idx]

	if !exists {
		var undef TuskType = TuskUndef{}
		return &undef
	}

	return v
}

//AtStr gets the value of a Tusk Hash based on a string
func (hash TuskHash) AtStr(idx string) *TuskType {
	var tuskstr TuskString
	tuskstr.FromGoType(idx)
	var tusktype TuskType = tuskstr
	return hash.At(&tusktype)
}

func (hash *TuskHash) Set(idx *TuskType, val TuskType) {

	if hash.Hash == nil {
		hash.Hash = make(map[*TuskType]*TuskType)
	}

	v, exists := hash.Hash[idx]

	if !exists {
		hash.keys = append(hash.keys, idx)
		hash.Length++
	}

	*v = val
}

func (hash TuskHash) Exists(idx *TuskType) bool {
	_, exists := hash.Hash[idx]
	return exists
}

func formathash(v *TuskType, pformatted *string) {
	if *pformatted != "[]" {
		newlineSplit := strings.Split(*pformatted, "\n")

		*pformatted = ""

		for _, val := range newlineSplit {
			*pformatted += strings.Repeat(" ", 2) + val + "\n"
		}

		*pformatted = strings.TrimSpace(*pformatted) //remove the trailing \n (because an extra was added) and the leading two spaces (because it will be on the same line)
	}
}

func (hash TuskHash) Format() string {

	return func() string {

		if len(hash.Hash) == 0 {
			return "[]"
		}

		var formatted = "["

		for k, v := range hash.Hash {
			kFormatted := (*k).Format()
			vFormatted := (*v).Format()

			switch (*v).(type) {
			case TuskHash: //if it is another hash, add the indents
				formathash(k, &kFormatted)
				formathash(v, &vFormatted)
			}

			formatted += "\n" + strings.Repeat(" ", 2) + kFormatted + " = " + vFormatted + ","
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

		ret := fn(k, v)

		if ret.Type == "break" {
			break
		} else if ret.Type == "return" {
			return &ret
		}
	}

	return nil
}
