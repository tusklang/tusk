package oatenc

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	. "github.com/omm-lang/omm/lang/types"
)

//export OatEncode
func OatEncode(filename string, data map[string][]Action) error {

	f, e := os.Create(filename)

	if e != nil {
		return e
	}

	//versioning and magic #
	fmt.Fprint(f, MAGIC)
	fmt.Fprintf(f, "%d.%d.%d\n", OMM_MAJOR, OMM_MINOR, OMM_BUG)
	////////////////////////

	for k, v := range data {

		var name = EncodeStr([]rune(k))
		var nameinter = make([]interface{}, len(name))
		var format = ""
		for k, v := range name {
			format += "%c"
			nameinter[k] = v
		}

		fmt.Fprintf(f, format, nameinter...)
		fmt.Fprintf(f, "%c", reserved["set global"])

		var encoded = EncodeActions(v)

		for _, v := range encoded {
			fmt.Fprintf(f, "%c", v)
		}

		fmt.Fprintf(f, "%c", reserved["new global"])
	}

	return nil
}

//export EncodeActions
func EncodeActions(data []Action) []rune {

	var final []rune

	for _, v := range data {
		fieldt := reflect.TypeOf(v)

		for i := 0; i < fieldt.NumField(); i++ {

			switch fieldt.Field(i).Name {

			case "File":

				final = append(final, EncodeStr([]rune(v.File))...)

			case "Line":

				final = append(final, reserved["escaper"], rune(v.Line))

			case "Type":

				final = append(final, reserved[v.Type])

			case "Name":

				final = append(final, EncodeStr([]rune(v.Name))...)

			case "Value":

				var putval func(v OmmType) []rune
				putval = func(v OmmType) []rune {

					var final []rune

					switch v.(type) {

					case OmmArray:

						final = append(final, reserved["make c-array"])

						for _, v := range v.(OmmArray).Array {
							final = append(final, putval(*v)...)
							final = append(final, reserved["value seperator"])
						}

					case OmmBool:

						final = append(final, reserved["make bool"], reserved["escaper"])
						if v.(OmmBool).ToGoType() {
							final = append(final, 1)
						} else {
							final = append(final, 0)
						}

					case OmmFunc:

						final = append(final, reserved["start function"])

						for _, v := range v.(OmmFunc).Overloads {
							for k := range v.Params {
								final = append(final, EncodeStr([]rune(v.Types[k]))...)
								final = append(final, reserved["seperate type-param"])
								final = append(final, EncodeStr([]rune(v.Params[k]))...)
								final = append(final, reserved["value seperator"])
							}
							final = append(final, reserved["param body split"])
							final = append(final, EncodeActions(v.Body)...)
							final = append(final, reserved["body var-ref split"])

							//list all of the variables that this function uses
							for _, v := range v.VarRefs {
								final = append(final, EncodeStr([]rune(v))...)
								final = append(final, reserved["value seperator"])
							}

							final = append(final, reserved["seperate overload"])
						}

						final = append(final, reserved["end function"])

					case OmmHash:

						final = append(final, reserved["make c-hash"])

						for k, v := range v.(OmmHash).Hash {
							final = append(final, EncodeStr([]rune(k))...)
							final = append(final, reserved["hash key seperator"])
							final = append(final, putval(*v)...)
							final = append(final, reserved["value seperator"])
						}

					case OmmNumber:

						final = append(final, reserved["start number"])

						if v.(OmmNumber).Integer != nil && len(*v.(OmmNumber).Integer) != 0 {
							for _, v := range *v.(OmmNumber).Integer {
								final = append(final, reserved["escaper"], rune(v))
							}
						}

						final = append(final, reserved["decimal spot"])

						if v.(OmmNumber).Decimal != nil && len(*v.(OmmNumber).Decimal) != 0 {
							for _, v := range *v.(OmmNumber).Decimal {
								final = append(final, reserved["escaper"], rune(v))
							}
						}

						final = append(final, reserved["end number"])

					case OmmProto:

						final = append(final, reserved["start proto"])

						//put the name
						final = append(final, EncodeStr([]rune(v.(OmmProto).ProtoName))...)
						final = append(final, reserved["seperate proto name"])
						//////////////

						for k, v := range v.(OmmProto).Static {
							final = append(final, EncodeStr([]rune(k))...)
							final = append(final, reserved["hash key seperator"])
							final = append(final, putval(*v)...)
							final = append(final, reserved["value seperator"])
						}

						final = append(final, reserved["seperate proto static instance"])

						for k, v := range v.(OmmProto).Instance {
							final = append(final, EncodeStr([]rune(k))...)
							final = append(final, reserved["hash key seperator"])
							final = append(final, putval(*v)...)
							final = append(final, reserved["value seperator"])
						}

						final = append(final, reserved["seperate proto static instance"])

						/*       put the access list       */
						for k, v := range v.(OmmProto).AccessList {
							final = append(final, EncodeStr([]rune(k))...)
							final = append(final, reserved["hash key seperator"])

							for _, vv := range v {
								final = append(final, EncodeStr([]rune(vv))...)
								final = append(final, reserved["sub value seperator"])
							}

							final = append(final, reserved["value seperator"])
						}
						/////////////////////////////////////

						final = append(final, reserved["seperate proto static instance"]) //also put the seperator here to denote the access list

						final = append(final, reserved["end proto"])

					case OmmRune:

						final = append(final, reserved["make rune"], reserved["escaper"])
						final = append(final, v.(OmmRune).ToGoType())

					case OmmString:

						final = append(final, reserved["make string"])
						final = append(final, EncodeStr(v.(OmmString).String)...)

					case OmmUndef:

						final = append(final, reserved["make undef"])

					}

					return final
				}

				final = append(final, putval(v.Value)...)

			case "ExpAct":

				if len(v.ExpAct) != 0 {
					final = append(final, reserved["start multi action"])
					final = append(final, EncodeActions(v.ExpAct)...)
					final = append(final, reserved["end multi action"])
				}

			case "First":

				if len(v.First) != 0 {
					final = append(final, reserved["start multi action"])
					final = append(final, EncodeActions(v.First)...)
					final = append(final, reserved["end multi action"])
				}

			case "Second":

				if len(v.Second) != 0 {
					final = append(final, reserved["start multi action"])
					final = append(final, EncodeActions(v.Second)...)
					final = append(final, reserved["end multi action"])
				}

			case "Array":

				final = append(final, reserved["start r-array"])

				for _, v := range v.Array {
					final = append(final, EncodeActions(v)...)
					final = append(final, reserved["value seperator"])
				}

				final = append(final, reserved["end r-array"])

			case "Hash":

				final = append(final, reserved["start r-hash"])

				for _, v := range v.Hash {
					final = append(final, reserved["start multi action"])
					final = append(final, EncodeActions(v[0])...)
					final = append(final, reserved["end multi action"])
					final = append(final, reserved["hash key seperator"])
					final = append(final, EncodeActions(v[1])...)
					final = append(final, reserved["value seperator"])
				}

				final = append(final, reserved["end r-hash"])

			}

			final = append(final, reserved["seperate "+strings.ToLower(fieldt.Field(i).Name)])

		}

		final = append(final, reserved["next action"])
	}

	return final
}

//export EncodeStr
func EncodeStr(splitted []rune) []rune {
	var slc []rune
	for _, v := range splitted {

		for _, vv := range reserved {
			if vv == v {
				slc = append(slc, reserved["escaper"])
				break
			}
		}

		slc = append(slc, v)
	}
	return slc
}
