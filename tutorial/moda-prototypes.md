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

  instance var somefunc: fn(self) {
    log self::a ;this would print a from the instance
  }
}
```

Now let us make a namespace in this prototype

```
var test_proto: proto {
  instance var a
  instance var b

  instance var initialize: fn(self, av, bv) {
    self::a: av
    self::b: bv
  }

  instance var somefunc: fn(self) {
    log self::a ;this would print a from the instance
  }

  static var c: 13
  static var d: "hello world"

  static var some_namespace_func: fn() {
    log "Coming from the namespace"
  }
}
```

Now we can use this prototype!

```clojure
var main: fn() {
  log test_proto::c ;would log 12
  obj := make[test_proto] ;create an object from the proto `test_proto`
  obj::initialize(obj, "value for a", "value for b") ;run the initialize function
  obj::somefunc(obj) ;would print the object's `a` variable
  obj::_private() ;would cause an error
}
```

In Omm, it is important to pass the `self` parameter to the instance to access the other variables in the instance. 