$ErrorActionPreference = "Stop"

$Repo = "f24aalam/dbmcp"
$Version = (Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest").tag_name
$InstallDir = "$env:USERPROFILE\bin"

$Arch = if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    "arm64"
} else {
    "amd64"
}

$Binary = "godbmcp_windows_$Arch.exe"
$Url = "https://github.com/$Repo/releases/download/$Version/$Binary"

Write-Host "📦 Installing godbmcp"
Write-Host "⬇️  $Url"

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
Invoke-WebRequest -Uri $Url -OutFile "$InstallDir\godbmcp.exe"

if ($env:PATH -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable(
        "PATH",
        "$env:PATH;$InstallDir",
        [EnvironmentVariableTarget]::User
    )
}

Write-Host "✅ Installed successfully"
Write-Host "Restart terminal and run: godbmcp --help"
