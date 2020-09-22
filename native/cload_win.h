#ifdef _WIN32

#ifndef TUSK_NATIVE_CLOAD_WIN_H_
#define TUSK_NATIVE_CLOAD_WIN_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <windows.h>
#include "cload.h"

struct CLib {
    HINSTANCE library;
    bool     err;
};

struct CProc {
    FARPROC proc;
    bool err;
};

static inline struct CLib cloadlib(char* name) {
    struct CLib ret;
    ret.err = false; //set it to none (by default)
    HINSTANCE lib = LoadLibraryA(name);

    if (!lib) ret.err = true;
    ret.library = lib;

    free(name);
    return ret;
}

static inline struct CProc cgetproc(struct CLib lib, char* name) {
    struct CProc ret;
    ret.err = false; //set it to none (by default)
    FARPROC proc = GetProcAddress(lib.library, name);

    if (!proc) ret.err = true;
    ret.proc = proc;

    free(name);
    return ret;
}

static inline double ccallproc(struct CProc proc, void** args, int argc) {
    FARPROC fproc = proc.proc;

    for (;argc < MAX_ARGC; ++argc)
        args[argc] = NULL; //set the length of args to be MAX_ARGC

    //just dump all of the args here
    return (double) fproc(
        args[0],
        args[1],
        args[2],
        args[3],
        args[4],
        args[5],
        args[6],
        args[7],
        args[8],
        args[9],
        args[10],
        args[11],
        args[12],
        args[13],
        args[14],
        args[15],
        args[16],
        args[17],
        args[18],
        args[19]
    );
}

//free a dll
static inline void freelib(struct CLib lib) {
    FreeLibrary(lib.library);
}

#ifdef __cplusplus
}
#endif

#endif

#endif