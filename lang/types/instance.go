package types

type Instance struct {
	Params CliParams
	vars   map[string]*KaVar
}

func (ins *Instance) Allocate(name string, value *KaType) {

	if ins.vars == nil {
		ins.vars = make(map[string]*KaVar)
	}

	ins.vars[name] = &KaVar{
		Name:  name,
		Value: value,
	}
}

func (ins *Instance) Deallocate(name string) {

	if ins.vars == nil {
		ins.vars = make(map[string]*KaVar)
	}

	//first do the complex dealloc of the type
	(*ins.vars[name].Value).Deallocate()

	delete(ins.vars, name)
}

func (ins *Instance) Fetch(name string) *KaVar {

	if ins.vars == nil {
		ins.vars = make(map[string]*KaVar)
	}

	return ins.vars[name]
}

func (ins Instance) Copy() *Instance {

	var nins Instance
	nins.Params = ins.Params
	nins.vars = make(map[string]*KaVar)

	for k, v := range ins.vars {
		nins.vars[k] = v
	}

	return &nins
}
