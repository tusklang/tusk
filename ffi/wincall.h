#ifndef WINDOWS_OMMFFI_H_
#define WINDOWS_OMMFFI_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32

#include <stdbool.h>
#include <windows.h>

struct LibraryRet {
    struct HINSTANCE__** module;
    bool                 error;
};

struct LibraryRet loadlib(char*);
void** callproc(struct HINSTANCE__**, char*, void**, int);

#endif

#ifdef __cplusplus
}
#endif

#endif