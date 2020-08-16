# Including Other Files

Omm provides the `include` keyword to use other files

```clojure
;test1.omm
include "test2.omm"
```

```clojure
;test2.omm
var main = fn() {
  log "hello world"
}
```

In addition, Omm automatically puts header guards on the files (header guards just prevent the developer from making an include cycle).

If you want to build the entire directory, you must provide a `main.omm` file. The `main.omm` file will include all of the other files in the directory if you run
```bash
cd /your-directory/
omm *
```

You can also include oat archives inside of omm scripts.

```clojure
include "test.oat"
```

To include from the Omm standard library, you can use backticks
```clojure
include `lib.oat` ;include lib.oat from the standard library
```

If you install from the MSI installer, `lib.oat` will be an archive of the entire stdlib. If you instaled from the source, you can browse [the Omm repository](https://github.com/omm-lang/omm/tree/master/src/stdlib) to see what the Omm standard library has and maybe even contribute to it.
