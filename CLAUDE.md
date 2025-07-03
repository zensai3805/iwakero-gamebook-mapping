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

This project is in its initial planning phase. No code has been written yet, and the technology stack needs to be chosen and set up.

## Development Guidelines

When implementing this project:
1. Start by setting up a Go project structure with appropriate module initialization
2. Consider using Markdown files for data storage to maintain simplicity and readability
3. For diagram generation, explore libraries that can work with Mermaid syntax or similar formats
4. Ensure cross-platform compatibility from the start
5. Design the data structure to efficiently use paragraph numbers as primary keys

## Document Maintenance Strategy

This CLAUDE.md file should be updated:
- After major feature implementations
- When project structure changes significantly
- When new technical decisions are made
- Monthly review (use `/init` command if comprehensive update needed)

### Update History
- 2025-07-04: Initial version created
- [Add new entries here as project evolves]

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