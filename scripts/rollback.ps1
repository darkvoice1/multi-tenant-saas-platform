param(
  [Parameter(Mandatory = $true)]
  [string]$Ref,
  [ValidateSet('dev','staging','prod')]
  [string]$Env = 'prod'
)

Write-Host "[rollback] checkout $Ref"

git checkout $Ref
if ($LASTEXITCODE -ne 0) {
  throw "git checkout failed"
}

$env:APP_ENV = $Env
Write-Host "[rollback] APP_ENV=$Env"
Write-Host "[rollback] docker compose up -d --build"

docker compose up -d --build