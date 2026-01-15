# File Structure Patterns - SP/3.0
# Lazy-loaded for project structure rules

STRUCTURE:UNIVERSAL
  domain/
    models/
      entity/
        types.{ext} (30 lines)
        impl.{ext} (50 lines)
      value_objects/
        id.{ext} (20 lines)
    services/
      usecase/
        handler.{ext} (50 lines)
        types.{ext} (30 lines)
    shared/
      utils/
        strings.{ext}
        errors.{ext}
      types/
        common.{ext}

STRUCTURE:GO
  cmd/{app}/main.go
  internal/
    domain/{entity}/
      types.go, impl.go
    service/{usecase}/
      handler.go, types.go
  pkg/shared/
    utils/, types/

STRUCTURE:RUST
  src/
    domain/{entity}/
      mod.rs, types.rs, impl.rs
    service/{usecase}/
      mod.rs, handler.rs
    shared/
      utils/, types/
  Cargo.toml

STRUCTURE:TYPESCRIPT
  src/
    domain/{entity}/
      types.ts, service.ts
    components/{feature}/
      Component.tsx, types.ts
    lib/
      utils/, types/
  package.json

STRUCTURE:DEPTH
  min: 3 levels (src/domain/entity/)
  avg: 5 levels (src/domain/entity/models/types.ts)
  max: 7 levels (avoid deeper)
