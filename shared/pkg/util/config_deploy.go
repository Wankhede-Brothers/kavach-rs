// Package util provides config deployment utilities.
// config_deploy.go: Ensures TOON config files are deployed to ~/.config/kavach/.
// Called by: session init on every session start.
package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureConfigDir creates ~/.config/kavach/ and deploys default configs if missing.
// Returns the config directory path.
func EnsureConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".config", "kavach")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return configDir, nil
}

// DeployDefaultConfig writes a config file if it doesn't exist.
func DeployDefaultConfig(configDir, filename, content string) error {
	path := filepath.Join(configDir, filename)
	if FileExists(path) {
		return nil // Already deployed
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// DeployAllConfigs deploys all TOON config files to ~/.config/kavach/.
// Only writes files that don't already exist (won't overwrite user customizations).
func DeployAllConfigs() error {
	configDir, err := EnsureConfigDir()
	if err != nil {
		return fmt.Errorf("config dir: %w", err)
	}

	configs := defaultConfigs()
	for filename, content := range configs {
		if err := DeployDefaultConfig(configDir, filename, content); err != nil {
			return fmt.Errorf("deploy %s: %w", filename, err)
		}
	}
	return nil
}

// defaultConfigs returns minimal default TOON configs for core functionality.
// Full configs should be deployed from the repo's config/ directory.
func defaultConfigs() map[string]string {
	return map[string]string{
		"nlu-patterns.toon": `# NLU Intent Patterns - SP/1.0
[DEBUG:PATTERNS]
keywords: fix,bug,error,crash,broken,not working,fail,debug,issue

[PERFORMANCE:PATTERNS]
keywords: slow,fast,optimize,performance,speed,latency,cache

[REFACTOR:PATTERNS]
keywords: refactor,clean,simplify,restructure,reorganize,modular

[RESEARCH:PATTERNS]
keywords: research,explore,investigate,find,search,look up,what is

[DOCS:PATTERNS]
keywords: document,readme,docs,comment,explain,describe

[AUDIT:PATTERNS]
keywords: audit,review,check,verify,inspect,scan,security audit

[IMPLEMENT:PATTERNS]
keywords: implement,create,build,add,make,write,develop,setup

[SECURITY:DOMAIN]
keywords: auth,login,password,token,encrypt,ssl,tls,oauth,jwt,xss,csrf,injection

[FRONTEND:DOMAIN]
keywords: ui,ux,css,html,react,vue,angular,tailwind,component,layout,responsive

[DATABASE:DOMAIN]
keywords: database,sql,query,migration,schema,index,table,postgresql,mysql

[INFRA:DOMAIN]
keywords: deploy,docker,kubernetes,terraform,cloudflare,aws,gcp,azure,ci,cd

[DEVOPS:DOMAIN]
keywords: pipeline,monitoring,logging,metrics,alert,grafana,prometheus

[TESTING:DOMAIN]
keywords: test,spec,assert,mock,coverage,unit test,integration test

[API:DOMAIN]
keywords: api,endpoint,route,handler,rest,graphql,grpc,middleware
`,
		"agent-mappings.toon": `# Agent Mappings - SP/1.0
[VALID:AGENTS]
nlu-intent-analyzer
ceo
research-director
backend-engineer
frontend-engineer
database-engineer
devops-engineer
security-engineer
qa-lead
aegis-guardian
code-reviewer
Explore
Plan
general-purpose
Bash

[ENGINEERS:LIST]
backend-engineer
frontend-engineer
database-engineer
devops-engineer
security-engineer
qa-lead
`,
		"skill-patterns.toon": `# Skill Patterns - SP/1.0
[debug-like-expert]
priority: 1
keywords: debug,fix,error,crash,bug,broken

[security]
priority: 2
keywords: auth,security,vulnerability,xss,csrf,injection

[frontend]
priority: 3
keywords: ui,ux,css,component,layout,responsive

[sql]
priority: 4
keywords: database,query,migration,schema,index

[rust]
priority: 5
keywords: rust,cargo,crate,lifetime,borrow

[api-design]
priority: 6
keywords: api,endpoint,rest,graphql,grpc

[testing]
priority: 7
keywords: test,spec,coverage,mock,assert

[arch]
priority: 8
keywords: architecture,design,system,scale,pattern

[dsa]
priority: 9
keywords: algorithm,data structure,sort,search,tree,graph

[cloud-infrastructure-mastery]
priority: 10
keywords: deploy,docker,kubernetes,terraform,cloudflare,aws
`,
	}
}
