package native

//all of the gofuncs
//functions written in go that are used by tusk

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unsafe"

	oatenc "github.com/tusklang/oat/format/encoding"
	. "github.com/tusklang/tusk/lang/types"
	tusksys "github.com/tusklang/tusk/native/systables"
)

//#include <stdbool.h>
//#include "syscall.h"
//#include "arrayc.h"
//#include "exec.h"
import "C"

//Native stores all of the native values. You can make your own by just putting it into this map
var Native = make(map[string]*TuskType)

func tusksprint(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) string {
	sprinted := ""

	for k, v := range args {
		sprinted += (*v).Format()
		if k+1 != len(args) {
			sprinted += " "
		}
	}

	return sprinted
}

//TuskPanic panics in an Tusk instance
func TuskPanic(err string, line uint64, file string, stacktrace []string) {
	fmt.Println("Panic on line", line, "file", file)
	fmt.Print(err)
	fmt.Print("\nWhen the error was thrown, this was the stack:\n")
	fmt.Println("  at line", line, "in file", file)
	for i := len(stacktrace) - 1; i >= 0; i-- { //print the stacktrace

		endl := "\n"
		if i == 0 {
			endl = ""
		}

		fmt.Print("  " + stacktrace[i] + endl)
	}
	fmt.Println()
	os.Exit(1)
}

func makectype(val *TuskType) (unsafe.Pointer, error) {
	switch (*val).(type) {
	case TuskNumber:
		cnum := C.int((*val).(TuskNumber).ToGoType())
		return C.makeunsafeint(cnum), nil
	case TuskString:
		cstr := C.CString((*val).(TuskString).ToGoType())
		return unsafe.Pointer(cstr), nil
	case TuskArray:
		carray := C.makecarray(C.long((*val).(TuskArray).Length))

		var err error

		(*val).(TuskArray).Range(func(k, v *TuskType) Returner {
			idx := (*k).(TuskNumber).ToGoType()
			cval, e := makectype(v)
			if e != nil {
				err = e //if there is an error, set it, then break the loop
				return Returner{
					Type: "break",
				}
			}
			C.setcarray(carray, C.int(idx), unsafe.Pointer(cval))
			return Returner{}
		})

		return unsafe.Pointer(carray), err
	default:
		return nil, fmt.Errorf("Cannot convert type %s to ctype", (*val).Type())
	}
}

func fromctype(val unsafe.Pointer, tuskarg *TuskType) TuskType {
	switch (*tuskarg).(type) {
	case TuskNumber:
		cnum := C.makeintfromunsafe(val)
		var tusknum TuskNumber
		tusknum.FromGoType(float64(cnum))
		return tusknum
	case TuskString:
		ccstr := (*C.char)(val)
		var tuskstr TuskString
		tuskstr.FromGoType(C.GoString(ccstr))
		return tuskstr
	case TuskArray:
		carray := (*unsafe.Pointer)(val)
		var tuskarray TuskArray
		(*tuskarg).(TuskArray).Range(func(k, v *TuskType) Returner {
			idx := (*k).(TuskNumber).ToGoType()
			safeval := fromctype(C.getcarray(carray, C.int(idx)), v)
			tuskarray.PushBack(safeval)
			return Returner{}
		})
		return tuskarray
	}

	return nil
}

//NativeFuncs are the native functions that are relatively simple to implement
var NativeFuncs = map[string]func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType{
	"log": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {
		fmt.Print(tusksprint(args, stacktrace, line, file, instance))
		fmt.Println() //print a newline at the end
		var tmpundef TuskType = TuskUndef{}
		return &tmpundef
	},
	"print": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {
		fmt.Print(tusksprint(args, stacktrace, line, file, instance))
		var tmpundef TuskType = TuskUndef{}
		return &tmpundef
	},
	"sprint": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {
		sprinted := tusksprint(args, stacktrace, line, file, instance)
		var tuskstr TuskString
		tuskstr.FromGoType(sprinted)
		var tusktype TuskType = tuskstr
		return &tusktype
	},
	"await": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) != 1 || (*args[0]).Type() != "thread" {
			TuskPanic("Function await requires (thread)", line, file, stacktrace)
		}

		interpreted := args[0]
		var awaited *TuskType

		switch (*interpreted).(type) {
		case TuskThread:

			//put the new value back into the given interpreted pointer
			thread := (*interpreted).(TuskThread)
			awaited = thread.Join()
			*interpreted = thread
			///////////////////////////////////////////////////////////

		default:
			awaited = interpreted
		}

		return awaited
	},
	"input": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		scanner := bufio.NewScanner(os.Stdin)

		if len(args) == 0 {
			//if it has 0 or 1 args, there is no error
		} else if len(args) == 1 {

			switch (*args[0]).(type) {
			case TuskString:
				str := (*args[0]).(TuskString).ToGoType()
				fmt.Print(str)
			default:
				TuskPanic("Expected a string as the argument to input[]", line, file, stacktrace)
			}

		} else {
			TuskPanic("Function input requires a parameter count of 0 or 1", line, file, stacktrace)
		}

		//get user input and convert it to TuskType
		scanner.Scan()
		input := scanner.Text()
		var inputTuskType TuskString
		inputTuskType.FromGoType(input)
		var inputType TuskType = inputTuskType

		return &inputType
	},
	"typeof": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) != 1 {
			TuskPanic("Function typeof requires a parameter count of 1", line, file, stacktrace)
		}

		typeof := (*args[0]).TypeOf()

		var str TuskString
		str.FromGoType(typeof)

		//convert to TuskType interface
		var tusktype TuskType = str
		return &tusktype
	},
	"append": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) != 2 {
			TuskPanic("Function append requires a parameter count of 2", line, file, stacktrace)
		}

		if (*args[0]).Type() != "array" {
			TuskPanic("Function append requires (array, any)", line, file, stacktrace)
		}

		a := (*args[0]).(TuskArray)
		a.PushBack(*args[1])
		var tusktype TuskType = a
		return &tusktype
	},
	"prepend": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) != 2 {
			TuskPanic("Function prepend requires a parameter count of 2", line, file, stacktrace)
		}

		if (*args[0]).Type() != "array" {
			TuskPanic("Function prepend requires the first argument to be an array", line, file, stacktrace)
		}

		a := (*args[0]).(TuskArray)
		a.PushFront(*args[1])
		var tusktype TuskType = a
		return &tusktype
	},
	"exit": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) == 1 {

			switch (*args[0]).(type) {
			case TuskNumber:

				var gonum = (*args[0]).(TuskNumber).ToGoType()
				os.Exit(int(gonum))

			case TuskBool:

				if (*args[0]).(TuskBool).ToGoType() == true {
					os.Exit(0)
				} else {
					os.Exit(1)
				}

			default:
				os.Exit(0)
			}

		} else if len(args) == 0 {
			os.Exit(0)
		} else {
			TuskPanic("Function exit requires a parameter count of 1 or 0", line, file, stacktrace)
		}

		var tmpundef TuskType = TuskUndef{}
		return &tmpundef
	},
	"make": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		errmsg := "Function make requires the signature (prototype)"

		if len(args) == 1 {

			switch (*args[0]).(type) {
			case TuskProto:
				var tusktype TuskType = (*args[0]).(TuskProto).New(*instance)
				return &tusktype
			default:
				TuskPanic(errmsg, line, file, stacktrace)
			}

		} else {
			TuskPanic(errmsg, line, file, stacktrace)
		}

		var tmpundef TuskType = TuskUndef{}
		return &tmpundef
	},
	"len": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) != 1 {
			TuskPanic("Function len requires the signature (any)", line, file, stacktrace)
		}

		var length uint64

		switch (*args[0]).(type) {
		case TuskString:
			length = (*args[0]).(TuskString).Length
		case TuskArray:
			length = (*args[0]).(TuskArray).Length
		case TuskHash:
			length = (*args[0]).(TuskHash).Length
		}

		var tusklen TuskNumber
		tusklen.FromGoType(float64(length))

		var tusktype TuskType = tusklen
		return &tusktype
	},
	"clone": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) != 1 {
			TuskPanic("Function clone requires the signature (any)", line, file, stacktrace)
		}

		val := *args[0]

		switch val.(type) {

		case TuskArray:

			var arr = val.(TuskArray).Array
			var cloned = append([]*TuskType{}, arr...) //append it to nothing (to clone it)
			var tusktype TuskType = TuskArray{
				Array:  cloned,
				Length: val.(TuskArray).Length,
			}

			return &tusktype

		case TuskBool:

			//take inderect of the bool, and place it in a temporary variable
			var tmp = *val.(TuskBool).Boolean

			var returner TuskType = TuskBool{
				Boolean: &tmp, //take address of tmp and place it into `Boolean` field of returner
			}

			return &returner

		case TuskHash:
			var hash = val.(TuskHash).Hash

			//clone it into `cloned`
			var cloned = make(map[string]*TuskType)
			for k, v := range hash {
				cloned[k] = v
			}
			////////////////////////

			var tusktype TuskType = TuskHash{
				Hash:   cloned,
				Length: val.(TuskHash).Length,
			}

			return &tusktype

		case TuskNumber:
			var number = val.(TuskNumber)

			//copy the integer and decimal
			var integer = append([]int64{}, *number.Integer...)

			var decimal []int64

			if number.Decimal != nil {
				decimal = append([]int64{}, *number.Decimal...)
			}
			//////////////////////////////

			var newnum TuskType = TuskNumber{
				Integer: &integer,
				Decimal: &decimal,
			}
			return &newnum

		case TuskRune:

			//take inderect of the rune, and place it in a temporary variable
			var tmp = *val.(TuskRune).Rune

			var returner TuskType = TuskRune{
				Rune: &tmp, //take address of tmp and place it into `Rune` field of returner
			}

			return &returner

		case TuskString:

			var tmp = val.(TuskString).ToRuneList() //convert it to a go type
			var kastr TuskString
			kastr.FromRuneList(append(tmp, []rune{}...)) //clone tmp
			var returner TuskType = kastr
			return &returner

		default:
			TuskPanic("Cannot clone type \""+val.Type()+"\"", line, file, stacktrace)
		}

		var tmpundef TuskType = TuskUndef{}
		return &tmpundef
	},
	"panic": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) != 1 || (*args[0]).Type() != "string" {
			TuskPanic("Function panic requires the signature (string)", line, file, stacktrace)
		}

		var err = (*args[0]).(TuskString)

		TuskPanic(err.ToGoType(), line, file, stacktrace)

		var tmpundef TuskType = TuskUndef{}
		return &tmpundef
	},
	"exec": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		errmsg := "Function exec requires the signature (string) or (string, string)"

		if len(args) == 1 {
			if (*args[0]).Type() != "string" {
				TuskPanic(errmsg, line, file, stacktrace)
			}

			var cmd = (*args[0]).(TuskString).ToGoType()

			var _execdir = C.getCmdExe()
			var _arg = C.getCmdOp()
			var execdir = C.GoString(_execdir)
			var arg = C.GoString(_arg)

			command := exec.Command(execdir, arg, cmd)
			out, _ := command.CombinedOutput()

			var stringValue TuskString
			stringValue.FromGoType(string(out))
			var tusktype TuskType = stringValue
			return &tusktype
		} else if len(args) == 2 {
			if (*args[0]).Type() != "string" || (*args[1]).Type() != "string" {
				TuskPanic(errmsg, line, file, stacktrace)
			}

			var cmd = (*args[0]).(TuskString).ToGoType()
			var stdin = (*args[1]).(TuskString).ToGoType()

			var _execdir = C.getCmdExe()
			var _arg = C.getCmdOp()
			var execdir = C.GoString(_execdir)
			var arg = C.GoString(_arg)

			command := exec.Command(execdir, arg, cmd)
			command.Stdin = strings.NewReader(stdin)
			out, _ := command.CombinedOutput()

			var stringValue TuskString
			stringValue.FromGoType(string(out))
			var tusktype TuskType = stringValue
			return &tusktype
		} else {
			TuskPanic(errmsg, line, file, stacktrace)
		}

		var tmpundef TuskType = TuskUndef{}
		return &tmpundef
	},
	"load": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) != 1 || (*args[0]).Type() != "string" {
			TuskPanic("Function load requires the signature (string)", line, file, stacktrace)
		}

		libname := (*args[0]).(TuskString).ToGoType()
		lib, e := oatenc.OatDecode(libname)

		if e != nil {
			TuskPanic(e.Error(), line, file, stacktrace)
		}

		var tuskhash TuskHash

		for k, v := range lib {
			tuskhash.SetStr(k, *v)
		}

		var tusktype TuskType = tuskhash
		return &tusktype
	},
	"syscall": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if (*args[0]).Type() != "number" {
			TuskPanic("Sysno must be a numeric value", line, file, stacktrace)
		}

		var cargs = make([]unsafe.Pointer, int(C.MAX_SYS_ARGC))

		i := 0

		for ; i < len(args) && i < int(C.MAX_SYS_ARGC); i++ {
			var e error
			cargs[i], e = makectype(args[i])
			if e != nil {
				TuskPanic(e.Error(), line, file, stacktrace)
			}
		}

		for ; i < int(C.MAX_SYS_ARGC); i++ {
			cnum := C.int(0)
			cargs[i] = C.makeunsafeint(cnum)
		}

		sysno := int((*args[0]).(TuskNumber).ToGoType())
		syscall, exists := tusksys.SysTable[sysno]

		if !exists {
			TuskPanic(fmt.Sprintf("Could not find syscall: %d", sysno), line, file, stacktrace)
		}

		called := C.callsys(
			syscall,
			cargs[1],
			cargs[2],
			cargs[3],
			cargs[4],
			cargs[5],
			cargs[6],
			cargs[7],
			cargs[8],
			cargs[9],
			cargs[10],
			cargs[11],
			cargs[12],
			cargs[13],
			cargs[14],
			cargs[15],
			cargs[16],
			cargs[17],
			cargs[18],
			cargs[19],
			cargs[20],
		)

		//give back the values to the original tusk pointers
		for i := 1; i < len(cargs); i++ {
			if i >= len(args) { //if there are no more args to give back, break
				break
			}
			(*args[i]) = fromctype(cargs[i], args[i])
		}

		var tusknum TuskNumber
		tusknum.FromGoType(float64(called))
		var tusktype TuskType = tusknum
		return &tusktype
	},
}
