# Terminal Tools - SP/1.0
# Lazy-loaded when Bash tool is used

TERMINAL:MODERN_TOOLS
  platform: Ubuntu 25.10 Questing Quokka
  principle: Rust/Zig tools > Legacy C tools

RUST_TOOLS{legacy,modern,binary}
  cat,bat,/usr/bin/batcat
  ls,eza,/usr/bin/eza
  find,fd,/usr/bin/fdfind
  grep,rg,/usr/bin/rg
  du,dust,~/.cargo/bin/dust
  ps,procs,~/.cargo/bin/procs
  sed,sd,~/.cargo/bin/sd
  diff,delta,~/.cargo/bin/delta
  cd,zoxide,~/.cargo/bin/zoxide
  top,bottom,~/.cargo/bin/btm

ZIG_TOOLS{tool,binary}
  bun,~/.bun/bin/bun
  zccinfo,~/.local/bin/zccinfo
  zig-cc,/usr/bin/zig

ALIASES{alias,command}
  cat,batcat --paging=never
  ls,eza --icons --group-directories-first
  ll,eza -la --icons --group-directories-first --git
  tree,eza --tree --level=3 --icons
  find,fdfind
  grep,rg --smart-case

RULE: Prefer Rust/Zig tools for better UX and performance
