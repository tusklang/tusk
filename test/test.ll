%Main = type { %Test }
%Test = type {}

@Main_main = global void ()* @f1
@Test_main = global void ()* @f2

define void @f1() {
0:
	ret void
}

define void @f2() {
0:
	ret void
}

define void @main() {
0:
	%1 = load void ()*, void ()** @Test_main
	call void %1()
	ret void
}
