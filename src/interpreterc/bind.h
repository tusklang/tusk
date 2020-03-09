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
  void bind(char *actions, char *calc_params, char *dir);
#ifdef __cplusplus
}
#endif
