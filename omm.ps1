#!/usr/bin/env pwsh
$cwd = ("$PWD").replace("\", "/")

Set-Location "$cwd"
& "$PSScriptRoot\omm_start.exe" $args