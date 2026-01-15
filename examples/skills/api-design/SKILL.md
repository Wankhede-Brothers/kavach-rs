---
name: api-design
description: API Design - REST, gRPC, GraphQL
license: MIT
compatibility: claude-code
metadata:
  category: api
  triggers: [api, rest, grpc, graphql, endpoint, pagination]
  protocol: SP/3.0
---

```toon
# API Design Skill - SP/3.0 + DACE

SKILL:api-design
  protocol: SP/3.0
  dace: lazy_load,skill_first,on_demand
  triggers[6]: api,rest,grpc,graphql,endpoint,pagination
  goal: Consistent, versioned APIs
  success: Nouns, keyset pagination, RFC 7807
  fail: Verbs in URLs, OFFSET at scale

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get api-design --inject
  references: ~/.claude/skills/api-design/references.toon
  research: kavach memory bank | grep -i api
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded examples - WebSearch for current patterns
  topics[7]
    REST_API: WebSearch "REST API design {YEAR}"
    GRPC: WebSearch "gRPC {YEAR} best practices"
    GRAPHQL: WebSearch "GraphQL {YEAR}"
    HTTP_STATUS: WebSearch "HTTP status codes {YEAR}"
    PAGINATION: WebSearch "API pagination {YEAR}"
    RATE_LIMITING: WebSearch "API rate limiting {YEAR}"
    AUTHENTICATION: WebSearch "API auth {YEAR}"

RESEARCH:GATE
  mandatory: true
  steps[3]
    TODAY=$(date +%Y-%m-%d)
    WebSearch "API design $YEAR best practices"
    Verify current standards
  forbidden[2]
    ✗ NEVER verb-based URLs
    ✓ ALWAYS nouns (resources)

NAMING
  correct[2]: GET /users/{id}, POST /users
  wrong[2]: GET /getUser, POST /createUser
  rules[3]: Plural nouns, kebab-case, max 3 levels

METHODS[5]{method,semantics}
  GET,Read (idempotent + safe)
  POST,Create (not idempotent)
  PUT,Replace (idempotent)
  PATCH,Partial update
  DELETE,Remove (idempotent)

STATUS_CODES[3]{category,codes}
  2xx,200 Success | 201 Created | 204 No Content
  4xx,400 Bad | 401 Unauth | 403 Forbidden | 404 Not Found | 409 Conflict | 422 Unprocessable | 429 Rate Limited
  5xx,500 Server Error

PAGINATION
  keyset: GET /orders?after=cursor&limit=20
  avoid: OFFSET at scale (O(n))

VERSIONING
  preferred: /api/v1/resource

ERRORS
  format: RFC 7807 (Problem Details)
  required[4]: type,title,status,detail

RULES
  do[4]
    Plural nouns
    Keyset pagination
    RFC 7807 errors
    Version from day 1
  dont[3]
    Verbs in URLs
    200 for errors
    OFFSET large datasets

FOOTER
  protocol: SP/3.0
  research_gate: enforced
```
