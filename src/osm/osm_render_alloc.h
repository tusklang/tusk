#ifndef OMM_OSM_RENDER_ALLOC_H_
#define OMM_OSM_RENDER_ALLOC_H_

#ifdef __cplusplus
extern "C" {
#endif
//

  #include <stdlib.h>

  typedef void* osmGoProc;
  typedef char* osmGoProcName;

  //function to alloc osm go procs
  static inline osmGoProc* allocOSM_GoProcs(size_t size) {

    return (osmGoProc*)(malloc(size * sizeof(osmGoProc)));
  }

  //function to alloc osm go proc names
  static inline osmGoProcName* allocOSM_GoProcNames(size_t size) {

    return (osmGoProcName*)(malloc(size * sizeof(osmGoProcName)));
  }

  //function to free a void** (osm go proc)
  static inline void freeOSM_GoProcs(osmGoProc* val) {
    free(val);
  }

  //function to free a char** (osm go proc name)
  static inline void freeOSM_GoProcNames(osmGoProcName* val) {
    free(val);
  }

#ifdef __cplusplus
}
#endif

#endif
