package interpreter

import (
	"github.com/tusklang/tusk/lang/types"
)

func strToBits(val1, val2 types.TuskType, operator func(a, b rune) rune) *types.TuskType {
	str1, str2 := val1.(types.TuskString).ToRuneList(), val2.(types.TuskString).ToRuneList()

	if len(str1) < len(str2) {
		str1, str2 = str2, str1
	}

	var fin = make([]rune, len(str1))

	for k, v := range str1 {
		fin[k] = operator(v, str2[k])
	}

	var finstr types.TuskString
	finstr.FromRuneList(fin)
	var tusktype types.TuskType = finstr
	return &tusktype
}

func strbitwiseAnd(val1, val2 types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return strToBits(val1, val2, func(a, b rune) rune {
		return a & b
	}), nil
}

func strbitwiseOr(val1, val2 types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return strToBits(val1, val2, func(a, b rune) rune {
		return a | b
	}), nil
}

func strbitwiseXor(val1, val2 types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {
	return strToBits(val1, val2, func(a, b rune) rune {
		return a ^ b
	}), nil
}

func strbitwiseNot(val1, val2 types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) (*types.TuskType, *types.TuskError) {

	val1 = types.TuskString{} //force the `none` type given to be an empty string (`b` stores the real string)

	return strToBits(val1, val2, func(a, b rune) rune {
		return ^b
	}), nil
}
