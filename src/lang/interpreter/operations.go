package interpreter

import "strconv"
import . "lang/types"

//list of operations
//export Operations
var Operations = map[string]func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {
	"number + number": number__plus__number,
	"number - number": number__minus__number,
	"number * number": number__times__number,
	"number / number": number__divide__number,
	"number % number": number__mod__number,
	"number ^ number": number__pow__number,
	"number == number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		var final = falsev

		if isEqual(val1.(OmmNumber), val2.(OmmNumber)) {
			final = truev
		}

		var finalType OmmType = final

		return &finalType
	},
	"number != number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		var final = truev

		if isEqual(val1.(OmmNumber), val2.(OmmNumber)) {
			final = falsev
		}

		var finalType OmmType = final

		return &finalType
	},
	"string == string": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		var isEqual OmmType = falsev

		if val1.(OmmString).ToGoType() == val2.(OmmString).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"string != string": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		var isEqual OmmType = truev

		if val1.(OmmString).ToGoType() == val2.(OmmString).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"bool == bool": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		var isEqual OmmType = falsev

		if val1.(OmmBool).ToGoType() == val2.(OmmBool).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"bool != bool": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		var isEqual OmmType = truev

		if val1.(OmmBool).ToGoType() == val2.(OmmBool).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"rune == rune": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		var isEqual OmmType = falsev

		if val1.(OmmRune).ToGoType() == val2.(OmmRune).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"rune != rune": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		var isEqual OmmType = truev

		if val1.(OmmBool).ToGoType() == val2.(OmmBool).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"undef == undef": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {
		var tmp OmmType = truev
		return &tmp
	},
	"undef != undef": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {
		var tmp OmmType = falsev
		return &tmp
	},
	"undef ! bool": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		boolean := !val2.(OmmBool).ToGoType()

		var converted OmmType = OmmBool{
			Boolean: &boolean,
		}

		return &converted
	},
	"number > number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		isGreaterv := !isLessOrEqual(val1.(OmmNumber), val2.(OmmNumber))
		var isGreaterType OmmType = falsev

		if isGreaterv {
			isGreaterType = truev
		}

		return &isGreaterType
	},
	"number >= number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		isGreaterOrEqualv := !isLess(val1.(OmmNumber), val2.(OmmNumber))
		var isGreaterOrEqualType OmmType = falsev

		if isGreaterOrEqualv {
			isGreaterOrEqualType = truev
		}

		return &isGreaterOrEqualType
	},
	"number < number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		isLessv := isLess(val1.(OmmNumber), val2.(OmmNumber))
		var isLessType OmmType = falsev

		if isLessv {
			isLessType = truev
		}

		return &isLessType
	},
	"number <= number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		isLessOrEqualv := isLessOrEqual(val1.(OmmNumber), val2.(OmmNumber))
		var isLessOrEqualType OmmType = falsev

		if isLessOrEqualv {
			isLessOrEqualType = truev
		}

		return &isLessOrEqualType
	},
	"array :: number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		//convert to int64
		idx := int64(val2.(OmmNumber).ToGoType())
		arr := val1.(OmmArray)

		if !arr.Exists(idx) {
			OmmPanic("Index " + strconv.FormatInt(idx, 10) + " out of range with length " + strconv.FormatUint(arr.Length, 10), line, file, stacktrace)
		}

		return arr.At(idx)
	},
	"string :: number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		//convert to int64
		idx := int64(val2.(OmmNumber).ToGoType())
		str := val1.(OmmString)

		if !str.Exists(idx) {
			OmmPanic("Index " + strconv.FormatInt(idx, 10) + " out of range with length " + strconv.FormatUint(str.Length, 10), line, file, stacktrace)
		}

		var ommtype OmmType = *str.At(idx)

		return &ommtype
	},
	"hash :: string": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		//convert index to go string
		gostr := val2.(OmmString).ToGoType()

		return val1.(OmmHash).At(gostr)
	},
	"proto :: string": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		//convert field to go string
		gostr := val2.(OmmString).ToGoType()

		if gostr[0] == '_' {
			OmmPanic("Cannot access private member: " + gostr, line, file, stacktrace)
		}
		
		field := val1.(OmmProto).GetStatic(gostr)

		if field == nil {
			OmmPanic("Class does not contain the field \"" + gostr + "\"", line, file, stacktrace)
		}

		return field
	},
	"object :: string": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		//convert field to go string
		gostr := val2.(OmmString).ToGoType()

		if gostr[0] == '_' {
			OmmPanic("Cannot access private member: " + gostr, line, file, stacktrace)
		}

		field := val1.(OmmObject).GetInstance(gostr)

		if field == nil {
			OmmPanic("Object does not contain the field \"" + gostr + "\"", line, file, stacktrace)
		}

		return field
	},
	"string + string": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		//alloc the space
		var space = make([]rune, val1.(OmmString).Length + val2.(OmmString).Length)

		var i uint
		var o uint

		val1l := val1.(OmmString).ToRuneList()
		val12 := val2.(OmmString).ToRuneList()

		for ;uint64(i) < val1.(OmmString).Length; i, o = i + 1, o + 1 {
			space[i] = val1l[i]
		}

		for ;uint64(i) < val1.(OmmString).Length; i, o = i + 1, o + 1 {
			space[i] = val12[i - o]
		}

		var ommstr OmmString
		ommstr.FromRuneList(space)
		var ommtype OmmType = ommstr
		return &ommtype
	},
	"string + rune": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {

		//alloc the space
		var space = make([]rune, val1.(OmmString).Length + 1)

		var i uint

		val1l := val1.(OmmString).ToRuneList()

		for ;uint64(i) < val1.(OmmString).Length; i++ {
			space[i] = val1l[i]
		}

		space[i] = val2.(OmmRune).ToGoType()

		var ommstr OmmString
		ommstr.FromRuneList(space)
		var ommtype OmmType = ommstr
		return &ommtype
	},
}
