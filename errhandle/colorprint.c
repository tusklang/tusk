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

void errprint(char *msg)
{
    fprintf(stderr, "%s%s%s\n", RED, msg, RESET);
}
#endif