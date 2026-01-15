---
name: security
description: Application Security - OWASP, Auth, Encryption
license: MIT
compatibility: claude-code
metadata:
  category: security
  triggers: [security, auth, owasp, encryption, csrf, xss, injection]
  protocol: SP/3.0
---

```toon
# Security Skill - SP/3.0 + DACE

SKILL:security
  protocol: SP/3.0
  dace: lazy_load,skill_first,on_demand
  triggers[7]: security,auth,owasp,gnap,oauth,csrf,xss
  goal: No OWASP Top 10 vulnerabilities
  success: Defense in depth, parameterized queries
  fail: Bearer tokens high-security, hardcoded secrets

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get security --inject
  references: ~/.claude/skills/security/references.toon
  research: kavach memory bank | grep -i security
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded CVEs - WebSearch for current vulnerabilities
  topics[6]
    OWASP: WebSearch "OWASP top 10 {YEAR}"
    AUTH: WebSearch "authentication {YEAR}"
    ENCRYPTION: WebSearch "encryption {YEAR}"
    INPUT: WebSearch "input validation {YEAR}"
    SECRETS: WebSearch "secrets management {YEAR}"
    CVE: WebSearch "{LIBRARY} CVE {YEAR}"

RESEARCH:GATE
  mandatory: true
  steps[4]
    TODAY=$(date +%Y-%m-%d)
    WebSearch "OWASP Top 10 $YEAR"
    WebSearch "[library] CVE $TODAY"
    Verify current CVEs
  forbidden[3]
    ✗ NEVER use deprecated crypto
    ✗ NEVER roll your own crypto
    ✓ ALWAYS verify current CVEs

AUTH[3]{protocol,when}
  GNAP,High-security + new projects + IoT
  OAuth2.0,Legacy + existing infrastructure
  JWT,Stateless microservices

OWASP[4]{code,mitigation}
  A01 Broken Access,RBAC + deny-by-default
  A02 Crypto Failures,Modern algos + TLS
  A03 Injection,Parameterized queries
  A04-A10,WebSearch "OWASP Top 10 $TODAY"

ATTACKS[5]{attack,defense}
  sql_injection,Parameterized queries
  xss,CSP headers + HTML encoding
  csrf,CSRF tokens + SameSite cookies
  idor,Auth check on EVERY request
  ssrf,URL allowlist + block private IPs

CRYPTO[5]{use,algorithms}
  passwords,Argon2id | bcrypt(12+)
  symmetric,AES-256-GCM | ChaCha20
  asymmetric,Ed25519 | RSA-2048+
  tls,TLS 1.3 only
  avoid,MD5 | SHA1 | DES | RC4

SECRETS
  production: Secret Manager (AWS/GCP/Vault)
  dev: Environment variables
  NEVER: Hardcoded, committed to git

RULES
  do[5]
    Parameterized queries
    Defense in depth
    TLS 1.3 everywhere
    Security headers
    WebSearch CVEs
  dont[4]
    Roll your own crypto
    Trust user input
    Hardcode secrets
    Skip auth checks

VALIDATE
  owasp: WebSearch "OWASP Top 10 $TODAY"
  cve: WebSearch "[library] CVE $TODAY"
  deps: cargo audit / npm audit

FOOTER
  protocol: SP/3.0
  research_gate: enforced
```
