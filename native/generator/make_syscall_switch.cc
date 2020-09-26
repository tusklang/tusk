#include <fstream>
#include <map>
#include <vector>
#include <math.h>
#include "../ctypes.h"

//this file generates the if-else statement to make a syscall
/*
    if (a0.type == 0 && a1.type == 0) {
        syscall(*((double*) a0), *((double*) a1))
    } else if (a0.type == 0 && a1.type == 1) {
        syscall(*((double*) a0), *((char**) a1))
    }
    ...
    //yeah its pretty long, so by handing it would be insane
*/

#define vvuc std::vector<std::vector<unsigned char>> //it is so long, so just #define it

vvuc appendnew(vvuc initial) {
    vvuc newcombs;

    for (std::vector<unsigned char> i : initial)
        for (int o = 0; o < TYPE_AMT; ++o) {
            std::vector<unsigned char> sub = i;
            sub.push_back(CTYPES[o]);
            newcombs.push_back(sub);
        }

    return newcombs;
}

std::string getTypeFromNum(int num) {
    switch (num) {
        case 1:
            return "char*";
    }

    return "double";
}

int main() {

    std::string fname = "native/syscall_switch.h";

    //delete the file first
    remove(fname.c_str());

    //then create the file
    std::ofstream generatorfile;
    generatorfile.open(fname);

    /*
        HOWTO:

            example: types are 0, 1, 2

                (0)         (1)         (2)
               / | \
              /  |  \ --> repeat
             /   |   \
          (0,0)(0,1)(0,2)

                 |
                 v
         repeat for max argc
    */

    vvuc combs(TYPE_AMT);

    for (int i = 0; i < MAX_SYSCALL_ARGC - 1; ++i)
        combs = appendnew(combs);

    //write the actual file:

    generatorfile
        << "#include <unistd.h>\n"
        << "#include \"ctypes.h\"\n"
        << "static inline long int tusksyscall(";

    for (int i = 0; i < MAX_SYSCALL_ARGC; ++i)
        generatorfile << "struct CType a" << i << (i + 1 == MAX_SYSCALL_ARGC ? "" : ","); //if it is the last one, do not put a comma

    generatorfile << "){";

    for (std::vector<unsigned char> i : combs) {
        generatorfile << "if(";
        for (int o = 0; o < i.size(); ++o)
            generatorfile
                << "a"
                << o
                << ".type"
                << "=="
                << std::to_string(int(i[o]))
                << (o + 1 == i.size() ? "" : "&&"); //if it is the last one, do not put an &&

        generatorfile << ") return syscall(";

        for (int o = 0; o < i.size(); ++o)
            for (int j = 0; j < MAX_SYSCALL_ARGC; ++j) generatorfile
                << "*(("
                << getTypeFromNum(i[o])
                << "*)"
                << "a"
                << j
                << ".val"
                << ")"
                << (j + 1 == MAX_SYSCALL_ARGC ? "" : ","); //if it is the last one, no trailing comma
        
        generatorfile << ");";
    }

    generatorfile << "}";

    generatorfile.close();
}