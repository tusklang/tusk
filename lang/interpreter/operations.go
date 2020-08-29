package interpreter

import (
	"strconv"

	. "github.com/omm-lang/omm/lang/types"
	. "github.com/omm-lang/omm/stdlib/native"
)

//list of operations
//export Operations
var Operations = map[string]func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType{
	"number + number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		return number__plus__number(val1, val2, instance, stacktrace, line, file)
	},
	"number - number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		return number__minus__number(val1, val2, instance, stacktrace, line, file)
	},
	"number * number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		return number__times__number(val1, val2, instance, stacktrace, line, file)
	},
	"number / number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		return number__divide__number(val1, val2, instance, stacktrace, line, file)
	},
	"number % number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		return number__mod__number(val1, val2, instance, stacktrace, line, file)
	},
	"number ^ number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		return number__pow__number(val1, val2, instance, stacktrace, line, file)
	},
	"number == number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		var final = falsev

		if isEqual(val1.(OmmNumber), val2.(OmmNumber)) {
			final = truev
		}

		var finalType OmmType = final

		return &finalType
	},
	"number != number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		var final = truev

		if isEqual(val1.(OmmNumber), val2.(OmmNumber)) {
			final = falsev
		}

		var finalType OmmType = final

		return &finalType
	},
	"string == string": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		var isEqual OmmType = falsev

		if val1.(OmmString).ToGoType() == val2.(OmmString).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"string != string": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		var isEqual OmmType = truev

		if val1.(OmmString).ToGoType() == val2.(OmmString).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"bool == bool": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		var isEqual OmmType = falsev

		if val1.(OmmBool).ToGoType() == val2.(OmmBool).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"bool != bool": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		var isEqual OmmType = truev

		if val1.(OmmBool).ToGoType() == val2.(OmmBool).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"rune == rune": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		var isEqual OmmType = falsev

		if val1.(OmmRune).ToGoType() == val2.(OmmRune).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"rune != rune": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		var isEqual OmmType = truev

		if val1.(OmmBool).ToGoType() == val2.(OmmBool).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"undef == undef": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		var tmp OmmType = truev
		return &tmp
	},
	"undef != undef": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		var tmp OmmType = falsev
		return &tmp
	},
	"undef ! bool": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		boolean := !val2.(OmmBool).ToGoType()

		var converted OmmType = OmmBool{
			Boolean: &boolean,
		}

		return &converted
	},
	"number > number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		isGreaterv := !isLessOrEqual(val1.(OmmNumber), val2.(OmmNumber))
		var isGreaterType OmmType = falsev

		if isGreaterv {
			isGreaterType = truev
		}

		return &isGreaterType
	},
	"number >= number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		isGreaterOrEqualv := !isLess(val1.(OmmNumber), val2.(OmmNumber))
		var isGreaterOrEqualType OmmType = falsev

		if isGreaterOrEqualv {
			isGreaterOrEqualType = truev
		}

		return &isGreaterOrEqualType
	},
	"number < number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		isLessv := isLess(val1.(OmmNumber), val2.(OmmNumber))
		var isLessType OmmType = falsev

		if isLessv {
			isLessType = truev
		}

		return &isLessType
	},
	"number <= number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		isLessOrEqualv := isLessOrEqual(val1.(OmmNumber), val2.(OmmNumber))
		var isLessOrEqualType OmmType = falsev

		if isLessOrEqualv {
			isLessOrEqualType = truev
		}

		return &isLessOrEqualType
	},
	"array :: number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		//convert to int64
		idx := int64(val2.(OmmNumber).ToGoType())
		arr := val1.(OmmArray)

		if !arr.Exists(idx) {
			OmmPanic("Index "+strconv.FormatInt(idx, 10)+" out of range with length "+strconv.FormatUint(arr.Length, 10), line, file, stacktrace)
		}

		return arr.At(idx)
	},
	"string :: number": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		//convert to int64
		idx := int64(val2.(OmmNumber).ToGoType())
		str := val1.(OmmString)

		if !str.Exists(idx) {
			OmmPanic("Index "+strconv.FormatInt(idx, 10)+" out of range with length "+strconv.FormatUint(str.Length, 10), line, file, stacktrace)
		}

		var ommtype OmmType = *str.At(idx)

		return &ommtype
	},
	"hash :: string": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		//convert index to go string
		gostr := val2.(OmmString).ToGoType()

		return val1.(OmmHash).At(gostr)
	},
	"proto :: string": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		//convert field to go string
		gostr := val2.(OmmString).ToGoType()

		if gostr[0] == '_' {
			OmmPanic("Cannot access private member: "+gostr, line, file, stacktrace)
		}

		//check for access (protected)

		if val1.(OmmProto).AccessList[gostr] == nil { //if it does not name any access, automatically make it public
			goto allowed
		}

		for _, v := range val1.(OmmProto).AccessList[gostr] {
			if file == v {
				goto allowed
			}
		}

		OmmPanic("File cannot acces field \""+gostr+"\"", line, file, stacktrace)

	allowed:
		field := val1.(OmmProto).GetStatic(gostr)

		if field == nil {
			OmmPanic("Prototype does not contain the field \""+gostr+"\"", line, file, stacktrace)
		}

		return field
	},
	"object :: string": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		//convert field to go string
		gostr := val2.(OmmString).ToGoType()

		if gostr[0] == '_' {
			OmmPanic("Cannot access private member: "+gostr, line, file, stacktrace)
		}

		//check for access (protected)

		if val1.(OmmObject).AccessList[gostr] == nil { //if it does not name any access, automatically make it public
			goto allowed
		}

		for _, v := range val1.(OmmObject).AccessList[gostr] {
			if file == v {
				goto allowed
			}
		}

		OmmPanic("File cannot acces field \""+gostr+"\"", line, file, stacktrace)

	allowed:
		field := val1.(OmmObject).GetInstance(gostr)

		if field == nil {
			OmmPanic("Object does not contain the field \""+gostr+"\"", line, file, stacktrace)
		}

		return field
	},
	"string + string": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		//alloc the space
		var space = make([]rune, val1.(OmmString).Length+val2.(OmmString).Length)

		var i uint
		var o uint

		val1l := val1.(OmmString).ToRuneList()
		val12 := val2.(OmmString).ToRuneList()

		for ; uint64(i) < val1.(OmmString).Length; i++ {
			space[i] = val1l[i]
		}

		for ; uint64(o) < val2.(OmmString).Length; i, o = i+1, o+1 {
			space[i] = val12[o]
		}

		var ommstr OmmString
		ommstr.FromRuneList(space)
		var ommtype OmmType = ommstr
		return &ommtype
	},
	"string + rune": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		//alloc the space
		var space = make([]rune, val1.(OmmString).Length+1)

		var i uint

		val1l := val1.(OmmString).ToRuneList()

		for ; uint64(i) < val1.(OmmString).Length; i++ {
			space[i] = val1l[i]
		}

		space[i] = val2.(OmmRune).ToGoType()

		var ommstr OmmString
		ommstr.FromRuneList(space)
		var ommtype OmmType = ommstr
		return &ommtype
	},
	"int + int": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		var (
			int1 = val1.(OmmInteger).Goint
			int2 = val2.(OmmInteger).Goint
		)

		var final OmmInteger
		final.Goint = int1 + int2
		var ret OmmType = final
		return &ret
	},
	"float + float": func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {

		var (
			float1 = val1.(OmmFloat).Gofloat
			float2 = val2.(OmmFloat).Gofloat
		)

		var final OmmFloat
		final.Gofloat = float1 + float2
		var ret OmmType = final
		return &ret
	},
}
