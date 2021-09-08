%tuskclass.test.Test2 = type {}
%tuskclass.Main = type {}
%tuskclass.Test = type {}

@Main_main = global void ()*

define void @_tusk_init() {
0:
	%1 = load void ()*, void ()* @tv_1
	store void ()* %1, void ()** @Main_main
	ret void
}

define void @tv_1() {
0:
	%1 = alloca %tuskclass.test.Test2
	ret void
}

declare void @tv_2()

define void @main() {
0:
	call void @_tusk_init()
	ret void
}
