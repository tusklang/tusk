package native

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"omm/lang/types"
)

type OmmFile struct {
	gofile *os.File
	Write  OmmGoFunc
	Read   OmmGoFunc
	Delete OmmGoFunc
	Close  OmmGoFunc
	Exists OmmGoFunc
}

var omm_newfile = func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

	if len(args) != 1 || (*args[0]).Type() != "string" {
		OmmPanic("Function files.open requires a parameter count of 1 with type string", line, file, stacktrace)
	}

	var name = (*args[0]).(types.OmmString).ToGoType()
	f := newfile(name)

	var ommtype types.OmmType = f
	return &ommtype
}

var omm_createfile = func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

	if len(args) != 1 || (*args[0]).Type() != "string" {
		OmmPanic("Function files.create requires a parameter count of 1 with type string", line, file, stacktrace)
	}

	var name = (*args[0]).(types.OmmString).ToGoType()
	f := createfile(name)

	var ommtype types.OmmType = f
	return &ommtype
}

func newfile(filename string) OmmFile {
	f, _ := os.Open(filename)
	var file OmmFile
	file.FromGoType(f)
	return file
}

func createfile(filename string) OmmFile {
	f, _ := os.Create(filename)
	var file OmmFile
	file.FromGoType(f)
	return file
}

func (f *OmmFile) FromGoType(osfile *os.File) {
	f.gofile = osfile

	f.Write = OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

			if len(args) != 1 || (*args[0]).Type() != "string" {
				OmmPanic("Function file::Write requires a parameter count of 1 with type string", line, file, stacktrace)
			}

			var data = (*args[0]).(types.OmmString).ToRuneList()

			for _, v := range data {
				fmt.Fprintf(osfile, "%c", v)
			}

			var undef types.OmmType = types.OmmUndef{}
			return &undef
		},
	}

	f.Read = OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {
			r := bufio.NewReader(osfile)
			var runelist []rune

			for { //read it as a rune list
				if c, _, err := r.ReadRune(); err != nil {
					if err == io.EOF {
						break
					} else {
						OmmPanic(err.Error(), line, file, stacktrace)
					}
				} else {
					runelist = append(runelist, c)
				}
			}

			var ommstr types.OmmString
			ommstr.FromRuneList(runelist)
			var ommtype types.OmmType = ommstr
			return &ommtype
		},
	}

	f.Delete = OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {
			os.Remove(osfile.Name())
			osfile.Close()
			var ommtype types.OmmType = types.OmmUndef{}
			return &ommtype
		},
	}

	f.Close = OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {
			osfile.Close()
			var ommtype types.OmmType = types.OmmUndef{}
			return &ommtype
		},
	}

	f.Exists = OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

			gotrue := true
			gofalse := false

			if _, e := os.Stat(osfile.Name()); os.IsNotExist(e) {
				var truev types.OmmType = types.OmmBool{
					Boolean: &gotrue,
				}
				return &truev
			}

			var falsev types.OmmType = types.OmmBool{
				Boolean: &gofalse,
			}
			return &falsev
		},
	}
}

func (f OmmFile) Format() string {
	return "{ file struct }"
}

func (f OmmFile) Type() string {
	return "file"
}

func (f OmmFile) TypeOf() string {
	return f.Type()
}

func (f OmmFile) Deallocate() {}

//Range ranges over a file
func (f OmmFile) Range(fn func(val1, val2 *types.OmmType) types.Returner) *types.Returner {
	return nil
}
