package interpreter

import (
	"strconv"

	. "ka/lang/types"
)

//list of operations
//export Operations
var Operations = map[string]func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType{
	"number + number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		return number__plus__number(val1, val2, instance, stacktrace, line, file)
	},
	"number - number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		return number__minus__number(val1, val2, instance, stacktrace, line, file)
	},
	"number * number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		return number__times__number(val1, val2, instance, stacktrace, line, file)
	},
	"number / number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		return number__divide__number(val1, val2, instance, stacktrace, line, file)
	},
	"number % number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		return number__mod__number(val1, val2, instance, stacktrace, line, file)
	},
	"number ^ number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		return number__pow__number(val1, val2, instance, stacktrace, line, file)
	},
	"number == number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		var final = falsev

		if isEqual(val1.(KaNumber), val2.(KaNumber)) {
			final = truev
		}

		var finalType KaType = final

		return &finalType
	},
	"number != number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		var final = truev

		if isEqual(val1.(KaNumber), val2.(KaNumber)) {
			final = falsev
		}

		var finalType KaType = final

		return &finalType
	},
	"string == string": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		var isEqual KaType = falsev

		if val1.(KaString).ToGoType() == val2.(KaString).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"string != string": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		var isEqual KaType = truev

		if val1.(KaString).ToGoType() == val2.(KaString).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"bool == bool": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		var isEqual KaType = falsev

		if val1.(KaBool).ToGoType() == val2.(KaBool).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"bool != bool": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		var isEqual KaType = truev

		if val1.(KaBool).ToGoType() == val2.(KaBool).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"rune == rune": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		var isEqual KaType = falsev

		if val1.(KaRune).ToGoType() == val2.(KaRune).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"rune != rune": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		var isEqual KaType = truev

		if val1.(KaBool).ToGoType() == val2.(KaBool).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"none == none": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		var tmp KaType = truev
		return &tmp
	},
	"none != none": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		var tmp KaType = falsev
		return &tmp
	},
	"none ! bool": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		boolean := !val2.(KaBool).ToGoType()

		var converted KaType = KaBool{
			Boolean: &boolean,
		}

		return &converted
	},
	"number > number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		isGreaterv := !isLessOrEqual(val1.(KaNumber), val2.(KaNumber))
		var isGreaterType KaType = falsev

		if isGreaterv {
			isGreaterType = truev
		}

		return &isGreaterType
	},
	"number >= number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		isGreaterOrEqualv := !isLess(val1.(KaNumber), val2.(KaNumber))
		var isGreaterOrEqualType KaType = falsev

		if isGreaterOrEqualv {
			isGreaterOrEqualType = truev
		}

		return &isGreaterOrEqualType
	},
	"number < number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		isLessv := isLess(val1.(KaNumber), val2.(KaNumber))
		var isLessType KaType = falsev

		if isLessv {
			isLessType = truev
		}

		return &isLessType
	},
	"number <= number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		isLessOrEqualv := isLessOrEqual(val1.(KaNumber), val2.(KaNumber))
		var isLessOrEqualType KaType = falsev

		if isLessOrEqualv {
			isLessOrEqualType = truev
		}

		return &isLessOrEqualType
	},
	"array :: number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		//convert to int64
		idx := int64(val2.(KaNumber).ToGoType())
		arr := val1.(KaArray)

		if !arr.Exists(idx) {
			KaPanic("Index "+strconv.FormatInt(idx, 10)+" out of range with length "+strconv.FormatUint(arr.Length, 10), line, file, stacktrace)
		}

		return arr.At(idx)
	},
	"string :: number": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		//convert to int64
		idx := int64(val2.(KaNumber).ToGoType())
		str := val1.(KaString)

		if !str.Exists(idx) {
			KaPanic("Index "+strconv.FormatInt(idx, 10)+" out of range with length "+strconv.FormatUint(str.Length, 10), line, file, stacktrace)
		}

		var katype KaType = *str.At(idx)

		return &katype
	},
	"hash :: string": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		//convert index to go string
		gostr := val2.(KaString).ToGoType()

		return val1.(KaHash).At(gostr)
	},
	"proto :: string": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		v, e := val1.(KaProto).Get(val2.(KaString).ToGoType(), file)

		if e != nil {
			KaPanic(e.Error(), line, file, stacktrace)
		}

		return v
	},
	"object :: string": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		v, e := val1.(KaObject).Get(val2.(KaString).ToGoType(), file)

		if e != nil {
			KaPanic(e.Error(), line, file, stacktrace)
		}

		return v
	},
	"string + string": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		//alloc the space
		var space = make([]rune, val1.(KaString).Length+val2.(KaString).Length)

		var i uint
		var o uint

		val1l := val1.(KaString).ToRuneList()
		val12 := val2.(KaString).ToRuneList()

		for ; uint64(i) < val1.(KaString).Length; i++ {
			space[i] = val1l[i]
		}

		for ; uint64(o) < val2.(KaString).Length; i, o = i+1, o+1 {
			space[i] = val12[o]
		}

		var kastr KaString
		kastr.FromRuneList(space)
		var katype KaType = kastr
		return &katype
	},
	"string + rune": func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {

		//alloc the space
		var space = make([]rune, val1.(KaString).Length+1)

		var i uint

		val1l := val1.(KaString).ToRuneList()

		for ; uint64(i) < val1.(KaString).Length; i++ {
			space[i] = val1l[i]
		}

		space[i] = val2.(KaRune).ToGoType()

		var kastr KaString
		kastr.FromRuneList(space)
		var katype KaType = kastr
		return &katype
	},
}
