package interpreter

//file that has all of the helper funcs for exponentiation

import . "ka/lang/types"

func number__pow__integer(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string) KaType {
	num1, num2 := val1.(KaNumber), val2.(KaNumber)
	ensurePrec(&num1, &num2, (*instance).Params)

	//using binary exponentiation
	//https://cp-algorithms.com/algebra/binary-exp.html#toc-tgt-1

	if isEqual(num2, zero) {
		return one
	}

	var two = zero
	two.Integer = &[]int64{2}

	divved := (*number__divide__number(num2, two, instance, stacktrace, line, file)).(KaNumber)
	divved.Decimal = &[]int64{} //round down to nearest whole

	res := number__pow__integer(num1, divved, instance, stacktrace, line, file).(KaNumber)

	resSquared := (*number__times__number(res, res, instance, stacktrace, line, file)).(KaNumber)

	modBy2 := (*number__mod__number(num2, two, instance, stacktrace, line, file)).(KaNumber)

	if isEqual(modBy2, one) {
		return (*number__times__number(resSquared, num1, instance, stacktrace, line, file)).(KaNumber)
	}

	return resSquared
}

func ln(val KaType, instance *Instance, stacktrace []string, line uint64, file string) KaType {
	x := val.(KaNumber)
	ensurePrec(&x, &KaNumber{}, (*instance).Params)

	//using taylor series expansion to calculate
	//found here https://www.efunda.com/math/taylor_series/logarithmic.cfm
	//algorithm 2

	var series = zero

	var two = zero
	two.Integer = &[]int64{2, 0}

	//calculate (x - 1) / (x + 1)
	xm1dxp1 := (*number__divide__number(*number__minus__number(x, one, instance, stacktrace, line, file), *number__plus__number(x, one, instance, stacktrace, line, file), instance, stacktrace, line, file)).(KaNumber)

	//convert precision to ka number
	kaNumberPrec := zero
	kaNumberPrec.FromGoType(float64((*instance).Params.Prec))

	//calculate taylor series to prec
	for i := one; isLess(i, kaNumberPrec); i = (*number__plus__number(i, two, instance, stacktrace, line, file)).(KaNumber) {

		//calculate 1/i
		onedi := *number__divide__number(one, i, instance, stacktrace, line, file)

		//calculate xm1dxp1 ^ i
		xm1dxp1pi := number__pow__integer(xm1dxp1, i, instance, stacktrace, line, file)

		//calculate onedi * xm1dxp1pi
		oneditxm1dxp1pi := *number__times__number(onedi, xm1dxp1pi, instance, stacktrace, line, file)

		//add to the series
		series = (*number__plus__number(series, oneditxm1dxp1pi, instance, stacktrace, line, file)).(KaNumber)
	}

	series = (*number__times__number(series, two, instance, stacktrace, line, file)).(KaNumber)
	return series
}

func fac(val KaType, instance *Instance, stacktrace []string, line uint64, file string) KaType {
	x := val.(KaNumber)
	ensurePrec(&x, &KaNumber{}, (*instance).Params)

	//factorial function for taylor series
	//using a naive method, but there is probably a faster method

	prod := one

	for i := one; isLessOrEqual(i, x); i = (*number__plus__number(i, one, instance, stacktrace, line, file)).(KaNumber) {
		prod = (*number__times__number(prod, i, instance, stacktrace, line, file)).(KaNumber)
	}

	return prod
}

func exp(val KaType, instance *Instance, stacktrace []string, line uint64, file string) KaType {
	x := val.(KaNumber)
	ensurePrec(&x, &KaNumber{}, (*instance).Params)

	//using taylor series expansion to calculate
	//found here https://www.efunda.com/math/taylor_series/exponential.cfm
	//algorithm 1

	var onePlaceholder = zero //temp value for one (because one wil get mutated if it is passed directly)
	onePlaceholder.Integer = &[]int64{1}

	var series KaNumber = one

	//convert precision to ka number
	kaNumberPrec := zero
	kaNumberPrec.FromGoType(float64((*instance).Params.Prec))

	for i := one; isLess(i, kaNumberPrec); i = (*number__plus__number(i, one, instance, stacktrace, line, file)).(KaNumber) {
		//calculate i!
		i_factorial := fac(i, instance, stacktrace, line, file)

		//calculate x^i
		xpi := number__pow__integer(x, i, instance, stacktrace, line, file)

		//calculate x ^ i / (i!)
		xpidifac := *number__divide__number(xpi, i_factorial, instance, stacktrace, line, file)

		//add x ^ i / (i!) to the series
		series = (*number__plus__number(series, xpidifac, instance, stacktrace, line, file)).(KaNumber)
	}

	return series
}
