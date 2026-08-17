---
name: oym-php
description: OYM PHP conventions. Use when writing, reviewing, or refactoring PHP in any On Your Marks repository - Craft CMS 5 sites and modules, Craft plugins, api-simple endpoints, or Symfony/API Platform. Covers the phpcs OYM standard, PHPStan and PHPUnit policy, and existing-vs-new project mode.
---

# OYM PHP

## Step 1: detect the flavour

Check in order, then read `${CLAUDE_PLUGIN_ROOT}/conventions/php-universal.md` plus exactly ONE flavour file. Never read all of them.

| Signal | Flavour | Read |
|---|---|---|
| composer.json `"type": "craft-plugin"` | Craft plugin | `conventions/php-craft-plugin.md` |
| composer.json requires `craftcms/cms` | Craft site | `conventions/php-craft.md` |
| File being edited is under a top-level `api-simple/` dir | api-simple action | `conventions/php-api-simple.md` |
| composer.json requires `symfony/framework-bundle` | Symfony | `conventions/php-symfony.md` |

A Craft site repo can contain both `modules/` (Craft rules) and `api-simple/` (api-simple rules); the rules apply per directory.

## Step 2: detect the mode

1. If the repo's `AGENTS.md` convention block header says `mode=existing` or `mode=new`, obey it.
2. Else: a `phpstan-baseline.neon` or substantial committed feature code means `existing`.
3. Default to `existing`. In existing mode the surrounding module outranks every rule below except the phpcs ruleset: mirror local idioms, do not retrofit patterns or tooling.

## Non-negotiables (both modes)

- `declare(strict_types=1);` in every file. Run `vendor/bin/phpcs` before finishing.
- Trailing commas everywhere multiline; arrow functions where possible; `static` closures when `$this` is unused; alphabetized `use` statements.
- Native types on all params, returns (`void`/`never` included), properties, constants (8.3+). Docblocks only for generics, array shapes, and `@property` maps.
- Backed enums for closed value sets. `match` over `switch` for values. Named arguments for boolean params. `JSON_THROW_ON_ERROR` always.
- No `@` suppression, no new `mixed`, no silently swallowed exceptions.
- Comments: non-obvious "why" only. No filler docblocks, no `@author`/`@package`/`@version`. Your default comment verbosity is too high; when in doubt, omit.
- Unit tests only: no DB, no container boot, mock at service boundaries. Add a test alongside any new service/handler with real logic; in existing repos without a suite, create the minimal config from `templates/php/` rather than skipping the test.
- PHPStan level max. New: green day one. Existing: max + generated baseline; the baseline only shrinks; never lower the level.

## New-project extras

`final` + `readonly` + constructor promotion by default; scaffold from `oym-craft5-boilerplate` (sites) and wire `templates/php/phpstan.neon.dist` + `templates/php/phpunit.xml.dist` before feature work. For a full greenfield checklist use the `oym-new-project` skill.
