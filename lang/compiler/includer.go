package compiler

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	oatenc "github.com/omm-lang/oat/format/encoding"

	. "github.com/omm-lang/omm/lang/types"
)

func includeSingle(filename string, line uint64, dir string, curdir string, includewhere int) ([]Action, error) {

	if strings.HasSuffix(filename, ".oat") {
		decoded, e := oatenc.OatDecode(filename)

		if e != nil {
			return nil, CompileError{
				Msg:   e.Error(),
				FName: filename,
				Line:  line,
			}
		}

		var actions []Action

		for k, v := range decoded {
			actions = append(actions, Action{
				Type: "var",
				Name: k,
				ExpAct: []Action{
					Action{
						Type:  (*v).Type(),
						Value: *v,
					},
				},
			})
		}

		return actions, nil
	}

	if strings.HasSuffix(filename, ".omm") {
		filename = strings.TrimSuffix(filename, ".omm")
	}

	filename += ".omm"

	for _, v := range included {
		if v == filename {
			return []Action{}, nil
		}
	}

	compiled, e := inclCompile(filename)

	if e != nil {
		return []Action{}, e
	}

	return compiled, nil
}

func includer(filename string, line uint64, dir string, curdir string, includewhere int) ([][]Action, error) {

	switch includewhere {
	case 1:
		filename = path.Join(path.Join(ommbasedir), filename)
	case 2:
		filename = path.Join(path.Join(curdir), filename)
	}

	stat, e := os.Stat(filename)

	if e != nil {
		return nil, makeCompilerErr("Could not open "+filename, dir, line)
	}

	if stat.IsDir() {

		files, _ := ioutil.ReadDir(filename)

		var actions [][]Action

		for _, v := range files {
			acts, e := includer(path.Join(filename, v.Name()), line, dir, curdir, 0)

			if e != nil {
				return nil, e
			}

			actions = append(actions, acts...)
		}

		return actions, nil
	}

	inc, e := includeSingle(filename, line, dir, curdir, includewhere)

	if e != nil {
		return [][]Action{}, e
	}

	return [][]Action{inc}, nil
}
