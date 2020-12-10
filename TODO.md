# High Priority:
- Allow for static methods to use private instance methods
  ```
  	var p = proto {
		var _t
	  
	  	static var new = fn() {
		  	var o = make:(p)
		  	o::_t = 1
		  	return o
		}
    }
  ```
  And allow static methods to call other private static methods, allow private instance methods of an object to access private instance variables in another object, and so on.

# Low Priority:
- Create date class
- Create a socket class, which provides higher level functions with sockets
- Add trig. functions