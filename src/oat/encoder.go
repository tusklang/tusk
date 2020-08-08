package oat

import "os"
import "reflect"
import . "lang/types"

func OatEncode(filename string, data Oat) error {

	acts := EncodeActions(data.Actions)

	f, e := os.Create(filename)

	if e != nil {
		return e
	}

	f.Write([]byte(string(acts)))
	f.Close()

	return nil
}

func EncodeActions(data []Action) []rune {

	var final []rune

	for _, v := range data {
		final = append(final, actionids[v.Type])
		final = append(final, '(')

		s := reflect.TypeOf(v)

		for i := 1; i < s.NumField(); i++ {
			field := s.Field(i)

			//add the tag to the final resp
			final = append(final, []rune{ []rune(field.Tag.Get("oat"))[0], ' ' }...)
		}

		final = append(final, ')')
	}

	return final
}