package types

import "strconv"

type OmmBool struct {
  Boolean *bool
}

func (b *OmmBool) FromGoType(val bool) {
  b.Boolean = &val
}

func (b OmmBool) ToGoType() bool {
  return *b.Boolean
}

func (b OmmBool) Format() string {
  return strconv.FormatBool(*b.Boolean)
}
