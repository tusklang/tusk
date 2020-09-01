package types

import (
	"errors"
)

type OmmProto struct {
	ProtoName  string
	Static     map[string]*OmmType
	Instance   map[string]*OmmType
	AccessList map[string][]string
}

func getfield(full map[string]*OmmType, field string, access map[string][]string, file string) (*OmmType, error) {
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

func (p OmmProto) Get(field string, file string) (*OmmType, error) {
	return getfield(p.Static, field, p.AccessList, file)
}

func (p OmmProto) Format() string {
	return "{" + p.ProtoName[1:] + "}"
}

func (p OmmProto) Type() string {
	return "proto"
}

func (p OmmProto) TypeOf() string {
	return p.ProtoName[1:] /* remove the leading $ */ + " prototype"
}

func (_ OmmProto) Deallocate() {}
