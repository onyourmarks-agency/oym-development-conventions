# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-08-17

### Added

- Convention documents for all stacks: shared, php-universal, php-craft, php-craft-plugin, php-api-simple, php-symfony, nestjs, svelte-universal, svelte-vite-craft, sveltekit.
- Claude Code plugin with four skills: oym-php, oym-nestjs, oym-svelte, oym-new-project.
- `sync-conventions.mjs` for vendoring convention cards into product repo `AGENTS.md` files, with drift checking.
- Templates: PHPStan (new and existing-with-baseline), PHPUnit, NestJS strictness deltas, Vitest configs per frontend family, `AGENTS.md`/`CLAUDE.md` starters.
- CI validation of card blocks and sync script self-test.
