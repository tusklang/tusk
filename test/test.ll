%Main = type {}

@Main_main = global void ()* null
@Main_a = global i32 0

define void @_init() {
0:
	store void ()* @tv_1, void ()** @Main_main
	store i32 43, i32* @Main_a
	ret void
}

define void @tv_1() {
0:
	%1 = add i32 32, 43
	ret void
}

define void @main() {
0:
	call void @_init()
	%1 = load void ()*, void ()** @Main_main
	call void %1()
	ret void
}
