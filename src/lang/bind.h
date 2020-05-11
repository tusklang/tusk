#ifdef __cplusplus
extern "C" {
#endif

  extern void Kill(void);
  extern char* AddC(char*, char*, char*);
  extern char* SubtractC(char*, char*, char*);
  extern char* MultiplyC(char*, char*, char*);
  extern char* DivisionC(char*, char*, char*);
  extern char* GetType(char*);
  extern char* ReturnInitC(char*);
  extern int IsLessC(char*, char*);
  extern char* GetOp(char*);

  void bindParser(char* actions, char* cli_params);
  void colorprint(char*, int);

#ifdef __cplusplus
}
#endif
