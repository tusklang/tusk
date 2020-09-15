package interpreter

//all of the gofuncs
//functions written in go that are used by ka

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	. "ka/lang/types"
)

//#include "exec.h"
import "C"

//Native stores all of the native values. You can make your own by just putting it into this map
var Native = make(map[string]*KaType)

func kaprint(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) {
	for k, v := range args {
		fmt.Print((*v).Format())
		if k+1 != len(args) {
			fmt.Print(" ")
		}
	}
}

//KaPanic panics in an Ka instance
func KaPanic(err string, line uint64, file string, stacktrace []string) {
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
var native = map[string]func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType{
	"log": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {
		kaprint(args, stacktrace, line, file, instance)
		fmt.Println() //print a newline at the end
		var tmpundef KaType = undef
		return &tmpundef
	},
	"print": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {
		kaprint(args, stacktrace, line, file, instance)
		var tmpundef KaType = undef
		return &tmpundef
	},
	"await": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		if len(args) != 1 || (*args[0]).Type() != "thread" {
			KaPanic("Function await requires (thread)", line, file, stacktrace)
		}

		interpreted := args[0]
		var awaited *KaType

		switch (*interpreted).(type) {
		case KaThread:

			//put the new value back into the given interpreted pointer
			thread := (*interpreted).(KaThread)
			awaited = thread.Join()
			*interpreted = thread
			///////////////////////////////////////////////////////////

		default:
			awaited = interpreted
		}

		return awaited
	},
	"input": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		scanner := bufio.NewScanner(os.Stdin)

		if len(args) == 0 {
			//if it has 0 or 1 args, there is no error
		} else if len(args) == 1 {

			switch (*args[0]).(type) {
			case KaString:
				str := (*args[0]).(KaString).ToGoType()
				fmt.Print(str)
			default:
				KaPanic("Expected a string as the argument to input[]", line, file, stacktrace)
			}

		} else {
			KaPanic("Function input requires a parameter count of 0 or 1", line, file, stacktrace)
		}

		//get user input and convert it to KaType
		scanner.Scan()
		input := scanner.Text()
		var inputKaType KaString
		inputKaType.FromGoType(input)
		var inputType KaType = inputKaType

		return &inputType
	},
	"typeof": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		if len(args) != 1 {
			KaPanic("Function typeof requires a parameter count of 1", line, file, stacktrace)
		}

		typeof := (*args[0]).TypeOf()

		var str KaString
		str.FromGoType(typeof)

		//convert to KaType interface
		var katype KaType = str
		return &katype
	},
	"append": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		if len(args) != 2 {
			KaPanic("Function append requires a parameter count of 2", line, file, stacktrace)
		}

		if (*args[0]).Type() != "array" {
			KaPanic("Function append requires (array, any)", line, file, stacktrace)
		}

		a := (*args[0]).(KaArray)
		a.PushBack(*args[1])
		var katype KaType = a
		return &katype
	},
	"prepend": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		if len(args) != 2 {
			KaPanic("Function prepend requires a parameter count of 2", line, file, stacktrace)
		}

		if (*args[0]).Type() != "array" {
			KaPanic("Function prepend requires the first argument to be an array", line, file, stacktrace)
		}

		a := (*args[0]).(KaArray)
		a.PushFront(*args[1])
		var katype KaType = a
		return &katype
	},
	"exit": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		if len(args) == 1 {

			switch (*args[0]).(type) {
			case KaNumber:

				var gonum = (*args[0]).(KaNumber).ToGoType()
				os.Exit(int(gonum))

			case KaBool:

				if (*args[0]).(KaBool).ToGoType() == true {
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
			KaPanic("Function exit requires a parameter count of 1 or 0", line, file, stacktrace)
		}

		var tmpundef KaType = undef
		return &tmpundef
	},
	"wait": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		if len(args) == 1 {

			if (*args[0]).Type() != "number" {
				KaPanic("Function wait requires a number as the argument", line, file, stacktrace)
			}

			var amt = (*args[0]).(KaNumber)

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

				for i := zero; isLess(i, amt); i = (*number__plus__number(i, n4294967295, &Instance{}, stacktrace, line, file)).(KaNumber) {
					time.Sleep(4294967295 * time.Millisecond)
				}
			}

		} else {
			KaPanic("Function wait requires a parameter count of 1", line, file, stacktrace)
		}

		var tmpundef KaType = undef
		return &tmpundef
	},
	"make": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		errmsg := "Function make requires the signature (prototype)"

		if len(args) == 1 {

			switch (*args[0]).(type) {
			case KaProto:
				var katype KaType = (*args[0]).(KaProto).New(*instance)
				return &katype
			default:
				KaPanic(errmsg, line, file, stacktrace)
			}

		} else {
			KaPanic(errmsg, line, file, stacktrace)
		}

		var tmpundef KaType = undef
		return &tmpundef
	},
	"len": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		if len(args) != 1 {
			KaPanic("Function len requires the signature (any)", line, file, stacktrace)
		}

		switch (*args[0]).(type) {
		case KaString:
			var length = (*args[0]).(KaString).Length

			var kanumber_length KaNumber
			kanumber_length.FromString(strconv.FormatUint(length, 10))
			var katype KaType = kanumber_length
			return &katype

		case KaArray:
			var length = (*args[0]).(KaArray).Length

			var kanumber_length KaNumber
			kanumber_length.FromString(strconv.FormatUint(length, 10))
			var katype KaType = kanumber_length
			return &katype

		case KaHash:
			var length = (*args[0]).(KaHash).Length

			var kanumber_length KaNumber
			kanumber_length.FromString(strconv.FormatUint(length, 10))
			var katype KaType = kanumber_length
			return &katype

		}

		var tmpzero KaType = zero
		return &tmpzero
	},
	"clone": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		if len(args) != 1 {
			KaPanic("Function clone requires the signature (any)", line, file, stacktrace)
		}

		val := *args[0]

		switch val.(type) {

		case KaArray:

			var arr = val.(KaArray).Array
			var cloned = append([]*KaType{}, arr...) //append it to nothing (to clone it)
			var katype KaType = KaArray{
				Array:  cloned,
				Length: val.(KaArray).Length,
			}

			return &katype

		case KaBool:

			//take inderect of the bool, and place it in a temporary variable
			var tmp = *val.(KaBool).Boolean

			var returner KaType = KaBool{
				Boolean: &tmp, //take address of tmp and place it into `Boolean` field of returner
			}

			return &returner

		case KaHash:
			var hash = val.(KaHash).Hash

			//clone it into `cloned`
			var cloned = make(map[string]*KaType)
			for k, v := range hash {
				cloned[k] = v
			}
			////////////////////////

			var katype KaType = KaHash{
				Hash:   cloned,
				Length: val.(KaHash).Length,
			}

			return &katype

		case KaNumber:
			var number = val.(KaNumber)

			//copy the integer and decimal
			var integer = append([]int64{}, *number.Integer...)

			var decimal []int64

			if number.Decimal != nil {
				decimal = append([]int64{}, *number.Decimal...)
			}
			//////////////////////////////

			var newnum KaType = KaNumber{
				Integer: &integer,
				Decimal: &decimal,
			}
			return &newnum

		case KaRune:

			//take inderect of the rune, and place it in a temporary variable
			var tmp = *val.(KaRune).Rune

			var returner KaType = KaRune{
				Rune: &tmp, //take address of tmp and place it into `Rune` field of returner
			}

			return &returner

		case KaString:

			var tmp = val.(KaString).ToRuneList() //convert it to a go type
			var kastr KaString
			kastr.FromRuneList(append(tmp, []rune{}...)) //clone tmp
			var returner KaType = kastr
			return &returner

		default:
			KaPanic("Cannot clone type \""+val.Type()+"\"", line, file, stacktrace)
		}

		var tmpundef KaType = undef
		return &tmpundef
	},
	"panic": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		if len(args) != 1 || (*args[0]).Type() != "string" {
			KaPanic("Function panic requires the signature (string)", line, file, stacktrace)
		}

		var err = (*args[0]).(KaString)

		KaPanic(err.ToGoType(), line, file, stacktrace)

		var tmpundef KaType = undef
		return &tmpundef
	},
	"exec": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		errmsg := "Function exec requires the signature (string) or (string, string)"

		if len(args) == 1 {
			if (*args[0]).Type() != "string" {
				KaPanic(errmsg, line, file, stacktrace)
			}

			var cmd = (*args[0]).(KaString).ToGoType()

			var _execdir = C.getCmdExe()
			var _arg = C.getCmdOp()
			var execdir = C.GoString(_execdir)
			var arg = C.GoString(_arg)

			command := exec.Command(execdir, arg, cmd)
			out, _ := command.CombinedOutput()

			var stringValue KaString
			stringValue.FromGoType(string(out))
			var katype KaType = stringValue
			return &katype
		} else if len(args) == 2 {
			if (*args[0]).Type() != "string" || (*args[1]).Type() != "string" {
				KaPanic(errmsg, line, file, stacktrace)
			}

			var cmd = (*args[0]).(KaString).ToGoType()
			var stdin = (*args[1]).(KaString).ToGoType()

			var _execdir = C.getCmdExe()
			var _arg = C.getCmdOp()
			var execdir = C.GoString(_execdir)
			var arg = C.GoString(_arg)

			command := exec.Command(execdir, arg, cmd)
			command.Stdin = strings.NewReader(stdin)
			out, _ := command.CombinedOutput()

			var stringValue KaString
			stringValue.FromGoType(string(out))
			var katype KaType = stringValue
			return &katype
		} else {
			KaPanic(errmsg, line, file, stacktrace)
		}

		var tmpundef KaType = undef
		return &tmpundef
	},
	"chdir": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		if len(args) != 1 || (*args[0]).Type() != "string" {
			KaPanic("Function chdir requires the signature (string)", line, file, stacktrace)
		}

		var dir = (*args[0]).(KaString).ToGoType()

		os.Chdir(dir)

		instance.Allocate("$__dirname", args[0])

		var tmpundef KaType = undef
		return &tmpundef
	},
	"getos": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {
		var os = runtime.GOOS
		var kastr KaString
		kastr.FromGoType(os)
		var katype KaType = kastr
		return &katype
	},
	"sprint": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		if len(args) != 1 {
			KaPanic("Function sprint requires the signature (any)", line, file, stacktrace)
		}

		var sprinted = (*args[0]).Format()
		var kastr KaString
		kastr.FromGoType(sprinted)
		var katype KaType = kastr
		return &katype
	},
	"syscall": func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType {

		const maxSyscallArgv = 18

		var fullargv [19]uintptr

		//convert to go types
		for k, v := range args {

			cur := &fullargv[k]

			switch (*v).(type) {
			case KaNumber:
				*cur = uintptr((*v).(KaNumber).ToGoType())
			case KaString:
				gostr := (*v).(KaString).ToGoType()
				*cur = uintptr(unsafe.Pointer(&gostr))
			default:
				KaPanic("Cannot use type "+(*v).Type()+" in a system call", line, file, stacktrace)
			}
		}

		c := makeSyscall(len(args)-1, fullargv)

		fmt.Println(c)

		return nil
	},
}
