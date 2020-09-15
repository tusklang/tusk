package types

import (
	"errors"
)

type KaProto struct {
	ProtoName  string
	Static     map[string]*KaType
	Instance   map[string]*KaType
	AccessList map[string][]string
}

func getfield(full map[string]*KaType, field string, access map[string][]string, file string) (*KaType, error) {
	if field[0] == '_' {
		return nil, errors.New("Cannot access private member: " + field)
	}

	//check for access (protected)

	if access[field] == nil || file == "" { //if it does not name any access, automatically make it public
		goto allowed
	}

	for _, v := range access[field] {
		if file == v {
			goto allowed
		}
	}

	return nil, errors.New("File cannot acces field \"" + field + "\"")

allowed:
	fieldv := full[field]

	if fieldv == nil {
		return nil, errors.New("Prototype does not contain the field \"" + field + "\"")
	}

	return fieldv, nil
}

func (p KaProto) Get(field string, file string) (*KaType, error) {
	return getfield(p.Static, field, p.AccessList, file)
}

func (p KaProto) Format() string {
	return "{" + p.ProtoName[1:] + "}"
}

func (p KaProto) Type() string {
	return "proto"
}

func (p KaProto) TypeOf() string {
	return p.ProtoName[1:] /* remove the leading $ */ + " prototype"
}

func (p KaProto) Deallocate() {}

//Range ranges over a prototype
func (p KaProto) Range(fn func(val1, val2 *KaType) Returner) *Returner {
	return nil
}
