package osm

import "net/http"
import "fmt"

//export CreateServer
func CreateServer(port string) {

  err := http.ListenAndServe("localhost:" + port, nil)

  if err != nil {
    fmt.Println("OSM warn\n", "Port", port, "is not a valid port! Defaulting to port 80")
    http.ListenAndServe("localhost:80", nil)
  }
}
