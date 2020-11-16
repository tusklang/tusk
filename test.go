package main

import (
	"bytes"
	"fmt"
	"strconv"
)

func chunks(s string, n int64) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%int(n) == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}

func rol(n, b int64) int64 {
	return ((n << b) | (n >> (32 - b))) & 0xffffffff
}

func main() {

	data := "abc"

	bytes := ""

	h0 := int64(0x67452301)
	h1 := int64(0xEFCDAB89)
	h2 := int64(0x98BADCFE)
	h3 := int64(0x10325476)
	h4 := int64(0xC3D2E1F0)

	for n := 0; n < len(data); n++ {
		bytes += fmt.Sprintf("%08b", data[n])
	}

	bits := bytes + "1"
	pBits := bits
	for len(pBits)%512 != 448 {
		pBits += "0"
	}

	pBits += fmt.Sprintf("%064b", len(bits)-1)

	for _, c := range chunks(pBits, 512) {
		words := chunks(c, 32)
		w := make([]int64, 80)
		for n := 0; n < 16; n++ {
			w[n], _ = strconv.ParseInt(words[n], 2, 64)
		}
		for i := 16; i < 80; i++ {
			w[i] = int64(rol(w[i-3]^w[i-8]^w[i-14]^w[i-16], 1))
		}

		a := h0
		b := h1
		c := h2
		d := h3
		e := h4

		for i := 0; i < 80; i++ {
			var f int64
			var k int64
			if i >= 0 && i <= 19 {
				f = (b & c) | ((^b) & d)
				k = 0x5A827999
			} else if i >= 20 && i <= 39 {
				f = b ^ c ^ d
				k = 0x6ED9EBA1
			} else if i >= 40 && i <= 59 {
				f = (b & c) | (b & d) | (c & d)
				k = 0x8F1BBCDC
			} else if i >= 60 && i <= 79 {
				f = b ^ c ^ d
				k = 0xCA62C1D6
			}

			temp := (rol(a, 5) + f + e + k + w[i]) & 0xffffffff
			e = d
			d = c
			c = rol(b, 30)
			b = a
			a = temp
		}

		h0 = (h0 + a) & 0xffffffff
		h1 = (h1 + b) & 0xffffffff
		h2 = (h2 + c) & 0xffffffff
		h3 = (h3 + d) & 0xffffffff
		h4 = (h4 + e) & 0xffffffff
	}

	fmt.Printf("%08x %08x %08x %08x %08x\n", h0, h1, h2, h3, h4)
}
