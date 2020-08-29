package oatenc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	. "github.com/omm-lang/omm/lang/types"
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
				vers += string(c)
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

	if LASTDEP[0] >= majorv && LASTDEP[1] >= minorv && LASTDEP[2] >= bugv {
		return nil, fmt.Errorf("Version %d.%d.%d has been depricated from your Omm version", majorv, minorv, bugv)
	}

	if majorv > OMM_MAJOR && minorv > OMM_MINOR && bugv > OMM_BUG {
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
			encodedVars[len(encodedVars)-1][0] = append(encodedVars[len(encodedVars)-1][0], v)
		} else {
			encodedVars[len(encodedVars)-1][1] = append(encodedVars[len(encodedVars)-1][1], v)
		}
	}

	//remove the trailing
	encodedVars = encodedVars[:len(encodedVars)-1]

	var decodedvars = make(map[string][]Action)

	for _, v := range encodedVars {
		decodedV, e := DecodeActions(v[1])

		if e != nil {
			return nil, e
		}

		decodedvars[string(DecodeStr(v[0]))] = decodedV
	}

	if mode == 0 {
		return decodedvars, nil
	}

	var all []Action
	for k, v := range decodedvars {
		all = append(all, Action{
			Type:   "var",
			Name:   k,
			ExpAct: v,
		})
	}
	return map[string][]Action{
		"$main": all,
	}, nil

}

func DecodeActions(encoded []rune) ([]Action, error) {
	var (
		decoded = make([]Action, 0)
		escaped = false
		curval  = make([]rune, 0)
		curact  = &Action{}
	)

	for i := 0; i < len(encoded); i++ {

		var matchers = make(map[string]int)

		for ; i < len(encoded); i++ {

			if escaped {
				escaped = false
				goto escaped
			}

			if encoded[i] == reserved["escaper"] {
				escaped = true
				goto escaped
			}

			if strings.HasPrefix(getReservedFromRune(encoded[i]), "start ") {
				matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "start ")]++
			}
			if strings.HasPrefix(getReservedFromRune(encoded[i]), "end ") {
				matchers[strings.TrimPrefix(getReservedFromRune(encoded[i]), "end ")]--
			}

			for _, v := range matchers {
				if v != 0 {
					goto escaped
				}
			}

			if strings.HasPrefix(getReservedFromRune(encoded[i]), "seperate ") {
				goto nextfield
			}

			if encoded[i] == reserved["next action"] {
				decoded = append(decoded, *curact)
				curact = &Action{}
				continue
			}

		escaped:
			curval = append(curval, encoded[i])
		}
		break

	nextfield:

		switch strings.TrimPrefix(getReservedFromRune(encoded[i]), "seperate ") {
		case "file":
			(*curact).File = string(DecodeStr(curval))

		case "line":
			if len(curval) != 2 {
				return nil, NOT_OAT
			}

			if curval[0] != reserved["escaper"] {
				return nil, NOT_OAT
			}

			(*curact).Line = uint64(curval[1])

		case "type":
			if len(curval) != 1 {
				return nil, NOT_OAT
			}

			(*curact).Type = getReservedFromRune(curval[0])

		case "name":

			(*curact).Name = string(DecodeStr(curval))

		case "value":

			var putval func([]rune) (OmmType, error)
			putval = func(cv []rune) (OmmType, error) {

				if len(cv) == 0 { //if it is nil return undef
					return OmmUndef{}, nil
				}

				switch cv[0] {
				case reserved["make c-array"]:

					cv = cv[1:]

					_arr := decode2d(cv, reserved["value seperator"])
					arr := make([]*OmmType, len(_arr))

					for k, v := range _arr {
						val, e := putval(v)
						if e != nil {
							return nil, e
						}
						arr[k] = &val
					}

					var ommarr OmmArray
					ommarr.Array = arr
					return ommarr, nil

				case reserved["make bool"]:

					if len(cv) != 3 {
						return nil, NOT_OAT
					}

					cv = cv[2:]         //remove the "make bool" and the escaper
					boolv := cv[0] == 1 //get the value as a bool (1 = true, 0 = false)

					var ommbool OmmBool
					ommbool.FromGoType(boolv)
					return ommbool, nil

				case reserved["start function"]:

					if len(cv) < 2 { //it must (at least) start with "start function" and end with "end function"
						return nil, NOT_OAT
					}

					cv = cv[1 : len(cv)-1]

					var encodedoverloads = decode2d(cv, reserved["seperate overload"])

					var overloads = make([]Overload, len(encodedoverloads)) //alloc the amount of overloads there are

					for kk, vv := range encodedoverloads {

						if len(vv) == 0 {
							continue
						}

						var parambodysplit = decode2d(vv, reserved["param body split"]) //seperate the params and the body

						if len(parambodysplit) != 2 {
							return nil, NOT_OAT
						}

						params, encbody := decode3d(parambodysplit[0], reserved["seperate type-param"], reserved["value seperator"]), decode2d(parambodysplit[1], reserved["body var-ref split"])

						var (
							pnames []string
							types  []string
						)

						for _, v := range params {

							typed := string(DecodeStr(v[0]))
							paramd := string(DecodeStr(v[1]))

							if paramd == "" || typed == "" { //skip empty params
								continue
							}

							types = append(types, typed)
							pnames = append(pnames, paramd)
						}

						decbody, e := DecodeActions(encbody[0]) //decode the body of the function

						if e != nil {
							return nil, e
						}

						_varrefs := decode2d(encbody[1], reserved["value seperator"]) //seperate all of the var references
						if len(_varrefs[len(_varrefs)-1]) == 0 {                      //remove the last empty one
							_varrefs = _varrefs[:len(_varrefs)-1]
						}
						var varrefs = make([]string, len(_varrefs)) //alloc the space of the []string var refs

						for k, v := range _varrefs {
							varrefs[k] = string(DecodeStr(v))
						}

						overloads[kk] = Overload{
							Params:  pnames,
							Types:   types,
							Body:    decbody,
							VarRefs: varrefs,
						}
					}

					var fn OmmFunc
					fn.Overloads = overloads

					return fn, nil

				case reserved["make c-hash"]:

					cv = cv[1:]

					_hash := decode3d(cv, reserved["hash key seperator"], reserved["value seperator"])
					hash := make(map[string]*OmmType)

					for _, v := range _hash {
						key := DecodeStr(v[0])
						val, e := putval(v[1])
						if e != nil {
							return nil, e
						}

						hash[string(key)] = &val
					}

					var ommhash OmmHash
					ommhash.Hash = hash
					return ommhash, nil

				case reserved["start number"]:

					if len(cv) < 2 { //numbers must at least have (start number - end number)
						return nil, NOT_OAT
					}

					cv = cv[1 : len(cv)-1] //remove the previously mentioned (start number and end number paddings)

					num := decode2d(cv, reserved["decimal spot"])

					if len(num) > 2 || len(num) == 0 {
						return nil, NOT_OAT
					}

					if len(num) == 1 {
						num = append(num, []rune{})
					}

					var str1 = DecodeStr(num[0])
					var str2 = DecodeStr(num[1])

					var integer = make([]int64, len(str1))
					var decimal = make([]int64, len(str2))

					for k, v := range str1 {
						integer[k] = int64(v)
					}
					for k, v := range str2 {
						decimal[k] = int64(v)
					}

					var ommnum OmmNumber
					ommnum.Integer = &integer
					ommnum.Decimal = &decimal
					return ommnum, nil

				case reserved["start proto"]:

					if len(cv) < 2 { //it must at least start with "start proto" and end with "end proto"
						return nil, NOT_OAT
					}

					cv = cv[1 : len(cv)-1]

					seperatedname := decode2d(cv, reserved["seperate proto name"])

					if len(seperatedname) != 2 { //it must have a name and a body
						return nil, NOT_OAT
					}

					name := string(DecodeStr(seperatedname[0]))

					body := decode2d(seperatedname[1], reserved["seperate proto static instance"])

					if len(body[len(body)-1]) == 0 {
						body = body[:len(body)-1]
					}

					if len(body) != 2 && len(body) != 3 { //it must have a static and an instance
						return nil, NOT_OAT
					}

					var parseprotobody = func(part []rune) (map[string]*OmmType, error) {
						decoded := decode3d(part, reserved["hash key seperator"], reserved["value seperator"])

						var protomap = make(map[string]*OmmType)

						for _, v := range decoded {

							if len(v[1]) == 0 { //if it is empty, skip it
								continue
							}

							tmp, e := putval(v[1]) //create a tmp to create a ptr
							if e != nil {
								return nil, e
							}
							protomap[string(DecodeStr(v[0]))] = &tmp
						}

						return protomap, nil
					}

					//parse the static and the instance
					static, e := parseprotobody(body[0])
					if e != nil {
						return nil, e
					}
					instance, e := parseprotobody(body[1])
					if e != nil {
						return nil, e
					}
					///////////////////////////////////

					//also decode the access list (if there is one)
					var accesslist = make(map[string][]string)

					if len(body) == 3 {
						decodedAccessList := decode3d(body[2], reserved["hash key seperator"], reserved["value seperator"])

						for _, vv := range decodedAccessList {

							if len(vv[0]) == 0 { //skip empty values
								continue
							}

							curlist := decode2d(vv[1], reserved["sub value seperator"])
							var parsedlist = make([]string, len(curlist))

							for kk, vvv := range curlist {
								parsedlist[kk] = string(DecodeStr(vvv))
							}

							accesslist[string(DecodeStr(vv[0]))] = parsedlist
						}

					}

					var proto = OmmProto{
						ProtoName:  name,
						Static:     static,
						Instance:   instance,
						AccessList: accesslist,
					}

					return proto, nil

				case reserved["make rune"]:

					cv = cv[1:]

					if len(cv) != 2 {
						return nil, NOT_OAT
					}

					//runes must have an escaper and a rune

					var ommrune OmmRune
					ommrune.FromGoType(cv[1])
					return ommrune, nil

				case reserved["make string"]:

					cv = cv[1:] //remove the "make string"

					runelist := DecodeStr(cv)
					var ommstr OmmString
					ommstr.FromRuneList(runelist)
					return ommstr, nil

				case reserved["make undef"]:

					if len(cv) != 1 {
						return nil, NOT_OAT
					}

					var undef OmmUndef //declare an undefined variable
					return undef, nil  //now return it

				}

				return nil, NOT_OAT
			}

			if len(curval) > 0 { //only if there is a value
				var e error                         //declare "e" here
				(*curact).Value, e = putval(curval) //and then put the value of "e" here
				if e != nil {                       //and now return the error if there was one
					return nil, e
				}
			}

		case "expact":
			if len(curval) > 2 {
				val, e := DecodeActions(curval[1 : len(curval)-1])

				if e != nil {
					return nil, e
				}

				(*curact).ExpAct = val
			}

		case "first":
			if len(curval) > 2 {
				val, e := DecodeActions(curval[1 : len(curval)-1])

				if e != nil {
					return nil, e
				}

				(*curact).First = val
			}

		case "second":
			if len(curval) > 2 {
				val, e := DecodeActions(curval[1 : len(curval)-1])

				if e != nil {
					return nil, e
				}

				(*curact).Second = val
			}

		case "array":

			if len(curval) > 2 {
				curval = curval[1 : len(curval)-1]
				var _arr = decode2d(curval, reserved["value seperator"])

				if len(_arr[len(_arr)-1]) == 0 { //remove the last empty value
					_arr = _arr[:len(_arr)-1]
				}

				var arr = make([][]Action, len(_arr))

				for kk, vv := range _arr {
					val, e := DecodeActions(vv)

					if e != nil {
						return nil, e
					}

					arr[kk] = val
				}

				(*curact).Array = arr
			}

		case "hash":

			if len(curval) > 2 {
				curval = curval[1 : len(curval)-1]

				_hash := decode3d(curval, reserved["hash key seperator"], reserved["value seperator"])

				if len(_hash[len(_hash)-1][0]) == 0 || len(_hash[len(_hash)-1][1]) == 0 { //remove the last empty value
					_hash = _hash[:len(_hash)-1]
				}

				var hash = make([][2][]Action, len(_hash))

				for kk, vv := range _hash {

					if len(vv[0]) < 2 { //it must at least start with "start multi action" and end with "end multi action"
						return nil, NOT_OAT
					}

					vv[0] = vv[0][1 : len(vv[0])-1]
					key, e := DecodeActions(vv[0])

					if e != nil {
						return nil, NOT_OAT
					}

					val, e := DecodeActions(vv[1])

					if e != nil {
						return nil, NOT_OAT
					}

					hash[kk] = [2][]Action{key, val}
				}

				(*curact).Hash = hash
			}

		}

		curval = nil
	}

	return decoded, nil
}

//DecodeStr Decode an oat string
func DecodeStr(str []rune) []rune {

	var final []rune
	var escaped = false

	for _, v := range str {

		decodedr := v

		if !escaped && decodedr == reserved["escaper"] {
			escaped = true
			continue
		}

		final = append(final, decodedr)

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

func decode2d(encoded []rune, seperator rune) [][]rune {

	var escaped bool
	var matchers = make(map[string]int)
	var seperated = make([][]rune, 1)

	for _, v := range encoded {

		if escaped {
			escaped = false
			goto escaped
		}

		if v == reserved["escaper"] {
			escaped = true
			goto escaped
		}

		if strings.HasPrefix(getReservedFromRune(v), "start ") {
			matchers[strings.TrimPrefix(getReservedFromRune(v), "start ")]++
		}
		if strings.HasPrefix(getReservedFromRune(v), "end ") {
			matchers[strings.TrimPrefix(getReservedFromRune(v), "end ")]--
		}

		for _, v := range matchers {
			if v != 0 {
				goto escaped
			}
		}

		if v == seperator {
			seperated = append(seperated, []rune{})
			continue
		}

	escaped:
		seperated[len(seperated)-1] = append(seperated[len(seperated)-1], v)
	}

	return seperated
}

func decode3d(encoded []rune, seperator1 rune, seperator2 rune) [][2][]rune {

	var escaped bool
	var matchers = make(map[string]int)
	var seperated = make([][2][]rune, 1)
	var cur bool

	for _, v := range encoded {

		if escaped {
			escaped = false
			goto escaped
		}

		if v == reserved["escaper"] {
			escaped = true
			goto escaped
		}

		if strings.HasPrefix(getReservedFromRune(v), "start ") {
			matchers[strings.TrimPrefix(getReservedFromRune(v), "start ")]++
		}
		if strings.HasPrefix(getReservedFromRune(v), "end ") {
			matchers[strings.TrimPrefix(getReservedFromRune(v), "end ")]--
		}

		for _, v := range matchers {
			if v != 0 {
				goto escaped
			}
		}

		if v == seperator1 {
			cur = true
			continue
		} else if v == seperator2 {
			cur = false
			seperated = append(seperated, [2][]rune{})
			continue
		}

	escaped:

		if !cur {
			seperated[len(seperated)-1][0] = append(seperated[len(seperated)-1][0], v)
		} else {
			seperated[len(seperated)-1][1] = append(seperated[len(seperated)-1][1], v)
		}

	}

	return seperated
}
