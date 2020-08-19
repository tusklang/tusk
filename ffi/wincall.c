#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32

#include <windows.h>
#include <stdlib.h>
#include "wincall.h"

//prevent linter from complaining
#if __has_include("callswitch.h")
#include "callswitch.h"
#else                     
long long int callswitch(FARPROC proc, int argc, void** argv) {
    return 0;
}
#endif
////////////////////////////////

struct LibraryRet loadlib(char* filepath) {
    struct HINSTANCE__* module = LoadLibrary(filepath); //load the library into an struct HINSTANCE__
    free(filepath); //free filepath (since it was c/malloc(ed) in C.CString)

    struct LibraryRet final;
    final.error = 0;

    if (module == NULL) { //if the module is NULL, return an error
        final.error = 1; //set the error to true
        return final; //return
    }

    final.module = calloc(1, sizeof(struct HINSTANCE__*)); //calloc the space
    
    *final.module = module; //set the module
    return final; //return
}

void** callproc(struct HINSTANCE__** module, char* fname, void** argv, int argc) {

    FARPROC proc = GetProcAddress(*module, fname); //get the function requested
    free(fname); //free values
    void** ret = (void**) calloc(1, sizeof(void*)); //calloc the values
    long long int called = callswitch(proc, 0, NULL); //call the proc
    *ret = (void*) &called;
    
    free(argv); //free the argv
    return ret; //return
}

#endif

#ifdef __cplusplus
}
#endif