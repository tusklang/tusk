package types

//GarbageCollectionThreshold stores the amount of variables that can be allocated until de-allocation starts
const GarbageCollectionThreshold = 1337000 //once `GarbageCollectionThreshold` variables are allocated, we should start deallocating

type Instance struct {
	Params   CliParams
	vars     map[string]*TuskVar
	varBinds map[string][]string
}

func (ins *Instance) Allocate(name string, value *TuskType) {

	if ins.vars == nil {
		ins.vars = make(map[string]*TuskVar)
	}

	if len(ins.vars) >= GarbageCollectionThreshold {
		for k, v := range ins.varBinds {
			if len(v) == 0 {
				ins.Deallocate(k)
			}
		}
	}

	ins.vars[name] = &TuskVar{
		Name:  name,
		Value: value,
	}
	ins.Bind(name, "scope") //add the scope bind
}

func (ins *Instance) Deallocate(name string) {
	if ins.vars == nil {
		ins.vars = make(map[string]*TuskVar)
	}

	//first do the complex dealloc of the type
	(*ins.vars[name].Value).Deallocate()

	delete(ins.vars, name)
	delete(ins.varBinds, name)
}

func (ins *Instance) Fetch(name string) *TuskVar {

	if ins.vars == nil {
		ins.vars = make(map[string]*TuskVar)
	}

	return ins.vars[name]
}

//Bind binds variable `var1` to variable `var2`
//once a variable has no binds left, it is deallocated
//every variable starts with a "scope" bind, but once the scope it is declared in ends, the "scope" bind is removed
//if a variable is returned, it is deallocated, but a return address is created with a "scope" bind
func (ins *Instance) Bind(var1, var2 string) {

	if ins.varBinds == nil {
		ins.varBinds = make(map[string][]string)
	}

	ins.varBinds[var1] = append(ins.varBinds[var1], var2)
}

//DeBind removed the bind created with `var1` to `var2`
func (ins *Instance) DeBind(var1, var2 string) {
	for k, v := range ins.varBinds[var1] {
		if v == var2 {

			//copy last element
			last := ins.varBinds[var1][len(ins.varBinds[var1])-1]
			//set the erased index to the last value
			ins.varBinds[var1][k] = last
			//remove the last index
			ins.varBinds[var1] = ins.varBinds[var1][:len(ins.varBinds[var1])-1]

			return
		}
	}
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
