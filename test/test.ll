%tuskclass.Main = type {}
%tuskclass.Test = type {}

@.Test_a = global i32 0
@.Main_main = global void ()* null

define void @_tusk_init() {
0:
	store void ()* @tv_2, void ()** @.Main_main
	store i32 43, i32* @.Test_a
	ret void
}

declare void @tv_1()

define void @tv_2() {
0:
	%1 = alloca i32
	%2 = load i32, i32* @.Test_a
	store i32 %2, i32* %1
	ret void
}

define void @main() {
0:
	call void @_tusk_init()
	ret void
}
