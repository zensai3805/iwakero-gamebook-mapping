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

**v0.1.0 Progress**: Core functionality implementation complete
- ✅ Basic recording functionality (paragraphs, choices) with TDD approach
- ✅ Markdown file persistence with Mermaid flow diagram generation  
- ✅ CLI interface with session management
- ✅ **Issue #1 RESOLVED**: Choice loading bug fixed - choices now save/load correctly
- ✅ GitHub Sub-Issue management established for milestone tracking
- 🔄 Working on error detection features (Issues #7, #8, #9)

## Development Guidelines

When implementing this project:
1. **Always check SPECIFICATION.md before starting any work** - this is the most important document
2. Follow Issue-driven development using GitHub Issues and PRs
3. Update CLAUDE.md in every PR
4. Use Test-Driven Development (TDD) approach
5. Ensure cross-platform compatibility from the start
6. Design the data structure to efficiently use paragraph numbers as primary keys

## Project Management

### Issue-Driven Development
- All work is managed through GitHub Issues
- Issues must start with SPECIFICATION.md confirmation
- Issues are closed via PRs that update CLAUDE.md
- Use Sub-Issues for complex features

### Milestones
- **v0.1.0**: Basic functionality (recording, flow diagrams, error detection)
- **v0.2.0**: Visualization features (2D text maps)
- **v1.0.0**: Voice input functionality

### Development Commands
```bash
# Test execution
go test ./...

# Build CLI
go build -o gamebook ./cmd/gamebook

# Run CLI
./gamebook [command]

# Available commands:
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

## Next Steps

1. Initialize a Go module for the project
2. Create the basic project structure
3. Design the data model for storing paragraph information
4. Implement core functionality for recording paragraph visits and choices
5. Add visualization capabilities for maps and flow diagrams

## Technology Stack (To be finalized)

- **Backend**: Go (Golang)
- **Data Storage**: Markdown files with structured format
- **Visualization**: Mermaid.js or similar for diagram generation
- **Cross-platform**: Standard Go build for multiple OS
- **Input**: CLI interface with potential voice input integration