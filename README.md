# roster

A Go CLI that analyses any repository and scaffolds an agent persona library with orchestration instructions for AI coding harnesses (Claude Code, Codex, Opencode).

Roster handles the scaffolding — the model selects, hydrates, and orchestrates agents.

## Install

```bash
brew install robinojw/roster/roster
```

Or download a binary from [Releases](https://github.com/robinojw/roster/releases).

## Usage

```bash
# Analyse repo and scaffold personas
roster bootstrap

# Analyse a specific path
roster bootstrap --path /path/to/repo

# Preview without writing files
roster bootstrap --dry-run
```

## What it does

`roster bootstrap` analyses your repo and writes:

| Output | Purpose |
|--------|---------|
| `.roster/signals.json` | Detected languages, frameworks, CI, test tools, design system |
| `.roster/personas/` | 11 agent persona files (architect, test, security, etc.) |
| `CLAUDE.md` | Managed section with signals, persona table, and orchestration instructions |
| `AGENTS.md` | Same managed section for Codex/Opencode compatibility |

For repos under 500 files, a file tree is included in the managed section to give the model additional context.

The managed section is delimited by `<!-- roster:start -->` / `<!-- roster:end -->`. Content outside these markers is never touched. Running bootstrap again updates the managed section idempotently.

## Bundled personas

| Persona | Focus |
|---------|-------|
| architect | System architecture, dependency management, API design |
| design | UI components, design system adherence, visual consistency |
| test | Test strategy, coverage, TDD |
| reviewer | Code review, standards enforcement |
| docs | Documentation, README, API docs |
| security | Auth, secrets, OWASP, input validation |
| performance | Profiling, caching, bundle size, query optimisation |
| accessibility | WCAG, screen readers, keyboard navigation |
| devops | CI/CD, Docker, deployment, infrastructure |
| data | Database schema, migrations, data modelling |
| api | REST/GraphQL design, versioning, error handling |

## How it works

1. The analyser walks the repo detecting languages, frameworks, test tools, CI, lint config, design systems, and monorepo patterns
2. The persona library (embedded via `go:embed`) is copied to `.roster/personas/`
3. A managed section is injected into `CLAUDE.md` and `AGENTS.md` containing the signals JSON, persona table, and step-by-step instructions for the model to hydrate and orchestrate agents

## Development

```bash
go build -o roster .
go test ./... -v -race
golangci-lint run
```
