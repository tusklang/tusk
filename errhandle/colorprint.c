#include "colorprint.h"

#ifdef _WIN32
void errprint(char *msg)
{
    // TODO
}
#else

#include <stdio.h>

#define RESET "\033[0m"
#define RED "\33[1;31m"
#define BLUE "\x1B[36m"

void compileErrorPrint(char *msg)
{
    fprintf(stderr, "%scompiler error: %s%s\n", RED, msg, RESET);
}

void parseErrorPrint(char *msg)
{
    fprintf(stderr, "%sparsing error: %s%s\n", BLUE, msg, RESET);
}
#endif