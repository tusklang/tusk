#!/usr/bin/env pwsh

$basedir = Split-Path $MyInvocation.MyCommand.Definition -Parent
$calldir = $PSScriptRoot.replace("\", "/") #convert to slash seperated

if ($args[0] -eq "get") {

  $_, $args_pop = $args

  if ($calldir -notmatch "\\$") {
    $calldir+="\"
  }

  & "$basedir/omm.exe" "$calldir" $args_pop "--mango-get"
} elseif ($args[0] -eq "rm" -or $args[0] -eq "remove") {

  $_, $args_pop = $args

  if ($calldir -notmatch "\\$") {
    $calldir+="\"
  }

  & "$basedir/omm.exe" "$calldir" $args_pop "--mango-rm"
} elseif ($args[0] -eq "wipe") {

  $_, $args_pop = $args

  if ($calldir -notmatch "\\$") {
    $calldir+="\"
  }

  & "$basedir/omm.exe" "$calldir" $args_pop "--mango-wipe"
}
