package compile

import "encoding/gob"
import "os"
import "lang/interpreter/bind"

import "lang" //omm language

//export Compile
func Compile(params bind.CliParams) {

  dir := params.GetFiles().GetDIR()
  fileName := params.GetFiles().GetNAME()

  file := lang.ReadFileJS(dir + fileName)[0]["Content"]

  lex := lang.Lexer(file, dir, fileName)
  acts := lang.Actionizer(lex, false, dir, fileName)

  if (IsAbsolute(params.GetCalc().GetO())) {

    writefile, _ := os.Create(params.GetCalc().GetO())

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(acts)
  } else {

    writefile, _ := os.Create(dir + params.GetCalc().GetO())

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(acts)
  }
}
