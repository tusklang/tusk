# Datatypes

The Omm languages comes default with 5 datatypes:

- Number: a number is a numeric value (floating point or integer) and can be written by just writing the number, e.g. 19, 12.2, 98.0
- Rune: a rune is Omm's version of a character, and can be written by using single quotes, e.g. 'a', 'b', 'c'
- String: a string is a series of runes, and can be written using either double quotes or single quotes, e.g. "hello world", \`hello world\`
- Boolean: a boolean is either `true` or `false` and can be written by just writing true or false, e.g. true, false
- Undefined: an undefined value is a value with no value, and can be written by just writing `undef`
- Array: an array is a list of values and can be written by writing [] with the values inside, e.g. [], [1, 2, 3], ["hello", 123, true, undef], [1 2 3]
- Hash: a hash is a list of values (like an array) but with key/value pairs, and can be written like: [: "a" = 1, "b" = 2, 3 = "c" :]
- Prototype: a prototype is Omm's for OOP, and you can learn about them in the prototypes module
- Object: an object is an instance of a prototype, and you can learn about them in the prototypes module
- Thread: a thread is the way to do asynchronous programming in Omm based on goroutines, and you can make them using the asynchronous operator (<~)
- Function: a function is used to store re-usable blocks of code, and you can learn more about them in the functions module
- Integer: 64 bit int
- Float: 64 bit float

The `typeof` function can be used to get the type of a value
To convert between types, you can use the `->` operator

```clojure
var converted = string -> 123 ;convert 123 to "123"
```
