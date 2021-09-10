%tuskclass.Main = type {}

@.Main_main = global void ()* null

define void @_tusk_init() {
0:
	store void ()* @tv_2, void ()** @.Main_main
	ret void
}

declare void @tv_1()

define void @tv_2() {
0:
	call void @tv_3()
	ret void
}

define void @tv_3() {
0:
	%1 = add i32 1, 1
	ret void
}

define void @main() {
0:
	call void @_tusk_init()
	ret void
}
