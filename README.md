# The Tusk Language

Tusk is a dynamicly typed programming language, which also supports many strictly typed features.

# Features

## Optional Strict Typing

Probably the most prominent feature is Tusk's optional strict typing. Tusk allows the developer to specify a function prototype argument types. This has no effect on compile-time error checking, but allowers for function overloading. This is already seen in Java, but JavaScript and Python both are known to not support function type overloading. For JavaScript, overloading is very messy:

```javascript
function f(a, b, c) {
    if (typeof(a) == "string" && typeof(b) == "string" && typeof(c) == "string") {

    } else if (typeof(a) == "string" && typeof(b) == "string" && c == undefined) {

    } //etc...
}
```

Python is a bit better, but type checking is still very messy

```python
def f(a, b, c):
    if isinstance(a, str) and isinstance(b, str) and isinstance(c, str):
        pass
    else:
        pass #error

def f(a, b):
    if isinstance(a, str) and isinstance(b, str):
        pass
    else:
        pass #error

```

Tusk subverts this problem by having an `ovld` keyword.

```clojure
var f = fn(string -> a, string -> b, string -> c) {
    
}

ovld f = fn(string -> a, string -> b) {

}

;if f is called and the signature is not (string, string, string) or (string, string) the program panics
```

This can be enourmously useful for code readability.

## Private and Public Fields

Maybe as useful as optional static typing is Tusk's support for private, public, and protected fields. Tusk, like many modern languages, has object oriented features. Classes in Tusk are defined with the `proto` keyword. In JavaScript, to create private variables, you must use something called a closure.

```javascript
function a() {
    var b = 5;
    return function() {
        //you can do things with b here
        return b;
    }
}
```

In Python, you can prefix your methods with a `_` to "mark them as private", but they are easily accessible from outside the class. Tusk has **real** private fields. Like in Python, you must prefix your field name with a `_` to mark it as private, but it actually works. 

```clojure
var a = proto {
    var publicmethod = fn() {
        ;public method
    }

    var _privatemethod = fn() {
        ;private method (will throw an error if it is accessed from outside this prototype)
    }
}
```
