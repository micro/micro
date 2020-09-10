$ErrorActionPreference = 'stop'

# GitHub Org and Repo to get archives from
$GitHubOrg="micro"
$GitHubRepo="micro"
$githubHeader = @{}

$MicroInstallDir="c:\micro"
$MicroCliName = "micro.exe"
$MicroCliPath = "${MicroInstallDir}\${MicroCliName}"

if((Get-ExecutionPolicy) -gt 'RemoteSigned' -or (Get-ExecutionPolicy) -eq 'ByPass') {
    Write-Output "PowerShell requires an execution policy of 'RemoteSigned'."
    Write-Output "To make this change please run:"
    Write-Output "'Set-ExecutionPolicy RemoteSigned -scope CurrentUser'"
    break
}

# Change security protocol to support TLS 1.2 / 1.1 / 1.0 - old powershell uses TLS 1.0 as a default protocol
[Net.ServicePointManager]::SecurityProtocol = "tls12, tls11, tls"

Write-Output "Installing micro..."

# Create micro install directory
Write-Output "Creating $MicroInstallDir directory"
New-Item -ErrorAction Ignore -Path $MicroInstallDir -ItemType "directory"
if (!(Test-Path $MicroInstallDir -PathType Container)) {
    throw "Could not create $MicroInstallDir"
}

# Get the list of releases from GitHub
Write-Output "Getting the latest micro release"
$releases = Invoke-RestMethod -Headers $githubHeader -Uri "https://api.github.com/repos/${GitHubOrg}/${GitHubRepo}/releases" -Method Get
if ($releases.Count -eq 0) {
    throw "No releases found in github.com/micro/micro repo"
}

# Filter windows binary and download archive
$windowsAsset = $releases[0].assets | where-object { $_.name -Like "*windows-amd64.zip" }
if (!$windowsAsset) {
    throw "Cannot find the windows micro archive"
}

$zipFilePath = $MicroInstallDir + "\" + $windowsAsset.name
Write-Output "Downloading $zipFilePath ..."

$githubHeader.Accept = "application/octet-stream"
Invoke-WebRequest -Headers $githubHeader -Uri $windowsAsset.url -OutFile $zipFilePath
if (!(Test-Path $zipFilePath -PathType Leaf)) {
    throw "Failed to download micro - $zipFilePath"
}

# Extract micro to ${MicroInstallDir}
Write-Output "Extracting $zipFilePath..."
Expand-Archive -Force -Path $zipFilePath -DestinationPath $MicroInstallDir
if (!(Test-Path $MicroCliPath -PathType Leaf)) {
    throw "Failed to download micro archive - $zipFilePath"
}

# Check the micro version
Invoke-Expression "$MicroCliPath --version"

# Clean up zipfile
Write-Output "Cleaning up $zipFilePath..."
Remove-Item $zipFilePath -Force

# Add MicroInstallDir directory to User Path environment variable
Write-Output "Attempting to add $MicroInstallDir to User Path Environment variable..."
$UserPathEnvionmentVar = [Environment]::GetEnvironmentVariable("PATH", "User")
if($UserPathEnvionmentVar -like "*$MicroInstallDir*") {
    Write-Output "Skipping to add $MicroInstallDir to User Path - $UserPathEnvionmentVar"
} else {
    [System.Environment]::SetEnvironmentVariable("PATH", $UserPathEnvionmentVar + ";$MicroInstallDir", "User")
    $UserPathEnvionmentVar = [Environment]::GetEnvironmentVariable("PATH", "User")
    Write-Output "Added $MicroInstallDir to User Path - $UserPathEnvionmentVar"
}

Write-Output "`r`nmicro has been installed successfully."
Write-Output "To start contributing to micro please visit https://github.com/micro"
Write-Output "Join micro community on slack https://micro.mu/slack"
