# Installing Omm

Omm can be installed 2 ways:

- Using the installer
- Building from the source

The Omm installations can be found [here](http://omm.zone/downloads.php) (Currently only have downloads for windows).

---

If you do not have windows, or if you prefer building from source, you can go to the [Omm github page](https://github.com/Ankizle/Omm), clone the repo, and run this in the cli.
```
cd src/
make
```
(Note that you must have [Go 1.14+](https://golang.org/doc/go1.14), and [GCC 64-bit 8.1+](http://mingw-w64.org/doku.php) to build from the source, and a version of GNU make)
Finally, you must add `'Directory to Omm installation'/src` to your `$PATH`.

Once you are able to install omm, you can use the `omm` command in the terminal to run omm scripts. Create a `test.omm` file and run `omm test.omm`.
