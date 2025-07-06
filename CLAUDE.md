# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 📚 Documentation Structure

### Core Documents
- **@SPECIFICATION.md** - Tool specifications (→ `docs/spec/` for details)
- **@DEVELOPMENT.md** - Action-based development guide

### Development Workflow
When starting ANY task, always check:
1. DEVELOPMENT.md for which documents to read
2. Relevant specification in `docs/spec/`
3. Implementation guidelines in `docs/dev/`

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
./scripts/claude-project-manager.sh feature "機能名" "説明" "v1.0.0"
./scripts/claude-project-manager.sh quality  # Run quality checks
```

### GitHub Templates
- `gh issue create --template ai_feature.md`
- `gh issue create --template ai_sub_issue.md`
- `gh issue create --template ai_bug_report.md`

## 📈 Current Status

- **Version**: v0.2.5
- **Latest**: Input assistance features
- **In Progress**: Issue #67 (Current location display)
- **History**: See `CHANGELOG.md`

## 🔄 Update History

- 2025-07-06: Document structure optimization - Reduced token consumption by splitting documents
- Earlier updates: See `CHANGELOG.md`

---
**Remember**: Always start by checking DEVELOPMENT.md to know which documents to read for your task!