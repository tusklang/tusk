package compiler

import (
	"os"
	"runtime"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/initialize"
	"github.com/tusklang/tusk/varprocessor"
)

var processor = varprocessor.NewProcessor()

var numtypes = map[string]data.Type{
	"i128": data.NewPrimitive(types.I128),
	"i64":  data.NewPrimitive(types.I64),
	"i32":  data.NewPrimitive(types.I32),
	"i16":  data.NewPrimitive(types.I16),
	"i8":   data.NewPrimitive(types.I8),
	"u128": data.NewNamedPrimitive(types.I128, "u128"),
	"u64":  data.NewNamedPrimitive(types.I64, "u64"),
	"u32":  data.NewNamedPrimitive(types.I32, "u32"),
	"u16":  data.NewNamedPrimitive(types.I16, "u16"),
	"u8":   data.NewNamedPrimitive(types.I8, "u8"),
	"f64":  data.NewNamedPrimitive(types.Double, "f64"),
	"f32":  data.NewNamedPrimitive(types.Float, "f32"),
}

//list of all the variables that are added by default
//has types to begin with, but it can store anything
//types are variables in tusk's parser so we need to add the default ones in like so
var prevars = map[string]data.Value{}

func Compile(prog *initialize.Program, outfile string) {

	var compiler ast.Compiler
	m := ir.NewModule() //create a new llvm module

	//set the module stuff in the compiler
	compiler.Module = m
	compiler.OS = runtime.GOOS
	compiler.ARCH = runtime.GOARCH

	compiler.VarMap = make(map[string]data.Value)
	compiler.LinkedFunctions = make(map[string]*ir.Func)

	initDefaultOps(&compiler)
	initDefaultCasts(&compiler)

	for k, v := range numtypes {
		prevars[k] = v.(data.Value)
	}

	//initialize the operations
	var initfunc = m.NewFunc("_tusk_init", types.Void) //initialize func ran before main
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
		//but we can't do this globally
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
			_, _ = k, v
		}

		processorCpy.ProcessVars(ic)

		_, _ = ic, c
	}

	//function used to malloc classes
	mallocf := m.NewFunc("malloc", types.I8Ptr, ir.NewParam("", types.I64))

	for ic, c := range cclasses {

		for _, v := range ic.Globals {

			if v.CRel == 2 {
				//it's a linked function
				v.Link.Compile(&compiler, c, nil, nil)
				continue
			}

			if v.Func != nil {
				v.Func.DeclareGlobal(&compiler, c, v.CRel == 1, v.Access)
				continue
			}

			v.Value.DeclareGlobal(c.ParentPackage.FullName+"."+c.Name+"_"+v.Value.Name, &compiler, c, v.CRel == 1, v.Access)
		}

		//if there is a constructor, compile the signature
		if ic.Constructor != nil {
			ic.Constructor.CompileSig(&compiler, c, c.Construct)
		}
	}

	for ic, c := range cclasses {

		var newAlloc = c.Construct.ActiveBlock.NewAlloca(types.NewPointer(c.SType))
		c.Construct.ActiveBlock.NewStore(
			c.Construct.ActiveBlock.NewBitCast(
				c.Construct.ActiveBlock.NewCall(mallocf, constant.NewInt(types.I64, int64(c.TypSiz))),
				types.NewPointer(c.SType),
			),
			newAlloc,
		)

		c.ConstructAlloc = c.Construct.ActiveBlock.NewLoad(types.NewPointer(c.SType), newAlloc)

		for _, v := range ic.Globals {

			if v.Value != nil {
				var compileToFn *data.Function

				switch v.CRel {
				case 0:
					//instance
					compileToFn = c.Construct
				case 1:
					//static
					compileToFn = compiler.InitFunc
				}

				v.Value.CompileGlobal(&compiler, c, compileToFn)
			} else if v.Func != nil {
				v.Func.CompileGlobal(&compiler, c, v.CRel == 1)
			}

		}

		if ic.Constructor != nil {
			//add the constructor into the mix of this jazz
			ic.Constructor.CompileConstructor(&compiler, c)
		}

		c.Construct.ActiveBlock.NewRet(c.ConstructAlloc) //return the allocated object at the end of the 'new' function
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
			mblock.NewCall(v.Static["main"].Value.LLVal(mblock))
		}
	}

	mblock.NewRet(constant.NewInt(types.I32, 0))

	f, _ := os.Create(outfile)
	f.Write([]byte(m.String()))
}
