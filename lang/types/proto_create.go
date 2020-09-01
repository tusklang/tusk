package types

//New creates a new OmmObject from an OmmProto
func (proto OmmProto) New(instance Instance) OmmObject {
	nins := instance.Copy() //copy the original vars

	for k, v := range proto.Instance {
		nins.Allocate(k, v)

		switch (*v).(type) {
		case OmmFunc: //if it is a function, change the instance
			tmp := (*v).(OmmFunc)
			tmp.Instance = nins
			*v = tmp
		}
	}

	var obj = OmmObject{
		Name:       proto.ProtoName,
		Instance:   *nins,
		AccessList: proto.AccessList,
	}
	return obj
}
