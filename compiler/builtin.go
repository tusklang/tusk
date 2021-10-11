package compiler

import (
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
)

func strarrIncludes(a []string, i string) bool {
	for _, v := range a {
		if v == i {
			return true
		}
	}
	return false
}

func addFuncIfndef(module *ir.Module, fn *ir.Func) {
	for _, v := range module.Funcs {
		if v.Name() == fn.Name() {
			return
		}
	}
	module.Funcs = append(module.Funcs, fn)
}

func ParseBuiltin(compiler *ast.Compiler, fs embed.FS) {
	bfs, e := fs.ReadDir("builtin")

	var builtins = make(map[string]builtin)

	if e != nil {
		log.Fatal(e)
	}

	//first we loop through all the config files
	for _, v := range bfs {

		if filepath.Ext(v.Name()) == ".config" {
			//detected a config

			//read the actual contents of the file
			cf, _ := fs.ReadFile("builtin/" + v.Name())

			//parse the contents
			cfs := string(cf)

			parts := strings.Split(cfs, "-")

			typ := strings.TrimSpace(parts[0])

			switch typ {
			case "class":

				if len(parts) != 4 {
					log.Fatal("Invalid builtin config file: " + v.Name())
				}

				var bclass = newBuiltinClass()
				cname := strings.TrimSpace(parts[2]) /* the third entry is the name */
				builtins[cname] = bclass
				bclass.name = cname

				//seperate each field in the config file
				var fieldsStr = strings.Split(strings.TrimSpace(parts[1]), "\n")

				for _, v := range fieldsStr {
					inf := strings.Split(v, " ")

					//it's in the format of <index> <field name> <accessibility>
					var (
						idx, _ = strconv.Atoi(inf[0])
						name   = inf[1]
						access = inf[2]
					)

					//only add it to the export if its public
					if access == "public" {
						bclass.exported[int64(idx)] = name
					}

				}

				//the fourth entry is the list of names to not mangle
				nomangle := strings.Split(strings.TrimSpace(parts[3]), " ")
				bclass.nomangle = nomangle

			case "func":

				if len(parts) != 3 {
					log.Fatal("Invalid builtin config file: " + v.Name())
				}

				bfunc := newBuiltinFunc()
				fname := strings.TrimSpace(parts[1]) /* the second entry is the name */
				builtins[fname] = bfunc

				bfunc.name = fname

				//the third entry is the list of names to not mangle
				nomangle := strings.Split(strings.TrimSpace(parts[2]), " ")
				bfunc.nomangle = nomangle

			}

		}
	}

	//now we loop through all the llvm ir implementations
	for _, v := range bfs {

		if filepath.Ext(v.Name()) == ".ll" {
			iname := strings.Split(v.Name(), ".")[0] //get the name of the builtin

			if obj, exists := builtins[iname]; !exists {
				log.Fatal("Missing config file for: " + iname)
			} else {

				//get the actual llvm ir
				ird, _ := fs.ReadFile("builtin/" + v.Name())

				//parse the ir
				parsed, e := asm.ParseBytes("", ird)

				if e != nil {
					log.Fatal(e)
				}

				switch o := obj.(type) {
				case *builtinClass:

					//look for the typedef that we defined in the config
					for _, vv := range parsed.TypeDefs {

						//if the current typedef is the one we defined in the config
						if vv.Name() == fmt.Sprintf("struct.%s", o.FetchName()) {
							vv.SetName("builtin." + o.FetchName())
							o.structT = vv.(*types.StructType)
							break
						}

					}

					//now look for the constructor
					for _, vv := range parsed.Funcs {

						//if the name of the functions is included in the nomangle list
						if strarrIncludes(o.nomangle, vv.Name()) {
							addFuncIfndef(compiler.Module, vv)
							continue
						}

						if vv.Name() == "construct" {
							vv.SetName("builtin." + o.FetchName() + ".construct")
							o.constructor = vv
						} else {
							vv.SetName("builtin." + o.FetchName() + ".helpers." + vv.Name())
							o.helpers = append(o.helpers, vv)
						}
					}

				case *builtinFunc:

					for _, v := range parsed.Funcs {

						if strarrIncludes(o.nomangle, v.Name()) {
							addFuncIfndef(compiler.Module, v)
							break
						}

						if v.Name() == o.FetchName() {
							v.SetName("builtin." + o.FetchName())
							o.exported = v
						} else {
							v.SetName("builtin." + o.FetchName() + ".helpers." + v.Name())
							o.helpers = append(o.helpers, v)
						}
					}

				}
			}
		}

	}

	//now we add the builtins we found to the build
	for _, v := range builtins {
		switch o := v.(type) {
		case *builtinClass:
			var tuskclass = data.NewClass(o.FetchName(), o.structT, nil)

			tuskclass.Construct = data.NewFunc(o.constructor, data.NewPointer(tuskclass))

			for k, vv := range o.exported {
				tuskclass.AddInstanceItem(vv, data.LLTypToTusk(o.structT.Fields[k]), k)
			}

			prevars[o.FetchName()] = tuskclass
			compiler.Module.Funcs = append(compiler.Module.Funcs, o.constructor)
			compiler.Module.TypeDefs = append(compiler.Module.TypeDefs, o.structT)

			//add the helper functions
			compiler.Module.Funcs = append(compiler.Module.Funcs, o.helpers...)
		case *builtinFunc:
			tuskfunc := data.NewFunc(o.exported, data.LLTypToTusk(o.exported.Sig.RetType))
			prevars[o.FetchName()] = tuskfunc

			compiler.Module.Funcs = append(compiler.Module.Funcs, o.exported)
			compiler.Module.Funcs = append(compiler.Module.Funcs, o.helpers...)
		}
	}

}
