package oat

import "os"
import . "lang/types"

func OatEncode(filename string, data Oat) {

	acts := EncodeActions(data.Actions)

	f, _ := os.Create(filename)
	f.Write(acts)

}

func EncodeActions(data []Action) []byte {

	var final string

	for _, v := range data {
		final+=string(rune(hash(v.Name))) + "("



		final+=")"
	}

	return []byte(final)
}

func hash(name string) int32 { //a hash function to hash an action type (because giving each one a unique id is so much work)

	var hash int32 = 1
	var c int32

	for c = 0; int(c) < len(name); c++ {
		hash = ((hash << 2) + hash) + int32(name[c])
	}

	return hash
}