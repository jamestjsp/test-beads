# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build -o todo .    # Build binary
./todo add "task"     # Add todo
./todo list           # List todos
./todo done 1         # Complete todo #1
```

## Architecture

Simple Cobra CLI app. Entry point `main.go` → `cmd.Execute()`. All commands in `cmd/root.go`:
- `addCmd`, `listCmd`, `doneCmd` registered via `rootCmd.AddCommand()`
- Todos stored in `~/.todos.json` as JSON array

## Issue Tracking

Uses **bd** (beads). See AGENTS.md for workflow. Key commands:
```bash
bd ready                              # Find work
bd update <id> --status in_progress   # Claim
bd close <id>                         # Complete
bd sync && git push                   # Always push when done
```
