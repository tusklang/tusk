package compiler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

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

	/*
		includewhere::
			0: relative to working directory
			1: relative to omm installation dir (stdlib)
			2: relative to current file
	*/

	var strfile = string(file)
	var includes [][]Action

	if strings.HasPrefix(strfile, ";include:") { //if it starts with an include list
		var istrfile = strings.TrimSpace(strings.TrimPrefix(strfile, ";include:")) //trim the include list identifier
		var includepaths [][2]string

		if len(istrfile) < 2 {
			return nil, errors.New("Must use an option for an include directory")
		}

		for istrfile[0] == ';' {
			istrfile = istrfile[1:]

			var dir = istrfile[0]

			istrfile = strings.TrimSpace(istrfile[1:])
			nline := strings.Index(istrfile, "\n")
			cip := strings.TrimSpace(istrfile[:nline])     //get the directory of the current include path
			istrfile = strings.TrimSpace(istrfile[nline:]) //remove the cip

			switch dir {
			case 'r': //relative directory
				includepaths = append(includepaths, [2]string{string(rune(0)), cip})
			case 's': //installation directory
				includepaths = append(includepaths, [2]string{string(rune(1)), cip})
			case 'w': //working directory
				includepaths = append(includepaths, [2]string{string(rune(2)), cip})
			default:
				return nil, fmt.Errorf("Unrecognized include path '%c", dir)
			}

		}

		for _, v := range includepaths {
			acts, e := includer(v[1], filename, int(v[0][0]))

			if e != nil {
				return nil, makeCompilerErr(e.Error(), v[1], 0)
			}

			includes = append(includes, acts...)
		}

	}

	lex, e := lexer(strfile, filename)

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

	for _, v := range includes {
		actions = append(v, actions...)
	}

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
