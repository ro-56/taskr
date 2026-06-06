$ErrorActionPreference = 'Stop'

$Repo = 'ro-56/taskr'
$InstallDir = "$env:LOCALAPPDATA\taskr"

# Determine version
if ($env:TASKR_VERSION) {
    $Version = $env:TASKR_VERSION
} else {
    $Release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
    $Version = $Release.tag_name
}

$VersionNum = $Version.TrimStart('v')
$Archive = "taskr_${VersionNum}_windows_amd64.zip"
$BaseUrl = "https://github.com/$Repo/releases/download/$Version"

$TmpDir = Join-Path $env:TEMP ([System.IO.Path]::GetRandomFileName())
New-Item -ItemType Directory -Path $TmpDir | Out-Null

try {
    Write-Host "Downloading taskr $Version (windows/amd64)..."
    Invoke-WebRequest "$BaseUrl/$Archive" -OutFile "$TmpDir\$Archive"
    Invoke-WebRequest "$BaseUrl/checksums.txt" -OutFile "$TmpDir\checksums.txt"

    # Verify SHA256
    $ChecksumLine = Get-Content "$TmpDir\checksums.txt" | Where-Object { $_ -match [regex]::Escape($Archive) }
    if (-not $ChecksumLine) {
        Write-Error "Could not find checksum for $Archive in checksums.txt"
        exit 1
    }
    $Expected = ($ChecksumLine -split '\s+')[0].ToLower()
    $Actual = (Get-FileHash "$TmpDir\$Archive" -Algorithm SHA256).Hash.ToLower()

    if ($Actual -ne $Expected) {
        Write-Error "SHA256 mismatch for ${Archive}:`n  expected: $Expected`n  got:      $Actual"
        exit 1
    }
    Write-Host "Checksum verified."

    # Extract and install
    Expand-Archive "$TmpDir\$Archive" -DestinationPath $TmpDir -Force
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Move-Item "$TmpDir\taskr.exe" "$InstallDir\taskr.exe" -Force
    Write-Host "Installed: $InstallDir\taskr.exe"

    # Add to user PATH via registry
    $RegPath = 'HKCU:\Environment'
    $CurrentPath = (Get-ItemProperty -Path $RegPath -Name PATH -ErrorAction SilentlyContinue).PATH

    if ($CurrentPath -and $CurrentPath -like "*$InstallDir*") {
        Write-Host "$InstallDir is already in your user PATH."
    } else {
        $NewPath = if ($CurrentPath) { "$CurrentPath;$InstallDir" } else { $InstallDir }
        Set-ItemProperty -Path $RegPath -Name PATH -Value $NewPath
        Write-Host "Added $InstallDir to user PATH via registry."
        Write-Host "Restart your terminal for the PATH change to take effect."
    }
} finally {
    Remove-Item -Recurse -Force $TmpDir -ErrorAction SilentlyContinue
}
