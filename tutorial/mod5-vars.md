# Variables

Variables are ways of storing values in a program. To declare a variable, you can use the `:=` operation, or the `var` keyword.

```clojure
var a = 0; declare variable 'a'
b := 0; declare variable 'b'
```

Both methods are exactly equivalent and can be used interchangeably, however it is common practice to use `var` when declaring globals, but use `:=` when declaring locals. This style is taken from the go language.

To set a new variable you can use the `=` operator

```clojure
b := 1
b = 2; is now 2
```

However, if you try and do this

```clojure
b = 2
```

Without declaring `b`, you would get a compiler error.

Lets take a look at this program:

```clojure
var main = fn() {
  local := 3
}
var anotherfn = fn() {
  global ;use the global variable
  local ;use the local variable (would cause an error)
}
var global = 0
```

In this example, a compiler error would be thrown because the `local` variable has not been declared in that scope. The solution would be to declare `local` in that scope, or make `local` a global variable.
