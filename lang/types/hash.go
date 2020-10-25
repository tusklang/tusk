package types

import (
	"strings"
)

type TuskHash struct {
	Hash   map[string]*TuskType
	keys   []*TuskType
	Length uint64
}

func (hash TuskHash) At(idx *TuskType) *TuskType {

	v, exists := hash.Hash[(*idx).Format()]

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
		hash.Hash = make(map[string]*TuskType)
	}

	v, exists := hash.Hash[(*idx).Format()]

	if !exists {
		hash.keys = append(hash.keys, idx)
		hash.Length++
		hash.Hash[(*idx).Format()] = &val
		return
	}

	*v = val
}

func (hash *TuskHash) SetStr(idx string, val TuskType) {
	var tuskstr TuskString
	tuskstr.FromGoType(idx)
	var tusktype TuskType = tuskstr
	hash.Set(&tusktype, val)
}

func (hash TuskHash) Exists(idx *TuskType) bool {
	_, exists := hash.Hash[(*idx).Format()]
	return exists
}

func formathash(pformatted *string) {
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
			return "[::]"
		}

		var formatted = "[:"

		for k, v := range hash.Hash {
			kFormatted := k
			vFormatted := (*v).Format()

			switch (*v).(type) {
			case TuskHash: //if it is another hash, add the indents
				formathash(&kFormatted)
				formathash(&vFormatted)
			}

			formatted += "\n" + strings.Repeat(" ", 2) + kFormatted + " = " + vFormatted + ","
		}

		return formatted + "\n:]"
	}() //staring with 2
}

func (hash TuskHash) Type() string {
	return "hash"
}

func (hash TuskHash) TypeOf() string {
	return hash.Type()
}

func (hash TuskHash) Deallocate() {}

func (hash TuskHash) Clone() *TuskType {
	var h = hash.Hash

	//clone it into `cloned`
	var cloned = make(map[string]*TuskType)
	for k, v := range h {
		cloned[k] = (*v).Clone()
	}
	////////////////////////

	var tusktype TuskType = TuskHash{
		Hash:   cloned,
		Length: hash.Length,
	}

	return &tusktype
}

//Range ranges over a hash
func (hash TuskHash) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {

	for _, keyi := range hash.keys {

		k, v := keyi, hash.Hash[(*keyi).Format()]

		ret, e := fn(k, v)

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
