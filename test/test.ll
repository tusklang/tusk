%tusk.string = type { i8*, i32 }
%tuskclass.Main = type {}

@.Main_main = global void ()* null

define %tusk.string @tusk.newstring(i8* %sptr, i32 %slen) {
0:
	%1 = alloca %tusk.string
	%2 = getelementptr %tusk.string, %tusk.string* %1, i32 0, i32 0
	store i8* %sptr, i8** %2
	%3 = getelementptr %tusk.string, %tusk.string* %1, i32 0, i32 1
	store i32 %slen, i32* %3
	%4 = load %tusk.string, %tusk.string* %1
	ret %tusk.string %4
}

define void @_tusk_init() {
0:
	store void ()* @1, void ()** @.Main_main
	ret void
}

define %tuskclass.Main* @tuskclass.new.Main() {
0:
	%1 = alloca %tuskclass.Main
	ret %tuskclass.Main* %1
}

declare void @0()

define void @1() {
0:
	%1 = alloca i32
	store i32 33, i32* %1
	%2 = alloca i32*
	store i32* %1, i32** %2
	%3 = alloca i32**
	store i32** %2, i32*** %3
	%4 = load i32**, i32*** %3
	%5 = load i32*, i32** %4
	%6 = alloca i32
	%7 = load i32, i32* %5
	store i32 %7, i32* %6
	ret void
}

define void @main() {
0:
	call void @_tusk_init()
	ret void
}
