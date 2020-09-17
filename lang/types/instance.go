package types

type Instance struct {
	Params CliParams
	vars   map[string]*TuskVar
}

func (ins *Instance) Allocate(name string, value *TuskType) {

	if ins.vars == nil {
		ins.vars = make(map[string]*TuskVar)
	}

	ins.vars[name] = &TuskVar{
		Name:  name,
		Value: value,
	}
}

func (ins *Instance) Deallocate(name string) {

	if ins.vars == nil {
		ins.vars = make(map[string]*TuskVar)
	}

	//first do the complex dealloc of the type
	(*ins.vars[name].Value).Deallocate()

	delete(ins.vars, name)
}

func (ins *Instance) Fetch(name string) *TuskVar {

	if ins.vars == nil {
		ins.vars = make(map[string]*TuskVar)
	}

	return ins.vars[name]
}

func (ins Instance) Copy() *Instance {

	var nins Instance
	nins.Params = ins.Params
	nins.vars = make(map[string]*TuskVar)

	for k, v := range ins.vars {
		nins.vars[k] = v
	}

	return &nins
}
