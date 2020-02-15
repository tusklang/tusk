#include <iostream>
#include <string>
#include <windows.h>
using namespace std;

void printer(int color, string val) {

  HANDLE hConsole;

  hConsole = GetStdHandle(STD_OUTPUT_HANDLE);
  SetConsoleTextAttribute(hConsole, color);
  cout << val << endl;

  SetConsoleTextAttribute(hConsole, 7);
}
