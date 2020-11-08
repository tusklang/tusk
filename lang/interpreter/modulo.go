package interpreter

import . "github.com/tusklang/tusk/lang/types"

func number__mod__number(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string) (*TuskType, *TuskError) {
	num1, num2 := val1.(TuskNumber), val2.(TuskNumber)
	ensurePrec(&num1, &num2, (*instance).Params)

	//ALGORITHM:
	//  num1 - num1 // num2 * num2
	//floor divide num1 // num2
	pdivided, e := number__floorDivide__number(num1, num2, instance, stacktrace, line, file)
	if e != nil {
		return nil, e
	}
	divided := (*pdivided).(TuskNumber)

	multiplied := (*number__times__number(divided, num2, instance, stacktrace, line, file)).(TuskNumber)
	return number__minus__number(num1, multiplied, instance, stacktrace, line, file), nil
}
