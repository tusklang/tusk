package types

import (
	"errors"
)

type TuskProto struct {
	ProtoName string
	Static    map[string]*TuskType
	Instance  map[string]*TuskType
}

func getfield(full map[string]*TuskType, field, file, namespace, protoname string) (*TuskType, error) {

	if namespace == protoname {
		goto has_access
	}

	if field[0] == '_' {
		return nil, errors.New("Cannot access private member: " + field)
	}

has_access:

	fieldv := full[field]

	if fieldv == nil {
		return nil, errors.New("Prototype does not contain the field \"" + field + "\"")
	}

	return fieldv, nil
}

func (p TuskProto) Get(field, file, namespace string) (*TuskType, error) {
	return getfield(p.Static, field, file, namespace, p.ProtoName)
}

func (p TuskProto) Format() string {
	return "{" + p.ProtoName + "}"
}

func (p TuskProto) Type() string {
	return "prototype"
}

func (p TuskProto) TypeOf() string {
	return p.ProtoName + " prototype"
}

func (p TuskProto) Deallocate() {}

func (p TuskProto) Clone() *TuskType {
	return nil
}

//Range ranges over a prototype
func (p TuskProto) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
