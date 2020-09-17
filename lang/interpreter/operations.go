package interpreter

import (
	"strconv"

	. "github.com/tusklang/tusk/lang/types"
)

//list of operations
//export Operations
var Operations = map[string]func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType{
	"number + number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {
		return number__plus__number(val1, val2, instance, stacktrace, line, file)
	},
	"number - number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {
		return number__minus__number(val1, val2, instance, stacktrace, line, file)
	},
	"number * number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {
		return number__times__number(val1, val2, instance, stacktrace, line, file)
	},
	"number / number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {
		return number__divide__number(val1, val2, instance, stacktrace, line, file)
	},
	"number % number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {
		return number__mod__number(val1, val2, instance, stacktrace, line, file)
	},
	"number ^ number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {
		return number__pow__number(val1, val2, instance, stacktrace, line, file)
	},
	"number == number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		var final = falsev

		if isEqual(val1.(TuskNumber), val2.(TuskNumber)) {
			final = truev
		}

		var finalType TuskType = final

		return &finalType
	},
	"number != number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		var final = truev

		if isEqual(val1.(TuskNumber), val2.(TuskNumber)) {
			final = falsev
		}

		var finalType TuskType = final

		return &finalType
	},
	"string == string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		var isEqual TuskType = falsev

		if val1.(TuskString).ToGoType() == val2.(TuskString).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"string != string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		var isEqual TuskType = truev

		if val1.(TuskString).ToGoType() == val2.(TuskString).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"bool == bool": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		var isEqual TuskType = falsev

		if val1.(TuskBool).ToGoType() == val2.(TuskBool).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"bool != bool": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		var isEqual TuskType = truev

		if val1.(TuskBool).ToGoType() == val2.(TuskBool).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"rune == rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		var isEqual TuskType = falsev

		if val1.(TuskRune).ToGoType() == val2.(TuskRune).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"rune != rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		var isEqual TuskType = truev

		if val1.(TuskBool).ToGoType() == val2.(TuskBool).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"none == none": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {
		var tmp TuskType = truev
		return &tmp
	},
	"none != none": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {
		var tmp TuskType = falsev
		return &tmp
	},
	"none ! bool": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		boolean := !val2.(TuskBool).ToGoType()

		var converted TuskType = TuskBool{
			Boolean: &boolean,
		}

		return &converted
	},
	"number > number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		isGreaterv := !isLessOrEqual(val1.(TuskNumber), val2.(TuskNumber))
		var isGreaterType TuskType = falsev

		if isGreaterv {
			isGreaterType = truev
		}

		return &isGreaterType
	},
	"number >= number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		isGreaterOrEqualv := !isLess(val1.(TuskNumber), val2.(TuskNumber))
		var isGreaterOrEqualType TuskType = falsev

		if isGreaterOrEqualv {
			isGreaterOrEqualType = truev
		}

		return &isGreaterOrEqualType
	},
	"number < number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		isLessv := isLess(val1.(TuskNumber), val2.(TuskNumber))
		var isLessType TuskType = falsev

		if isLessv {
			isLessType = truev
		}

		return &isLessType
	},
	"number <= number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		isLessOrEqualv := isLessOrEqual(val1.(TuskNumber), val2.(TuskNumber))
		var isLessOrEqualType TuskType = falsev

		if isLessOrEqualv {
			isLessOrEqualType = truev
		}

		return &isLessOrEqualType
	},
	"array :: number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		//convert to int64
		idx := int64(val2.(TuskNumber).ToGoType())
		arr := val1.(TuskArray)

		if !arr.Exists(idx) {
			TuskPanic("Index "+strconv.FormatInt(idx, 10)+" out of range with length "+strconv.FormatUint(arr.Length, 10), line, file, stacktrace)
		}

		return arr.At(idx)
	},
	"string :: number": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		//convert to int64
		idx := int64(val2.(TuskNumber).ToGoType())
		str := val1.(TuskString)

		if !str.Exists(idx) {
			TuskPanic("Index "+strconv.FormatInt(idx, 10)+" out of range with length "+strconv.FormatUint(str.Length, 10), line, file, stacktrace)
		}

		var tusktype TuskType = *str.At(idx)

		return &tusktype
	},
	"hash :: string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		//convert index to go string
		gostr := val2.(TuskString).ToGoType()

		return val1.(TuskHash).At(gostr)
	},
	"proto :: string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {
		v, e := val1.(TuskProto).Get(val2.(TuskString).ToGoType(), file)

		if e != nil {
			TuskPanic(e.Error(), line, file, stacktrace)
		}

		return v
	},
	"object :: string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {
		v, e := val1.(TuskObject).Get(val2.(TuskString).ToGoType(), file)

		if e != nil {
			TuskPanic(e.Error(), line, file, stacktrace)
		}

		return v
	},
	"string + string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

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

		var kastr TuskString
		kastr.FromRuneList(space)
		var tusktype TuskType = kastr
		return &tusktype
	},
	"string + rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *TuskType {

		//alloc the space
		var space = make([]rune, val1.(TuskString).Length+1)

		var i uint

		val1l := val1.(TuskString).ToRuneList()

		for ; uint64(i) < val1.(TuskString).Length; i++ {
			space[i] = val1l[i]
		}

		space[i] = val2.(TuskRune).ToGoType()

		var kastr TuskString
		kastr.FromRuneList(space)
		var tusktype TuskType = kastr
		return &tusktype
	},
}
