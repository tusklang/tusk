package interpreter

import . "ka/lang/types"

func number__minus__number(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string) *KaType {
	num1, num2 := val1.(KaNumber), val2.(KaNumber)
	ensurePrec(&num1, &num2, (*instance).Params)

	num2Placeholder := zero                                                                //create a placeholder for num2 (so it wont mutate)
	tmpInt, tmpDec := make([]int64, len(*num2.Integer)), make([]int64, len(*num2.Decimal)) //allocate the length
	num2Placeholder.Integer, num2Placeholder.Decimal = &tmpInt, &tmpDec

	//looks like this
	// a - b = a + -b

	//invert the decimal
	for k, v := range *num2.Decimal {
		(*num2Placeholder.Decimal)[k] = -1 * v
	}

	//invert the integer
	for k, v := range *num2.Integer {
		(*num2Placeholder.Integer)[k] = -1 * v
	}

	return number__plus__number(num1, num2Placeholder, instance, stacktrace, line, file)
}
