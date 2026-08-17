---
name: oym-new-project
description: Scaffold a new On Your Marks project. Use when starting a new repo, Craft site, Craft plugin, api-simple service, Symfony API, NestJS service, or Svelte/SvelteKit app from scratch, or when the user asks to "set up a new X". Composes the right convention layers and wires PHPStan, tests, AGENTS.md, and CI from day one.
---

# OYM new project

New projects always run in `new` mode: the full standard applies from the first commit. No baselines, no legacy allowances. Current version targets: `${CLAUDE_PLUGIN_ROOT}/VERSIONS.md` (never hardcode versions from memory).

## Step 1: pick the composition

| Building | Scaffold | Convention layers to read |
|---|---|---|
| Craft CMS site | `oym-craft5-boilerplate` repo | shared + php + framework-craft |
| Craft plugin | structure in `framework-craft-plugin.md` | shared + php + framework-craft-plugin |
| High-load endpoint in a Craft site | `oym/api-simple` package | shared + php + framework-api-simple |
| Symfony API | Symfony skeleton + structure in `framework-symfony.md` | shared + php + framework-symfony |
| NestJS service | structure in `framework-nestjs.md` | shared + javascript + framework-nestjs |
| SvelteKit app | structure in `tool-sveltekit.md` | shared + javascript + framework-svelte + tool-sveltekit |
| Frontend inside a Craft site | `oym-vite-boilerplate` repo | shared + javascript + framework-svelte + tool-vite |

Webpack is never chosen for new projects. Conventions live at `${CLAUDE_PLUGIN_ROOT}/conventions/`, templates at `${CLAUDE_PLUGIN_ROOT}/templates/`.

## Step 2: checklist (all stacks)

1. Scaffold from the boilerplate where one exists; otherwise build the structure described in the framework/tool file. Never hand-roll a different structure.
2. Wire conventions with the CLI: `oym-conventions init --mode new --stacks <stacks from the table> --yes` (writes AGENTS.md with vendored cards, CLAUDE.md, and offers the quality-gate templates in one go). If the CLI is not installed: `curl -fsSL https://raw.githubusercontent.com/onyourmarks-agency/oym-development-conventions/main/install.sh | bash`.
3. Fill in the AGENTS.md project sections (what the project is, project-specific rules).
4. Quality gates before feature work:
   - PHP: `templates/php/phpcs.xml` (adjust scan targets) + `templates/php/phpstan.neon.dist` (level max) + `templates/php/phpunit.xml.dist`; require-dev `squizlabs/php_codesniffer` + `slevomat/coding-standard` + `dealerdirect/phpcodesniffer-composer-installer` + `phpstan/phpstan` (+ `phpstan/phpstan-symfony`/`-doctrine` for Symfony).
   - NestJS: full `strict: true` (`templates/nestjs/tsconfig.strict.delta.jsonc`), `no-explicit-any: error` (`templates/nestjs/eslint.any-rules.snippet.mjs`), Jest wired.
   - Frontend: Vitest config from `templates/svelte/` (pick the tool variant), husky + lint-staged, `frontend-cs` aggregate script.
5. GitLab CI from the `oym-gitlab-deploy-templates` catalog: the PHP deploy flow components for PHP/Craft, the docker flow components (`stage-build-image`, `stage-tag-environment`) for Node apps. Copy the include block from the nearest comparable repo; component versions in `VERSIONS.md`.
6. Conventional Commits from the first commit (`git cz`); node projects pin `engines.node` + `packageManager` (pnpm) per `VERSIONS.md`; PHP projects point composer at the Repman registry and include `oym/sentry-logger`.
7. Keep a Changelog `CHANGELOG.md` for anything published (plugins, packages).

## Step 3: verify

Run every gate once before the first feature commit: phpcs + phpstan + phpunit green, or `frontend-cs` green, or lint + test green. A new repo that starts red stays red.
