package interpreter

func num_normalize(num Action) string {

  integer, decimal := sliceToBigInt(num.Integer), sliceToBigInt(num.Decimal)

  return integer.String() + "." + decimal.String()
}
