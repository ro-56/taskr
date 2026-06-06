$ErrorActionPreference = 'Stop'

$InstallDir = "$env:LOCALAPPDATA\taskr"

if (Test-Path $InstallDir) {
    Remove-Item -Recurse -Force $InstallDir
    Write-Host "Removed: $InstallDir"
} else {
    Write-Host "taskr not found at $InstallDir"
}

# Remove from user PATH via registry
$RegPath = 'HKCU:\Environment'
$CurrentPath = (Get-ItemProperty -Path $RegPath -Name PATH -ErrorAction SilentlyContinue).PATH

if ($CurrentPath -and $CurrentPath -like "*$InstallDir*") {
    $NewPath = ($CurrentPath -split ';' | Where-Object { $_ -ne $InstallDir }) -join ';'
    Set-ItemProperty -Path $RegPath -Name PATH -Value $NewPath
    Write-Host "Removed $InstallDir from user PATH."
    Write-Host "Restart your terminal for the change to take effect."
} else {
    Write-Host "$InstallDir was not in user PATH."
}
