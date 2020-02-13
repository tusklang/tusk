package main

import "strings"
import "strconv"
import "math"
import "fmt"

func subtract(_num1 string, _num2 string, calc_params paramCalcOpts, line uint64, functions []Funcs) string {

  num1_, num2_ := initAdd(_num1, _num2)

  decPlace := getDec(num1_)

  var final = ""

  if !strings.HasPrefix(num1_, "-") && !strings.HasPrefix(num2_, "-") {

    _num1_ := addDec(strings.Replace(num1_, "-", "0", 1))
    _num2_ := addDec(strings.Replace(num2_, "-", "0", 1))

    num1, num2 := initAdd(_num1_, _num2_)

    decPlace = getDec(num1)

    num1 = strings.Replace(num1, ".", "", 1)
    num2 = strings.Replace(num2, ".", "", 1)

    if !isLess(num1, num2) {

      for i := len(num1) - 1; i >= 0; i-- {
        n1, _ := strconv.ParseUint(string(num1[i]), 10, 64)
        n2, _ := strconv.ParseUint(string(num2[i]), 10, 64)

        if n1 < n2 {
          n1+=10

          o := i - 1
          o1 := []rune(num1)

          if o1[o] == '.' {
            o--
          }

          for ;o1[o] == '0'; {
            o1[o] = '9'
            o--
          }
          cur, _ := strconv.ParseInt(string(o1[o]), 10, 64)
          o1[o] = []rune(strconv.Itoa(int(cur - 1)))[0]

          num1 = string(o1)
        }

        sum := strconv.Itoa(int(n1 - n2))

        final = sum + final

        if calc_params.logger {
          fmt.Println("Omm Logger ~ Subtraction: " + final)
        }
      }

      final = final[:decPlace] + "." + final[decPlace:]
    } else {
      switchOpts := []string{ num1, num2 }

      num1 = switchOpts[1]
      num2 = switchOpts[0]

      for i := len(num1) - 1; i >= 0; i-- {
        n1, _ := strconv.ParseUint(string(num1[i]), 10, 64)
        n2, _ := strconv.ParseUint(string(num2[i]), 10, 64)

        if n1 < n2 {
          n1+=10

          o := i - 1
          o1 := []rune(num1)

          if o1[o] == '.' {
            o--
          }

          for ;o1[o] == '0'; {
            o1[o] = '9'
            o--
          }
          cur, _ := strconv.ParseInt(string(o1[o]), 10, 64)
          o1[o] = []rune(strconv.Itoa(int(cur - 1)))[0]

          num1 = string(o1)
        }

        sum := strconv.Itoa(int(n1 - n2))

        final = sum + final

        if calc_params.logger {
          fmt.Println("Omm Logger ~ Subtraction: " + final)
        }
      }

      final = "-" + final[:decPlace] + "." + final[decPlace:]
    }

  } else if strings.HasPrefix(num1_, "-") && !strings.HasPrefix(num2_, "-") {
    var carry = 0

    _num1_ := addDec(strings.Replace(num1_, "-", "0", 1))
    _num2_ := addDec(strings.Replace(num2_, "-", "0", 1))

    num1__, num2__ := initAdd(_num1_, _num2_)

    decPlace = getDec(num1__)

    num1 := Chunk(strings.ReplaceAll(num1__, ".", ""), 9)
    num2 := Chunk(strings.ReplaceAll(num2__, ".", ""), 9)

    for i := len(num1) - 1; i >= 0; i-- {
      n1, _ := strconv.ParseUint(num1[i], 10, 64)
      n2, _ := strconv.ParseUint(num2[i], 10, 64)

      sum := strconv.Itoa(int(n1 + n2 + uint64(carry)))

      carry = 0
      if len(sum) > len(num1[i]) {
        sum = sum[1:]
        carry = 1
      }

      nL := math.Max(float64(len(num1[i])), float64(len(num2[i])))

      for ;float64(len(sum)) < nL; {
        sum = "0" + sum
      }

      final = sum + final;

      if calc_params.logger {
        fmt.Println("Omm Logger ~ Subtraction: " + final)
      }
    }

    final = "-" + final[:decPlace] + "." + final[decPlace:]

  } else if !strings.HasPrefix(num1_, "-") && strings.HasPrefix(num2_, "-") {
    var carry = 0

    num1_ := addDec(strings.Replace(num1_, "-", "0", 1))
    num2_ := addDec(strings.Replace(num2_, "-", "0", 1))

    num1__, num2__ := initAdd(num1_, num2_)

    decPlace = getDec(num1__)

    num1 := Chunk(strings.ReplaceAll(num1__, ".", ""), 9)
    num2 := Chunk(strings.ReplaceAll(num2__, ".", ""), 9)

    for i := len(num1) - 1; i >= 0; i-- {
      n1, _ := strconv.ParseUint(num1[i], 10, 64)
      n2, _ := strconv.ParseUint(num2[i], 10, 64)

      sum := strconv.Itoa(int(n1 + n2 + uint64(carry)))

      carry = 0
      if len(sum) > len(num1[i]) {
        sum = sum[1:]
        carry = 1
      }

      nL := math.Max(float64(len(num1[i])), float64(len(num2[i])))

      for ;float64(len(sum)) < nL; {
        sum = "0" + sum
      }

      final = sum + final;

      if calc_params.logger {
        fmt.Println("Omm Logger ~ Subtraction: " + final)
      }
    }

    final = final[:decPlace] + "." + final[decPlace:]


  } else if strings.HasPrefix(num1_, "-") && strings.HasPrefix(num2_, "-") {

    num1_ = num1_[1:]
    num2_ = num2_[1:]

    num1 := addDec(strings.Replace(num1_, "-", "0", 1))
    num2 := addDec(strings.Replace(num2_, "-", "0", 1))

    num1, num2 = initAdd(num1, num2)

    decPlace = getDec(num1)

    num1 = strings.Replace(num1, ".", "", 1)
    num2 = strings.Replace(num2, ".", "", 1)

    if !isLess(num1, num2) {

      for i := len(num1) - 1; i >= 0; i-- {
        n1, _ := strconv.ParseUint(string(num1[i]), 10, 64)
        n2, _ := strconv.ParseUint(string(num2[i]), 10, 64)

        if n1 < n2 {
          n1+=10

          o := i - 1
          o1 := []rune(num1)

          if o1[o] == '.' {
            o--
          }

          for ;o1[o] == '0'; {
            o1[o] = '9'
            o--
          }
          cur, _ := strconv.ParseInt(string(o1[o]), 10, 64)
          o1[o] = []rune(strconv.Itoa(int(cur - 1)))[0]

          num1 = string(o1)
        }

        sum := strconv.Itoa(int(n1 - n2))

        final = sum + final

        if calc_params.logger {
          fmt.Println("Omm Logger ~ Subtraction: " + final)
        }
      }

      final = "-" + final[:decPlace] + "." + final[decPlace:]
    } else {
      switchOpts := []string{ num1, num2 }

      num1 = switchOpts[1]
      num2 = switchOpts[0]

      for i := len(num1) - 1; i >= 0; i-- {
        n1, _ := strconv.ParseUint(string(num1[i]), 10, 64)
        n2, _ := strconv.ParseUint(string(num2[i]), 10, 64)

        if n1 < n2 {
          n1+=10

          o := i - 1
          o1 := []rune(num1)

          if o1[o] == '.' {
            o--
          }

          for ;o1[o] == '0'; {
            o1[o] = '9'
            o--
          }
          cur, _ := strconv.ParseInt(string(o1[o]), 10, 64)
          o1[o] = []rune(strconv.Itoa(int(cur - 1)))[0]

          num1 = string(o1)
        }

        sum := strconv.Itoa(int(n1 - n2))

        final = sum + final

        if calc_params.logger {
          fmt.Println("Omm Logger ~ Subtraction: " + final)
        }
      }

      final = final[:decPlace] + "." + final[decPlace:]
    }
  }

  return returnInit(final)
}
