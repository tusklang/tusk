package types

import (
	"errors"
)

type TuskProto struct {
	ProtoName string
	Static    map[string]*TuskType
	Instance  map[string]*TuskType
}

func getfield(full map[string]*TuskType, field string, file string) (*TuskType, error) {
	if field[0] == '_' {
		return nil, errors.New("Cannot access private member: " + field)
	}

	fieldv := full[field]

	if fieldv == nil {
		return nil, errors.New("Prototype does not contain the field \"" + field + "\"")
	}

	return fieldv, nil
}

func (p TuskProto) Get(field string, file string) (*TuskType, error) {
	return getfield(p.Static, field, file)
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
