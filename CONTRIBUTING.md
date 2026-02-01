# Contributing to Kavach

Thank you for your interest in contributing to Kavach! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful and constructive. We welcome contributors of all experience levels.

## Getting Started

### Prerequisites

#### Build Requirements

| Tool | Version | Purpose |
|------|---------|---------|
| Rust | stable | Build and test |
| Git | 2.0+ | Version control |

#### Rust CLI Tools (Required)

Kavach enforces modern Rust CLI tools. Install before contributing:

```bash
# Install all required tools
cargo install bat eza fd-find ripgrep

# Or with Homebrew (macOS)
brew install bat eza fd ripgrep

# Or with Scoop (Windows)
scoop install bat eza fd ripgrep
```

| Tool | Replaces | Why Required |
|------|----------|--------------|
| `bat` | `cat` | Syntax highlighting in tests |
| `eza` | `ls` | Directory listings with icons |
| `fd` | `find` | Fast file discovery |
| `rg` | `grep` | Fast pattern searching |

**Note**: Running tests or development commands with legacy tools (`cat`, `ls`, `find`, `grep`) will be blocked by Kavach's bash gate.

### Development Setup

```bash
# 1. Fork the repository on GitHub

# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/kavach-rs.git
cd kavach-rs

# 3. Add upstream remote
git remote add upstream https://github.com/Wankhede-Brothers/kavach-rs.git

# 4. Create a feature branch
git checkout -b feature/my-feature

# 5. Build and test
just build
just test

# 6. Verify binary works
./target/release/kavach status
```

## Project Structure

```
kavach/
├── cmd/kavach/                # Main binary
│   └── internal/commands/     # Cobra command tree
│       ├── gates/             # Enforcement gates
│       ├── memory/            # Memory Bank operations
│       ├── session/           # Session management
│       ├── agents/            # Agent management
│       └── skills/            # Skill management
├── shared/                    # Shared packages
│   └── pkg/
│       ├── hook/              # Hook I/O utilities
│       ├── toon/              # TOON parser
│       ├── patterns/          # Pattern matching
│       ├── agentic/           # Research gate (TABULA_RASA)
│       └── util/              # Utilities
├── examples/                  # Example agents and skills
├── configs/                   # Platform configurations
├── install/                   # Installation scripts
└── docs/                      # Documentation
```

## Code Style

### Rust Standards

- **Formatting**: Run `just fmt` before committing
- **Linting**: Code must pass `just lint` (cargo clippy)
- **Errors**: Use `anyhow` for error handling with context

### DACE Principles

Kavach follows **Dynamic Agentic Context Engineering (DACE)**:

- **Max 100 lines per file** - Split by concern, not line count
- **Single responsibility** - One struct/impl per file
- **Lazy loading** - Load context on-demand
- **No hardcoding** - Use patterns for dynamic detection

## How to Contribute

### Adding a New Gate

1. Create gate file in `crates/kavach-cli/src/commands/gates/`:

2. Implement the gate command using clap `Args` and the `hook` module

3. Register in `mod.rs`

### Adding a New Skill

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

### Adding a New Agent

1. Create agent file in `examples/agents/`:

```markdown
# examples/agents/my-agent.md

---
name: my-agent
level: 1
model: sonnet
domain: my-domain
---

AGENT:my-agent
  level: 1
  model: sonnet
  domain: Description of domain
```

## Testing

### Run Tests

```bash
# All tests
just test

# With verbose output
cd crates/kavach-cli && cargo test -- --nocapture
```

### Manual Testing

```bash
# Test session init
kavach session init

# Test gates (pipe JSON input)
echo '{"tool_name":"Write","tool_input":{"file_path":"test.go"}}' | kavach gates enforcer --hook

# Test memory bank
kavach memory bank
```

## Submitting Changes

### Pull Request Process

1. **Update your fork**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Push your branch**:
   ```bash
   git push origin feature/my-feature
   ```

3. **Create Pull Request** on GitHub

4. **CI Pipeline runs automatically**:
   - Build verification
   - Test suite
   - Lint checks
   - Cross-compilation

5. **Address review feedback**

6. **Merge** after approval

### PR Checklist

Before submitting, ensure:

- [ ] `just test` passes
- [ ] `just fmt` applied
- [ ] `just lint` passes
- [ ] Max 100 lines per file (DACE compliance)
- [ ] Tests added for new functionality
- [ ] Documentation updated if needed
- [ ] Commit messages are descriptive

### Commit Message Format

```
type(scope): short description

Longer description if needed.

Co-Authored-By: Your Name <email@example.com>
```

Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`

Examples:
- `feat(gates): Add new validation gate`
- `fix(memory): Handle concurrent write race condition`
- `docs(readme): Update installation instructions`

## Release Process

Releases are automated via GitHub Actions:

1. Maintainer tags a release:
   ```bash
   git tag v0.2.0
   git push origin v0.2.0
   ```

2. CI builds binaries for all platforms:
   - Linux (amd64, arm64)
   - macOS (Intel, Apple Silicon)
   - Windows (amd64)

3. GitHub Release created with:
   - Compiled binaries
   - SHA256 checksums
   - Auto-generated release notes

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/Wankhede-Brothers/kavach-rs/issues)
- **Documentation**: [docs/](docs/) or [Website](https://wankhedebrothers.com/docs/kavach/)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
