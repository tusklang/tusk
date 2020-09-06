package compiler

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/omm-lang/omm/lang/interpreter"

	. "github.com/omm-lang/omm/lang/types"
)

type CompileError struct {
	Msg   string
	FName string
	Line  uint64
}

func (e CompileError) Error() string {
	f := ""
	f += fmt.Sprintln("Error while compiling", e.FName, "on line", e.Line)
	f += fmt.Sprint(e.Msg)
	return f
}

func makeCompilerErr(msg, fname string, line uint64) error {
	return CompileError{
		Msg:   msg,
		FName: fname,
		Line:  line,
	}
}

func inclCompile(filename string) ([]Action, error) {

	file, e := ioutil.ReadFile(filename)

	if e != nil {
		return nil, errors.New("Could not open file: " + filename)
	}

	lex, e := lexer(string(file), filename)

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

func Compile(params CliParams) (map[string]*OmmType, error) {

	ommbasedir = params.OmmDirname

	var e error

	actions, e := inclCompile(params.Name)

	if e != nil {
		return nil, e
	}

	getnativetypes()

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

		if vars[k] == nil { //skip for declares
			varnames[k] = k
			continue
		}

		varnames[k] = k
	}

	//also account for the gofuncs
	for k := range interpreter.Native {
		varnames[k] = k
	}

	for _, v := range vars {

		var tmp = make([]Action, 1) //create a temporary action slice to pass to changevarnames
		tmp[0].Type = (*v).Type()
		tmp[0].Value = *v

		_, e = changevarnames(tmp, varnames)
		if e != nil {
			return nil, e
		}

		putFunctionVarRefs(tmp)
		*v = tmp[0].Value
	}

	return vars, nil
}
