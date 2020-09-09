package interpreter

import . "github.com/omm-lang/omm/lang/types"

func naive_mul(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {
	num1, num2 := val1.(OmmNumber), val2.(OmmNumber)
	ensurePrec(&num1, &num2, (*instance).Params)

	var multFin [][]int64 //store the final values that were multiplied
	trailingZeroCount := 0

	//amount of decimal places there are
	decPlaceCount := len(*num1.Decimal) + len(*num2.Decimal)

	//actual numbers for num1 and num2
	num1n := append(*num1.Decimal, *num1.Integer...)
	num2n := append(*num2.Decimal, *num2.Integer...)

	if len(num1n) < len(num2n) { //swap num1n and num2n if num2n is greater (improves performance)
		num1n, num2n = num2n, num1n
	}

	for _, v := range num1n {
		//add a new row to the multiplied values
		multFin = append(multFin, []int64{})

		var carry int64 = 0 //variable to carry over overflowed numbers

		for i := 0; i < trailingZeroCount; i++ { //insert the trailing zeros
			multFin[len(multFin)-1] = append(multFin[len(multFin)-1], 0)
		}

		for _, sv := range num2n {

			var product int64 = v*sv + carry
			carry = 0 //reverrt carry after it was factored in

			if product > MAX_DIGIT || product < MIN_DIGIT {
				var rounded int64 = product / (MAX_DIGIT + 1) * (MAX_DIGIT + 1) //round down by `MAX_DIGIT`
				carry = product / (MAX_DIGIT + 1)                               //divide by the max digit to get the carry
				product -= rounded
			}

			multFin[len(multFin)-1] = append(multFin[len(multFin)-1], product)
		}

		multFin[len(multFin)-1] = append(multFin[len(multFin)-1], carry)
		trailingZeroCount++
	}

	totalSum := []int64{0}

	//multiply the values
	for _, v := range multFin {
		totalSumAct := zero //placeholder number to pass into add
		totalSumAct.Integer = &totalSum

		multFinAct := zero
		multFinAct.Integer = &v

		totalSum = *(*number__plus__number(totalSumAct, multFinAct, instance, stacktrace, line, file)).(OmmNumber).Integer
	}

	decimalRet := totalSum[:decPlaceCount]
	integerRet := totalSum[decPlaceCount:]

	returner := zero
	returner.Integer, returner.Decimal = &integerRet, &decimalRet

	var returnerType OmmType = returner

	return &returnerType
}

func number__times__number(num1, num2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

	//maybe switch to karatsuba later?
	//look into this: http://www.cburch.com/proj/karat/karat.txt

	return naive_mul(num1, num2, instance, stacktrace, line, file)
}
