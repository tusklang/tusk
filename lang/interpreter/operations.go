package interpreter

import (
	"math"
	"strconv"

	. "github.com/tusklang/tusk/lang/types"
	"github.com/tusklang/tusk/native"
)

//list of operations
//export Operations
var Operations = map[string]func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string){

	"int + int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var intt TuskInt
		intt.FromGoType(val1.(TuskInt).Int + val2.(TuskInt).Int)
		var tusktype TuskType = intt
		return &tusktype, nil, ""
	},
	"int - int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var intt TuskInt
		intt.FromGoType(val1.(TuskInt).Int - val2.(TuskInt).Int)
		var tusktype TuskType = intt
		return &tusktype, nil, ""
	},
	"int * int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var intt TuskInt
		intt.FromGoType(val1.(TuskInt).Int * val2.(TuskInt).Int)
		var tusktype TuskType = intt
		return &tusktype, nil, ""
	},
	"int / int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		if val2.(TuskInt).Int == 0 {
			return nil, native.TuskPanic("Divide By Zero Error", line, file, stacktrace, native.ErrCodes["DBZ"]), ""
		}

		var intt TuskInt
		intt.FromGoType(val1.(TuskInt).Int / val2.(TuskInt).Int)
		var tusktype TuskType = intt
		return &tusktype, nil, ""
	},
	"int % int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		if val2.(TuskInt).Int == 0 {
			return nil, native.TuskPanic("Divide By Zero Error", line, file, stacktrace, native.ErrCodes["DBZ"]), ""
		}

		var intt TuskInt
		intt.FromGoType(val1.(TuskInt).Int % val2.(TuskInt).Int)
		var tusktype TuskType = intt
		return &tusktype, nil, ""
	},
	"int ** int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var intt TuskInt
		intt.FromGoType(int64(math.Pow(float64(val1.(TuskInt).Int), float64(val2.(TuskInt).Int))))
		var tusktype TuskType = intt
		return &tusktype, nil, ""
	},

	"float + float": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var floatt TuskFloat
		floatt.FromGoType(val1.(TuskFloat).Float + val2.(TuskFloat).Float)
		var tusktype TuskType = floatt
		return &tusktype, nil, ""
	},
	"float - float": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var floatt TuskFloat
		floatt.FromGoType(val1.(TuskFloat).Float - val2.(TuskFloat).Float)
		var tusktype TuskType = floatt
		return &tusktype, nil, ""
	},
	"float * float": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var floatt TuskFloat
		floatt.FromGoType(val1.(TuskFloat).Float * val2.(TuskFloat).Float)
		var tusktype TuskType = floatt
		return &tusktype, nil, ""
	},
	"float / float": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		if val2.(TuskFloat).Float == 0 {
			return nil, native.TuskPanic("Divide By Zero Error", line, file, stacktrace, native.ErrCodes["DBZ"]), ""
		}

		var floatt TuskFloat
		floatt.FromGoType(val1.(TuskFloat).Float / val2.(TuskFloat).Float)
		var tusktype TuskType = floatt
		return &tusktype, nil, ""
	},
	"float ** float": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var floatt TuskFloat
		floatt.FromGoType(math.Pow(val1.(TuskFloat).Float, val2.(TuskFloat).Float))
		var tusktype TuskType = floatt
		return &tusktype, nil, ""
	},

	"big + big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		return number__plus__number(val1, val2, instance, stacktrace, line, file), nil, ""
	},
	"big - big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		return number__minus__number(val1, val2, instance, stacktrace, line, file), nil, ""
	},
	"big * big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		return number__times__number(val1, val2, instance, stacktrace, line, file), nil, ""
	},
	"big / big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		tmp, e := number__divide__number(val1, val2, instance, stacktrace, line, file)
		return tmp, e, ""
	},
	"big // big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		tmp, e := number__floorDivide__number(val1, val2, instance, stacktrace, line, file)
		return tmp, e, ""
	},
	"big % big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		tmp, e := number__mod__number(val1, val2, instance, stacktrace, line, file)
		return tmp, e, ""
	},
	"big ** big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		tmp, e := number__pow__number(val1, val2, instance, stacktrace, line, file)
		return tmp, e, ""
	},

	"int & int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var t TuskType = TuskInt{
			Int: val1.(TuskInt).Int & val2.(TuskInt).Int,
		}
		return &t, nil, ""
	},
	"int | int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var t TuskType = TuskInt{
			Int: val1.(TuskInt).Int | val2.(TuskInt).Int,
		}
		return &t, nil, ""
	},
	"int ^ int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var t TuskType = TuskInt{
			Int: val1.(TuskInt).Int ^ val2.(TuskInt).Int,
		}
		return &t, nil, ""
	},
	"none ~ int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var t TuskType = TuskInt{
			Int: ^val2.(TuskInt).Int,
		}
		return &t, nil, ""
	},
	"int >> int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var t TuskType = TuskInt{
			Int: val1.(TuskInt).Int >> val2.(TuskInt).Int,
		}
		return &t, nil, ""
	},
	"int << int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var t TuskType = TuskInt{
			Int: val1.(TuskInt).Int << val2.(TuskInt).Int,
		}
		return &t, nil, ""
	},

	"big & big":  bitwiseAnd,
	"big | big":  bitwiseOr,
	"big ^ big":  bitwiseXor,
	"none ~ big": bitwiseNot,
	"big >> big": bitwiseRShift,
	"big << big": bitwiseLShift,

	"string & string": strbitwiseAnd,
	"string | string": strbitwiseOr,
	"string ^ string": strbitwiseXor,
	"none ~ string":   strbitwiseNot,

	//arithmetic operators for runes
	"rune + rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		v := val1.(TuskRune).ToGoType() + val2.(TuskRune).ToGoType() //perform operation on runes, as go int32
		var tuskrune TuskRune
		tuskrune.FromGoType(v)
		var tusktype TuskType = tuskrune
		return &tusktype, nil, ""
	},
	"rune - rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		v := val1.(TuskRune).ToGoType() - val2.(TuskRune).ToGoType() //perform operation on runes, as go int32
		var tuskrune TuskRune
		tuskrune.FromGoType(v)
		var tusktype TuskType = tuskrune
		return &tusktype, nil, ""
	},
	"rune * rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		v := val1.(TuskRune).ToGoType() * val2.(TuskRune).ToGoType() //perform operation on runes, as go int32
		var tuskrune TuskRune
		tuskrune.FromGoType(v)
		var tusktype TuskType = tuskrune
		return &tusktype, nil, ""
	},
	"rune / rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		v := val1.(TuskRune).ToGoType() / val2.(TuskRune).ToGoType() //perform operation on runes, as go int32
		var tuskrune TuskRune
		tuskrune.FromGoType(v)
		var tusktype TuskType = tuskrune
		return &tusktype, nil, ""
	},
	"rune % rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		v := val1.(TuskRune).ToGoType() % val2.(TuskRune).ToGoType() //perform operation on runes, as go int32
		var tuskrune TuskRune
		tuskrune.FromGoType(v)
		var tusktype TuskType = tuskrune
		return &tusktype, nil, ""
	},
	"rune ** rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		v := rune(math.Pow(float64(val1.(TuskRune).ToGoType()), float64(val2.(TuskRune).ToGoType()))) //perform operation on runes, as go int32
		var tuskrune TuskRune
		tuskrune.FromGoType(v)
		var tusktype TuskType = tuskrune
		return &tusktype, nil, ""
	},
	////////////////////////////////

	"int == int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var boolt TuskBool
		boolt.FromGoType(val1.(TuskInt).Int == val2.(TuskInt).Int)
		var tusktype TuskType = boolt
		return &tusktype, nil, ""
	},
	"int != int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var boolt TuskBool
		boolt.FromGoType(val1.(TuskInt).Int != val2.(TuskInt).Int)
		var tusktype TuskType = boolt
		return &tusktype, nil, ""
	},

	"float == float": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var boolt TuskBool
		boolt.FromGoType(val1.(TuskInt).Int == val2.(TuskInt).Int)
		var tusktype TuskType = boolt
		return &tusktype, nil, ""
	},
	"float != float": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var boolt TuskBool
		boolt.FromGoType(val1.(TuskFloat).Float != val2.(TuskFloat).Float)
		var tusktype TuskType = boolt
		return &tusktype, nil, ""
	},

	"big == big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		var final = falsev

		if isEqual(val1.(TuskNumber), val2.(TuskNumber)) {
			final = truev
		}

		var finalType TuskType = final

		return &finalType, nil, ""
	},
	"big != big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		var final = truev

		if isEqual(val1.(TuskNumber), val2.(TuskNumber)) {
			final = falsev
		}

		var finalType TuskType = final

		return &finalType, nil, ""
	},
	"string == string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		var isEqual TuskType = falsev

		if val1.(TuskString).ToGoType() == val2.(TuskString).ToGoType() {
			isEqual = truev
		}

		return &isEqual, nil, ""
	},
	"string != string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		var isEqual TuskType = truev

		if val1.(TuskString).ToGoType() == val2.(TuskString).ToGoType() {
			isEqual = falsev
		}

		return &isEqual, nil, ""
	},
	"bool == bool": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		var isEqual TuskType = falsev

		if val1.(TuskBool).ToGoType() == val2.(TuskBool).ToGoType() {
			isEqual = truev
		}

		return &isEqual, nil, ""
	},
	"bool != bool": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		var isEqual TuskType = truev

		if val1.(TuskBool).ToGoType() == val2.(TuskBool).ToGoType() {
			isEqual = falsev
		}

		return &isEqual, nil, ""
	},
	"rune == rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		var isEqual TuskType = falsev

		if val1.(TuskRune).ToGoType() == val2.(TuskRune).ToGoType() {
			isEqual = truev
		}

		return &isEqual, nil, ""
	},
	"rune != rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		var isEqual TuskType = truev

		if val1.(TuskRune).ToGoType() == val2.(TuskRune).ToGoType() {
			isEqual = falsev
		}

		return &isEqual, nil, ""
	},
	"none == none": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var tmp TuskType = truev
		return &tmp, nil, ""
	},
	"none != none": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var tmp TuskType = falsev
		return &tmp, nil, ""
	},
	"none ! bool": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		boolean := !val2.(TuskBool).ToGoType()

		var converted TuskType = TuskBool{
			Boolean: &boolean,
		}

		return &converted, nil, ""
	},
	"none ! none": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var tmp TuskType = truev //!undef is always true
		return &tmp, nil, ""
	},
	"int > int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var boolt TuskBool
		boolt.FromGoType(val1.(TuskInt).Int > val2.(TuskInt).Int)
		var tusktype TuskType = boolt
		return &tusktype, nil, ""
	},
	"float > float": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var boolt TuskBool
		boolt.FromGoType(val1.(TuskFloat).Float > val2.(TuskFloat).Float)
		var tusktype TuskType = boolt
		return &tusktype, nil, ""
	},
	"big > big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		isGreaterv := !isLessOrEqual(val1.(TuskNumber), val2.(TuskNumber))
		var isGreaterType TuskType = falsev

		if isGreaterv {
			isGreaterType = truev
		}

		return &isGreaterType, nil, ""
	},
	"rune > rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		v := val1.(TuskRune).ToGoType() > val2.(TuskRune).ToGoType() //if the first arg (as int32) is greater than the second

		//value to return (return true if true and false if falsse)
		var ret TuskType = falsev
		if v {
			ret = truev
		}

		return &ret, nil, ""
	},
	"int >= int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var boolt TuskBool
		boolt.FromGoType(val1.(TuskInt).Int >= val2.(TuskInt).Int)
		var tusktype TuskType = boolt
		return &tusktype, nil, ""
	},
	"float >= float": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var boolt TuskBool
		boolt.FromGoType(val1.(TuskFloat).Float >= val2.(TuskFloat).Float)
		var tusktype TuskType = boolt
		return &tusktype, nil, ""
	},
	"big >= big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		isGreaterOrEqualv := !isLess(val1.(TuskNumber), val2.(TuskNumber))
		var isGreaterOrEqualType TuskType = falsev

		if isGreaterOrEqualv {
			isGreaterOrEqualType = truev
		}

		return &isGreaterOrEqualType, nil, ""
	},
	"rune >= rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		v := val1.(TuskRune).ToGoType() >= val2.(TuskRune).ToGoType() //if the first arg (as int32) is greater than or equal to the second

		//value to return (return true if true and false if falsse)
		var ret TuskType = falsev
		if v {
			ret = truev
		}

		return &ret, nil, ""
	},
	"int < int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var boolt TuskBool
		boolt.FromGoType(val1.(TuskInt).Int < val2.(TuskInt).Int)
		var tusktype TuskType = boolt
		return &tusktype, nil, ""
	},
	"float < float": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var boolt TuskBool
		boolt.FromGoType(val1.(TuskFloat).Float < val2.(TuskFloat).Float)
		var tusktype TuskType = boolt
		return &tusktype, nil, ""
	},
	"big < big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		isLessv := isLess(val1.(TuskNumber), val2.(TuskNumber))
		var isLessType TuskType = falsev

		if isLessv {
			isLessType = truev
		}

		return &isLessType, nil, ""
	},
	"rune < rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		v := val1.(TuskRune).ToGoType() < val2.(TuskRune).ToGoType() //if the first arg (as int32) is less than the second

		//value to return (return true if true and false if falsse)
		var ret TuskType = falsev
		if v {
			ret = truev
		}

		return &ret, nil, ""
	},
	"int <= int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var boolt TuskBool
		boolt.FromGoType(val1.(TuskInt).Int <= val2.(TuskInt).Int)
		var tusktype TuskType = boolt
		return &tusktype, nil, ""
	},
	"float <= float": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		var boolt TuskBool
		boolt.FromGoType(val1.(TuskFloat).Float <= val2.(TuskFloat).Float)
		var tusktype TuskType = boolt
		return &tusktype, nil, ""
	},
	"big <= big": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		isLessOrEqualv := isLessOrEqual(val1.(TuskNumber), val2.(TuskNumber))
		var isLessOrEqualType TuskType = falsev

		if isLessOrEqualv {
			isLessOrEqualType = truev
		}

		return &isLessOrEqualType, nil, ""
	},
	"rune <= rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		v := val1.(TuskRune).ToGoType() <= val2.(TuskRune).ToGoType() //if the first arg (as int32) is less than or equal to the second

		//value to return (return true if true and false if falsse)
		var ret TuskType = falsev
		if v {
			ret = truev
		}

		return &ret, nil, ""
	},
	"array :: int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		idx := val2.(TuskInt).Int
		arr := val1.(TuskArray)

		if !arr.Exists(idx) {
			native.TuskPanic("Index "+strconv.FormatInt(idx, 10)+" out of range with length "+strconv.FormatUint(arr.Length(), 10), line, file, stacktrace, native.ErrCodes["IOB"])
		}

		return arr.At(idx), nil, ""
	},
	"string :: int": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		idx := val2.(TuskInt).Int
		str := val1.(TuskString)

		if !str.Exists(idx) {
			return nil, native.TuskPanic("Index "+strconv.FormatInt(idx, 10)+" out of range with length "+strconv.FormatUint(str.Length(), 10), line, file, stacktrace, native.ErrCodes["IOB"]), ""
		}

		var tusktype TuskType = *str.At(idx)

		return &tusktype, nil, ""
	},
	"hash :: string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		return val1.(TuskHash).At(&val2), nil, ""
	},
	"prototype :: string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		v, e := val1.(TuskProto).Get(val2.(TuskString).ToGoType(), file, namespace)

		if e != nil {
			return nil, native.TuskPanic(e.Error(), line, file, stacktrace, native.ErrCodes["ITEMNOTFOUND"]), ""
		}

		return v, nil, ""
	},
	"object :: string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		v, e := val1.(TuskObject).Get(val2.(TuskString).ToGoType(), file, namespace)

		if e != nil {
			return nil, native.TuskPanic(e.Error(), line, file, stacktrace, native.ErrCodes["ITEMNOTFOUND"]), ""
		}

		return v, nil, ""
	},
	"string + string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		//alloc the space
		var space = make([]rune, val1.(TuskString).Length()+val2.(TuskString).Length())

		var i uint
		var o uint

		val1l := val1.(TuskString).ToRuneList()
		val12 := val2.(TuskString).ToRuneList()

		for ; uint64(i) < val1.(TuskString).Length(); i++ {
			space[i] = val1l[i]
		}

		for ; uint64(o) < val2.(TuskString).Length(); i, o = i+1, o+1 {
			space[i] = val12[o]
		}

		var tuskstr TuskString
		tuskstr.FromRuneList(space)
		var tusktype TuskType = tuskstr
		return &tusktype, nil, ""
	},
	"string + rune": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		//alloc the space
		var space = make([]rune, val1.(TuskString).Length()+1)

		var i uint

		val1l := val1.(TuskString).ToRuneList()

		for ; uint64(i) < val1.(TuskString).Length(); i++ {
			space[i] = val1l[i]
		}

		space[i] = val2.(TuskRune).ToGoType()

		var tuskstr TuskString
		tuskstr.FromRuneList(space)
		var tusktype TuskType = tuskstr
		return &tusktype, nil, ""
	},
	"clib :: string": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {
		proc := val1.(native.CLibrary).GetProc(val2.(TuskString).ToGoType())
		var tusktype TuskType = proc
		return &tusktype, nil, ""
	},
	"cproc : array": func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint, namespace string) (*TuskType, *TuskError, string) {

		args := make([]*TuskType, val2.(TuskArray).Length())
		val2.(TuskArray).Range(func(k, v *TuskType) (Returner, *TuskError) {
			args[(*k).(TuskInt).ToGoType()] = v
			return Returner{}, nil
		})
		ret, e := val1.(native.CProc).Call(args)

		if e != nil {
			return nil, native.TuskPanic(e.Error(), line, file, stacktrace, native.ErrCodes["INVALIDARG"]), ""
		}

		var tuskint TuskInt
		tuskint.FromGoType(ret)
		var tusktype TuskType = tuskint
		return &tusktype, nil, ""
	},
}
