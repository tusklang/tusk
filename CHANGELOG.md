# Omm Beta August 16, 2020

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
