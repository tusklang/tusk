#ifndef PRINT_C_
#define PRINT_C_

#ifdef __cplusplus
extern "C" {
#endif
  #include <windows.h>
  #include <stdio.h>

  void colorprint(char* str, int color) {

    HANDLE hConsole = GetStdHandle(STD_OUTPUT_HANDLE);

    CONSOLE_SCREEN_BUFFER_INFO console_info;
    GetConsoleScreenBufferInfo(GetStdHandle(STD_OUTPUT_HANDLE), &console_info);

    //get current color (so we can change it back after printing)
    short cur_color = console_info.wAttributes;

    SetConsoleTextAttribute(hConsole, color);
    printf("%s", str);

    SetConsoleTextAttribute(hConsole, cur_color);
  }

#ifdef __cplusplus
}
#endif

#endif
