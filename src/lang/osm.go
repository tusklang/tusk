package lang

import "fmt"
import "net/http"
import "unsafe"

// #cgo CFLAGS: -std=c99
// #include "bind.h"
// #include "../osm/osm_render_alloc.h"
import "C"

type OsmPath struct {
  Url          string
  RequestType  string
  Directory    string
}

var paths []OsmPath

//struct to hold a process
//this struct will be converted to a void*
//then when it is passed back to go it will re convert back into an OsmProc
type OsmProc struct {
  GoProc func(params unsafe.Pointer, cli_params unsafe.Pointer, vars unsafe.Pointer, this_vals unsafe.Pointer, dir unsafe.Pointer) *C.char
}

//export CallOSMProc
func CallOSMProc(fnPtr unsafe.Pointer, args unsafe.Pointer, cli_params unsafe.Pointer, vars unsafe.Pointer, this_vals unsafe.Pointer, dir unsafe.Pointer) *C.char { //func to call goprocs for OSM
  //osm go procs can only return char*

  fn := (*((*OsmProc)(fnPtr))).GoProc //convert the void* to a go function

  return fn(args, cli_params, vars, this_vals, dir) //call the function
}

func handler(res http.ResponseWriter, req *http.Request) {

  //get the url
  url := req.URL.Path

  for k, v := range paths {

    if v.Url == url {

      //map of all the goprocs to send to the c++ interpreter
      goprocesses := map[string](func(params unsafe.Pointer, cli_params unsafe.Pointer, vars unsafe.Pointer, this_vals unsafe.Pointer, dir unsafe.Pointer) *C.char){
        "render": func(params unsafe.Pointer, cli_params unsafe.Pointer, vars unsafe.Pointer, this_vals unsafe.Pointer, dir unsafe.Pointer) *C.char {
          res.Header().Set("Content-Type", "text/html")

          //call ombrBind to template the ombr (json) to html
          html := C.GoString(C.ombrBind(params, cli_params, vars, this_vals, dir))
          res.Write([]byte(html))
          return C.CString(html) //return the converted ombr (html)
        },
      }

      const size = 1 //if you add more goprocesses, increase this

      //alloc the goprocNames array
      goprocNamesArr := C.allocOSM_GoProcNames(C.size_t(size))
      goprocNamesArrPass := (*[100000]C.osmGoProcName)(unsafe.Pointer(goprocNamesArr))[:size:size]

      //alloc the goproc array
      goprocArr := C.allocOSM_GoProcs(C.size_t(size))
      goprocArrPass := (*[size]C.osmGoProc)(unsafe.Pointer(goprocArr))[:size:size]

      //declare the current index of the array
      index := 0

      for k, v := range goprocesses {
        goprocArrPass[index] = (C.osmGoProc)(unsafe.Pointer(&OsmProc{ v })) //set the current index to be the current value
        goprocNamesArrPass[index] = (C.osmGoProcName)(unsafe.Pointer(C.CString(k)))

        index++ //increment the index
      }

      C.bindOsm(C.int(k), /* pass the url of the request */ C.CString(url), goprocArr, goprocNamesArr, C.int(size))

      //free the goproc and goprocName arrays
      C.freeOSM_GoProcs(goprocArr)
      C.freeOSM_GoProcNames(goprocNamesArr)
      break
    }

  }
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
func NewPath(pathC, requestTypeC, dirC *C.char) {

  path, requestType, dir := C.GoString(pathC), C.GoString(requestTypeC), C.GoString(dirC)

  paths = append(paths, OsmPath{ path, requestType, dir })

}
