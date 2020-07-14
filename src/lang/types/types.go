package types

//package to store all of the datatypes in omm

type CliParams map[string]map[string]interface{}

type OmmType interface {
  Format()  string

  //difference between type and typeof is for prototypes
  //Type() for a proto will return "prototype"
  //TypeOf() will return the proto's name
  Type()    string
  TypeOf()  string
  //////////////////////////////////////////////////////
}
