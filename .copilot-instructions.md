# GitHub Copilot Instructions for binarytime

## Project Overview
This repository contains both Go and TypeScript/NPM implementations of the `binarytime` library, which provides high-precision time calculations using 128-bit fixed-point arithmetic.

## Repository Structure
- **Go Implementation**: Original library in `/pkg/`, `/cmd/`, `/tools/`
- **TypeScript/NPM Implementation**: New port in `/src/`, `/dist/`
- **Documentation**: `/docs/` and `/html/`

## Technology-Specific Instructions

### For TypeScript/NPM Development
See: [`.copilot-instructions-npm.md`](./.copilot-instructions-npm.md)

### For Go Development  
See: [`.copilot-instructions-go.md`](./.copilot-instructions-go.md)

## General Development Guidelines
- **Incremental Development**: Don't try to do everything at once
- **Step-by-Step Approval**: Explain strategy and get approval at each step
- **Ask Clarifying Questions**: When appropriate, especially for design decisions
- **Documentation**: Keep relevant docs updated with process and choices
- **CLI Documentation**: Use and document appropriate command-line tool usage

## Core Library Concepts
- **Fixed128**: 128-bit fixed-point arithmetic (Go: big.Int, TS: BigInt)
- **Date/BinaryDate**: High-precision timestamp representation
- **Duration/BinaryDuration**: High-precision time span representation
- **Constants**: `DAY_NANOSECONDS = 86_400_000_000_000` (Go: dayNs, TS: DAY_NANOSECONDS)

## Current Focus
The primary development focus is on the **TypeScript/NPM implementation** - porting the Go library to JavaScript/TypeScript for web and Node.js use.