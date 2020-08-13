package interpreter

//osm is omm's server manager

import "strconv"
import "net/http"
import "time"
import . "lang/types"

func OSM_start(args []*OmmType, stacktrace []string, line uint64, file string, instance *Instance) *OmmType {

  if len(args) != 1 {
    OmmPanic("Function osm.start requires a parameter count of 1", line, file, stacktrace)
  }

  var port = (*cast(*args[0], "string", stacktrace, line, file)).(OmmString).ToGoType()

  http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) { //whenever a request is made, send the corresponding html file in the "public" folder
    http.ServeFile(res, req, "public/" + req.URL.Path[1:])
  })

  if err := http.ListenAndServe(":" + port, nil); err != nil {
    OmmPanic("Cannot use port " + port, line, file, stacktrace)
  }

  //return undef
  var tmpundef OmmType = undef
  return &tmpundef
}

func OSM_handle(args []*OmmType, stacktrace []string, line uint64, file string, instance *Instance) *OmmType {

  if len(args) != 2 {
    OmmPanic("Function osm.handle requires a parameter count of 2", line, file, stacktrace)
  }

  if (*args[0]).Type() != "string" || (*args[1]).Type() != "function" {
    OmmPanic("Function osm.handle requires (string, function)", line, file, stacktrace)
  }

  var path = (*args[0]).(OmmString).ToGoType()
  var cb = (*args[1]).(OmmFunc)

  if len(cb.Overloads[0].Params) != 2 {
    OmmPanic("Function osm.handle requires a callback with 2 parameters", line, file, stacktrace)
  }

  http.HandleFunc(path, func(res http.ResponseWriter, req *http.Request) {
    var (oreq, ores OmmHash) //make request and response hashes (hashes because everything is instance)

    oreq.Hash = make(map[string]*OmmType)
    ores.Hash = make(map[string]*OmmType)

    //in the response hash, put these
    ores.Set("send", OmmGoFunc{
      Function: func(args []*OmmType, stacktrace []string, line uint64, file string, instance *Instance) *OmmType {
        if len(args) != 1 {
          OmmPanic("Expected 1 argument to send a text response", line, file, stacktrace)
        }

        var resp = (*cast(*args[0], "string", stacktrace, line, file)).(OmmString).ToGoType()
        res.Write([]byte(resp))

        //return undef
        var tmpundef OmmType = undef
        return &tmpundef
      },
    })
    ores.Set("sendfile", OmmGoFunc{
      Function: func(args []*OmmType, stacktrace []string, line uint64, file string, instance *Instance) *OmmType {
        if len(args) != 1 {
          OmmPanic("Expected 1 argument to send a html response", line, file, stacktrace)
        }

        var htmlf = (*cast(*args[0], "string", stacktrace, line, file)).(OmmString).ToGoType()
        http.ServeFile(res, req, "public/" + htmlf)

        //return undef
        var tmpundef OmmType = undef
        return &tmpundef
      },
    })
    ores.Set("setheader", OmmGoFunc{
      Function: func(args []*OmmType, stacktrace []string, line uint64, file string, instance *Instance) *OmmType {
        if len(args) != 2 {
          OmmPanic("Expected 2 arguments to set a header", line, file, stacktrace)
        }

        var htype = (*cast(*args[0], "string", stacktrace, line, file)).(OmmString).ToGoType()
        var val = (*cast(*args[1], "string", stacktrace, line, file)).(OmmString).ToGoType()
        res.Header().Set(htype, val)

        //return undef
        var tmpundef OmmType = undef
        return &tmpundef
      },
    })
    ores.Set("cookie", OmmGoFunc{
      Function: func(args []*OmmType, stacktrace []string, line uint64, file string, instance *Instance) *OmmType {

        if len(args) != 2 {
          OmmPanic("Expected 2 arguments to set a cookie", line, file, stacktrace)
        }

        var name = (*args[0]).(OmmString).ToGoType()
        var opts = (*args[1]).(OmmHash)
        var cookie http.Cookie
        cookie.Name = name

        if opts.Exists("value") && (*opts.At("value")).Type() == "string" {
          cookie.Value = (*opts.At("value")).(OmmString).ToGoType()
        }
        if opts.Exists("path") && (*opts.At("path")).Type() == "string" {
          cookie.Path = (*opts.At("path")).(OmmString).ToGoType()
        }
        if opts.Exists("domain") && (*opts.At("domain")).Type() == "string" {
          cookie.Domain = (*opts.At("domain")).(OmmString).ToGoType()
        }
        if opts.Exists("expires") && (*opts.At("expires")).Type() == "number" {
          cookie.Expires = time.Now().Add(time.Duration((*opts.At("expires")).(OmmNumber).ToGoType()) * time.Millisecond)
        }
        if opts.Exists("maxage") && (*opts.At("maxage")).Type() == "number" {
          cookie.MaxAge = int((*opts.At("maxage")).(OmmNumber).ToGoType())
        }
        if opts.Exists("secure") && (*opts.At("secure")).Type() == "bool" {
          cookie.Secure = (*opts.At("secure")).(OmmBool).ToGoType()
        }
        if opts.Exists("httponly") && (*opts.At("httponly")).Type() == "bool" {
          cookie.HttpOnly = (*opts.At("httponly")).(OmmBool).ToGoType()
        }
        if opts.Exists("samesite") && (*opts.At("samesite")).Type() == "string" {

          var ss = 1

          switch (*opts.At("samesite")).(OmmString).ToGoType() {
            case "lax":
              ss = 2
            case "strict":
              ss = 3
            case "none":
              ss = 4
          }

          cookie.SameSite = http.SameSite(ss)
        }

        http.SetCookie(res, &cookie)

        var ommtype OmmType = undef
        return &ommtype
      },
    })
    /////////////////////////////////

    //in the request hash put this
    oreq.Set("cookies", OmmGoFunc{
      Function: func(args []*OmmType, stacktrace []string, line uint64, file string, instance *Instance) *OmmType {
        var cookies = req.Cookies()
        var ommhash OmmHash
        ommhash.Hash = make(map[string]*OmmType)

        for _, v := range cookies { //map each cookie to it's name
          var cookie OmmHash

          var value OmmString
          value.FromGoType(v.Value)
          cookie.Set("value", value)

          var ommtype OmmType = value
          ommhash.Hash[v.Name] = &ommtype
        }

        var ommtype OmmType = ommhash
        return &ommtype
      },
    })
    //////////////////////////////

    var (
      ommtype_req OmmType = oreq
      ommtype_res OmmType = ores
    )

    instance.Allocate(cb.Overloads[0].Params[0], &ommtype_req)
    instance.Allocate(cb.Overloads[0].Params[1], &ommtype_res)
    Interpreter(instance, cb.Overloads[0].Body, append(stacktrace, "osm handler callback at line " + strconv.FormatUint(line, 10) + " in file " + file))
  })

  //return undef
  var tmpundef OmmType = undef
  return &tmpundef
}
