package lang

import "fmt"
import "net/http"

// #cgo CFLAGS: -std=c99
// #include "bind.h"
import "C"

type OsmPath struct {
  Url          string
  RequestType  string
}

var paths []OsmPath

func handler(res http.ResponseWriter, req *http.Request) {

  //get the url
  url := req.URL.Path

  for k, v := range paths {

    if v.Url == url {
      C.bindOsm(C.int(k), C.CString(url))
    }

  }

  fmt.Fprintln(res, "hello")
}

//export OSM_create_server
func OSM_create_server(portC *C.char) {

  //get the port as a go string
  port := C.GoString(portC)

  http.HandleFunc("/", handler)

  //create the server
  err := http.ListenAndServe("localhost:" + port, nil)

  //f the server is invalid
  if err != nil {
    fmt.Println("Port", port, "is invalid. Defaulting to port 80")
    http.ListenAndServe("localhost:80", nil)
  }

}

//export NewPath
func NewPath(pathC, requestTypeC *C.char) {

  path, requestType := C.GoString(pathC), C.GoString(requestTypeC)

  paths = append(paths, OsmPath{ path, requestType })

}
