# Omm Beta August 16, 2020

## Beta 1.1.4 - August 29, 2020

Fixed oat powershell and shell scripts. They stopped working because of the last update where the `omml.exe` and `omml` binaries were replaced with `omm.exe` and `omm`. The oat ps1 and sh scripts were still calling the non-existant `omml.exe` and `omml` binaries.

## Beta 1.1.3 - August 28, 2020

- Fixed the string concatenation operator
- Removed `omm.ps1` and `omm.sh` in favor for `omm.exe` and `omm` binaries.
- Fixed including from the standard library 

## Beta 1.1.2 - August 26, 2020

This beta release targets the bugs in the oat format. The previous oat format given in b1.1.1 has now been deprecated, and a new, more robust, oat version has been implemented. Previously, functions would not be properly parsed in oat. In addition, prototype access lists would be ignored by oat. A few more bugs have been fixed with the decoder as well, but these are just bugs where the decoder would not properly read strings in a few fields.

## Beta 1.1.1 - August 21, 2020

- When a function was called, the arguments would not be garbage collected. This causes memory leaks when arguments are passed.
- Languages that support closures allow something like this:
  ```js
  function test() {
    var a = 1;
    return () => {
      return a; //a is allowed to be used in this lambda
    }
  }
  ```
  Omm previously allowed that, but a bug with overloading may have undid that. This update patches that bug. 

## Beta 1.1.0 - August 20, 2020

- Added protected fields inside prototypes
- Allowed for functions to be overloaded inside in prototypes
  Previously, if you tried to overload a function inside a prototype, it would say that you must declare that function outside of the prototype, then use it. Now, you do all of the work inside the prototype. 

## Beta 1.0.1 - August 19, 2020

- Fixed a garbage collection bug
  Previously, garbage collection would only occur when `return` was called. This causes memory leaks when `return` was not used

- Made threads rely on c threads (pthreads/win32 threads) instead of goroutines
  Goroutines have less versatility than c threads. 

## Beta 1.0.0 - August 16, 2020

- Initial Beta Release
