---
name: oym-new-project
description: Scaffold a new On Your Marks project. Use when starting a new repo, Craft site, Craft plugin, api-simple service, Symfony API, NestJS service, or Svelte/SvelteKit app from scratch, or when the user asks to "set up a new X". Selects the canonical boilerplate and wires PHPStan, tests, AGENTS.md, and CI from day one.
---

# OYM new project

New projects always run in `new` mode: the full standard applies from the first commit. No baselines, no legacy allowances.

## Step 1: pick the base

| Building | Start from | Conventions to read |
|---|---|---|
| Craft CMS site | `oym-craft5-boilerplate` | `php-universal.md` + `php-craft.md` |
| Craft plugin | mirror `craft-cacheable` layout | `php-universal.md` + `php-craft-plugin.md` |
| High-load endpoint in a Craft site | `oym/api-simple` package | `php-api-simple.md` |
| Symfony API | mirror `nevobo-trainersplatform-api` | `php-universal.md` + `php-symfony.md` |
| NestJS service | mirror `unisport-agent` | `nestjs.md` |
| SvelteKit app | mirror `oym-chant-fe` (forms/CSP: `oym-paddock-fe`) | `svelte-universal.md` + `sveltekit.md` |
| Frontend inside a Craft site | `oym-vite-boilerplate` | `svelte-universal.md` + `svelte-vite-craft.md` |

Conventions live at `${CLAUDE_PLUGIN_ROOT}/conventions/`, templates at `${CLAUDE_PLUGIN_ROOT}/templates/`.

## Step 2: checklist (all stacks)

1. Scaffold from the base above. Never hand-roll structure.
2. `AGENTS.md` from `templates/AGENTS.template.md` with the relevant stack blocks, header `mode=new`. Vendor the cards: `node scripts/sync-conventions.mjs --stacks <stacks> --target .` (from a checkout of this repo, or `npx github:onyourmarks-agency/oym-development-conventions`).
3. `CLAUDE.md` from `templates/CLAUDE.template.md` (one line: `@AGENTS.md`).
4. Quality gates before feature work:
   - PHP: `templates/php/phpstan.neon.dist` (level max) + `templates/php/phpunit.xml.dist`; phpcs OYM ruleset from the boilerplate; require-dev `phpstan/phpstan` (+ `phpstan/phpstan-symfony`/`-doctrine` for Symfony).
   - NestJS: full `strict: true` (`templates/nestjs/tsconfig.strict.delta.jsonc`), `no-explicit-any: error` (`templates/nestjs/eslint.any-rules.snippet.mjs`), Jest wired.
   - Svelte: Vitest config from `templates/svelte/` (pick the family variant), husky + lint-staged, `frontend-cs` aggregate script.
5. GitLab CI from the `oym-gitlab-deploy-templates` catalog: `@4.x` components for PHP/Craft (build-vendors, test-phpcs, deploy), `@1.x` components for SvelteKit/Nest Docker apps (build-image, tag-environment). Copy the include block from the nearest comparable repo.
6. Conventional Commits from the first commit (`git cz`); node projects pin `engines.node: "24"` + `packageManager` (pnpm); PHP projects point composer at the Repman registry and include `oym/sentry-logger`.
7. Keep a Changelog `CHANGELOG.md` for anything published (plugins, packages).

## Step 3: verify

Run every gate once before the first feature commit: phpcs + phpstan + phpunit green, or `frontend-cs` green, or lint + test green. A new repo that starts red stays red.
