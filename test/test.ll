%tuskclass.Main = type {}

@.Main_test = global i32 0
@.Main_test2 = global i32 0
@.Main_main = global void ()* null

define void @_tusk_init() {
0:
	%1 = load i32, i32* @.Main_test2
	store i32 %1, i32* @.Main_test
	%2 = load i32, i32* @.Main_test
	store i32 %2, i32* @.Main_test2
	store void ()* @tv_2, void ()** @.Main_main
	ret void
}

declare void @tv_1()

define void @tv_2() {
0:
	%1 = load i32, i32* @.Main_test
	%2 = add i32 %1, 1
	%3 = alloca i32
	store i32 %2, i32* %3
	ret void
}

define void @main() {
0:
	call void @_tusk_init()
	ret void
}
