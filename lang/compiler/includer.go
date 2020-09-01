package compiler

import (
	"io/ioutil"
	"path"
	"strings"

	oatenc "github.com/omm-lang/oat/format/encoding"

	. "github.com/omm-lang/omm/lang/types"
)

func includeSingle(filename string, line uint64, dir string, curdir string, includewhere int) ([]Action, error) {

	switch includewhere {
	case 1:
		filename = path.Join(ommbasedir, filename)
	case 2:
		filename = path.Join(curdir, filename)
	}

	if strings.HasSuffix(filename, ".oat") {
		decoded, e := oatenc.OatDecode(filename, 1)

		if e != nil {
			return nil, CompileError{
				Msg:   e.Error(),
				FName: filename,
				Line:  line,
			}
		}

		return decoded["$main"], nil
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

	if strings.HasSuffix(filename, "*") {

		switch includewhere {
		case 1:
			filename = path.Join(path.Join(ommbasedir), filename)
		case 2:
			filename = path.Join(path.Join(curdir), filename)
		}

		files, e := ioutil.ReadDir(strings.TrimSuffix(filename, "*"))

		if e != nil {
			return [][]Action{}, makeCompilerErr("Could not find directory: "+filename, dir, line)
		}

		var actions [][]Action

		for _, v := range files {

			if !strings.HasSuffix(v.Name(), ".omm") && !strings.HasSuffix(v.Name(), ".oat") { //if it is not an omm or an oat file, skip it
				continue
			}

			if v.IsDir() {
				inc, e := includer(path.Join(strings.TrimSuffix(filename, "*"), v.Name()+"/*"), line, dir, curdir, 0)

				if e != nil {
					return [][]Action{}, e
				}

				actions = append(actions, inc...)
			} else {
				inc, e := includeSingle(path.Join(strings.TrimSuffix(filename, "*"), v.Name()), line, dir, curdir, 0)

				if e != nil {
					return [][]Action{}, e
				}

				actions = append(actions, inc)
			}
		}

		return actions, nil
	}

	inc, e := includeSingle(filename, line, dir, curdir, includewhere)

	if e != nil {
		return [][]Action{}, e
	}

	return [][]Action{inc}, nil
}
