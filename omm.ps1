#!/usr/bin/env pwsh
$cwd = ("$PWD").replace("\", "/")

chdir "$cwd"
& "$PSScriptRoot\omm_start.exe" $args