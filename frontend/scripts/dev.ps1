$ErrorActionPreference = 'Stop'

$frontendDir = Split-Path -Parent $PSScriptRoot
Set-Location $frontendDir

npm.cmd run dev

