package interpreter

import . "tusk/lang/types"

func number__pow__number(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string) *TuskType {
	num1, num2 := val1.(TuskNumber), val2.(TuskNumber)
	ensurePrec(&num1, &num2, (*instance).Params)

	expNeg := false
	if isLess(num2, zero) { //account for negative exponents
		expNeg = true
		num2 = (*number__times__number(num2, neg_one, instance, stacktrace, line, file)).(TuskNumber)
	}

	if len(*num2.Decimal) == 0 { //if the exponent is an integer, use binary exponentiation for an O(log n) solution

		powwed := number__pow__integer(num1, num2, instance, stacktrace, line, file)

		if expNeg {
			powwed = *number__divide__number(one, powwed, instance, stacktrace, line, file)
		}

		return &powwed
	}

	var two = zero
	two.Integer = &[]int64{2}

	neg := false
	if isLess(num1, zero) && isEqual((*number__mod__number(num1, two, instance, stacktrace, line, file)).(TuskNumber), zero) { //because ln (n < 0) is undefined
		neg = true
	}

	num1 = abs(num1, stacktrace, (*instance).Params).(TuskNumber)

	powwed := exp((*number__times__number(num2, ln(num1, instance, stacktrace, line, file), instance, stacktrace, line, file)), instance, stacktrace, line, file)

	if expNeg {
		powwed = *number__divide__number(one, powwed, instance, stacktrace, line, file)
	}
	if neg {
		powwed = *number__times__number(powwed, neg_one, instance, stacktrace, line, file)
	}

	//algorithm is exp(n2 ln n1)
	return &powwed
}
