package interpreter

import "github.com/tusklang/tusk/lang/types"

func bitwise(f func(a, b int64) int64, a, b types.TuskType) (*types.TuskType, *types.TuskError) {
	gon := f(int64(a.(types.TuskNumber).ToGoType()), int64(b.(types.TuskNumber).ToGoType()))
	var tusknumber types.TuskNumber
	tusknumber.FromGoType(float64(gon))
	var tusktype types.TuskType = tusknumber
	return &tusktype, nil
}

func bitwiseNot(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {

	//force a to be a dummy tusk number (to not have a type assertion error)
	a = types.TuskNumber{}

	return bitwise(func(a, b int64) int64 {
		return ^b
	}, a, b)
}

func bitwiseAnd(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return bitwise(func(a, b int64) int64 {
		return a & b
	}, a, b)
}

func bitwiseOr(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return bitwise(func(a, b int64) int64 {
		return a | b
	}, a, b)
}

func bitwiseXor(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return bitwise(func(a, b int64) int64 {
		return a ^ b
	}, a, b)
}

func bitwiseRShift(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return bitwise(func(a, b int64) int64 {
		return a >> uint64(b)
	}, a, b)
}

func bitwiseLShift(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return bitwise(func(a, b int64) int64 {
		return a << uint64(b)
	}, a, b)
}

func bitwiseURShift(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return bitwise(func(a, b int64) int64 {
		return int64(uint64(a) >> uint64(b))
	}, a, b)
}
