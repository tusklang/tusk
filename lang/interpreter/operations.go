package interpreter

import (
	"math"
	"strconv"

	. "github.com/tusklang/tusk/lang/types"
	. "github.com/tusklang/tusk/native"
)

//list of operations
//export Operations
var Operations = map[string]func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError){
	"number + number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		return number__plus__number(val1, val2, instance, stacktrace, line, file), nil
	},
	"number - number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		return number__minus__number(val1, val2, instance, stacktrace, line, file), nil
	},
	"number * number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		return number__times__number(val1, val2, instance, stacktrace, line, file), nil
	},
	"number / number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		return number__divide__number(val1, val2, instance, stacktrace, line, file)
	},
	"number // number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		return number__floorDivide__number(val1, val2, instance, stacktrace, line, file)
	},
	"number % number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		return number__mod__number(val1, val2, instance, stacktrace, line, file)
	},
	"number ** number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		return number__pow__number(val1, val2, instance, stacktrace, line, file)
	},
	"number & number":  bitwiseAnd,
	"number | number":  bitwiseOr,
	"number ^ number":  bitwiseXor,
	"none ~ number":    bitwiseNot,
	"number >> number": bitwiseRShift,
	"number << number": bitwiseLShift,

	//arithmetic operators for runes
	"rune + rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		v := val1.(TuskRune).ToGoType() + val2.(TuskRune).ToGoType() //perform operation on runes, as go int32
		var tuskrune TuskRune
		tuskrune.FromGoType(v)
		var tusktype TuskType = tuskrune
		return &tusktype, nil
	},
	"rune - rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		v := val1.(TuskRune).ToGoType() - val2.(TuskRune).ToGoType() //perform operation on runes, as go int32
		var tuskrune TuskRune
		tuskrune.FromGoType(v)
		var tusktype TuskType = tuskrune
		return &tusktype, nil
	},
	"rune * rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		v := val1.(TuskRune).ToGoType() * val2.(TuskRune).ToGoType() //perform operation on runes, as go int32
		var tuskrune TuskRune
		tuskrune.FromGoType(v)
		var tusktype TuskType = tuskrune
		return &tusktype, nil
	},
	"rune / rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		v := val1.(TuskRune).ToGoType() / val2.(TuskRune).ToGoType() //perform operation on runes, as go int32
		var tuskrune TuskRune
		tuskrune.FromGoType(v)
		var tusktype TuskType = tuskrune
		return &tusktype, nil
	},
	"rune % rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		v := val1.(TuskRune).ToGoType() % val2.(TuskRune).ToGoType() //perform operation on runes, as go int32
		var tuskrune TuskRune
		tuskrune.FromGoType(v)
		var tusktype TuskType = tuskrune
		return &tusktype, nil
	},
	"rune ** rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		v := rune(math.Pow(float64(val1.(TuskRune).ToGoType()), float64(val2.(TuskRune).ToGoType()))) //perform operation on runes, as go int32
		var tuskrune TuskRune
		tuskrune.FromGoType(v)
		var tusktype TuskType = tuskrune
		return &tusktype, nil
	},
	////////////////////////////////

	"number == number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		var final = falsev

		if isEqual(val1.(TuskNumber), val2.(TuskNumber)) {
			final = truev
		}

		var finalType TuskType = final

		return &finalType, nil
	},
	"number != number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		var final = truev

		if isEqual(val1.(TuskNumber), val2.(TuskNumber)) {
			final = falsev
		}

		var finalType TuskType = final

		return &finalType, nil
	},
	"string == string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		var isEqual TuskType = falsev

		if val1.(TuskString).ToGoType() == val2.(TuskString).ToGoType() {
			isEqual = truev
		}

		return &isEqual, nil
	},
	"string != string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		var isEqual TuskType = truev

		if val1.(TuskString).ToGoType() == val2.(TuskString).ToGoType() {
			isEqual = falsev
		}

		return &isEqual, nil
	},
	"bool == bool": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		var isEqual TuskType = falsev

		if val1.(TuskBool).ToGoType() == val2.(TuskBool).ToGoType() {
			isEqual = truev
		}

		return &isEqual, nil
	},
	"bool != bool": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		var isEqual TuskType = truev

		if val1.(TuskBool).ToGoType() == val2.(TuskBool).ToGoType() {
			isEqual = falsev
		}

		return &isEqual, nil
	},
	"rune == rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		var isEqual TuskType = falsev

		if val1.(TuskRune).ToGoType() == val2.(TuskRune).ToGoType() {
			isEqual = truev
		}

		return &isEqual, nil
	},
	"rune != rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		var isEqual TuskType = truev

		if val1.(TuskRune).ToGoType() == val2.(TuskRune).ToGoType() {
			isEqual = falsev
		}

		return &isEqual, nil
	},
	"none == none": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		var tmp TuskType = truev
		return &tmp, nil
	},
	"none != none": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		var tmp TuskType = falsev
		return &tmp, nil
	},
	"none ! bool": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		boolean := !val2.(TuskBool).ToGoType()

		var converted TuskType = TuskBool{
			Boolean: &boolean,
		}

		return &converted, nil
	},
	"number > number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		isGreaterv := !isLessOrEqual(val1.(TuskNumber), val2.(TuskNumber))
		var isGreaterType TuskType = falsev

		if isGreaterv {
			isGreaterType = truev
		}

		return &isGreaterType, nil
	},
	"rune > rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		v := val1.(TuskRune).ToGoType() > val2.(TuskRune).ToGoType() //if the first arg (as int32) is greater than the second

		//value to return (return true if true and false if falsse)
		var ret TuskType = falsev
		if v {
			ret = truev
		}

		return &ret, nil
	},
	"number >= number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		isGreaterOrEqualv := !isLess(val1.(TuskNumber), val2.(TuskNumber))
		var isGreaterOrEqualType TuskType = falsev

		if isGreaterOrEqualv {
			isGreaterOrEqualType = truev
		}

		return &isGreaterOrEqualType, nil
	},
	"rune >= rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		v := val1.(TuskRune).ToGoType() >= val2.(TuskRune).ToGoType() //if the first arg (as int32) is greater than or equal to the second

		//value to return (return true if true and false if falsse)
		var ret TuskType = falsev
		if v {
			ret = truev
		}

		return &ret, nil
	},
	"number < number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		isLessv := isLess(val1.(TuskNumber), val2.(TuskNumber))
		var isLessType TuskType = falsev

		if isLessv {
			isLessType = truev
		}

		return &isLessType, nil
	},
	"rune < rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		v := val1.(TuskRune).ToGoType() < val2.(TuskRune).ToGoType() //if the first arg (as int32) is less than the second

		//value to return (return true if true and false if falsse)
		var ret TuskType = falsev
		if v {
			ret = truev
		}

		return &ret, nil
	},
	"number <= number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		isLessOrEqualv := isLessOrEqual(val1.(TuskNumber), val2.(TuskNumber))
		var isLessOrEqualType TuskType = falsev

		if isLessOrEqualv {
			isLessOrEqualType = truev
		}

		return &isLessOrEqualType, nil
	},
	"rune <= rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		v := val1.(TuskRune).ToGoType() <= val2.(TuskRune).ToGoType() //if the first arg (as int32) is less than or equal to the second

		//value to return (return true if true and false if falsse)
		var ret TuskType = falsev
		if v {
			ret = truev
		}

		return &ret, nil
	},
	"array :: number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		//convert to int64
		idx := int64(val2.(TuskNumber).ToGoType())
		arr := val1.(TuskArray)

		if !arr.Exists(idx) {
			TuskPanic("Index "+strconv.FormatInt(idx, 10)+" out of range with length "+strconv.FormatUint(arr.Length, 10), line, file, stacktrace)
		}

		return arr.At(idx), nil
	},
	"string :: number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		//convert to int64
		idx := int64(val2.(TuskNumber).ToGoType())
		str := val1.(TuskString)

		if !str.Exists(idx) {
			return nil, TuskPanic("Index "+strconv.FormatInt(idx, 10)+" out of range with length "+strconv.FormatUint(str.Length, 10), line, file, stacktrace)
		}

		var tusktype TuskType = *str.At(idx)

		return &tusktype, nil
	},
	"hash :: string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		return val1.(TuskHash).At(&val2), nil
	},
	"prototype :: string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		v, e := val1.(TuskProto).Get(val2.(TuskString).ToGoType(), file)

		if e != nil {
			return nil, TuskPanic(e.Error(), line, file, stacktrace)
		}

		return v, nil
	},
	"object :: string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		v, e := val1.(TuskObject).Get(val2.(TuskString).ToGoType(), file)

		if e != nil {
			return nil, TuskPanic(e.Error(), line, file, stacktrace)
		}

		return v, nil
	},
	"string + string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		//alloc the space
		var space = make([]rune, val1.(TuskString).Length+val2.(TuskString).Length)

		var i uint
		var o uint

		val1l := val1.(TuskString).ToRuneList()
		val12 := val2.(TuskString).ToRuneList()

		for ; uint64(i) < val1.(TuskString).Length; i++ {
			space[i] = val1l[i]
		}

		for ; uint64(o) < val2.(TuskString).Length; i, o = i+1, o+1 {
			space[i] = val12[o]
		}

		var tuskstr TuskString
		tuskstr.FromRuneList(space)
		var tusktype TuskType = tuskstr
		return &tusktype, nil
	},
	"string + rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {

		//alloc the space
		var space = make([]rune, val1.(TuskString).Length+1)

		var i uint

		val1l := val1.(TuskString).ToRuneList()

		for ; uint64(i) < val1.(TuskString).Length; i++ {
			space[i] = val1l[i]
		}

		space[i] = val2.(TuskRune).ToGoType()

		var tuskstr TuskString
		tuskstr.FromRuneList(space)
		var tusktype TuskType = tuskstr
		return &tusktype, nil
	},
}
