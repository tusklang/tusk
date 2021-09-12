%tusk.string = type { i8*, i32 }
%tuskclass.Main = type {}

@.Main_a = global %tusk.string ()* null
@.Main_b = global i32 ()* null
@.Main_main = global void ()* null
@tv_5 = global [5 x i8] c"hello"

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
	store %tusk.string ()* @tv_4, %tusk.string ()** @.Main_a
	store i32 ()* @tv_6, i32 ()** @.Main_b
	store void ()* @tv_7, void ()** @.Main_main
	ret void
}

declare %tusk.string @tv_1()

declare i32 @tv_2()

declare void @tv_3()

define %tusk.string @tv_4() {
0:
	%1 = call %tusk.string @tusk.newstring(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @tv_5, i32 0, i32 0), i32 5)
	ret %tusk.string %1
}

define i32 @tv_6() {
0:
	ret i32 33
}

define void @tv_7() {
0:
	ret void
}

define void @main() {
0:
	call void @_tusk_init()
	ret void
}
