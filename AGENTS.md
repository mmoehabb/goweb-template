# AGENTS.md

This file provides AI agents with information about working on this project.

## Quick Start

1. **Read this file** - You're reading it now.
2. **Read memory** - Start with `.agents/memory.md` to understand project context.
3. **Check todos** - Review `.agents/todos.md` to see what needs doing.
4. **Follow guide** - Use `.agents/guide.md` for code conventions.

## Files

| File | Description |
|------|-------------|
| [`.agents/architecture.md`](.agents/architecture.md) | Package structure and dependencies |
| [`.agents/guide.md`](.agents/guide.md) | How to extend the template |
| [`.agents/todos.md`](.agents/todos.md) | Task tracking |
| [`.agents/instructions.md`](.agents/instructions.md) | Session rules |
| [`.agents/memory.md`](.agents/memory.md) | Shared memory |

## How to Use

1. At session start, read `.agents/memory.md` and `.agents/instructions.md`
2. Review `.agents/todos.md` to see pending tasks
3. Reference `.agents/guide.md` and `.agents/architecture.md` as needed
4. After each session, update `.agents/memory.md` with important context

## Conventions

All code conventions are documented in `.agents/guide.md`. Key patterns:

- Handlers go in `handlers/<feature>/`
- Database layer in `db/<feature>/`
- UI components in `ui/forms/`, `ui/components/`, `ui/layouts/`
- Pages in `pages/<page>_templ.go`
- Routes registered in `main.go`