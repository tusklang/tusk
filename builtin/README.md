# Built-in Functions and Classes

## How does this work?

This folder is home to quite a few C files. Each of these files stores a non-primitive class that is present in Tusk. We define these classes as C structures, and they're later converted to LLVM-IR and parsed by the compiler to include in a Tusk build.

## What are these .config files?

LLVM does not provide all the data needed to make a complete Tusk class. For example, LLVM-IR doesn't include field names within structs. Tusk accesses fields by name, so this is an issue. Another is private variables. Some built-in classes require private instance variables (e.g. the string class prevents you from editing the raw char pointer). LLVM and C both have no concept of encapsulation, so this config file is used to define that.

In general, a config file for classes are structured like so:
```
class
-
<field index> <field name> <field accessibility>
-
<class name>
-
<names to not mangle>
```

And for functions:
```
func
-
<function name>
-
<names to not mangle>
```