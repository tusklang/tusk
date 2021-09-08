%Main = type {}
%Test = type { i32 }

@Main_main = global void ()* null

define void @_tusk_init() {
0:
	store void ()* @tv_1, void ()** @Main_main
	ret void
}

define void @tv_1() {
0:
	%1 = alloca %Main
	ret void
}

declare void @tv_2()

define void @main() {
0:
	call void @_tusk_init()
	%1 = load void ()*, void ()** @Main_main
	call void %1()
	ret void
}
