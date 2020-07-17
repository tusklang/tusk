#!/usr/bin/env pwsh
$basedir = Split-Path $MyInvocation.MyCommand.Definition -Parent
$calldir = $PSScriptRoot.replace("\", "/") #convert to slash seperated

function CallScript {

  Param (
    [string]$type,
    [array]$arguments
  )

  #command to call script
  & "$basedir/omm.exe" "$calldir\" $arguments "$type"
}

#if it starts with oat run or oat compile, then perform the given task.
#it is pretty much an alias for oat file.oat -r and oat file.omm -c

if ($args[0] -eq "compile") {

  #remove first element from array
  $_, $firstpop = $args

  CallScript -type --compile -arguments $firstpop
  exit 0

} elseif ($args[0] -eq "run") {

  #remove first element from array
  $_, $firstpop = $args

  CallScript -type --run -arguments $firstpop
  exit 0

}

foreach ($cmd in $args) {
  if ($cmd -eq "--compile" -or $cmd -eq "-c") {

    #function to call script
    CallScript -type --compile -arguments $args

    #exit the script
    exit 0

  } elseif ($cmd -eq "--run" -or $cmd -eq "-r") {

    #function to call script
    CallScript -type --run -arguments $args

    #exit the script
    exit 0

  }
}

#if no cli param was given, just call with --run
CallScript -type --run -arguments $args
