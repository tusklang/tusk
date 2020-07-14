package types

type OmmGoFunc struct {
  Function func(args []*OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType
}

func (ogf OmmGoFunc) Format() string {
  return "{native gofunc}"
}

func (ogf OmmGoFunc) Type() string {
  return "gofunc"
}

func (ogf OmmGoFunc) TypeOf() string {
  return ogf.Type()
}
