package native

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/omm-lang/omm/lang/types"
	. "github.com/omm-lang/omm/lang/types"
)

//OmmURLResp is a url response
type OmmURLResp struct {
	Status     OmmString
	StatusCode OmmNumber
	Error      OmmString
	Protocol   OmmString
	Header     OmmHash
	Body       OmmString
}

func (r *OmmURLResp) FromGoType(res *http.Response, e error) {

	if res == nil { //do not proceed if res is nil (just in case)
		return
	}

	if e == nil { //only proceed if error is nil
		var rstatus OmmString
		rstatus.FromGoType(res.Status)
		var rstatuscode OmmNumber
		rstatuscode.FromGoType(float64(res.StatusCode))
		var proto OmmString
		proto.FromGoType(res.Proto)

		head := res.Header
		var ommhead OmmHash

		for k, v := range head {
			var ommval OmmArray

			for _, vv := range v {
				var ommsubval OmmString
				ommsubval.FromGoType(vv)
				ommval.PushBack(ommsubval)
			}

			ommhead.Set(k, ommval)
		}

		body, _ := ioutil.ReadAll(res.Body)
		var ommbody OmmString
		ommbody.FromGoType(string(body))
		res.Body.Close()

		r.Status = rstatus
		r.StatusCode = rstatuscode
		r.Protocol = proto
		r.Header = ommhead
		r.Body = ommbody
	} else {
		var err OmmString
		err.FromGoType(e.Error())
		r.Error = err
	}
}

func (r OmmURLResp) Format() string {
	return "{ http response }"
}

func (r OmmURLResp) Type() string {
	return "http_response"
}

func (r OmmURLResp) TypeOf() string {
	return r.Type()
}

func (r OmmURLResp) Deallocate() {}

//Range ranges over a url response
func (r OmmURLResp) Range(fn func(val1, val2 *types.OmmType) types.Returner) *types.Returner {
	return nil
}

func urlget(args []*OmmType, stacktrace []string, line uint64, file string, instance *Instance) *OmmType {

	if len(args) != 1 || (*args[0]).Type() != "string" {
		OmmPanic("Function url.get requires an argument count of 1 which must be of type string", line, file, stacktrace)
	}

	res, e := http.Get((*args[0]).(OmmString).ToGoType())

	var final OmmURLResp
	final.FromGoType(res, e)
	var ommtype OmmType = final
	return &ommtype
}

var urlrequest = func(args []*OmmType, stacktrace []string, line uint64, file string, instance *Instance) *OmmType {

	var res *http.Response
	var e error

	/*
		Overloads are

		(string) => get request
		(string, string) => post with json data
		(string, string, hash) => post with data and header
	*/
	if len(args) == 1 && (*args[0]).Type() == "string" {
		return urlget(args, stacktrace, line, file, instance)
	} else if len(args) == 2 && (*args[0]).Type() == "string" && (*args[1]).Type() == "string" {
		res, e = http.Post((*args[0]).(OmmString).ToGoType(), "application/json", bytes.NewBuffer([]byte((*args[1]).(OmmString).ToGoType())))
	} else if len(args) == 3 && (*args[0]).Type() == "string" && (*args[1]).Type() == "string" && (*args[2]).Type() == "hash" {

		var req *http.Request
		req, e = http.NewRequest("post", (*args[0]).(OmmString).ToGoType(), bytes.NewBuffer([]byte((*args[1]).(OmmString).ToGoType())))

		if e == nil { //only if there is no error, continue
			for k, v := range (*args[2]).(OmmHash).Hash {

				if (*v).Type() != "string" {
					OmmPanic("Given an invalid header to url.post", line, file, stacktrace)
				}

				headerv := (*v).(OmmString).ToGoType()
				req.Header.Set(k, headerv)
			}

			var client = &http.Client{} //create a new client
			res, e = client.Do(req)
		}
	} else {
		OmmPanic("Function url.post requires an argument count of 1 or 2 with the types of string and (optional) hash", line, file, stacktrace)
	}

	var final OmmURLResp
	final.FromGoType(res, e)
	var ommtype OmmType = final
	return &ommtype

}
