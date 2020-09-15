package types

//New creates a new KaObject from an KaProto
func (proto KaProto) New(instance Instance) KaObject {
	nins := instance.Copy() //copy the original vars

	for k, v := range proto.Instance {
		nins.Allocate(k, v)

		switch (*v).(type) {
		case KaFunc: //if it is a function, change the instance
			tmp := (*v).(KaFunc)
			tmp.Instance = nins
			*v = tmp
		}
	}

	var obj = KaObject{
		Name:       proto.ProtoName,
		Instance:   *nins,
		AccessList: proto.AccessList,
	}
	return obj
}
