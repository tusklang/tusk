package types

import (
	"math"
	"strconv"
)

//export NumNormalize
func NumNormalize(num KaNumber) string {

	/*
	  ALGORITHM TO NORMALIZE:
	      Starting with this number:
	        [3412, -9912, 0001]
	      STEP 1:
	        loop through each number (from from decimal to integer)
	        in each iteration of the loop, if the number is the opposite of `isNeg` (meaning if `isNeg` is false, then the current value should be positive and vice versa)
	        use the following expression to get the complement
	          `MAX_DIGIT` + 1 - |`current num`|
	        replace the `current num` with this new value.
	        next, if `isNeg`, the digit to the right should be added by one, otherwise, subtract it by 1.
	        go to the next value and repeat
	      STEP 2:
	        join the vector of the integer and decimal with '.', then join each digit with ''
	        if `isNeg` then precede the string with a '-'
	        finally, return the result
	*/

	tmp := make([]int64, 1)

	if num.Integer == nil {
		num.Integer = &tmp
	}

	if num.Decimal == nil {
		num.Decimal = &tmp
	}

	integer := *num.Integer
	decimal := *num.Decimal

	//the first digit is actually the last index
	//because ka numbers are stored as so [1234, 5678, 9101] = 910, 156, 781, 234

	//alloc amounts into the copies
	var decimalCopy = make([]int64, len(decimal))
	var integerCopy = make([]int64, len(integer))

	//determine if it is less than 0

	combined := append(decimal, integer...)

	var negative int = 0 //0 means 0, -1 means negative, 1 means positive

	for i := len(combined) - 1; i >= 0; i-- {
		if combined[i] == 0 {
			continue
		}

		if combined[i] < 0 {
			negative = -1
		} else {
			negative = 1
		}

		break
	}
	////////////////////////////////

	var isNeg bool

	if negative == 0 {
		return "0"
	} else if negative == -1 {
		isNeg = true
	} else {
		isNeg = false
	}

	var carry int64 = 0

	for k := range decimal {
		decimalCopy[k] = decimal[k]
		decimalCopy[k] += carry
		curIsNeg := decimalCopy[k] < 0
		carry = 0

		if decimalCopy[k] == 0 {
			continue
		}

		if curIsNeg != isNeg {
			complement := MAX_DIGIT + 1 - int64(math.Abs(float64(decimalCopy[k])))
			decimalCopy[k] = complement

			if isNeg {
				carry = 1
			} else {
				carry = -1
			}
		}

		decimalCopy[k] = int64(math.Abs(float64(decimalCopy[k]))) //set the current digit to |current digit|
	}

	for k := range integer {
		integerCopy[k] = integer[k]
		integerCopy[k] += carry
		curIsNeg := integerCopy[k] < 0
		carry = 0

		if integerCopy[k] == 0 {
			continue
		}

		if curIsNeg != isNeg {
			complement := MAX_DIGIT + 1 - int64(math.Abs(float64(integerCopy[k])))
			integerCopy[k] = complement

			if isNeg {
				carry = 1
			} else {
				carry = -1
			}
		}

		integerCopy[k] = int64(math.Abs(float64(integerCopy[k]))) //set the current digit to |current digit|
	}

	for len(integerCopy) != 0 && integerCopy[len(integerCopy)-1] == 0 {
		integerCopy = integerCopy[:len(integerCopy)-1]
	}
	for len(decimalCopy) != 0 && decimalCopy[0] == 0 {
		decimalCopy = decimalCopy[1:]
	}

	var joined = ""

	for _, v := range decimalCopy {
		joined = strconv.FormatInt(v, 10) + joined
	}
	if len(decimalCopy) != 0 {
		joined = "." + joined
	}
	for _, v := range integerCopy {
		joined = strconv.FormatInt(v, 10) + joined
	}

	if len(joined) == 0 { //just in case it is all zero
		return "0"
	}

	if joined[0] == '.' {
		/*
		   if joined looks like: .123
		   convert it to 0.123
		*/
		joined = "0" + joined
	}
	if isNeg {
		joined = "-" + joined
	}

	return joined
}
