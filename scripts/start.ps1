param([switch]$Demo, [switch]$Showcase)
$ErrorActionPreference = 'Stop'
Push-Location (Join-Path $PSScriptRoot '..')
try {
    docker compose build
    if ($LASTEXITCODE -ne 0) { throw 'The Talknet build failed.' }
    if ($Demo) {
        docker compose run --rm --no-deps talknet -seed-demo
        if ($LASTEXITCODE -ne 0) { throw 'Sample data setup failed.' }
    }
    if ($Showcase) {
        docker compose run --rm --no-deps talknet -seed-showcase
        if ($LASTEXITCODE -ne 0) { throw 'Showcase sample data setup failed.' }
    }
    docker compose up -d --wait
    if ($LASTEXITCODE -ne 0) { throw 'Talknet did not become healthy.' }
    Write-Host 'Talknet is ready. Default address: http://localhost:8088'
} finally { Pop-Location }
