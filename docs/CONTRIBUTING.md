# Contributing

Thank you for your interest in contributing to Kavach!

## Development Setup

```bash
# Clone repository
git clone https://github.com/Wankhede-Brothers/kavach-go.git
cd kavach-go

# Build locally
go build -o kavach ./cmd/kavach

# Run tests
go test ./...

# Install
cp kavach ~/.local/bin/
```

## Code Style

- **Go:** Standard `go fmt` formatting
- **Comments:** Exported functions must have godoc comments
- **Errors:** Always wrap errors with context: `fmt.Errorf("...: %w", err)`
- **Logging:** Use stderr for debug output, stdout for decision JSON
- **DACE:** Max 100 lines per file, single responsibility

## Project Structure

```
kavach/
├── cmd/kavach/           # Main binary
│   └── internal/
│       └── commands/     # Cobra command tree
│           ├── gates/    # Gate commands
│           ├── memory/   # Memory commands
│           ├── session/  # Session commands
│           ├── agents/   # Agent commands
│           └── skills/   # Skill commands
├── shared/               # Shared packages
│   └── pkg/
│       ├── hook/         # Hook I/O utilities
│       ├── toon/         # TOON parser
│       ├── patterns/     # Pattern matching
│       └── util/         # Utilities
├── examples/             # Example agents and skills
├── configs/              # Platform configurations
├── install/              # Installation scripts
└── docs/                 # Documentation
```

## Adding a New Gate

1. Create gate file in `cmd/kavach/internal/commands/gates/`:
   ```bash
   touch cmd/kavach/internal/commands/gates/my-gate.go
   ```

2. Implement gate:
   ```go
   package gates

   import (
       "github.com/Wankhede-Brothers/kavach-go/shared/pkg/hook"
       "github.com/spf13/cobra"
   )

   var myGateCmd = &cobra.Command{
       Use:   "my-gate",
       Short: "My custom gate",
       RunE: func(cmd *cobra.Command, args []string) error {
           input, err := hook.ReadInput()
           if err != nil {
               return hook.Block("Failed to read input")
           }
           // Gate logic here
           return hook.Approve()
       },
   }
   ```

3. Register in `register.go`:
   ```go
   gatesCmd.AddCommand(myGateCmd)
   ```

## Adding a New Skill

1. Create skill directory in `examples/skills/`:
   ```
   examples/skills/my-skill/
   ├── SKILL.md           # Skill definition
   └── references.toon    # Dynamic WebSearch queries
   ```

2. SKILL.md format:
   ```markdown
   ---
   name: my-skill
   description: What this skill does
   trigger: /my-skill
   ---

   # My Skill

   SKILL:my-skill
     protocol: SP/1.0
     dace: lazy_load,skill_first,on_demand

   KAVACH:DYNAMIC
     context: kavach skills --get my-skill --inject
     status: kavach status
   ```

3. Register in `builtin.go` (optional for built-in skills)

## Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./cmd/kavach/internal/commands/gates/...

# Run with verbose output
go test -v ./...

# Test coverage
go test -cover ./...
```

## Manual Testing

```bash
# Test session init
kavach session init

# Test gates
echo '{"tool_name":"Write","tool_input":{"file_path":"test.go"}}' | kavach gates enforcer --hook

# Test memory
kavach memory bank
```

## Submitting Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Add tests for new functionality
5. Run `go fmt ./...` to format code
6. Submit a pull request

## Pull Request Checklist

- [ ] Tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Documentation updated if needed
- [ ] README updated if adding new command
- [ ] Follows DACE principles (max 100 lines per file)

## Release Process

Releases are automated via GitHub Actions:

1. Push a tag: `git tag v0.1.0 && git push origin v0.1.0`
2. GitHub Actions builds binaries for all platforms (Linux, macOS, Windows)
3. Release is created with compiled binaries
4. Install scripts automatically fetch latest release

## Support

- Issues: [GitHub Issues](https://github.com/Wankhede-Brothers/kavach-go/issues)
- Documentation: [docs/](../docs/)
