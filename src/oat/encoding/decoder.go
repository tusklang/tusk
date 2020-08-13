package oatenc

import (
	"errors"
	"bufio"
	"os"
	"regexp"
	"strings"
	"strconv"
	"io"
	"reflect"
	. "lang/types"
)

var NOT_OAT error

func OatDecode(filename string, mode int) (map[string][]Action, error) {

	NOT_OAT = errors.New("Given file is not an oat: " + filename) 

	file, e := os.Open(filename)

	if e != nil {
		return nil, errors.New("Could not find file: " + filename)
	}

	reader := bufio.NewReader(file)

	var vers = "" //store the omm version
	var data []rune
	var cur bool = true

	for {

		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, e
			}
		} else {

			if cur && c == '\n' { //on newline, switch to data
				cur = false
				continue
			}

			if cur {
				vers+=string(c)
			} else {
				data = append(data, c)
			}

		}
	}

	if !strings.HasPrefix(vers, MAGIC) {
		return nil, NOT_OAT
	}

	vers = strings.TrimPrefix(vers, MAGIC)

	m, _ := regexp.MatchString("\\d+\\.\\d+\\.\\d+", vers)

	if !m {
		return nil, NOT_OAT
	}

	version_spl := strings.Split(vers, ".")

	majorv, _ := strconv.Atoi(version_spl[0])
	minorv, _ := strconv.Atoi(version_spl[1])
	bugv, _ := strconv.Atoi(version_spl[2])

	if !(majorv >= OMM_MAJOR && minorv >= OMM_MINOR && bugv >= OMM_BUG) {
		return nil, errors.New("Please upgrade your omm version to " + vers + " in order use this oat")
	}

	var encodedVars = make([][2][]rune, 1)

	escaped := false
	isOnName := true

	for _, v := range data {

		if escaped {
			escaped = false
			goto escaped
		}

		if v == reserved["escaper"] {
			escaped = true
			goto escaped
		}

		if v == reserved["set global"] {
			isOnName = false
			continue
		} else if v == reserved["new global"] {
			encodedVars = append(encodedVars, [2][]rune{})
			isOnName = true
			continue
		}

		escaped:
		if isOnName {
			encodedVars[len(encodedVars) - 1][0] = append(encodedVars[len(encodedVars) - 1][0], v)
		} else {
			encodedVars[len(encodedVars) - 1][1] = append(encodedVars[len(encodedVars) - 1][1], v)
		}
	}

	//remove the trailing
	encodedVars = encodedVars[:len(encodedVars) - 1]

	var decodedvars = make(map[string][]Action)

	for _, v := range encodedVars {
		decodedV, e := DecodeVariable(v[1])

		if e != nil {
			return nil, e
		}

		decodedvars[string(DecodeStr(v[0]))] = decodedV
	}

	if mode == 0 {
		return decodedvars, nil
	} else {
		var all []Action
		for k, v := range decodedvars {
			all = append(all, Action{
				Type: "var",
				Name: k,
				ExpAct: v,
			})
		}
		return map[string][]Action{
			"all": all,
		}, nil
	}

}

func DecodeVariable(encoded []rune) ([]Action, error) {
	var (
		decoded = make([]Action, 1)
		curpos  = 0
		escaped = false
		curval  = make([]rune, 0)
	)

	for i := 0; i < len(encoded); i++ {

		curact := &decoded[len(decoded) - 1]
		s := reflect.ValueOf(curact)

		if curpos >= s.Elem().NumField() {
			return nil, NOT_OAT
		}

		field := s.Elem().Field(curpos)

		if !field.IsValid() || !field.CanSet() {
			return nil, NOT_OAT
		}

		var matchers = map[string]int{}

		for ;i < len(encoded); i++ {

			if escaped {
				curval = append(curval, encoded[i])
				escaped = false
				continue
			}

			if encoded[i] == reserved["escaper"] {
				escaped = true
				curval = append(curval, encoded[i])
				continue
			}

			if strings.HasPrefix(getReservedFromRune(encoded[i]), "start ") {
				matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "start ")]++
			} else if strings.HasPrefix(getReservedFromRune(encoded[i]), "end ") {

				if matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "end ")] <= 0 {
					return nil, NOT_OAT
				}

				matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "end ")]--
			}

			for _, v := range matchers {
				if v != 0 {
					goto skipchecks
				}
			}

			if encoded[i] == reserved["seperate field"] {
				goto nextfield
			}

			skipchecks:
			curval = append(curval, encoded[i])
		}

		return nil, NOT_OAT

		nextfield:

		switch s.Elem().Type().Field(curpos).Name {
			case "Type":

				if len(curval) != 1 {
					return nil, NOT_OAT
				}

				(*curact).Type = getReservedFromRune(curval[0])
				
			case "Name":
				(*curact).Name = string(DecodeStr(curval))

			case "Value":

				var decval func(cv []rune) (OmmType, error)
				decval = func(cv []rune) (OmmType, error) {

					if len(cv) == 0 {
						return nil, NOT_OAT
					}

					switch cv[0] {
						case reserved["make c-array"]:
							cv = cv[1:]

							escaped := false
							var arr = make([][]rune, 1)

							for _, v := range cv {

								if escaped {
									escaped = false
									goto escape
								}

								if v == reserved["escaper"] {
									escaped = true
									goto escape
								}

								if v == reserved["value seperator"] {
									arr = append(arr, []rune{})
									continue
								}

								escape:
								arr[len(arr) - 1] = append(arr[len(arr) - 1], v)
							}

							//remove the trailing
							arr = arr[:len(arr) - 1]

							var oarr []*OmmType
							
							for _, v := range arr {
								val, e := decval(v)
								if e != nil {
									return nil, e
								}
								oarr = append(oarr, &val)
							}

							return OmmArray{
								Array: oarr,
								Length: uint64(len(oarr)),
							}, nil

						case reserved["make bool"]:
							if len(cv) != 3 {
								return nil, NOT_OAT
							}

							cv = cv[2:]

							var ommbool bool

							if cv[0] == 1 {
								ommbool = true
							} else if cv[0] == 0 {
								ommbool = false
							} else {
								return nil, NOT_OAT
							}

							var boolean OmmType = OmmBool{
								Boolean: &ommbool,
							}
							return boolean, nil

						case reserved["make c-hash"]:
							cv = cv[1:]

							var hash = make([][2][]rune, 1)
							var cur = true
							var escaped = false

							for _, v := range cv {
								
								if escaped {
									escaped = false
									goto escp
								}

								if v == reserved["escaper"] {
									escaped = true
									goto escp
								}

								if v == reserved["hash key seperator"] {
									cur = false
									continue
								}

								if v == reserved["value seperator"] {
									cur = true
									continue
								}

								escp:
								if cur {
									hash[len(hash) - 1][0] = append(hash[len(hash) - 1][0], v)
								} else {
									hash[len(hash) - 1][1] = append(hash[len(hash) - 1][1], v)
								}
							}

							//remove the trailing
							hash = hash[:len(hash) - 1]

							var ohash map[string]*OmmType

							for _, v := range hash {
								cdecoded, e := decval(v[1])
								if e != nil {
									return nil, e
								}
								ohash[string(DecodeStr(v[0]))] = &cdecoded
							}

							return OmmHash{
								Hash: ohash,
								Length: uint64(len(ohash)),
							}, nil

						case reserved["start number"]:
							if len(cv) < 2 {
								return nil, NOT_OAT
							}

							cv = cv[1:len(cv) - 1]

							var inte []int64
							var deci []int64
							var cur = true
							var escaped = false

							for _, v := range cv {
								if v == reserved["escaper"] {
									escaped = true
									continue
								}

								if escaped {
									escaped = false
									goto esc
								}

								if v == reserved["decimal place"] {
									cur = false
									continue
								}

								esc:
								if cur {
									inte = append(inte, int64(v))
								} else {
									deci = append(deci, int64(v))
								}
							}

							return OmmNumber{
								Integer: &inte,
								Decimal: &deci,
							}, nil

						case reserved["start proto"]:

							if len(cv) < 3 {
								return nil, NOT_OAT
							}

							cv = cv[1:len(cv) - 1]

							if cv[0] != reserved["start proto name"] {
								return nil, NOT_OAT
							}

							cv = cv[1:]

							var escaped = false
							var name string
							var k int
							var v rune

							var matchers = make(map[string]int)

							for k, v = range cv {

								if escaped {
									escaped = false
									goto protoname_esc
								}

								if v == reserved["escaper"] {
									escaped = true
									continue
								}

								if strings.HasPrefix(getReservedFromRune(v), "start ") {
									matchers[strings.TrimPrefix(getReservedFromRune(v), "start ")]++
								} else if strings.HasPrefix(getReservedFromRune(v), "end ") {
					
									if matchers[strings.TrimPrefix(getReservedFromRune(v), "end ")] <= 0 {

										if v == reserved["end proto name"] {
											k++
											break
										}
										
										return nil, NOT_OAT
									}
					
									matchers[strings.TrimPrefix(getReservedFromRune(v), "end ")]--
								}
					
								for _, v := range matchers {
									if v != 0 {
										goto protoname_esc
									}
								}
								
								protoname_esc:
								name+=string(v)
							}

							cv = cv[k:]

							if len(cv) < 4 || cv[0] != reserved["start proto static"] || cv[len(cv) - 1] != reserved["end proto instance"] {
								return nil, NOT_OAT
							}

							cv = cv[1:len(cv) - 1]

							matchers = make(map[string]int)

							var staticparts = make([][2][]rune, 1)
							var instanceparts = make([][2][]rune, 1)
							var curp = true

							for k, v = range cv {

								if escaped {
									escaped = false
									goto protostatic_esc
								}

								if v == reserved["escaper"] {
									escaped = true
									goto protostatic_esc
								}

								if strings.HasPrefix(getReservedFromRune(v), "start ") {
									matchers[strings.TrimPrefix(getReservedFromRune(v), "start ")]++
								} else if strings.HasPrefix(getReservedFromRune(v), "end ") {
					
									if matchers[strings.TrimPrefix(getReservedFromRune(v), "end ")] <= 0 {

										if v == reserved["end proto static"] {
											k++
											break
										}

										return nil, NOT_OAT
									}
					
									matchers[strings.TrimPrefix(getReservedFromRune(v), "end ")]--
								}
					
								for _, v := range matchers {
									if v != 0 {
										goto protostatic_esc
									}
								}

								if v == reserved["hash key seperator"] {
									curp = false
									continue
								}
								if v == reserved["value seperator"] {
									staticparts = append(staticparts, [2][]rune{})
									curp = true
									continue
								}

								protostatic_esc:
								if curp {
									staticparts[len(staticparts) - 1][0] = append(staticparts[len(staticparts) - 1][0], v)
								} else {
									staticparts[len(staticparts) - 1][1] = append(staticparts[len(staticparts) - 1][1], v)
								}
							}
							staticparts = staticparts[:len(staticparts) - 1]

							if k + 1 >= len(cv) {
								goto noins
							}

							cv = cv[k + 1:]

							matchers = make(map[string]int)
							curp = true

							for k, v = range cv {

								if escaped {
									escaped = false
									goto protoinstance_esc
								}

								if v == reserved["escaper"] {
									escaped = true
									goto protoinstance_esc
								}

								if strings.HasPrefix(getReservedFromRune(v), "start ") {
									matchers[strings.TrimPrefix(getReservedFromRune(v), "start ")]++
								} else if strings.HasPrefix(getReservedFromRune(v), "end ") {
					
									if matchers[strings.TrimPrefix(getReservedFromRune(v), "end ")] <= 0 {

										if v == reserved["end proto instance"] {
											k++
											break
										}
										
										return nil, NOT_OAT
									}
					
									matchers[strings.TrimPrefix(getReservedFromRune(v), "end ")]--
								}
					
								for _, v := range matchers {
									if v != 0 {
										goto protoinstance_esc
									}
								}

								if v == reserved["hash key seperator"] {
									curp = false
									continue
								}
								if v == reserved["value seperator"] {
									instanceparts = append(instanceparts, [2][]rune{})
									curp = true
									continue
								}

								protoinstance_esc:
								if curp {
									instanceparts[len(instanceparts) - 1][0] = append(instanceparts[len(instanceparts) - 1][0], v)
								} else {
									instanceparts[len(instanceparts) - 1][1] = append(instanceparts[len(instanceparts) - 1][1], v)
								}
							}

							noins:
							instanceparts = instanceparts[:len(instanceparts) - 1]
							var nstaticp = make(map[string]*OmmType)
							var ninstancep = make(map[string]*OmmType)

							for _, v := range staticparts {
								value, e := decval(v[1])
								if e != nil {
									return nil, e
								}
								nstaticp[string(DecodeStr(v[0]))] = &value
							}

							for _, v := range instanceparts {
								value, e := decval(v[1])
								if e != nil {
									return nil, e
								}
								ninstancep[string(DecodeStr(v[0]))] = &value
							}
							
							return OmmProto{
								ProtoName: name,
								Static: nstaticp,
								Instance: ninstancep,
							}, nil

						case reserved["make rune"]:
							if len(cv) != 3 {
								return nil, NOT_OAT
							}

							return OmmRune{
								Rune: &cv[2],
							}, nil

						case reserved["make string"]:
							cv = cv[1:]

							return OmmString{
								String: DecodeStr(cv),
							}, nil

						case reserved["make undef"]:
							return OmmUndef{}, nil

						case reserved["start function"]:
							cv = cv[1:]

							if cv[len(cv) - 1] != reserved["end function"] {
								return nil, NOT_OAT
							}

							cv = cv[:len(cv) - 1]

							var params = make([][2]string, 1)

							var curp = false
							var o int
							var escaped = false

							var matchers = make(map[string]int)

							for o = 0; o < len(cv); o++ {

								if escaped {
									goto isescaped
								}

								if cv[o] == reserved["escaper"] {
									escaped = true
									continue
								}

								if strings.HasPrefix(getReservedFromRune(encoded[i]), "start ") {
									matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "start ")]++
								} else if strings.HasPrefix(getReservedFromRune(encoded[i]), "end ") {
					
									if matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "end ")] <= 0 {
										return nil, NOT_OAT
									}
					
									matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "end ")]--
								}
					
								for _, v := range matchers {
									if v != 0 {
										goto isescaped
									}
								}

								if cv[o] == reserved["param body split"] {
									o++
									break
								}

								isescaped:
								if curp {
									params[len(params) - 1][0]+=string(cv[o])
								} else {
									params[len(params) - 1][1]+=string(cv[o])
								}
								escaped = false
							}

							params = params[:len(params) - 1] //remove the trailing [2]string
							cv = cv[o:]

							var fn OmmFunc
							fn.Overloads = make([]Overload, 1)

							fn.Overloads[0].Types = make([]string, len(params))
							fn.Overloads[0].Params = make([]string, len(params))
							
							for k, v := range params {
								fn.Overloads[0].Types[k] = v[0]
								fn.Overloads[0].Params[k] = v[1]
							}

							bodyv, e := DecodeVariable(cv)

							if e != nil {
								return nil, e
							}

							fn.Overloads[0].Body = bodyv

							return fn, nil

						default:
							return nil, NOT_OAT
					}
				}

				if len(curval) != 0 {
					calc, e := decval(curval)
					if e != nil {
						return nil, e
					}
					(*curact).Value = calc
				}

			case "ExpAct":
				if len(curval) != 0 {
					curval = curval[1:len(curval) - 1]
					val, e := DecodeVariable(curval)
					if e != nil {
						return nil, e
					}
					(*curact).ExpAct = val
				}

			case "First":
				if len(curval) != 0 {
					curval = curval[1:len(curval) - 1]
					val, e := DecodeVariable(curval)
					if e != nil {
						return nil, e
					}
					(*curact).First = val
				}

			case "Second":
				if len(curval) != 0 {
					curval = curval[1:len(curval) - 1]
					val, e := DecodeVariable(curval)
					if e != nil {
						return nil, e
					}
					(*curact).Second = val
				}

			case "Array":

				var arr [][]Action
				var current []rune

				escaped := false

				if len(curval) < 2 {
					return nil, NOT_OAT
				}

				var matchers = make(map[string]int)

				for o := 1; o < len(curval) - 1; o++ {

					for ;o < len(curval); o++ {

						if escaped {
							escaped = false
							goto arr_esc
						}

						if curval[o] == reserved["escaper"] {
							escaped = true
							goto arr_esc
						}

						if strings.HasPrefix(getReservedFromRune(encoded[i]), "start ") {
							matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "start ")]++
						} else if strings.HasPrefix(getReservedFromRune(encoded[i]), "end ") {
			
							if matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "end ")] <= 0 {
								return nil, NOT_OAT
							}
			
							matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "end ")]--
						}
			
						for _, v := range matchers {
							if v != 0 {
								goto arr_esc
							}
						}

						if getReservedFromRune(curval[o]) == "value seperator" {
							break
						}

						arr_esc:
						current = append(current, curval[o])
					}

					dec, e := DecodeVariable(current)
					if e != nil {
						return nil, e
					}
					arr = append(arr, dec)
					current = nil //clear current
				}

				(*curact).Array = arr

			case "Hash":

				var hash = make([][2][]Action, 0)
				var currentk []rune
				var currentv []rune

				curpoint := true
				escaped := false

				if len(curval) < 2 {
					return nil, NOT_OAT
				}

				var matchers = make(map[string]int)

				for o := 1; o < len(curval) - 1; o++ {

					for ;o < len(curval); o++ {

						if escaped {
							escaped = false
							goto hash_esc
						}

						if curval[o] == reserved["escaper"] {
							escaped = true
							goto hash_esc
						}

						if strings.HasPrefix(getReservedFromRune(encoded[i]), "start ") {
							matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "start ")]++
						} else if strings.HasPrefix(getReservedFromRune(encoded[i]), "end ") {
			
							if matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "end ")] <= 0 {
								return nil, NOT_OAT
							}
			
							matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "end ")]--
						}
			
						for _, v := range matchers {
							if v != 0 {
								goto hash_esc
							}
						}

						if getReservedFromRune(curval[o]) == "hash key seperator" {
							if !curpoint {
								return nil, NOT_OAT
							}
							curpoint = false
							continue
						}

						if getReservedFromRune(curval[o]) == "value seperator" {
							if curpoint {
								return nil, NOT_OAT
							}
							curpoint = true
							continue
						}

						hash_esc:
						if curpoint {
							currentk = append(currentk, curval[o])
						} else {
							currentv = append(currentv, curval[o])
						}
					}

					deckey, e :=  DecodeVariable(currentk[1:len(currentk) - 2])
					if e != nil {
						return nil, e
					}
					decval, e :=  DecodeVariable(currentv)
					if e != nil {
						return nil, e
					}

					hash = append(hash, [2][]Action{
						deckey,
						decval,
					})

					//clear both currents
					currentk = nil
					currentv = nil
				}

				(*curact).Hash = hash

			case "File":
				(*curact).File = string(DecodeStr(curval))

			case "Line":
				if len(curval) != 2 {
					return nil, NOT_OAT
				}

				if curval[0] != reserved["escaper"] {
					return nil, NOT_OAT
				}

				(*curact).Line = uint64(curval[1])
		}

		if encoded[i + 1] == reserved["next action"] {
			i++
			curpos = -1
			decoded = append(decoded, Action{})
		}

		curval = nil //clear curval
		curpos++
	}

	//remove the trailing one
	decoded = decoded[:len(decoded) - 1]

	return decoded, nil
}

func DecodeStr(str []rune) []rune {

	var final []rune
	var escaped = false

	for _, v := range str {

		if !escaped && v == reserved["escaper"] {
			escaped = true
			continue
		}

		final = append(final, v)

		escaped = false
	}

	return final
}

func getReservedFromRune(val rune) string {

	for k, v := range reserved {
		if v == val {
			return k
		}
	}

	return ""
}
