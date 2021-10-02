%tusk.string = type { i8*, i32 }
%tuskclass.Main = type { i32, void ()* }

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
	store void ()* @3, void ()** @.Main_main
	ret void
}

define %tuskclass.Main* @tuskclass.new.Main() {
0:
	%1 = alloca %tuskclass.Main
	%2 = getelementptr %tuskclass.Main, %tuskclass.Main* %1, i32 0, i32 0
	%3 = getelementptr %tuskclass.Main, %tuskclass.Main* %1, i32 0, i32 1
	store i32 4, i32* %2
	store void ()* @2, void ()** %3
	ret %tuskclass.Main* %1
}

declare void @0()

declare void @1()

define void @2() {
0:
	ret void
}

define void @3() {
0:
	%1 = call %tuskclass.Main* @tuskclass.new.Main()
	%2 = alloca %tuskclass.Main
	%3 = load %tuskclass.Main, %tuskclass.Main* %1
	store %tuskclass.Main %3, %tuskclass.Main* %2
	%4 = call %tuskclass.Main* @tuskclass.new.Main()
	%5 = load %tuskclass.Main, %tuskclass.Main* %4
	store %tuskclass.Main %5, %tuskclass.Main* %2
	ret void
}

define void @main() {
0:
	call void @_tusk_init()
	ret void
}
