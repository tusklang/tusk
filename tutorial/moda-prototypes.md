# Prototypes and Objects

Prototypes are classes in Omm. Prototypes allow programs to create new datatypes. Prototypes are made up of instance and static variables. Instance variables are stored in an *object* of that prototype, while the static is stored in the namespace of that prototype.

We can make a new prototype like this

```
var test_proto: proto {
  ; TODO: make this prototype
}
```

Now let us fill this prototype with an instance

```
var test_proto: proto {
  instance var a
  instance var b

  instance var initialize: fn(av, bv) {
    a = av
    b = bv
  }

  instance var somefunc: fn() {
    log a ;this would print a from the instance
  }
}
```

If you come from object oriented statically typed languages, like Java, C++, Go, C#, etc, then you will know about public and private variables. Omm, despite being dynammically typed, also has this feature. To create a private variable in a prototype, you can prefix the variable name with an uderscore.

```
var test_proto = proto {
  instance var a
  instance var b

  instance var initialize = fn(av, bv) {
    a: av
    b: bv
  }

  instance var somefunc = fn() {
    log a ;this would print a from the instance
  }

  instance var _private = fn() { ;this method is private
    log "Coming from a private method"
  }
}
```

Now let us make a namespace in this prototype

```
var test_proto = proto {
  instance var a
  instance var b

  instance var initialize = fn(av, bv) {
    a: av
    b: bv
  }

  instance var somefunc = fn() {
    log a ;this would print a from the instance
  }

  instance var _private = fn() { ;this method is private
    log "Coming from a private method"
  }

  static var c = 13
  static var d = "hello world"

  static var some_namespace_func = fn() {
    log "Coming from the namespace"
  }
}
```

Now we can use this prototype!

```clojure
var main = fn() {
  log test_proto::c ;would log 12
  obj := make[test_proto] ;create an object from the proto `test_proto`
  obj::initialize("value for a", "value for b") ;run the initialize function
  obj::somefunc() ;would print the object's `a` variable
  obj::_private() ;would cause an error
}
```