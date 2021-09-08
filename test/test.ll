%Main = type {}

@Main_main = global void ()* null

define void @_init() {
0:
	store void ()* @tv_1, void ()** @Main_main
	ret void
}

define void @tv_1() {
0:
	%1 = alloca void ()*
	store void ()* @tv_2, void ()** %1
	ret void
}

define void @tv_2() {
0:
	ret void
}

define void @main() {
0:
	call void @_init()
	%1 = load void ()*, void ()** @Main_main
	call void %1()
	ret void
}
