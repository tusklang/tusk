# First Steps

Every Omm program must have a main function. This style is shown in languages like, C, Go, Java, Rust, and many more. The main function looks like:

```clojure
var main = fn() {
  ; program goes here
}
```

Now, lets break apart this simple program.
First we can see: `var main: `. This declares a **global variable** in the program. See the variables module for more information about variables. Next we can see `fn() {}`. This denotes a function. Lastly, we can see `; program goes here`. Many languages use `//` as a comment, some use `#`. However, assembly languages, lisps, and omm use the semicolon to comment. You may be asking "What do we use to terminate a statement?" Omm does not require the developer to use terminators because they are automatically inserted. If you must insert a terminator, you can use the comma. Again, they are not required.
