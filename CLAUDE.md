# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 🚨 CRITICAL: Required Documentation Check

### Absolute Reference Documents
- **@SPECIFICATION.md** - Absolute specification standard
- **@DEVELOPMENT.md** - Absolute development flow standard

### Most Important Policy
**MUST check before starting any work:**
1. Check **@DEVELOPMENT.md** for action-specific references
2. Read specified specialized documents (`docs/spec/` or `docs/dev/`)
3. Execute procedures strictly

**Even if these documents are included in shared memory, always verify the content before starting work.**

## 🎯 Project Overview

**iwakero-gamebook-mapping** - A gamebook playing assistant tool for tracking progress through paragraphs and visualizing the journey.

### Key Features
- Record paragraph visits and choices
- Generate flow diagrams and 2D maps
- Rich terminal UI with PTerm
- Support multiple gamebooks

### Technical Stack
- **Language**: Go (Golang)
- **UI**: PTerm (rich terminal UI)
- **Data**: Markdown format
- **Platform**: Cross-platform (Windows/macOS/Linux)

## 🔧 Development Principles

### Must Follow
1. **TDD (Test-Driven Development)** - Write tests first, always
2. **Issue-Driven Development** - All work tracked via GitHub Issues
3. **Clean Architecture** - 4-layer separation (see `docs/dev/ai-optimization-strategy.md`)
4. **Update CLAUDE.md in every PR**

### AI Optimization
- Use Sub-Issues to limit scope and reduce token consumption
- Read only relevant layer documentation when implementing
- Batch tool calls for better performance

## 🚀 Quick Start

```bash
# Interactive mode (recommended)
./gamebook

# CLI commands
./gamebook new "GameTitle"      # Create new gamebook
./gamebook add 1 "Description"  # Add paragraph
./gamebook show                 # Display current state
```

## 📂 Project Structure

```
├── cmd/gamebook/         # Frameworks Layer (UI, CLI)
├── internal/
│   ├── domain/          # Entities Layer (pure business logic)
│   ├── usecase/         # Usecase Layer (application logic)
│   ├── infrastructure/  # Interface Adapters (repository, presenter)
│   └── interface/       # Interface Adapters (controllers)
├── docs/
│   ├── dev/            # Development guidelines (6 files)
│   └── spec/           # Specifications (4 layers)
└── data/               # Game data (Markdown files)
```

## 🛠️ AI Tools

### Project Management Script
```bash
./scripts/claude-project-manager.sh feature "FeatureName" "Description" "v1.0.0"
./scripts/claude-project-manager.sh quality  # Run quality checks
```

### GitHub Templates
- `gh issue create --template ai_feature.md`
- `gh issue create --template ai_sub_issue.md`
- `gh issue create --template ai_bug_report.md`

## 📈 Current Status

- **Version**: v0.2.9
- **Latest**: Issue #93 (Interactive mode logging conflict resolution)
- **Completed**: Added automatic log output control for interactive mode to prevent UI corruption
- **Previous**: Issue #92 (CLI logging integration completion)
- **History**: See `CHANGELOG.md`

## 🔄 Update History

- 2025-07-07: Issue #93 completed - Interactive mode logging conflict resolution with dynamic output switching and level filtering
- 2025-07-07: Issue #92 completed - CLI logging system integration (all commands now use CLIExecutor)  
- 2025-07-07: Logger integration in development flow - Updated DEVELOPMENT.md and related docs with Logger usage guidelines
- 2025-07-07: Development flow improvement - Added mandatory logging verification procedures to prevent false "verification complete" reports
- 2025-07-06: Document structure optimization - Reduced token consumption by splitting documents
- Earlier updates: See `CHANGELOG.md`

---
**Remember**: Always start by checking DEVELOPMENT.md to know which documents to read for your task!