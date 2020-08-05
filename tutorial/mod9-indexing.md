# Indexing

Omm has an operator to index types (::). We can use this operator to index strings, arrays, and hashes.

```clojure
;example for string
v1 := "test"
log v1::0 ;logs 't' (as a rune)
;;;;;;;;;;;;;;;;;;;

;example for array
v2 := [1, 2, 3]
log v2::2 ;logs 3
;;;;;;;;;;;;;;;;;;

;example for hash
v3 := [:
  key1 = "hello",
  "key2" = "world"
:]
log v3::key1 ;logs "hello"
;;;;;;;;;;;;;;;;;
```

To index using a variable, you can wrap the variable in parenthesis.

```clojure
testindex := "another_index"
testvalue := [:
  testindex = "hello",
  another_index = "world",
:]
log testvalue::(testindex) ;logs "world"
log testvalue::testindex ;logs "hello"
```
