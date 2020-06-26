package interpreter

func num_normalize(num Action) string {

  integer, decimal := sliceToBig(num.Integer), sliceToBig(num.Decimal)

  return integer.String() + "." + decimal.String()
}
