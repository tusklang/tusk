#ifndef SYSTABLES_SYSCALLS_GETDIR_H_
#define SYSTABLES_SYSCALLS_GETDIR_H_

#ifdef __cplusplus
extern "C"
{
#endif

#include <stdio.h>
#include <unistd.h>
#include <dirent.h>
#include <limits.h>
#include <sys/stat.h>

#ifndef _WIN32
#define MAX_PATH FILENAME_MAX
#endif

    long long int syslsdir(long long int loc, void **subnames, void **fsnumbers, int len)
    {
        DIR *dir = (DIR *)loc;

        for (int i = 0; i < len; i++)
        {
            struct dirent *d = readdir(dir); //get the next file
            if (d == NULL)
                break; //it does not exist
            subnames[i] = (void *)d->d_name;
            fsnumbers[i] = (void *)((long long int)d->d_ino);
        }

        return 0;
    }

    long long int syssizedir(long long int loc)
    {
        int file_count = 0;
        DIR *dirp = (DIR *)loc;
        struct dirent *entry;

        for (; (entry = readdir(dirp)) != NULL; file_count++)
            ;
        return 0;
    }

    long long int sysclosedir(long long int loc)
    {
        return closedir((DIR *)loc);
    }

    long long int sysgetcwd(char *buf)
    {
        buf = (char *)realloc(buf, sizeof(char) * MAX_PATH);
        char *_ = getcwd(buf, MAX_PATH);
        return 0;
    }

    long long int syschdir(char *path)
    {
        return chdir(path);
    }

    long long int sysrename(char *oldp, char *newp)
    {
        return rename(oldp, newp);
    }

    long long int sysmkdir(char *path)
    {
        return mkdir(path
//all perms
#ifndef _WIN32
                     ,
                     S_IRWXU | S_IRWXG | S_IROTH | S_IXOTH
#endif
        );
    }

    long long int sysrmdir(char *path)
    {
        return rmdir(path);
    }

#ifdef __cplusplus
}
#endif

#endif