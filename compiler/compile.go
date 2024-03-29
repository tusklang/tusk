package compiler

import (
	"os"
	"runtime"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/parser"
	"github.com/tusklang/tusk/varprocessor"
)

//list of all the variables that are added by default
//has types to begin with, but it can store anything
//types are variables in tusk's parser so we need to add the default ones in like so
var prevars = map[string]data.Value{}

func Compile(prog *parser.Program, outfile string) {

	var compiler ast.Compiler
	m := ir.NewModule() //create a new llvm module

	//set the module stuff in the compiler
	compiler.Module = m
	compiler.OS = runtime.GOOS
	compiler.ARCH = runtime.GOARCH

	compiler.VarMap = make(map[string]data.Value)
	compiler.LinkedFunctions = make(map[string]*ir.Func)

	for k, v := range ast.InttypeV {
		ast.Numtypes[k] = v
	}

	for k, v := range ast.UinttypeV {
		ast.Numtypes[k] = v
	}

	for k, v := range ast.FloattypeV {
		ast.Numtypes[k] = v
	}

	for k, v := range ast.Numtypes {
		prevars[k] = v.(data.Value)
	}

	initDefaultOps(&compiler)
	initDefaultCasts(&compiler)

	var processor = varprocessor.NewProcessor(&compiler)

	//initialize the operations
	var initfunc = m.NewFunc(".tusk.initfunc", types.Void) //initialize func ran before main
	compiler.InitFunc = data.NewFunc(initfunc, data.NewPrimitive(types.Void))
	compiler.InitFunc.ActiveBlock = initfunc.NewBlock("")

	cclasses := parseProjStructure(&compiler, prog)

	//add all the prevars as variables to the compiler
	for k, v := range prevars {
		processor.AddPreDecl(k)
		compiler.AddVar(k, v)
	}

	//process all the variables

	for ic, c := range cclasses {

		//create a new processor specific to this class
		//imagine a project like this:
		/*
			package1/
				file1.tusk
				file2.tusk
			file3.tusk
		*/
		//now say the file1 class references the file2 class
		//and file3 also references file2
		//file1 could just use `file2.prop`
		//but file3 would have to use `package1.file2.prop`
		//because file3 is not in the same package
		//to have this, we need to include the files/packages nested within the same package in the variable processor
		//but we can't do this globally otherwise file3 could also just use `file1` or `file2` to access the classes
		var processorCpy = varprocessor.CloneProcessor(processor)

		classpack := c.ParentPackage

		for k, v := range classpack.ChildPacks {

			allparents := v.ReferenceFromStart()

			var (
				operationMac = new(ast.ASTNode)
				curRef       = operationMac
			)

			for kk, vv := range allparents {

				if kk+1 == len(allparents) {
					curRef.Group = &ast.VarRef{
						Name: vv.PackageName,
					}
					break
				}

				curRef.Left = []*ast.ASTNode{{
					Group: &ast.VarRef{
						Name: vv.PackageName,
					},
				}}
				curRef.Group = &ast.Operation{
					OpType: ".",
				}
				curRef.Right = make([]*ast.ASTNode, 1)
				curRef.Right[0] = new(ast.ASTNode)
				curRef = curRef.Right[0]
			}

			processorCpy.AddMacro(k, operationMac)
		}
		for k, v := range classpack.Classes {

			parents := v.ParentPackage.ReferenceFromStart()

			var (
				operationMac = new(ast.ASTNode)
				curRef       = operationMac
			)

			for _, vv := range parents {
				curRef.Left = []*ast.ASTNode{{
					Group: &ast.VarRef{
						Name: vv.PackageName,
					},
				}}
				curRef.Group = &ast.Operation{
					OpType: ".",
				}
				curRef.Right = make([]*ast.ASTNode, 1)
				curRef.Right[0] = new(ast.ASTNode)
				curRef = curRef.Right[0]
			}

			curRef.Group = &ast.VarRef{
				Name: v.CName,
			}

			processorCpy.AddMacro(k, operationMac)

		}

		processorCpy.ProcessVars(ic)
	}

	//use the malloc function to allocate instances/objects
	mallocf := m.NewFunc("malloc", types.I64Ptr)
	mallocf.Sig.Variadic = true
	compiler.LinkedFunctions["malloc"] = mallocf

	//use the free function to deallocate instances/objects
	freef := m.NewFunc("free", types.I64Ptr)
	freef.Sig.Variadic = true
	compiler.LinkedFunctions["free"] = freef

	for ic, c := range cclasses {

		for _, v := range ic.Globals {

			if v.CRel == 2 {
				//it's a linked function
				v.Link.Compile(&compiler, c, nil, nil)
				continue
			}

			if v.Func != nil {
				v.Func.DeclareGlobal(&compiler, c, parser.IsStatic(v), parser.IsPure(v), v.Access)
				continue
			}

			v.Value.DeclareGlobal(c.ParentPackage.FullName+"."+c.Name+"_"+v.Value.Name, &compiler, c, parser.IsStatic(v), v.Access)
		}

		//if there is a constructor, compile the signature
		if ic.Constructor != nil {
			ic.Constructor.CompileSig(&compiler, c, c.Construct)
		}
	}

	for ic, c := range cclasses {

		var newAlloc = c.Construct.ActiveBlock.NewAlloca(c.SType)
		c.ConstructAlloc = newAlloc

		for _, v := range ic.Globals {

			if v.Value != nil {
				var compileToFn *data.Function

				switch v.CRel {
				case 0:
					//instance
					compileToFn = c.Construct
				case 1:
					//static
					fallthrough
				case 3:
					compileToFn = compiler.InitFunc
				}

				v.Value.CompileGlobal(&compiler, c, compileToFn)
			} else if v.Func != nil {
				v.Func.CompileGlobal(&compiler, c, parser.IsStatic(v))
			}

		}

		if ic.Constructor != nil {
			//add the constructor into the mix of this jazz
			ic.Constructor.CompileConstructor(&compiler, c)
		}

		c.Construct.ActiveBlock.NewRet(c.Construct.ActiveBlock.NewLoad(c.SType, c.ConstructAlloc)) //return the allocated object at the end of the 'new' function
	}

	if compiler.Errors != nil {
		//there were some compile-time errors
		for _, v := range compiler.Errors {
			v.Print()
		}
		errhandle.PKill()
		return
	}

	//declare the llvm entry function
	mfunc := m.NewFunc("main", types.I32)
	mblock := mfunc.NewBlock("")

	compiler.InitFunc.ActiveBlock.NewRet(nil) //append a `return void` to the init function

	mblock.NewCall(initfunc) //call the initialize function

	//load and run the main function

	for _, v := range cclasses {
		if v.Name == prog.Config.Entry {
			//entry class
			mblock.NewCall(v.Static["main"].Value.LLVal(nil))
		}
	}

	mblock.NewRet(constant.NewInt(types.I32, 0))

	f, _ := os.Create(outfile)
	f.Write([]byte(m.String()))
}
