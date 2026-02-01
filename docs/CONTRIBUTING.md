# Contributing

Thank you for your interest in contributing to Kavach!

## Development Setup

```bash
# Clone repository
git clone https://github.com/Wankhede-Brothers/kavach-rs.git
cd kavach-rs

# Build locally
just build

# Run tests
just test

# Install
just install
```

## Code Style

- **Rust:** Standard `cargo fmt` formatting
- **Linting:** Code must pass `cargo clippy`
- **Errors:** Use `anyhow` for error handling with context
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

1. Create gate file in `crates/kavach-cli/src/commands/gates/`

2. Implement the gate command using clap `Args` and the `hook` module

3. Register in `mod.rs`

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
just test

# Run with verbose output
cd crates/kavach-cli && cargo test -- --nocapture
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
5. Run `just fmt` to format code
6. Submit a pull request

## Pull Request Checklist

- [ ] Tests pass (`just test`)
- [ ] Code is formatted (`just fmt`)
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

- Issues: [GitHub Issues](https://github.com/Wankhede-Brothers/kavach-rs/issues)
- Documentation: [docs/](../docs/)
