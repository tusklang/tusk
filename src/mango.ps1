#!/usr/bin/env pwsh

$basedir = Split-Path $MyInvocation.MyCommand.Definition -Parent
$calldir = $PSScriptRoot

if ($args[0] -eq "get") {
  & "$basedir/omm.exe" "$calldir" $args "--mango-get"
} elseif ($args[0] -eq "rm" -or $args[0] -eq "remove") {
  & "$basedir/omm.exe" "$calldir" $args "--mango-rm"
} elseif ($args[0] -eq "wipe") {
  & "$basedir/omm.exe" "$calldir" $args "--mango-wipe"
}
