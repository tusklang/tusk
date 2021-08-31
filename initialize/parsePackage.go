package initialize

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

func parsePackage(dir string, pkg *Package, prog *Program) error {

	fsinfo, e := ioutil.ReadDir(dir)

	if e != nil {
		return e
	}

	for _, v := range fsinfo {

		//joined path (of the parent directories and current one)
		jpth := path.Join(dir, v.Name())

		if v.IsDir() {
			//a new package

			var spkg Package

			//because the name is an array (see `Package.go`) we want to get the package names of all the parents
			spkg.Name = append(spkg.Name, strings.Split(jpth, "/")...)

			e = parsePackage(jpth, &spkg, prog)

			if e != nil {
				return e
			}

			continue
		}

		//only append a new class if it's a tusk file
		if filepath.Ext(v.Name()) != ".tusk" {
			continue
		}

		//a new class in the package
		pf, e := parseFile(jpth)

		if e != nil {
			return e
		}

		pkg.Files = append(pkg.Files, pf)
	}

	prog.Packages = append(prog.Packages, pkg)
	return nil
}
