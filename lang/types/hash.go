package types

import (
	"strings"
)

//TuskHash represents a hash in tusk
type TuskHash struct {
	hash   map[string]*TuskType
	keys   []*TuskType
	length uint64
}

//At gets the value of a tusk hash based on a tusk type
func (hash TuskHash) At(idx *TuskType) *TuskType {

	k := (*idx).Format()
	v, exists := hash.hash[k]

	if !exists {
		var undef TuskType = TuskUndef{}
		var pundef *TuskType = &undef
		hash.hash[k] = pundef
		return hash.hash[k]
	}

	return v
}

func (hash TuskHash) Length() uint64 {
	return hash.length
}

//AtStr gets the value of a tusk hash based on a string
func (hash TuskHash) AtStr(idx string) *TuskType {
	var tuskstr TuskString
	tuskstr.FromGoType(idx)
	var tusktype TuskType = tuskstr
	return hash.At(&tusktype)
}

//Set sets a value of a hash given a key and a value
func (hash *TuskHash) Set(idx *TuskType, val TuskType) {

	if hash.hash == nil {
		hash.hash = make(map[string]*TuskType)
	}

	v, exists := hash.hash[(*idx).Format()]

	if !exists {
		hash.keys = append(hash.keys, idx)
		hash.length++
		hash.hash[(*idx).Format()] = &val
		return
	}

	*v = val
}

//SetStr sets a value of a hash given a key and a value, where the key is given as a go string
func (hash *TuskHash) SetStr(idx string, val TuskType) {
	var tuskstr TuskString
	tuskstr.FromGoType(idx)
	var tusktype TuskType = tuskstr
	hash.Set(&tusktype, val)
}

//Exists determines if a key exists in a hash
func (hash TuskHash) Exists(idx *TuskType) bool {
	_, exists := hash.hash[(*idx).Format()]
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

//Format formats a hash
func (hash TuskHash) Format() string {

	return func() string {

		if len(hash.hash) == 0 {
			return "[::]"
		}

		var formatted = "[:"

		for k, v := range hash.hash {
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

//Type returns the type of a value
func (hash TuskHash) Type() string {
	return "hash"
}

//TypeOf returns the type of a value's object
func (hash TuskHash) TypeOf() string {
	return hash.Type()
}

//Deallocate deallocates any hanging memory in a hash
func (hash TuskHash) Deallocate() {
	for _, v := range hash.hash {
		(*v).Deallocate()
	}
}

func (hash TuskHash) Clone() *TuskType {
	var h = hash.hash

	//clone it into `cloned`
	var cloned = make(map[string]*TuskType)
	for k, v := range h {
		cloned[k] = (*v).Clone()
	}
	////////////////////////

	var tusktype TuskType = TuskHash{
		hash:   cloned,
		keys:   hash.keys,
		length: hash.length,
	}

	return &tusktype
}

//Range ranges over a hash
func (hash TuskHash) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {

	for _, keyi := range hash.keys {

		k, v := keyi, hash.hash[(*keyi).Format()]

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
