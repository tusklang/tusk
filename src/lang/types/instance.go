package types

type Instance struct {
  Params   CliParams
  vars     map[string]*OmmVar
}

func (ins *Instance) Allocate(name string, value *OmmType) {

  if ins.vars == nil {
    ins.vars = make(map[string]*OmmVar)
  }

  ins.vars[name] = &OmmVar{
    Name: name,
    Value: value,
  }
}

func (ins *Instance) Deallocate(name string) {
  delete(ins.vars, name)
}

func (ins *Instance) Fetch(name string) *OmmVar {
	return ins.vars[name]
}
