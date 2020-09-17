package interpreter

//all of the gofuncs
//functions written in go that are used by tusk

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	. "tusk/lang/types"
)

//#include "exec.h"
import "C"

//Native stores all of the native values. You can make your own by just putting it into this map
var Native = make(map[string]*TuskType)

func kaprint(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) {
	for k, v := range args {
		fmt.Print((*v).Format())
		if k+1 != len(args) {
			fmt.Print(" ")
		}
	}
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
}

//these are the native functions that are relatively simple to implement
var native = map[string]func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType{
	"log": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {
		kaprint(args, stacktrace, line, file, instance)
		fmt.Println() //print a newline at the end
		var tmpundef TuskType = undef
		return &tmpundef
	},
	"print": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {
		kaprint(args, stacktrace, line, file, instance)
		var tmpundef TuskType = undef
		return &tmpundef
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

		var tmpundef TuskType = undef
		return &tmpundef
	},
	"wait": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) == 1 {

			if (*args[0]).Type() != "number" {
				TuskPanic("Function wait requires a number as the argument", line, file, stacktrace)
			}

			var amt = (*args[0]).(TuskNumber)

			var n4294967295 = zero
			n4294967295.Integer = &[]int64{5, 9, 2, 7, 6, 9, 4, 9, 2, 4}

			//if amt is less than 2 ^ 32 - 1, just convert to a go int
			if isLess(amt, n4294967295) {
				gonum := amt.ToGoType()

				time.Sleep(time.Duration(gonum) * time.Millisecond)
			} else {
				//this is how this works
				/*
				   the loop starts at 0 with an increment of 2 ^ 32 - 1
				   in each iteration, it will wait for 2 ^ 32 - 1 milliseconds
				*/

				for i := zero; isLess(i, amt); i = (*number__plus__number(i, n4294967295, &Instance{}, stacktrace, line, file)).(TuskNumber) {
					time.Sleep(4294967295 * time.Millisecond)
				}
			}

		} else {
			TuskPanic("Function wait requires a parameter count of 1", line, file, stacktrace)
		}

		var tmpundef TuskType = undef
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

		var tmpundef TuskType = undef
		return &tmpundef
	},
	"len": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) != 1 {
			TuskPanic("Function len requires the signature (any)", line, file, stacktrace)
		}

		switch (*args[0]).(type) {
		case TuskString:
			var length = (*args[0]).(TuskString).Length

			var kanumber_length TuskNumber
			kanumber_length.FromString(strconv.FormatUint(length, 10))
			var tusktype TuskType = kanumber_length
			return &tusktype

		case TuskArray:
			var length = (*args[0]).(TuskArray).Length

			var kanumber_length TuskNumber
			kanumber_length.FromString(strconv.FormatUint(length, 10))
			var tusktype TuskType = kanumber_length
			return &tusktype

		case TuskHash:
			var length = (*args[0]).(TuskHash).Length

			var kanumber_length TuskNumber
			kanumber_length.FromString(strconv.FormatUint(length, 10))
			var tusktype TuskType = kanumber_length
			return &tusktype

		}

		var tmpzero TuskType = zero
		return &tmpzero
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

		var tmpundef TuskType = undef
		return &tmpundef
	},
	"panic": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) != 1 || (*args[0]).Type() != "string" {
			TuskPanic("Function panic requires the signature (string)", line, file, stacktrace)
		}

		var err = (*args[0]).(TuskString)

		TuskPanic(err.ToGoType(), line, file, stacktrace)

		var tmpundef TuskType = undef
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

		var tmpundef TuskType = undef
		return &tmpundef
	},
	"chdir": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) != 1 || (*args[0]).Type() != "string" {
			TuskPanic("Function chdir requires the signature (string)", line, file, stacktrace)
		}

		var dir = (*args[0]).(TuskString).ToGoType()

		os.Chdir(dir)

		instance.Allocate("$__dirname", args[0])

		var tmpundef TuskType = undef
		return &tmpundef
	},
	"getos": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {
		var os = runtime.GOOS
		var kastr TuskString
		kastr.FromGoType(os)
		var tusktype TuskType = kastr
		return &tusktype
	},
	"sprint": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {

		if len(args) != 1 {
			TuskPanic("Function sprint requires the signature (any)", line, file, stacktrace)
		}

		var sprinted = (*args[0]).Format()
		var kastr TuskString
		kastr.FromGoType(sprinted)
		var tusktype TuskType = kastr
		return &tusktype
	},
	"syscall": func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType {
		return nil
	},
}
