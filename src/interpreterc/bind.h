#ifdef __cplusplus
extern "C" {
#endif
  extern void Kill(void);
  extern char* Add(char*, char*, char*, int);
  extern char* Subtract(char*, char*, char*, int);
  extern char* Multiply(char*, char*, char*, int);
  extern char* Division(char*, char*, char*, int);
  extern char* Modulo(char*, char*, char*, int);
  extern char* Exponentiate(char*, char*, char*, int);
  extern char* Cactions(char*);
  extern char* GetType(char*);
  extern char* ReturnInitC(char*);
  extern int IsLessC(char*, char*);
  extern char* CReadFile(char*, char*, int);
  extern char* CLex(char*);
  extern char* NQReplaceC(char*);
  extern int GetActNumC(char*);
  extern char* AddC(char*, char*);
  extern char* GetOp(char*);
  extern char* Similar(char*, char*, char*, char*, int, char*, char*);
  extern char* AddStrings(char*, char*, char*, int);
  extern char* SubtractStrings(char*, char*, char*, int);
  void bind(char *actions, char *calc_params, char *dir);
  char* parser_exp(const char* actionsP, const char* calc_paramsP, const char* varsP, const char* dirP, const int groupReturn, int line, const int expReturn);
#ifdef __cplusplus
}
#endif
