# Functions

Functions are ways re-using blocks of code. This tutorial has already shown you a few examples of functions, but lets look at a more complex example

```
var main: fn() {
  val := test[1, 2, 3]
}

var test: fn(a, b, c) {
  return a + b + c
}
```

Now you may notice one thing in Omm that is significantly different than in other languages. To call a function, we use `func_name[]`, instead of `func_name()`. The reason for this is because omm inserts an operator between `func_name` and `[]` during compilation. This operation is `<-` which means synchronous call. The above example would be translated from `test[1, 2, 3]` to `test <- [1, 2, 3]`. However, this change happens automatically during compilation. Omm also has asynchronous functions. To call an function asynchronously, you can use the <~ operator as mentioned before.

```
test <~ [1, 2, 3]
; instead of
test <- [1, 2, 3]
```

The `<~` operator returns a thread type. If you want to get the value of a thread type, you can use the `await` keyword

```
value := await test <~ [1, 2, 3]
; or you can do
val1 := test <~ [1, 2, 3]
; do stuff
val2 := await val1
```
