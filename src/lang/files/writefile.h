#ifndef WRITEFILE_H_
#define WRITEFILE_H_

#ifdef __cplusplus
extern "C" {
#endif;
  #include <stdio.h>

  void writefile(char* dir, char* content) {

    FILE* f = fopen(dir, "w+");
    fputs(content, f);

    fclose(f);
  }

#ifdef __cplusplus
}
#endif

#endif
