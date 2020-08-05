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
