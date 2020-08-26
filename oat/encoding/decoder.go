package oatenc

import (
	"bufio"
	"encoding/json"
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
		decodedV, e := DecodeAction(v[1])

		if e != nil {
			return nil, e
		}

		decodedvars[string(DecodeStr(v[0]))] = decodedV
	}

	j, _ := json.MarshalIndent(decodedvars, "", "  ")
	fmt.Println(string(j))

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

func DecodeAction(encoded []rune) ([]Action, error) {
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

			var putval = func(cv []rune) (OmmType, error) {

				if len(cv) == 0 {
					return nil, NOT_OAT
				}

				switch cv[0] {
				case reserved["make c-array"]:
				case reserved["make bool"]:

					if len(cv) != 3 {
						return nil, NOT_OAT
					}

					cv = cv[2:]
					boolv := cv[0] == 1 //get the value as a bool (1 = true, 0 = false)

					var ommbool OmmBool
					ommbool.FromGoType(boolv)
					return ommbool, nil

				case reserved["start function"]:
				case reserved["make c-hash"]:
				case reserved["start number"]:

					if len(cv) < 2 {
						return nil, NOT_OAT
					}

					cv = cv[1 : len(cv)-1]

					num := decode2d(cv, reserved["decimal spot"])

					if len(num) > 2 || len(num) == 0 {
						return nil, NOT_OAT
					}

					if len(num) == 1 {
						num = append(num, []rune{})
					}

					var str1 = DecodeRaw(num[0])
					var str2 = DecodeRaw(num[1])

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

				case reserved["make rune"]:

					if len(cv) != 2 {
						return nil, NOT_OAT
					}

					var ommrune OmmRune
					ommrune.FromGoType(cv[1])
					return ommrune, nil

				case reserved["make string"]:

					cv = cv[1:]

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

			var e error                         //declare "e" here
			(*curact).Value, e = putval(curval) //and then put the value of "e" here
			if e != nil {                       //and now return the error if there was one
				return nil, e
			}

		case "expact":
			if len(curval) > 2 {
				val, e := DecodeAction(curval[1 : len(curval)-1])

				if e != nil {
					return nil, e
				}

				(*curact).ExpAct = val
			}

		case "first":
			if len(curval) > 2 {
				val, e := DecodeAction(curval[1 : len(curval)-1])

				if e != nil {
					return nil, e
				}

				(*curact).First = val
			}

		case "second":
			if len(curval) > 2 {
				val, e := DecodeAction(curval[1 : len(curval)-1])

				if e != nil {
					return nil, e
				}

				(*curact).Second = val
			}

		case "array":

			if len(curval) < 2 {
				return nil, NOT_OAT
			}

			curval = curval[1 : len(curval)-1]
			var _arr = decode2d(curval, reserved["value seperator"])

			var arr = make([][]Action, len(_arr))

			for kk, vv := range _arr {
				val, e := DecodeAction(vv)

				if e != nil {
					return nil, e
				}

				arr[kk] = val
			}

			(*curact).Array = arr

		case "hash":

			if len(curval) < 2 {
				return nil, NOT_OAT
			}

			curval = curval[1 : len(curval)-1]
			_hash := decode3d(curval, reserved["hash key seperator"], reserved["value seperator"])

			var hash = make([][2][]Action, len(_hash))

			for kk, vv := range _hash {
				key, e := DecodeAction(vv[0])

				if e != nil {
					return nil, NOT_OAT
				}

				val, e := DecodeAction(vv[1])

				if e != nil {
					return nil, NOT_OAT
				}

				hash[kk] = [2][]Action{key, val}
			}

			(*curact).Hash = hash

		}

		curval = nil
	}

	return decoded, nil
}

//DecodeStr Decode an oat string
func DecodeStr(str []rune) []rune {

	rawd := DecodeRaw(str)

	for k := range rawd {
		rawd[k] -= 5000
	}

	return rawd
}

//DecodeRaw Decode an oat raw string
func DecodeRaw(str []rune) []rune {

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
	var matchers map[string]int
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

	return seperated[:len(seperated)-1]
}

func decode3d(encoded []rune, seperator1 rune, seperator2 rune) [][][]rune {

	var escaped bool
	var matchers map[string]int
	var seperated = make([][][]rune, 1)
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
			seperated = append(seperated, [][]rune{})
			continue
		}

	escaped:

		if !cur {
			seperated[len(seperated)-1][0] = append(seperated[len(seperated)-1][0], v)
		} else {
			seperated[len(seperated)-1][1] = append(seperated[len(seperated)-1][1], v)
		}

	}

	return seperated[:len(seperated)-1]
}
