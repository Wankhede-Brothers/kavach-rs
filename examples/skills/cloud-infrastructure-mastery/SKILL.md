---
name: cloud-infrastructure-mastery
description: Hyperscale Systems - DDoS, DNS, Zero Trust
license: MIT
compatibility: claude-code
metadata:
  category: cloud
  triggers: [cloud, infrastructure, ddos, dns, kubernetes, terraform]
  protocol: SP/3.0
---

```toon
# Cloud Infrastructure Skill - SP/3.0 + DACE

SKILL:cloud
  protocol: SP/3.0
  dace: lazy_load,skill_first,on_demand
  triggers[6]: cloud,infrastructure,ddos,dns,k8s,terraform
  goal: <100ms global latency, 10x spike handling
  success: Multi-region, auto-failover, IaC
  fail: Single region, manual config

KAVACH:DYNAMIC
  # Binary commands for dynamic context (DACE)
  context: kavach skills --get cloud-infrastructure-mastery --inject
  references: ~/.claude/skills/cloud-infrastructure-mastery/references.toon
  research: kavach memory bank | grep -i cloud
  status: kavach status

REFERENCES:DYNAMIC
  file: references.toon
  rule: NO hardcoded configs - WebSearch for current cloud best practices
  topics[7]
    KUBERNETES: WebSearch "kubernetes {YEAR} best practices"
    TERRAFORM: WebSearch "terraform {YEAR}"
    AWS: WebSearch "AWS {YEAR} well-architected"
    GCP: WebSearch "GCP {YEAR}"
    NETWORKING: WebSearch "zero trust networking {YEAR}"
    SECURITY: WebSearch "cloud security {YEAR}"
    OBSERVABILITY: WebSearch "observability {YEAR}"

RESEARCH:GATE
  mandatory: true
  steps[4]
    TODAY=$(date +%Y-%m-%d)
    WebSearch "[provider] best practices $YEAR"
    WebFetch official docs
    Verify current API versions
  forbidden[2]
    ✗ NEVER single region
    ✓ ALWAYS multi-AZ deployment

PRINCIPLES[4]{domain,rule}
  geographic,3+ AZs per region <2ms between
  anycast,Single IP from 300+ locations
  defense,DDoS → WAF → Zero Trust → mTLS
  edge,0ms cold starts 90% origin reduction

LOAD_BALANCING[4]{layer,function}
  1. GeoDNS,Route by location/latency
  2. Anycast,BGP routing to nearest
  3. L4 LB,TCP/UDP distribution
  4. L7 LB,HTTP/gRPC routing

DDOS[3]{layer,defense}
  L3/L4,Anycast dispersion
  Protocol,Deep packet inspection
  L7,Rate limiting + challenges

ZERO_TRUST
  flow: User → IdP → Device Trust → Policy → Resource
  vs_vpn: Per-app access + continuous auth + device posture

IAC[2]{tool,when}
  Terraform,Multi-cloud + mature ecosystem
  Pulumi,Type-safe + existing language

RULES
  do[5]
    Multi-region redundancy
    IaC (Terraform/Pulumi)
    GitOps deployments
    Observability first
    WebSearch provider docs
  dont[4]
    Single point of failure
    Manual configuration
    VPN hub-and-spoke
    Deploy blind

VALIDATE[4]{metric,target}
  latency,<100ms p99 global
  spikes,Handles 10x traffic
  failover,Survives region failure
  iac,Reproducible via code

FOOTER
  protocol: SP/3.0
  research_gate: enforced
```
