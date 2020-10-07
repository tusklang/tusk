# The Tusk Language

Tusk is a dynamicly typed programming language, which also supports many strictly typed features.

# Features

## Optional Strict Typing

Probably the most prominent feature is Tusk's optional strict typing. Tusk allows the developer to specify a function prototype argument types. This has no effect on compile-time error checking, but allowers for function overloading. This is already seen in Java, but JavaScript and Python both are famously known to not support function overloading. For JavaScript, overloading is very messy:

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

;if f is called and the prototype is not (string, string, string) or (string, string) the program panics
```

This can be enourmously useful for both efficiency, and code readability.

## Private, Public, and Protected Fields

Maybe as useful as optional static typing is Tusk's support for private, public, and protected fields. Tusk, like many modern languages, has object oriented features. Classes in Tusk are defined with the `proto` keyword. In JavaScript, to crrate private variables, you must use something called a closure.

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

##### Protected Fields

In Java, a protected field is a field in a class that can only be accessed from nested classes, or from classes within the same package. Tusk has a similar idea, but Tusk does not have a package structure. Instead, you can use the `access` keyword to specify which files can use the given field. This allows for the developer to specify access, even though Tusk is a dynamic language.

```clojure
var a = proto {
    ;thisf is a macro for the current file
    access ("thisf", "anotherfile.tusk")
    var protectedfield = fn() {
        ;this method can only be accessed in the current file `thisf` and anotherfile.tusk
    }
}
```
