package lang

import "strings"
import "strconv"
import "lang/interpreter/bind"

//export DigitSize
var DigitSize = 1;

//export BigNumConverter
func BigNumConverter(num string) ([]int64, []int64) {

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
      var put []int64

      for _, v := range chunked { //convert all numbers to a negative

        if v == "" {
          continue
        }

        num, e := strconv.ParseInt(v, 10, 64)

        if neg {
          num*=-1
        }

        if e != nil {
          return []int64{}, []int64{}
        } else {
          put = append(put, int64(num))
        }
      }

      return put, []int64{}
    case 2:

      chunked_integer := Chunk(splitted[0], DigitSize)
      var put_int []int64

      for _, v := range chunked_integer { //convert all numbers to a negative

        if v == "" {
          continue
        }

        num, e := strconv.ParseInt(v, 10, 64)

        if neg {
          num*=-1
        }

        if e != nil {
          return []int64{}, []int64{}
        } else {
          put_int = append(put_int, int64(num))
        }
      }

      chunked_decimal := Chunk(splitted[1], DigitSize)
      var put_dec []int64

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
          return []int64{}, []int64{}
        } else {
          put_dec = append(put_dec, int64(num))
        }

      }

      return put_int, put_dec
    default:
      return []int64{ 0 }, []int64{}
  }
}
