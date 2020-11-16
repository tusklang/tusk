package interpreter

import (
	"math/big"

	"github.com/tusklang/tusk/lang/types"
)

//using the math/big library to implement bitwise operators

func bitwise(f func(a, b *big.Int) *big.Int, a, b types.TuskType) (*types.TuskType, *types.TuskError) {
	gobig := f(a.(types.TuskNumber).ToBigInt(), b.(types.TuskNumber).ToBigInt())
	var tusknumber types.TuskNumber
	tusknumber.FromBigInt(gobig)
	var tusktype types.TuskType = tusknumber
	return &tusktype, nil
}

func bitwiseNot(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {

	//force a to be a dummy tusk number (to not have a type assertion error)
	a = types.TuskNumber{}

	return bitwise(func(a, b *big.Int) *big.Int {
		return b.Not(b)
	}, a, b)
}

func bitwiseAnd(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return bitwise(func(a, b *big.Int) *big.Int {
		return a.And(a, b)
	}, a, b)
}

func bitwiseOr(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return bitwise(func(a, b *big.Int) *big.Int {
		return a.Or(a, b)
	}, a, b)
}

func bitwiseXor(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return bitwise(func(a, b *big.Int) *big.Int {
		return a.Xor(a, b)
	}, a, b)
}

func bitwiseRShift(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return bitwise(func(a, b *big.Int) *big.Int {
		return a.Rsh(a, uint(b.Int64()))
	}, a, b)
}

func bitwiseLShift(a, b types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return bitwise(func(a, b *big.Int) *big.Int {
		return a.Lsh(a, uint(b.Int64()))
	}, a, b)
}
