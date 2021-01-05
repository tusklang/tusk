#include <dlfcn.h>

int main()
{
    void *handle = dlopen("./test.so", RTLD_NOW);
    int (*a)();
    a = dlsym(handle, "test");
    printf("%d\n", a());
}