#!/usr/bin/env pwsh

$basedir = Split-Path $MyInvocation.MyCommand.Definition -Parent
$calldir = ("$PWD").replace("\", "/") #convert to slash seperated

#call the actual script
& "$basedir/omml.exe" "$calldir/" $args
exit 0
