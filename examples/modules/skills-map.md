# Skills Mapping - SP/3.0
# Maps task types to recommended skills

SKILLS:BACKEND
  "API endpoint"         → /api-design
  "REST API"             → /api-design
  "gRPC service"         → /api-design
  "GraphQL"              → /api-design
  "Database schema"      → /sql
  "SQL query"            → /sql
  "GIN index"            → /sql
  "Full-text search"     → /sql
  "PostgreSQL"           → /sql
  "Rust implementation"  → /rust
  "Axum handler"         → /rust
  "Backend service"      → /rust
  "Rate limiting"        → /rust + /security
  "Session management"   → /rust + /security
  "Token refresh"        → /security
  "Webhook handler"      → /rust + /api-design

SKILLS:FRONTEND
  "React component"      → /frontend
  "UI component"         → /frontend
  "Form validation"      → /frontend
  "Loading states"       → /frontend
  "Skeletons"            → /frontend
  "Toast notifications"  → /frontend
  "Navigation"           → /frontend
  "Responsive design"    → /frontend
  "Dashboard"            → /frontend
  "TypeScript"           → /frontend

SKILLS:INFRASTRUCTURE
  "Deploy"               → /cloud-infrastructure-mastery
  "DigitalOcean"         → /cloud-infrastructure-mastery
  "Kubernetes"           → /cloud-infrastructure-mastery
  "Health endpoints"     → /cloud-infrastructure-mastery + /api-design
  "Cloudflare"           → /cloud-infrastructure-mastery
  "Docker"               → /cloud-infrastructure-mastery
  "CI/CD"                → /cloud-infrastructure-mastery
  "Monitoring"           → /cloud-infrastructure-mastery

SKILLS:SECURITY
  "Authentication"       → /security
  "Authorization"        → /security
  "MFA"                  → /security
  "RSA signature"        → /security
  "Encryption"           → /security
  "OWASP"                → /security
  "Escrow"               → /security + /sql

SKILLS:DATA
  "Data pipeline"        → /high-performance-data-processing
  "Parquet"              → /high-performance-data-processing
  "Arrow"                → /high-performance-data-processing
  "Polars"               → /high-performance-data-processing
  "ETL"                  → /high-performance-data-processing

SKILLS:QUALITY
  "Debug"                → /debug-like-expert
  "Bug fix"              → /debug-like-expert
  "Investigation"        → /debug-like-expert
  "Test"                 → /testing
  "Unit test"            → /testing
  "Integration test"     → /testing
  "Code analysis"        → /heal
  "Refactor"             → /heal

SKILLS:ARCHITECTURE
  "System design"        → /arch
  "Architecture"         → /arch
  "Scalability"          → /arch
  "Performance"          → /arch + /dsa
  "Algorithm"            → /dsa
  "Data structure"       → /dsa
  "Optimization"         → /dsa

TASK_OUTPUT:FORMAT
  ┌──────────┬────────────────────────┬──────────┬─────────────────┐
  │ Priority │ Task                   │ Type     │ Skill           │
  ├──────────┼────────────────────────┼──────────┼─────────────────┤
  │ 🔴       │ GIN Index for search   │ Backend  │ /sql            │
  │ 🔴       │ Deploy Gatus           │ Infra    │ /cloud-infra... │
  │ 🔴       │ Health endpoints       │ Backend  │ /api-design     │
  │ 🟠       │ Redis sessions         │ Backend  │ /rust + /sec    │
  │ 🟠       │ Rate limiting          │ Backend  │ /rust + /sec    │
  │ 🟠       │ Token refresh          │ Backend  │ /security       │
  └──────────┴────────────────────────┴──────────┴─────────────────┘
