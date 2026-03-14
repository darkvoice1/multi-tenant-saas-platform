param(
  [ValidateSet('dev','staging','prod')]
  [string]$Env = 'staging'
)

$env:APP_ENV = $Env
Write-Host "[rollout] APP_ENV=$Env"
Write-Host "[rollout] docker compose up -d --build"

docker compose up -d --build