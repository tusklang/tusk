#ifndef OMM_CGO_LANG_BIND_H_
#define OMM_CGO_LANG_BIND_H_

#ifdef __cplusplus
extern "C" {
#endif
//

  extern void Kill(void);
  extern char* AddC(char*, char*, char*);
  extern char* SubtractC(char*, char*, char*);
  extern char* MultiplyC(char*, char*, char*);
  extern char* DivisionC(char*, char*, char*);
  extern char* ReturnInitC(char*);
  extern int IsLessC(char*, char*);
  extern char* GetOp(char*); //get operation alias, e.g. add --> +, subtract --> -, etc.
  extern char* ExecCmd(char*, char*, char*); //funtion to execute 

  void bindParser(char* actions, char* cli_params, char* dir, int argc, char ** argv);
  void colorprint(char*, int);

#ifdef __cplusplus
}
#endif

#endif
