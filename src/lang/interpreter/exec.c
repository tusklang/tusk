#ifdef __cplusplus
extern "C" {
#endif

#include <string.h>
#include <stdlib.h>

char* getCmdExe() {
    //function to get the executable to execute a command

    #ifdef _WIN32
        char* plc = "\\windows\\system32\\cmd.exe"; //the location (not in any drive yet)

        //add the drive
        const char* drive = getenv("SystemDrive"); //get the default drive
        char* total = (char*) calloc(strlen(drive) + strlen(plc) - 1, sizeof(char)); //alocate the space

        strcpy(total, drive);
        strcat(total, plc);

        //add it to the heap
        char* heap_added = (char*) calloc(strlen(total), sizeof(char*));
        strcpy(heap_added, total);

        return heap_added;
    #else
        return "/bin/sh";
    #endif

}

char* getCmdOp() {
    //on windows it must be /C, but on unix/linux, it should be -c

    #ifdef _WIN32
        return "/C";
    #else
        return "-c";
    #endif

}

#ifdef __cplusplus
}
#endif