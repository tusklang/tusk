#ifdef __cplusplus
extern "C" {
#endif
  #include "bind.h"
  #include <stdio.h>
  #include <string.h>
  #include <stdlib.h>

  void write(char* dir, char* name, char* contents) {

    if (IsAbsolute(name)) {

      FILE* f = fopen(dir, "w");
      fprintf(f, contents);

      fclose(f);
    } else {

      char* filepath = (char*) malloc(1 + strlen(dir) + strlen(name));

      strcpy(filepath, dir);
      strcat(filepath, name);

      FILE* f = fopen(filepath, "w");
      fprintf(f, contents);

      fclose(f);
    }
  }

#ifdef __cplusplus
}
#endif
