package oat

import "os"
import "reflect"
import "fmt"
import "strings"
import "strconv"
import . "lang/types"
import . "lang/interpreter"

//export OatEncode
func OatEncode(filename string, data Oat) error {

	acts := EncodeActions(data.Actions)

	f, e := os.Create(filename)

	if e != nil {
		return e
	}

	//versioning
	fmt.Fprint(f, "OAT")
	fmt.Fprintln(f, string(reserved["MAJOR"]) + OMM_MAJOR)
	fmt.Fprintln(f, string(reserved["MINOR"]) + OMM_MINOR)
	fmt.Fprintln(f, string(reserved["BUG"]) + OMM_BUG)
	////////////

	f.Write([]byte(string(acts)))
	f.Close()

	return nil
}

//export EncodeActions
func EncodeActions(data []Action) []rune {

	var final []rune

	for _, v := range data {

		// fieldv := reflect.ValueOf(v)
		fieldt := reflect.TypeOf(v)

		for i := 0; i < fieldt.NumField(); i++ {

			switch fieldt.Field(i).Name {
				case "Type":

					final = append(final, reserved[v.Type])

				case "Name":

					if strings.HasPrefix(v.Name, "v ") {
						n, _ := strconv.Atoi(strings.TrimPrefix(v.Name, "v "))
						final = append(final, reserved["varname start"])
						final = append(final, rune(n))
					} else {
						final = append(final, EncodeStr([]rune(v.Name))...)
					}

				case "Value":

					var putval func(v OmmType) []rune
					putval = func(v OmmType) []rune {

						var final []rune

						switch v.(type) {

							case OmmArray:

								final = append(final, reserved["start c-array"])

								for _, v := range v.(OmmArray).Array {
									final = append(final, putval(*v)...)
									final = append(final, reserved["value seperator"])
								}

								final = append(final, reserved["end c-array"])

							case OmmBool:

								final = append(final, reserved["make bool"])
								if v.(OmmBool).ToGoType() {
									final = append(final, 1)
								} else {
									final = append(final, 0)
								}

							case OmmHash:

								final = append(final, reserved["start c-hash"])

								for k, v := range v.(OmmHash).Hash {
									final = append(final, EncodeStr([]rune(k))...)
									final = append(final, reserved["hash key seperator"])
									final = append(final, putval(*v)...)
									final = append(final, reserved["value seperator"])
								}

								final = append(final, reserved["start c-hash"])

							case OmmNumber:

								final = append(final, reserved["start number"])

								if v.(OmmNumber).Integer != nil {
									for _, v := range *v.(OmmNumber).Integer {
										final = append(final, rune(v))
									}
								}

								final = append(final, reserved["decimal place"])

								if v.(OmmNumber).Decimal != nil {
									for _, v := range *v.(OmmNumber).Decimal {
										final = append(final, rune(v))
									}
								}

								final = append(final, reserved["end number"])

							case OmmProto:

								final = append(final, reserved["start proto"])
								final = append(final, reserved["start proto name"])

								//put the name
								final = append(final, EncodeStr([]rune(v.(OmmProto).ProtoName))...)

								final = append(final, reserved["start proto static"])
								for k, v := range v.(OmmProto).Static {
									final = append(final, EncodeStr([]rune(k))...)
									final = append(final, reserved["hash key seperator"])
									final = append(final, putval(*v)...)
								}
								final = append(final, reserved["end proto static"])

								final = append(final, reserved["start proto instance"])
								for k, v := range v.(OmmProto).Instance {
									final = append(final, EncodeStr([]rune(k))...)
									final = append(final, reserved["hash key seperator"])
									final = append(final, putval(*v)...)
								}
								final = append(final, reserved["end proto instance"])

								final = append(final, reserved["start proto"])

							case OmmRune:

								final = append(final, reserved["make rune"])
								final = append(final, v.(OmmRune).ToGoType())

							case OmmString:

								final = append(final, reserved["make string"])
								final = append(final, EncodeStr(v.(OmmString).String)...)

							case OmmUndef:

								final = append(final, reserved["make undef"])

							case OmmFunc:

								final = append(final, reserved["start function"])

								for _, v := range v.(OmmFunc).Overloads {
									final = append(final, reserved["start overload"])

									final = append(final, reserved["start params"])
									for k := range v.Params {
										final = append(final, EncodeStr([]rune(v.Types[k]))...)
										final = append(final, reserved["seperate type-param"])
										final = append(final, EncodeStr([]rune(v.Params[k]))...)
									}
									final = append(final, reserved["end params"])
									final = append(final, EncodeActions(v.Body)...)
									final = append(final, reserved["end overload"])
								}

								final = append(final, reserved["end function"])

						}

						return final
					}

					final = append(final, putval(v.Value)...)

				case "ExpAct":

					final = append(final, reserved["start multi action"])
					final = append(final, EncodeActions(v.ExpAct)...)
					final = append(final, reserved["end multi action"])

				case "First":

					final = append(final, reserved["start multi action"])
					final = append(final, EncodeActions(v.ExpAct)...)
					final = append(final, reserved["end multi action"])

				case "Second":

					final = append(final, reserved["start multi action"])
					final = append(final, EncodeActions(v.ExpAct)...)
					final = append(final, reserved["end multi action"])

				case "Hash":

					final = append(final, reserved["start r-hash"])

					for k, v := range v.Hash {
						final = append(final, EncodeStr([]rune(k))...)
						final = append(final, reserved["hash key seperator"])
						final = append(final, EncodeActions(v)...)
						final = append(final, reserved["value seperator"])
					}

					final = append(final, reserved["end r-hash"])

				case "Array":

					final = append(final, reserved["start r-array"])

					for _, v := range v.Hash {
						final = append(final, EncodeActions(v)...)
						final = append(final, reserved["value seperator"])
					}

					final = append(final, reserved["end r-array"])

				case "File":

					final = append(final, EncodeStr([]rune(v.File))...)

				case "Line":

					final = append(final, rune(v.Line))

			}

			final = append(final, reserved["seperate field"])

		}

		final = append(final, reserved["new action"])
	}

	return final
}

//export EncodeStr
func EncodeStr(splitted []rune) []rune {
	var neg1_multiplied []rune
	for _, v := range splitted {

		negated := -1 * v

		for _, v := range reserved {
			if v == negated {
				neg1_multiplied = append(neg1_multiplied, reserved["escaper"])
			}
		}


		neg1_multiplied = append(neg1_multiplied, negated)
	}
	return neg1_multiplied
}