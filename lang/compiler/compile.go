package compiler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	. "github.com/tusklang/tusk/lang/types"
	"github.com/tusklang/tusk/native"
)

//CompileError represents a compile-time error in Tusk, and it implements the `error` interface
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

var tuskbasedir string

func inclCompile(filename string) ([]Action, error) {

	file, e := ioutil.ReadFile(filename)

	var strfile = string(file)
	var includes []Action

	if strings.HasPrefix(strfile, ";include:") { //if it starts with an include list
		var istrfile = strings.TrimSpace(strings.TrimPrefix(strfile, ";include:")) //trim the include list identifier
		var includepaths []string

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
				includepaths = append(includepaths, path.Join(path.Dir(filename), cip))
			case 'w': //working directory
				includepaths = append(includepaths, cip)
			case 's': //installation directory
				includepaths = append(includepaths, path.Join(tuskbasedir, cip))
			default:
				return nil, fmt.Errorf("Unrecognized include path '%c", dir)
			}

		}

		for _, v := range includepaths {
			acts, e := includer(v)

			if e != nil {
				return nil, makeCompilerErr(e.Error(), filename, 0)
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

	actions = append(includes, actions...) //append the includes

	return actions, e
}

func Compile(params CliParams) (map[string]*TuskType, error) {

	tuskbasedir = params.TuskDirname

	var e error

	actions, e := includer(params.Name)

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

	actions, access := globalAccess(actions)

	vars, e := getvars(actions)
	if e != nil {
		return nil, e
	}

	//make each var have only it's name
	var varnames = make(map[string]string)

	varnames["__dirname"] = "__dirname"
	varnames["argv"] = "argv"

	for k := range vars {

		if vars[k] == nil { //skip for declares
			varnames[k] = k
			continue
		}

		varnames[k] = k
	}

	//also account for the gofuncs
	for k := range native.Native {
		varnames[k] = k
	}

	for _, v := range vars {
		var tmp = make([]Action, 1) //create a temporary action slice to pass to changevarnames
		tmp[0].Type = (*v).Type()
		tmp[0].Value = *v

		_, e = changevarnames(tmp, varnames, access)
		if e != nil {
			return nil, e
		}

		*v = tmp[0].Value
	}

	return vars, nil
}
