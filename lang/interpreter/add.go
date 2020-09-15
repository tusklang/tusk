package interpreter

import . "ka/lang/types"

func number__plus__number(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string) *KaType {
	num1, num2 := val1.(KaNumber), val2.(KaNumber)
	ensurePrec(&num1, &num2, (*instance).Params)

	int1 := *num1.Integer
	int2 := *num2.Integer
	dec1 := *num1.Decimal
	dec2 := *num2.Decimal

	//ensure that the lengths of num1 are larger than the lengths of num2
	if len(int1) < len(int2) {
		int1, int2 = int2, int1
	}
	if len(dec1) < len(dec2) {
		dec1, dec2 = dec2, dec1
	}

	//got this from the python source code
	//https://github.com/python/cpython/blob/master/Objects/longobject.c#L3003
	//basically, it does not require an "if" statement to see if num2 has an element at an index

	var carry int64 = 0

	//do the decimal
	newDec := make([]int64, len(dec1))
	dec1Len := len(dec1)
	dec2Len := len(dec2)

	var deci1 int
	var deci2 int

	for ; deci1 < dec1Len-dec2Len; deci1++ {
		added := dec1[deci1] + carry
		carry = 0

		if added > MAX_DIGIT {
			carry = 1
			added -= MAX_DIGIT + 1
		}
		if added < MIN_DIGIT {
			carry = -1
			added -= MIN_DIGIT - 1
		}

		newDec[deci1] = added
	}
	for ; deci1 < dec1Len; deci1, deci2 = deci1+1, deci2+1 {
		added := dec1[deci1] + dec2[deci2] + carry
		carry = 0

		if added > MAX_DIGIT {
			carry = 1
			added -= MAX_DIGIT + 1
		}
		if added < MIN_DIGIT {
			carry = -1
			added -= MIN_DIGIT - 1
		}

		newDec[deci1] = added
	}

	//do the integer
	newInt := make([]int64, len(int1))
	int1Len := len(int1)
	int2Len := len(int2)

	var inti int

	for ; inti < int2Len; inti++ {
		added := int1[inti] + int2[inti] + carry
		carry = 0

		if added > MAX_DIGIT {
			carry = 1
			added -= MAX_DIGIT + 1
		}
		if added < MIN_DIGIT {
			carry = -1
			added -= MIN_DIGIT - 1
		}

		newInt[inti] = added
	}
	for ; inti < int1Len; inti++ {
		added := int1[inti] + carry
		carry = 0

		if added > MAX_DIGIT {
			carry = 1
			added -= MAX_DIGIT + 1
		}
		if added < MIN_DIGIT {
			carry = -1
			added -= MIN_DIGIT - 1
		}

		newInt[inti] = added
	}

	if carry != 0 {
		newInt = append(newInt, carry) //append the final carry to the new integer
	}

	number := zero
	number.Integer = &newInt
	number.Decimal = &newDec

	var numtype KaType = number

	return &numtype
}
