package interpreter

import . "github.com/tusklang/tusk/lang/types"

func number__mod__number(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string) *TuskType {
	num1, num2 := val1.(TuskNumber), val2.(TuskNumber)
	ensurePrec(&num1, &num2, (*instance).Params)

	//ALGORITHM:
	//  num1 - floor(num1 / num2) * num2

	if num2.Decimal == nil { //ensure that a nil pointer reference doesnt happen
		num2.Decimal = &[]int64{0}
	}

	num2P := zero //create a placeholder for num2 (because it will get mutated)
	tmpInt, tmpDec := append([]int64{}, *num2.Integer...), append([]int64{}, *num2.Decimal...)
	num2P.Integer, num2P.Decimal = &tmpInt, &tmpDec

	//if you set the prec to 0 here, it will mutate it
	divided := (*number__divide__number(num1, num2, instance, stacktrace, line, file)).(TuskNumber)
	*divided.Decimal = nil //round down

	multiplied := (*number__times__number(divided, num2, instance, stacktrace, line, file)).(TuskNumber)
	return number__minus__number(num1, multiplied, instance, stacktrace, line, file)
}
