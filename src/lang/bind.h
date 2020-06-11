#ifndef OMM_CGO_LANG_BIND_H_
#define OMM_CGO_LANG_BIND_H_

#ifdef __cplusplus
extern "C" {
#endif
//

  extern void Kill(void);
  extern char* GetOp(char*); //get operation alias, e.g. add --> +, subtract --> -, etc.
  extern char* ExecCmd(char*, char*, char*); //funtion to execute

  void bindParser(char* actions, char* cli_params, char* dir, int argc, char ** argv);
  void colorprint(char*, int);

#ifdef __cplusplus
}
#endif

#endif
