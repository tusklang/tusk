# Installing Omm

Omm can be installed 2 ways:

- Using the installer
- Building from the source

The Omm installations can be found [here](http://omm.zone/downloads.php) (Currently only have downloads for windows).

---

If you do not have windows, or if you prefer building from source, you can go to the [Omm github page](https://github.com/Ankizle/Omm), clone the repo, and run this in the cli.
```
cd `Directory to Omm installation`/src
make
```
(Note that you must have [Go 1.14+](https://golang.org/doc/go1.14) to build from the source, and a version of GNU make)
Finally, you must add `'Directory to Omm installation'/src` to your `$PATH`.