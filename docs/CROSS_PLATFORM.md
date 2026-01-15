# Cross-Platform Support

Kavach supports **Linux**, **macOS**, and **Windows** with automatic path detection.

## Directory Locations

### Linux

| Component | Path |
|-----------|------|
| Binary | `~/.local/bin/kavach` |
| Memory Bank | `~/.local/share/shared-ai/memory/` |
| Config | `~/.claude/` |
| System Prompt | `~/.claude/CLAUDE.md` |

### macOS

| Component | Path |
|-----------|------|
| Binary | `~/.local/bin/kavach` |
| Memory Bank | `~/Library/Application Support/shared-ai/memory/` |
| Config | `~/.claude/` |
| System Prompt | `~/.claude/CLAUDE.md` |

### Windows

| Component | Path |
|-----------|------|
| Binary | `%USERPROFILE%\.local\bin\kavach.exe` |
| Memory Bank | `%APPDATA%\shared-ai\memory\` |
| Config | `%USERPROFILE%\.claude\` |
| System Prompt | `%USERPROFILE%\.claude\CLAUDE.md` |

## Installation

### Linux/macOS

```bash
# One-line install
curl -fsSL https://raw.githubusercontent.com/Wankhede-Brothers/kavach-go/main/install/install.sh | bash

# Or from source
git clone https://github.com/Wankhede-Brothers/kavach-go.git
cd kavach-go
go build -o kavach ./cmd/kavach
cp kavach ~/.local/bin/

# Add to PATH (add to ~/.bashrc, ~/.zshrc, etc.)
export PATH="$HOME/.local/bin:$PATH"
```

### Windows (PowerShell)

```powershell
# One-line install
irm https://raw.githubusercontent.com/Wankhede-Brothers/kavach-go/main/install/install.ps1 | iex

# Or from source
git clone https://github.com/Wankhede-Brothers/kavach-go.git
cd kavach-go
go build -o kavach.exe .\cmd\kavach
Copy-Item kavach.exe "$env:USERPROFILE\.local\bin\"

# Add to PATH (run in PowerShell as Admin)
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$env:USERPROFILE\.local\bin",
    "User"
)
```

### Windows (CMD)

```cmd
REM Clone repository
git clone https://github.com/Wankhede-Brothers/kavach-go.git
cd kavach-go

REM Build
go build -o kavach.exe .\cmd\kavach

REM Add to PATH (requires restart)
setx PATH "%PATH%;%USERPROFILE%\.local\bin"
```

## Binary Names

| Platform | Binary Name |
|----------|-------------|
| Linux | `kavach` |
| macOS | `kavach` |
| Windows | `kavach.exe` |

## Memory Bank Structure

All platforms use the same Memory Bank structure:

```
memory/
├── GOVERNANCE.toon
├── index.toon
├── volatile.toon
├── decisions/
├── patterns/
├── research/
├── kanban/
├── proposals/
├── roadmaps/
├── graph/
└── STM/
```

## Environment Variables

### Platform Detection

```bash
# Linux/macOS
uname -s  # Linux or Darwin

# Windows
echo %OS%  # Windows_NT
```

### Path Detection

| Platform | Home | Data |
|----------|------|------|
| Linux | `$HOME` | `$XDG_DATA_HOME` or `~/.local/share` |
| macOS | `$HOME` | `~/Library/Application Support` |
| Windows | `%USERPROFILE%` | `%APPDATA%` |

## Claude Code Configuration

### Linux/macOS

```bash
# settings.json location
~/.claude/settings.json

# Hooks use kavach directly (in PATH)
"command": "kavach session init"
```

### Windows

```powershell
# settings.json location
%USERPROFILE%\.claude\settings.json

# Hooks use full path or PATH
"command": "kavach session init"
```

## Shell Scripts vs Batch Files

### Linux/macOS

Shell scripts use `#!/bin/bash`:
```bash
#!/bin/bash
kavach session init
```

### Windows

Batch files (`.cmd`) or PowerShell:
```cmd
@echo off
kavach.exe session init
```

## Testing Cross-Platform

```bash
# Linux/macOS
kavach status
kavach memory bank

# Windows (PowerShell)
kavach.exe status
kavach.exe memory bank

# Windows (CMD)
kavach.exe status
kavach.exe memory bank
```

## Troubleshooting

### "command not found" (Linux/macOS)

```bash
# Check if in PATH
which kavach

# Add to PATH
export PATH="$HOME/.local/bin:$PATH"

# Check permissions
chmod +x ~/.local/bin/kavach
```

### "not recognized" (Windows)

```cmd
# Check if in PATH
where kavach.exe

# Check PATH
echo %PATH%

# Restart terminal after adding to PATH
```

### Path Detection Issues

If paths are not detected correctly:

```bash
# Linux/macOS
echo $HOME
echo $XDG_DATA_HOME

# Windows
echo %USERPROFILE%
echo %APPDATA%
```

## CI/CD

GitHub Actions builds for all platforms:

```yaml
strategy:
  matrix:
    include:
      - goos: linux
        goarch: amd64
      - goos: darwin
        goarch: amd64
      - goos: darwin
        goarch: arm64
      - goos: windows
        goarch: amd64
```

See `.github/workflows/release.yml` for details.
