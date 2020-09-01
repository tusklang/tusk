package types

type Instance struct {
	Params CliParams
	vars   map[string]*OmmVar
}

func (ins *Instance) Allocate(name string, value *OmmType) {

	if ins.vars == nil {
		ins.vars = make(map[string]*OmmVar)
	}

	ins.vars[name] = &OmmVar{
		Name:  name,
		Value: value,
	}
}

func (ins *Instance) Deallocate(name string) {

	if ins.vars == nil {
		ins.vars = make(map[string]*OmmVar)
	}

	//first do the complex dealloc of the type
	(*ins.vars[name].Value).Deallocate()

	delete(ins.vars, name)
}

func (ins *Instance) Fetch(name string) *OmmVar {

	if ins.vars == nil {
		ins.vars = make(map[string]*OmmVar)
	}

	return ins.vars[name]
}

func (ins Instance) Copy() *Instance {

	var nins Instance
	nins.Params = ins.Params
	nins.vars = make(map[string]*OmmVar)

	for k, v := range ins.vars {
		nins.vars[k] = v
	}

	return &nins
}
