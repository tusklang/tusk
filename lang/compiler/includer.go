package compiler

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	oatenc "oat/format/encoding"

	. "omm/lang/types"
)

var included = []string{} //list of the included files from omm

func includeSingle(filename string) ([]Action, error) {

	for _, v := range included { //ensure includes are not duplicated (header guards)
		if v == filename {
			return nil, nil
		}
	}

	if strings.HasSuffix(filename, ".oat") {
		decoded, e := oatenc.OatDecode(filename)

		if e != nil {
			return nil, e
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

	included = append(included, filename)

	compiled, e := inclCompile(filename)

	if e != nil {
		return []Action{}, e
	}

	return compiled, nil
}

func includer(filename string) ([]Action, error) {

	stat, e := os.Stat(filename)

	if e != nil {
		return nil, errors.New("Could not open " + filename)
	}

	if stat.IsDir() {

		files, _ := ioutil.ReadDir(filename)

		var actions []Action

		for _, v := range files {
			acts, e := includer(path.Join(filename, v.Name()))

			if e != nil {
				return nil, e
			}

			actions = append(actions, acts...)
		}

		return actions, nil
	}

	inc, e := includeSingle(filename)

	if e != nil {
		return nil, e
	}

	return inc, nil
}
