# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- SvelteKit conventions follow SvelteKit 3: configuration through the Vite plugin in `sveltekit.config.js` (no `svelte.config.js`), `tsconfig.json` extending `$app/tsconfig`, own modules as `#lib/...` with file extensions, environment variables declared in `src/env.ts` and read from `$app/env/private`, and hook types from `@sveltejs/kit/hooks`. SvelteKit 2 is no longer supported.
- `VERSIONS.md`: SvelteKit 3, which needs Vite 8 and TypeScript 6.
- `templates/svelte/vitest.config.sveltekit.ts` passes the shared options to the SvelteKit plugin.

### Added

- `templates/svelte/sveltekit.config.js`: the SvelteKit options, with runes forced and a CSP to start from.
- `oym-seed` skill: gives a Craft site seed data with `oym/craft-seed` (generated content, no copy of staging or production), or repairs the seed when the drift check reports that the content model changed.
- `templates/ddev/seed`: the `ddev seed` command the skill adds to a project.

## [1.2.0] - 2026-11-11

### Added

- Project manifest rules in the shared layer: every repository has a `.project.yml`, it is the source of truth for project metadata, and it is updated in the same change that affects it.
- `templates/project.template.yml` as the starter manifest for new repositories; the `oym-new-project` skill creates it.

## [1.1.0] - 2026-09-02

### Added

- Formatting rule in the shared layer: every file ends with a single trailing newline.

## [1.0.0] - 2026-08-17

### Added

- Layered convention documents: `shared`, language layers (`php`, `javascript`), framework layers (`framework-craft`, `framework-craft-plugin`, `framework-api-simple`, `framework-symfony`, `framework-svelte`, `framework-nestjs`), and tool layers (`tool-vite`, `tool-webpack`, `tool-sveltekit`). NestJS composes the shared `javascript` layer with the front-ends.
- `VERSIONS.md` as the single location for version targets; convention rules are version-free.
- Claude Code plugin with four skills: `oym-php`, `oym-nestjs`, `oym-frontend`, `oym-new-project`.
- The `oym-conventions` Go CLI (`init`, `sync`, `check`, `update`, `version`): detects a project's stacks, wires AGENTS.md/CLAUDE.md/quality-gate templates, and keeps the vendored cards current. Rules and templates are embedded per release; installed via `install.sh` from GitHub releases.
- GitLab CI snippet for monthly freshness (`ci/gitlab-conventions-sync.yml`): scheduled check plus opt-in auto-MR sync.
- Templates: PHPStan (new and existing-with-baseline), PHPUnit, NestJS strictness deltas, Vitest configs per frontend tool, `AGENTS.md`/`CLAUDE.md` starters.
- CI validation of card blocks and the sync lifecycle as Go tests; GoReleaser release workflow on version tags.
