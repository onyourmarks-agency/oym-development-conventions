# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-08-17

### Added

- Layered convention documents: `shared`, language layers (`php`, `javascript`), framework layers (`framework-craft`, `framework-craft-plugin`, `framework-api-simple`, `framework-symfony`, `framework-svelte`, `framework-nestjs`), and tool layers (`tool-vite`, `tool-webpack`, `tool-sveltekit`). NestJS composes the shared `javascript` layer with the front-ends.
- `VERSIONS.md` as the single location for version targets; convention rules are version-free.
- Claude Code plugin with four skills: `oym-php`, `oym-nestjs`, `oym-frontend`, `oym-new-project`.
- `sync-conventions.mjs` for vendoring convention cards into product repo `AGENTS.md` files, with drift checking.
- Templates: PHPStan (new and existing-with-baseline), PHPUnit, NestJS strictness deltas, Vitest configs per frontend tool, `AGENTS.md`/`CLAUDE.md` starters.
- CI validation of card blocks and sync script self-test.
