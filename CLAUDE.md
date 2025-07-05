# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a gamebook (choose-your-own-adventure) playing assistant tool called "iwakero-gamebook-mapping". The project aims to help gamebook players track their progress through paragraphs and visualize their journey.

## Project Goals

### Primary Goals
- Create a tool to record paragraph summaries and player choices during gamebook sessions
- Use paragraph numbers as primary keys for all records
- Dynamically generate dungeon maps or flow diagrams based on recorded data

### Technical Preferences
- **Language**: Go (Golang) preferred
- **Platform**: Cross-platform support for Windows, macOS, and Linux
- **Data Format**: Prefer common formats like Markdown (potentially with Mermaid syntax for diagrams)
- **Input Method**: Voice input support desired for easier recording during gameplay

## Current Status

**v0.1.5 Complete**: Interactive mode with rich terminal UI
- ✅ **v0.1.0**: Core functionality (recording, flow diagrams, error detection)
- ✅ **v0.1.5**: Interactive mode implementation complete
  - PTerm-based rich terminal UI with colors, tables, and menus
  - CLI wrapper architecture avoiding code duplication
  - Auto-display of current game status
  - Optimized menu ordering based on usage frequency
  - Test-driven development (TDD) workflow restored
- 🔄 **v0.2.0**: PTerm visualization features planned

## Development Guidelines

When implementing this project:
1. **Always check SPECIFICATION.md and DEVELOPMENT.md before starting any work** - these are the most important documents
2. Follow Issue-driven development using GitHub Issues and PRs
3. Update CLAUDE.md in every PR
4. Use Test-Driven Development (TDD) approach - this is the highest priority
5. Ensure cross-platform compatibility from the start
6. Design the data structure to efficiently use paragraph numbers as primary keys

## Project Management

### Issue-Driven Development
- All work is managed through GitHub Issues
- Issues must start with SPECIFICATION.md confirmation
- Issues are closed via PRs that update CLAUDE.md
- Use Sub-Issues for complex features

### Milestones
- **v0.1.0**: Basic functionality (recording, flow diagrams, error detection) ✅ **Complete**
- **v0.1.5**: Interactive mode with PTerm rich UI ✅ **Complete**
- **v0.2.0**: PTerm visualization features (dynamic flow, interactive 2D maps)
- **v1.0.0**: Voice input functionality

### Development Commands
```bash
# Test execution
go test ./...

# Build CLI
go build -o gamebook ./cmd/gamebook

# Run CLI
./gamebook [command]

# Interactive mode (no arguments):
./gamebook                      # Launch rich PTerm interactive mode

# CLI commands:
./gamebook new "GameTitle"      # Create new gamebook
./gamebook add 1 "Description"  # Add paragraph
./gamebook choice 1 "Go north" 2  # Add choice
./gamebook select 1 1           # Select choice and move
./gamebook show                 # Display current state
```

## Document Maintenance Strategy

This CLAUDE.md file should be updated:
- After major feature implementations
- When project structure changes significantly
- When new technical decisions are made
- Monthly review (use `/init` command if comprehensive update needed)

### Update History
- 2025-07-04: Initial version created
- 2025-07-03: Issue #1 resolved - Fixed critical choice loading bug in markdown_repository.go
- 2025-07-03: v0.1.5 completed - PTerm interactive mode implementation with CLI wrapper architecture

## Next Steps

1. Initialize a Go module for the project
2. Create the basic project structure
3. Design the data model for storing paragraph information
4. Implement core functionality for recording paragraph visits and choices
5. Add visualization capabilities for maps and flow diagrams

## Technology Stack

- **Backend**: Go (Golang)
- **Data Storage**: Markdown files with structured format
- **UI**: PTerm (modern Go terminal UI library)
- **Visualization**: Mermaid.js for diagram generation + PTerm for interactive displays
- **Cross-platform**: Standard Go build for multiple OS
- **Interface**: CLI + Rich interactive mode
- **Future**: Voice input integration planned