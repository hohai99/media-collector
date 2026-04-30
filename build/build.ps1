# Media Collector Build Script (Windows)
# Usage: .\build\build.ps1

Write-Host "Building Media Collector..." -ForegroundColor Cyan

wails build -platform windows/amd64

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!" -ForegroundColor Green
    Write-Host "Output: build\bin\Media Collector.exe" -ForegroundColor Yellow
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}
