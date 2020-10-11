package interpreter

//file that has all of the helper funcs for exponentiation

import . "github.com/tusklang/tusk/lang/types"

func number__pow__integer(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string) (TuskType, *TuskError) {
	num1, num2 := val1.(TuskNumber), val2.(TuskNumber)
	ensurePrec(&num1, &num2, (*instance).Params)

	//using binary exponentiation
	//https://cp-algorithms.com/algebra/binary-exp.html#toc-tgt-1

	if isEqual(num2, zero) {
		return one, nil
	}

	var two = zero
	two.Integer = &[]int64{2}

	pdivved, e := number__divide__number(num2, two, instance, stacktrace, line, file)
	if e != nil {
		return nil, e
	}
	divved := (*pdivved).(TuskNumber)
	divved.Decimal = &[]int64{} //round down to nearest whole

	pres, e := number__pow__integer(num1, divved, instance, stacktrace, line, file)
	if e != nil {
		return nil, e
	}
	res := pres.(TuskNumber)

	resSquared := (*number__times__number(res, res, instance, stacktrace, line, file)).(TuskNumber)

	pmodBy2, e := number__mod__number(num2, two, instance, stacktrace, line, file)
	if e != nil {
		return nil, e
	}
	modBy2 := (*pmodBy2).(TuskNumber)

	if isEqual(modBy2, one) {
		return (*number__times__number(resSquared, num1, instance, stacktrace, line, file)).(TuskNumber), nil
	}

	return resSquared, nil
}

func ln(val TuskType, instance *Instance, stacktrace []string, line uint64, file string) (TuskType, *TuskError) {
	x := val.(TuskNumber)
	ensurePrec(&x, &TuskNumber{}, (*instance).Params)

	//using taylor series expansion to calculate
	//found here https://www.efunda.com/math/taylor_series/logarithmic.cfm
	//algorithm 2

	var series = zero

	var two = zero
	two.Integer = &[]int64{2, 0}

	//calculate (x - 1) / (x + 1)
	pxm1dxp1, e := number__divide__number(*number__minus__number(x, one, instance, stacktrace, line, file), *number__plus__number(x, one, instance, stacktrace, line, file), instance, stacktrace, line, file)
	if e != nil {
		return nil, e
	}
	xm1dxp1 := *pxm1dxp1

	//convert precision to tusk number
	kaNumberPrec := zero
	kaNumberPrec.FromGoType(float64((*instance).Params.Prec))

	//calculate taylor series to prec
	for i := one; isLess(i, kaNumberPrec); i = (*number__plus__number(i, two, instance, stacktrace, line, file)).(TuskNumber) {

		//calculate 1/i
		ponedi, e := number__divide__number(one, i, instance, stacktrace, line, file)
		if e != nil {
			return nil, e
		}
		onedi := *ponedi

		//calculate xm1dxp1 ^ i
		xm1dxp1pi, e := number__pow__integer(xm1dxp1, i, instance, stacktrace, line, file)

		//calculate onedi * xm1dxp1pi
		oneditxm1dxp1pi := *number__times__number(onedi, xm1dxp1pi, instance, stacktrace, line, file)

		//add to the series
		series = (*number__plus__number(series, oneditxm1dxp1pi, instance, stacktrace, line, file)).(TuskNumber)
	}

	series = (*number__times__number(series, two, instance, stacktrace, line, file)).(TuskNumber)
	return series, nil
}

func fac(val TuskType, instance *Instance, stacktrace []string, line uint64, file string) TuskType {
	x := val.(TuskNumber)
	ensurePrec(&x, &TuskNumber{}, (*instance).Params)

	//factorial function for taylor series
	//using a naive method, but there is probably a faster method

	prod := one

	for i := one; isLessOrEqual(i, x); i = (*number__plus__number(i, one, instance, stacktrace, line, file)).(TuskNumber) {
		prod = (*number__times__number(prod, i, instance, stacktrace, line, file)).(TuskNumber)
	}

	return prod
}

func exp(val TuskType, instance *Instance, stacktrace []string, line uint64, file string) (TuskType, *TuskError) {
	x := val.(TuskNumber)
	ensurePrec(&x, &TuskNumber{}, (*instance).Params)

	//using taylor series expansion to calculate
	//found here https://www.efunda.com/math/taylor_series/exponential.cfm
	//algorithm 1

	var onePlaceholder = zero //temp value for one (because one wil get mutated if it is passed directly)
	onePlaceholder.Integer = &[]int64{1}

	var series TuskNumber = one

	//convert precision to tusk number
	kaNumberPrec := zero
	kaNumberPrec.FromGoType(float64((*instance).Params.Prec))

	for i := one; isLess(i, kaNumberPrec); i = (*number__plus__number(i, one, instance, stacktrace, line, file)).(TuskNumber) {
		//calculate i!
		i_factorial := fac(i, instance, stacktrace, line, file)

		//calculate x^i
		xpi, e := number__pow__integer(x, i, instance, stacktrace, line, file)
		if e != nil {
			return nil, e
		}

		//calculate x ^ i / (i!)
		pxpidifac, e := number__divide__number(xpi, i_factorial, instance, stacktrace, line, file)
		if e != nil {
			return nil, e
		}
		xpidifac := *pxpidifac

		//add x ^ i / (i!) to the series
		series = (*number__plus__number(series, xpidifac, instance, stacktrace, line, file)).(TuskNumber)
	}

	return series, nil
}
