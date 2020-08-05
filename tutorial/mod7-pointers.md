# Pointer Variables

In Omm, every variable is a reference (pointer). Lets take an example

```clojure
a := 1
b := a
b = 5
;a is is now 5
```

This is called a pointer, and many languages (C, C++, D, and Go) have this feature. To make `a` stay as `1`, we need to perform an indirect of `a`. Here is an example from C

```c
int _a = 5;
int* a = &_a;
int b = *a;
```

Omm can do this with the `clone[]` built in function.

```clojure
a := 1
b := clone[a]
b = 4
;a stays as 1
```
