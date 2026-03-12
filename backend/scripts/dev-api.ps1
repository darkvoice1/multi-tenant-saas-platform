$ErrorActionPreference = 'Stop'

$backendDir = Split-Path -Parent $PSScriptRoot
Set-Location $backendDir

$cacheDir = Join-Path $backendDir '.cache\gocache'
$modDir = Join-Path $backendDir '.cache\gomod'
$tmpDir = Join-Path $backendDir '.cache\gotmp'
New-Item -ItemType Directory -Force -Path $cacheDir, $modDir, $tmpDir | Out-Null

$env:GOCACHE = $cacheDir
$env:GOMODCACHE = $modDir
$env:GOTMPDIR = $tmpDir
$env:GOPROXY = 'https://goproxy.cn,https://proxy.golang.com.cn,direct'

go run ./cmd/api
