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

  instance var initialize: fn(self, av, bv) {
    self::a: av
    self::b: bv
  }
}
```

Now let us make a namespace in this prototype

```
var test_proto: proto {
  instance var a
  instance var b
  instance var self

  instance var initialize: fn(self, av, bv) {
    self::self: self
    self::a: av
    self::b: bv
  }

  instance var somefunc: fn() {
    log self::a ;would log the value of the proto's `a`
    log a ;would also work
  }

  static var c: 12
  static var d: "hi there"
}
```

Now we can use this prototype!

```
var main: fn() {
  log test_proto::c ;would log 12
  obj := make[test_proto] ;create an object from the proto `test_proto`
  obj::initialize(obj, "value for a", "value for b") ;run the initialize function and pass the object itself, "value for a", and "value for b"
  obj::somefunc()
}
```
