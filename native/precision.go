package native

import "omm/lang/types"

var setprec = func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

	if len(args) != 1 || (*args[0]).Type() != "number" {
		OmmPanic("Function instance.setprec requires an argument count of 1 with the type of number", line, file, stacktrace)
	}

	var goprec = uint64((*args[0]).(types.OmmNumber).ToGoType())
	instance.Params.Prec = goprec

	var undef types.OmmType = types.OmmUndef{}
	return &undef
}
