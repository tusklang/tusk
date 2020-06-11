package lang

import "strings"
import "strconv"

// #cgo CFLAGS: -std=c99
import "C"

//see interpreter/run.hpp for more info about DigitSize
var DigitSize = 1; //if you change this, change it also in interpreter/run.hpp

func bigNumConverter(num string) ([]int, []int) {

  //this converts a numeric value, e.g. 123456
  //to [1, 2, 3, 4, 4, 5, 6]
  //for negatives it converts -123456
  //to [-1, -2, -3, -4, -5, -6]

  neg := false

  if strings.HasPrefix(num, "-") {
    neg = true
    num = num[1:]
  }

  splitted := strings.Split(num, ".")

  switch len(splitted) {
    case 1:

      chunked := Chunk(splitted[0], DigitSize)
      reverseStringSlice(chunked)
      var put []int

      for _, v := range chunked { //convert all numbers to a negative

        if v == "" {
          continue
        }

        num, e := strconv.ParseInt(v, 10, 64)

        if neg {
          num*=-1
        }

        if e != nil {
          return []int{}, []int{}
        } else {
          put = append(put, int(num))
        }
      }

      return put, []int{}
    case 2:

      chunked_integer := Chunk(splitted[0], DigitSize)
      reverseStringSlice(chunked_integer)
      var put_int []int

      for _, v := range chunked_integer { //convert all numbers to a negative

        if v == "" {
          continue
        }

        num, e := strconv.ParseInt(v, 10, 64)

        if neg {
          num*=-1
        }

        if e != nil {
          return []int{}, []int{}
        } else {
          put_int = append(put_int, int(num))
        }
      }

      chunked_decimal := Chunk(splitted[1], DigitSize)
      reverseStringSlice(chunked_decimal)
      var put_dec []int

      for _, v := range chunked_decimal { //convert all numbers to a negative

        if v == "" {
          continue
        }

        for ;len(v) != DigitSize; { //add ending zeros to the decimal .4 -> .4000
          v+="0"
        }

        num, e := strconv.ParseInt(v, 10, 64)

        if neg {
          num*=-1
        }

        if e != nil {
          return []int{}, []int{}
        } else {
          put_dec = append(put_dec, int(num))
        }

      }

      return put_int, put_dec
    default:
      return []int{ 0 }, []int{}
  }
}
