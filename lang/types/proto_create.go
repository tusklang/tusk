package types

//New creates a new TuskObject from an TuskProto
func (proto TuskProto) New(instance Instance) TuskObject {
	nins := instance.Copy() //copy the original vars

	for k, v := range proto.Instance {
		nins.Allocate(k, v)

		switch (*v).(type) {
		case TuskFunc: //if it is a function, change the instance
			tmp := (*v).(TuskFunc)
			tmp.Instance = nins
			*v = tmp
		}
	}

	var obj = TuskObject{
		Name:       proto.ProtoName,
		Instance:   *nins,
		AccessList: proto.AccessList,
	}
	return obj
}
