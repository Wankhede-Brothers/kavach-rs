# Kavach Installation Script - Windows (Brahmastra Stack)
# Usage: irm https://raw.githubusercontent.com/Wankhede-Brothers/kavach-rs/main/install/install.ps1 | iex
# Or: irm ... | iex -CLI opencode

$ErrorActionPreference = "Stop"
$Repo = "Wankhede-Brothers/kavach-rs"

function Write-Header {
    Write-Host "`n============================================" -ForegroundColor Blue
    Write-Host "  KAVACH - Brahmastra Stack Installer" -ForegroundColor Blue
    Write-Host "  Protocol: SP/1.0 (Sutra Protocol)" -ForegroundColor Blue
    Write-Host "============================================`n" -ForegroundColor Blue
}

function Detect-CLI {
    param([string]$CLIArg)
    if ($CLIArg) { $script:CLI = $CLIArg }
    elseif (Test-Path "$env:USERPROFILE\.claude") { $script:CLI = "claude-code" }
    elseif ($env:OPENCODE_HOME) { $script:CLI = "opencode" }
    else { $script:CLI = "claude-code" }
    Write-Host "[DETECT] cli: $CLI"
}

function Set-Paths {
    $script:MemoryDir = "$env:LOCALAPPDATA\shared-ai\memory"
    if ($CLI -eq "claude-code") {
        $script:BinDir = "$env:USERPROFILE\.claude\bin"
        $script:SettingsPath = "$env:USERPROFILE\.claude\settings.json"
        $script:PromptPath = "$env:USERPROFILE\.claude\CLAUDE.md"
    } else {
        $script:BinDir = "$env:LOCALAPPDATA\kavach"
        $script:SettingsPath = "$env:LOCALAPPDATA\opencode\settings.json"
        $script:PromptPath = "$env:LOCALAPPDATA\opencode\AGENTS.md"
    }
    Write-Host "[PATHS]"
    Write-Host "  bin: $BinDir\kavach.exe"
    Write-Host "  memory: $MemoryDir"
    Write-Host "  prompt: $PromptPath"
}

function Create-Directories {
    $dirs = @($BinDir, "$MemoryDir\decisions", "$MemoryDir\graph", "$MemoryDir\kanban",
              "$MemoryDir\patterns", "$MemoryDir\proposals", "$MemoryDir\research",
              "$MemoryDir\roadmaps", "$MemoryDir\STM", (Split-Path $SettingsPath))
    foreach ($d in $dirs) { if (!(Test-Path $d)) { New-Item -ItemType Directory -Path $d -Force | Out-Null } }
}

function Download-Binary {
    Write-Host "[DOWNLOAD]"
    $arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
    $url = "https://github.com/$Repo/releases/latest/download/kavach-windows-$arch.exe"
    Write-Host "  url: $url"
    try {
        Invoke-WebRequest -Uri $url -OutFile "$BinDir\kavach.exe" -UseBasicParsing
        Write-Host "  status: ok"
        return $true
    } catch {
        Write-Host "  status: failed (will build)"
        return $false
    }
}

function Build-FromSource {
    Write-Host "[BUILD]"
    if (!(Get-Command cargo -ErrorAction SilentlyContinue)) {
        Write-Host "  error: Rust not installed (install via https://rustup.rs)"; exit 1
    }
    $tmp = Join-Path $env:TEMP "kavach-build"
    if (Test-Path $tmp) { Remove-Item $tmp -Recurse -Force }
    New-Item -ItemType Directory -Path $tmp | Out-Null
    Set-Location $tmp
    git clone --depth 1 "https://github.com/$Repo.git" .
    Set-Location "crates\kavach-cli"
    cargo build --release
    Copy-Item "target\release\kavach.exe" "$BinDir\kavach.exe"
    Set-Location $env:USERPROFILE
    Remove-Item $tmp -Recurse -Force
    Write-Host "  status: ok"
}

function Install-Prompt {
    Write-Host "[PROMPT]"
    $file = if ($CLI -eq "claude-code") { "CLAUDE.md" } else { "AGENTS.md" }
    $url = "https://raw.githubusercontent.com/$Repo/main/configs/windows/$file"
    try {
        Invoke-WebRequest -Uri $url -OutFile $PromptPath -UseBasicParsing
        Write-Host "  installed: $PromptPath"
    } catch { Write-Host "  status: using default" }
}

function Install-Memory {
    Write-Host "[MEMORY]"
    $date = Get-Date -Format "yyyy-MM-dd"
    $index = "# Memory Bank Index - SP/1.0`nINDEX:memory-bank`n  version: 1.0`n  created: $date"
    Set-Content -Path "$MemoryDir\index.toon" -Value $index
    Set-Content -Path "$MemoryDir\volatile.toon" -Value "# Volatile - SP/1.0`nVOLATILE:session`n  created: $date"
    Write-Host "  initialized: $MemoryDir"
}

function Create-Aliases {
    $aliases = @("ceo-gate","bash-sanitizer","read-blocker","session-init","memory-bank")
    foreach ($a in $aliases) {
        Set-Content -Path "$BinDir\$a.cmd" -Value "@echo off`r`n`"$BinDir\kavach.exe`" %*"
    }
    Write-Host "  aliases: created"
}

function Update-Path {
    $p = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($p -notlike "*$BinDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$p;$BinDir", "User")
        $env:Path = "$env:Path;$BinDir"
    }
}

function Write-Success {
    Write-Host "`n[COMPLETE]"
    Write-Host "  binary: $BinDir\kavach.exe"
    Write-Host "  memory: $MemoryDir"
    Write-Host "  prompt: $PromptPath"
    Write-Host "`n[NEXT]"
    Write-Host "  1. Restart terminal"
    Write-Host "  2. Run: kavach status"
    Write-Host "  3. Run: kavach memory bank`n"
}

# Main
param([string]$CLI)
Write-Header
Detect-CLI -CLIArg $CLI
Set-Paths
Create-Directories
if (!(Download-Binary)) { Build-FromSource }
Install-Prompt
Install-Memory
Create-Aliases
Update-Path
Write-Success
