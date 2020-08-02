# Control Flow

Control Flow is how you control the order of statements in your program. There are three ways to do this

- Conditionals
- While loops
- Each loops

Conditionals are if, elif, and else statements. An example would look like this

```clojure
if (true) log "abc"
elif (false) { ;'elif' is just shorthand for 'else if'
  log "def"
} else if (true) log "ghi"
else log "jkl"
```

While loops are just statements that execute *while* a statement is true

```clojure
i := 0
while (true) {
  i++
  log i
}
```

Each loops are range based

```clojure
iterator := [1, 2, 3]
each (iterator, k, v) {
  print k
  print " "
  print v
  log ""
}
; this program would print
; 0 1
; 1 2
; 2 3
```

To skip the loop, you can use the `continue,` statement, and to break it, you can use the `break,` statement.

```clojure
while (true) {
  if (false) break,
  log "hello"
  continue,
}
; you can also use continue and break in each loops
```
