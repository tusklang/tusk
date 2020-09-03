#!/usr/bin/env pwsh
$cwd = ("$PWD").replace("\", "/")
& "$PSScriptRoot\omm_start.exe" $args -cwd="$cwd"