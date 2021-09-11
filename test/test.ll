%tusk.string = type { i8*, i32 }
%tuskclass.main = type {}

@.main_main = global void ()* null
@tv_3 = global [11 x i8] c"hello world"

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
	store void ()* @tv_2, void ()** @.main_main
	ret void
}

declare void @tv_1()

define void @tv_2() {
0:
	%1 = alloca %tusk.string
	%2 = call %tusk.string @tusk.newstring(i8* getelementptr inbounds ([11 x i8], [11 x i8]* @tv_3, i32 0, i32 0), i32 11)
	store %tusk.string %2, %tusk.string* %1
	ret void
}

define void @main() {
0:
	call void @_tusk_init()
	ret void
}
