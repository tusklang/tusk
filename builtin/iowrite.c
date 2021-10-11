#ifdef _WIN32
#include <io.h>
#else
#include <unistd.h>
#endif

int iowrite(int fd, char *s, int len)
{
    return write(fd, s, len);
}