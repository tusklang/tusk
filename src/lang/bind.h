#ifdef __cplusplus
extern "C" {
#endif
  extern void Kill(void);
  extern char* Add(char*, char*, char*);
  extern char* Subtract(char*, char*, char*);
  extern char* Multiply(char*, char*, char*);
  extern char* Division(char*, char*, char*);
  extern char* Modulo(char*, char*, char*);
  extern char* Exponentiate(char*, char*, char*);
  extern char* GetType(char*);
  extern char* ReturnInitC(char*);
  extern int IsLessC(char*, char*);
  extern int GetActNumC(char*);
  extern char* AddC(char*, char*);
  extern char* GetOp(char*);
  extern char* Similar(char*, char*, char*, char*, int, char*, char*);
  extern char* AddStrings(char*, char*, char*);
  extern char* SubtractStrings(char*, char*, char*);
  void bindCgo(char* actions, char* cli_params);
  void colorprint(char*, int);
#ifdef __cplusplus
}
#endif
