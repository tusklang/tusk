#ifndef SYSTABLES_SYSCALLS_OPEN_H_
#define SYSTABLES_SYSCALLS_OPEN_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdio.h>

#ifdef _WIN32
#include <io.h>
#else
#include <fcntl.h>
#endif

#define special_open(name, mode) fileno(fopen(name, mode))

long int sysopen(char* name, int mode) {

    switch (mode) {

        //0 + 1 are the same
        case 0:
        case 1:
        ////////////////////

        //2 is append
        case 2:
            return special_open(name, "a");

        //3 is read + write
        case 3:
            return special_open(name, "r+");

        //4 is create empty file (for w+)
        case 4:
            return special_open(name, "w+");
        
        //5 is create empty file (for a+)
        case 5:
            return special_open(name, "a+");
    }

    return open(name, mode);
};

#ifdef __cplusplus
}
#endif

#endif