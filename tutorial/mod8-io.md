# Input and Output

Omm has two ways of outputting to the console

- `log`
- `print`

```clojure
print "hello"
log "world"
```

In the console, you would see

```
helloworld
```

The difference between `log` and `print` is that `log` adds a newline to the output, while `print` does not

To read text from the console, you can use the built in `input` function.

```clojure
input_value := input["hello world"]
```

And you would see

```
hello world: 
```

and you can input a value to the console
