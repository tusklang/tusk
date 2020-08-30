package compiler

import (
	"fmt"
	"os"

	"github.com/omm-lang/omm/lang/interpreter"

	. "github.com/omm-lang/omm/lang/types"
)

type _CompileErr struct {
	Msg   string
	FName string
	Line  uint64
}

func (e _CompileErr) Print() {
	fmt.Println("Error while compiling", e.FName, "on line", e.Line)
	fmt.Println(e.Msg)
	os.Exit(1)
}

type CompileErr interface {
	Print()
}

func makeCompilerErr(msg, fname string, line uint64) CompileErr {
	return _CompileErr{
		Msg:   msg,
		FName: fname,
		Line:  line,
	}
}

func inclCompile(file, filename string) ([]Action, CompileErr) {

	lex, e := lexer(file, filename)

	if e != nil {
		return []Action{}, e
	}

	groups, e := makeGroups(lex)

	if e != nil {
		return []Action{}, e
	}

	operations, e := makeOperations(groups)

	if e != nil {
		return []Action{}, e
	}

	actions, e := actionizer(operations)

	return actions, e
}

var ommbasedir string

func Compile(file, filename string, params CliParams) (map[string][]Action, CompileErr) {

	ommbasedir = params.OmmDirname

	var e CompileErr

	actions, e := inclCompile(file, filename)

	if e != nil {
		return nil, e
	}

	if e != nil {
		return nil, e
	}

	//a bunch of validations and initializers
	e = has_non_global_prototypes(actions, true)
	if e != nil {
		return nil, e
	}
	put_proto_types(actions)
	e = validate_types(actions)
	if e != nil {
		return nil, e
	}
	/////////////////////////////////////////

	vars, e := getvars(actions)
	if e != nil {
		return nil, e
	}

	//make each var have only it's name
	var varnames = make(map[string]string)

	varnames["$__dirname"] = "$__dirname"
	varnames["$argv"] = "$argv"

	for k := range vars {

		if len(vars[k]) == 0 { //skip for declares
			varnames[k] = k
			continue
		}

		//ensure that the globals do not have any compound types (such as operations)
		if vars[k][0].Value == nil {
			return nil, makeCompilerErr("Cannot have compound types at the global scope", vars[k][0].File, vars[k][0].Line)
		}

		varnames[k] = k
	}

	//also account for the gofuncs
	for k := range interpreter.Native {
		varnames[k] = k
	}

	for k := range vars {
		_, e = changevarnames(vars[k], varnames)
		if e != nil {
			return nil, e
		}

		putFunctionVarRefs(vars[k])
	}

	return vars, nil
}
