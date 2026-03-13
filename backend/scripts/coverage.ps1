param(
  [int]$MinCoverage = 20
)

$ErrorActionPreference = "Stop"

Write-Host "Running tests with coverage..."
$coverFile = "coverage.out"

go test ./... -coverprofile $coverFile -covermode atomic | Write-Host

if (-not (Test-Path $coverFile)) {
  throw "coverage.out not generated"
}

$line = go tool cover -func $coverFile | Select-String -Pattern "total:" | Select-Object -First 1
if (-not $line) {
  throw "Failed to read total coverage"
}

$parts = $line.ToString().Split()
$percentText = $parts[-1].TrimEnd('%')
$percent = [double]$percentText

Write-Host ("Total coverage: {0:N2}%" -f $percent)

if ($percent -lt $MinCoverage) {
  throw ("Coverage {0:N2}% is below threshold {1}%" -f $percent, $MinCoverage)
}

Write-Host "Coverage gate passed."
