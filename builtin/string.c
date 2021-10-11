#include <stdlib.h>

struct string
{
    int _length;
    char *_data;
    int (*length)();
};

int string_length(struct string *self)
{
    return (*self)._length;
}

struct string *construct(char *d, int len)
{
    struct string *s = (struct string *)calloc(1, sizeof(struct string));
    (*s)._data = d;
    (*s)._length = len;
    (*s).length = &string_length;
    return s;
}
