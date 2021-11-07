# TODO

- Fix linked functions
Linked functions in different files will have the same name in the generated .ll file, giving an error
And we can't just make all these functions connected to the same function in IR
Consider this example

```
//test1.tusk

link var printf: fn(*i8, i32) -> printf; //link to printf, where we can print one integer
```

```
//test2.tusk

link var printf: fn(*i8, *i8) -> printf; //link to printf, where we can print one string
```

This would compile to an .ll file with two definitions of printf, where each has a different signature.