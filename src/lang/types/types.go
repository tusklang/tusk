package types

//package to store all of the datatypes in omm

type CliParams map[string]map[string]interface{}

type OmmType interface {
  ValueFunc()
}
