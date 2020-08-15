# Functions

Functions are ways re-using blocks of code. This tutorial has already shown you a few examples of functions, but lets look at a more complex example

```clojure
var main: fn() {
  val := test[1, 2, 3]
}

var test: fn(number -> a, number -> b, number -> c) {
  return a + b + c
}
```

Now you may notice two things in Omm that are significantly different than in other languages. First, To call a function, we use `func_name[]`, instead of `func_name()`. The reason for this is because Omm inserts an operator between `func_name` and `[]` during compilation. This operation is `<-` which means synchronous call. The above example would be translated from `test[1, 2, 3]` to `test <- [1, 2, 3]`. However, this change happens automatically during compilation. Omm also has asynchronous functions. To call an function asynchronously, you can use the <~ operator as mentioned before.

```clojure
test <~ [1, 2, 3]
; instead of
test <- [1, 2, 3]
```

The `<~` operator returns a thread type. If you want to get the value of a thread type, you can use the `await` keyword

```clojure
value := await test <~ [1, 2, 3]
; or you can do
val1 := test <~ [1, 2, 3]
; do stuff
val2 := await val1
```

But, the omm compiler is smart, and can automatically replace the [] with (), so writing

```clojure
func_name()
;instead of
func_name[]
```

Will also work (even inserting the <- operator automatically).

Second, parameter lists in omm look like
```
(string -> a, string -> b)
```
Whereas in javascript, and python they look more like
```
(a, b)
```
The key difference is that omm names a type before an argument. This is because omm allows something called overloading. If you want an argument to have any type you can use `any ->`.
```clojure
fn(any -> a, any -> b) ;accepts any type for a and b
```

But again, the omm compiler can automatically assume the type:

```clojure
fn(a, b) ;automatically gets converted to fn(any -> a, any -> b)
```

Now lets look at the previous example, but with wrong types

```clojure
var main = fn() {
  ;pass bool, string, and array
  val := test[true, "hi", [1, 2, 3]]
}

var test: fn(number -> a, number -> b, number -> c) {
  return a + b + c
}
```

This would cause an error, where the function with that type list does not exist. Now let's overload this function!

To overload, you can use the built in `ovld` keyword.

```clojure
var main = fn() {
  ;pass bool, string, and array
  val := test[true, "hi", [1, 2, 3]]
}

var test = fn(number -> a, number -> b, number -> c) {
  return a + b + c
}

;create an overloaded `test` with the type list of bool, string, and array
ovld test = fn(bool -> a, string -> b, array -> c) {

  log a
  log b
  log c

  return false
}
```
