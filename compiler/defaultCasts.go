package compiler

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
)

var (
	inttypes   = []string{"i8", "i16", "i32", "i64", "i128"}
	uinttypes  = []string{"u8", "u16", "u32", "u64", "u128"}
	floattypes = []string{"f32", "f64"}
)

func addCastArray(compiler *ast.Compiler, typArr []string, fromType string, fn func(tname string, fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value) {
	for _, _v := range typArr {
		v := _v
		compiler.CastStore.NewCast(true, v, fromType, func(toType data.Type, fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return fn(v, fromData, compiler, function, class)
		})
	}
}

func addXCasts2(auto, slice bool, compiler *ast.Compiler, fromArr []string, toArr []string, fn func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value) {
	for k, _v := range fromArr {
		var v = _v

		sl := 0

		if slice {
			sl = k
		}

		for _, _vv := range toArr[sl:] {
			var vv = _vv

			if v == vv {
				continue
			}

			compiler.CastStore.NewCast(auto, vv, v, func(toType data.Type, fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
				return data.NewInstVariable(fn(fromData, compiler, function, class, numtypes[vv].Type()), numtypes[vv])
			})

		}
	}
}

func addXCasts(auto bool, compiler *ast.Compiler, typArr []string, fn func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value) {
	addXCasts2(auto, true, compiler, typArr, typArr, fn) //just do this function, but the outer and inner loops params are equal
}

//reverse a string array (type arrays)
func reverseStrArr(a []string) []string {
	var fin = make([]string, len(a))
	for k, v := range a {
		fin[len(a)-k-1] = v
	}
	return fin
}

func initDefaultCasts(compiler *ast.Compiler) {
	compiler.CastStore = ast.NewCastStore()

	//add upcasts
	addXCasts(true, compiler, inttypes, func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewSExt(fromData.LLVal(function), typ)
	})
	addXCasts(true, compiler, uinttypes, func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewZExt(fromData.LLVal(function), typ)
	})
	addXCasts(true, compiler, floattypes, func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewFPExt(fromData.LLVal(function), typ)
	})

	//add downcasts
	addXCasts(false, compiler, reverseStrArr(inttypes), func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewTrunc(fromData.LLVal(function), typ)
	})
	addXCasts(false, compiler, reverseStrArr(uinttypes), func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewTrunc(fromData.LLVal(function), typ)
	})
	addXCasts(false, compiler, reverseStrArr(floattypes), func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewFPTrunc(fromData.LLVal(function), typ)
	})

	//add casts between int/uint/float types
	addXCasts2(false, true, compiler, inttypes, uinttypes, func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewSExt(fromData.LLVal(function), typ)
	})
	addXCasts2(false, false, compiler, inttypes, floattypes, func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewSIToFP(fromData.LLVal(function), typ)
	})
	addXCasts2(false, true, compiler, uinttypes, inttypes, func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewZExt(fromData.LLVal(function), typ)
	})
	addXCasts2(false, false, compiler, uinttypes, floattypes, func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewUIToFP(fromData.LLVal(function), typ)
	})
	addXCasts2(false, false, compiler, floattypes, inttypes, func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewFPToSI(fromData.LLVal(function), typ)
	})
	addXCasts2(false, false, compiler, floattypes, uinttypes, func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewFPToUI(fromData.LLVal(function), typ)
	})

	//and also for downcasts
	addXCasts2(false, true, compiler, reverseStrArr(inttypes), reverseStrArr(uinttypes), func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewTrunc(fromData.LLVal(function), typ)
	})
	addXCasts2(false, true, compiler, reverseStrArr(uinttypes), reverseStrArr(inttypes), func(fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class, typ types.Type) value.Value {
		return function.ActiveBlock.NewTrunc(fromData.LLVal(function), typ)
	})

	//add casts from untyped numeric vals
	addCastArray(compiler, append(inttypes, uinttypes...), "untypedint", func(tname string, fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(constant.NewInt(numtypes[tname].Type().(*types.IntType), fromData.(*data.Integer).UTypVal), numtypes[tname])
	})

	addCastArray(compiler, floattypes, "untypedint", func(tname string, fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(constant.NewFloat(numtypes[tname].Type().(*types.FloatType), float64(fromData.(*data.Integer).UTypVal)), numtypes[tname])
	})

	addCastArray(compiler, append(inttypes, uinttypes...), "untypedfloat", func(tname string, fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(constant.NewInt(numtypes[tname].Type().(*types.IntType), int64(fromData.(*data.Float).UTypVal)), numtypes[tname])
	})

	addCastArray(compiler, floattypes, "untypedfloat", func(tname string, fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(constant.NewFloat(numtypes[tname].Type().(*types.FloatType), fromData.(*data.Float).UTypVal), numtypes[tname])
	})

	//other casts
	compiler.CastStore.NewCast(true, "slice", "fixed", func(toType data.Type, fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {

		if !fromData.TypeData().GetOtherDat("valtyp").(data.Type).Equals(toType.TypeData().GetOtherDat("valtyp").(data.Type)) {
			//error
			//slice type and fixed type don't match
			//this is someting like:
			//	[]i32 -> []i64{};
			//most likely
		}

		toTypet := toType.(*data.SliceArray)

		malloc := compiler.LinkedFunctions["malloc"]        //fetch the malloc function
		length := fromData.TypeData().GetOtherDat("length") //fetch the length of the fixed arr

		ftyp := types.NewPointer(toTypet.ValType().Type())

		alc := function.ActiveBlock.NewAlloca(ftyp)
		alc.Align = ir.Align(8)

		function.ActiveBlock.NewStore(
			function.ActiveBlock.NewBitCast(
				function.ActiveBlock.NewCall(
					malloc,
					function.ActiveBlock.NewMul(
						length.LLVal(function),
						constant.NewInt(length.Type().(*types.IntType), int64(toTypet.ValType().Alignment())),
					),
				),
				ftyp,
			),
			alc,
		)

		loaded := function.ActiveBlock.NewLoad(ftyp, alc)

		lengthVal := length.(*data.Integer).GetInt()

		//see the indexing operation source for why we allocate it again
		fromAlc := function.ActiveBlock.NewAlloca(fromData.Type())
		function.ActiveBlock.NewStore(fromData.LLVal(function), fromAlc)

		for i := 0; i < int(lengthVal); i++ {
			gepTo := function.ActiveBlock.NewGetElementPtr(toTypet.ValType().Type(), loaded, constant.NewInt(types.I32, int64(i)))
			gepFrom := function.ActiveBlock.NewGetElementPtr(fromData.Type(), fromAlc, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(i)))
			function.ActiveBlock.NewStore(function.ActiveBlock.NewLoad(toTypet.ValType().Type(), gepFrom), gepTo)
		}

		return data.NewInstVariable(loaded, data.NewPointer(toTypet.ValType()))
	})

	compiler.CastStore.NewCast(true, "slice", "varied", func(toType data.Type, fromData data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {

		toTypet := toType.(*data.SliceArray)                //convert the toType to a slice array (we know it is one)
		malloc := compiler.LinkedFunctions["malloc"]        //fetch the malloc function
		length := fromData.TypeData().GetOtherDat("length") //fetch the length of the fixed arr
		lllength := length.LLVal(function)                  //get the llvm value of length

		ftyp := types.NewPointer(toTypet.ValType().Type())
		alc := function.ActiveBlock.NewAlloca(ftyp)
		alc.Align = ir.Align(8)

		function.ActiveBlock.NewStore(
			function.ActiveBlock.NewBitCast(
				function.ActiveBlock.NewCall(
					malloc,
					function.ActiveBlock.NewMul(
						lllength,
						constant.NewInt(length.Type().(*types.IntType), int64(toTypet.ValType().Alignment())),
					),
				), ftyp,
			),
			alc,
		)

		loaded := function.ActiveBlock.NewLoad(ftyp, alc)

		//get the old body block
		oldb := function.ActiveBlock

		counter := oldb.NewAlloca(length.Type())
		ltyp := length.Type().(*types.IntType)
		oldb.NewStore(constant.NewInt(ltyp, 0), counter)

		//create a new block
		//this block will check if the current counter is less than the length of the varied arr
		condb := function.LLFunc.NewBlock("")
		oldb.NewBr(condb)

		//create a new block
		//this block will loop through all the values of the varied array
		//and put them into the slice array
		bodb := function.LLFunc.NewBlock("")
		bodb.NewBr(condb)

		//create a new block
		//this block will contain all further instructions
		restb := function.LLFunc.NewBlock("")

		loadedCond := condb.NewLoad(length.Type(), counter)
		condb.NewCondBr(
			condb.NewICmp(enum.IPredULT, loadedCond, lllength),
			bodb,
			restb,
		)

		function.ActiveBlock = bodb

		gepTo := bodb.NewGetElementPtr(toTypet.ValType().Type(), loaded, loadedCond)
		gepFrom := bodb.NewGetElementPtr(toTypet.ValType().Type(), fromData.LLVal(function), loadedCond)
		bodb.NewStore(bodb.NewLoad(toTypet.ValType().Type(), gepFrom), gepTo)

		bodb.NewStore(bodb.NewAdd(loadedCond, constant.NewInt(ltyp, 1)), counter)

		function.ActiveBlock = restb

		return data.NewInstVariable(loaded, data.NewPointer(toTypet.TType()))
	})

}
