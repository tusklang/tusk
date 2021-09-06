%Main = type {}

@Main_main = global void ()* null
@Main_a = global i64 0

define void @_init() {
0:
	store void ()* @tv_1, void ()** @Main_main
	%1 = sext i32 32 to i64
	store i64 %1, i64* @Main_a
	ret void
}

define void @tv_1() {
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
