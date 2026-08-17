# PHP conventions (back-end language layer)

Applies to every PHP file in every OYM repository. Framework files (`framework-craft.md`, `framework-craft-plugin.md`, `framework-api-simple.md`, `framework-symfony.md`) build on this.

## Framework detection

Check in order:

1. `composer.json` has `"type": "craft-plugin"` (or `src/*.php` extends `craft\base\Plugin`): read `framework-craft-plugin.md`.
2. `composer.json` requires `craftcms/cms`: read `framework-craft.md`. If a top-level `api-simple/` directory exists, files under it follow `framework-api-simple.md` instead.
3. `composer.json` requires `symfony/framework-bundle`: read `framework-symfony.md`.

## Coding standard: phpcs OYM ruleset

- `[always]` The law is the `phpcs.xml` "OYM" ruleset: PSR-12 (line length excluded) plus the Slevomat rules below. Run `vendor/bin/phpcs` before finishing any PHP change.
- The ruleset enforces:
  - `declare(strict_types=1);` in every file, one blank line after the open tag.
  - Trailing commas in every multiline call, declaration, closure use, and array.
  - Arrow functions where possible (`RequireArrowFunction`); `static` closures when `$this` is unused (`StaticClosure`).
  - Alphabetically sorted `use` statements.
  - Explicit visibility on class constants.
  - Nullable type for parameters with a `null` default.
  - No useless docblocks, no `Class Foo`/`Interface Foo` filler comments, no `@author`/`@package`/`@version` style annotations.
- `[existing]` Older repos may carry the minimal ruleset (PSR-12 only, ruleset name "TDE"). Follow whatever the repo's `phpcs.xml` enforces; write to the full standard anyway, it is a superset.

## Language and style

- `[always]` Native types on everything: parameters, returns (including `: void` and `: never`), properties, and class constants where the PHP version allows typed constants.
- `[always]` Docblock types only where native types fall short: generics (`@template`, `@param array<string, Team>`), array shapes (`@return array{id: int, name: string}`), `@property` maps for magic accessors.
- `[new]` Default to `final` classes, `readonly` where the PHP version allows, and constructor property promotion. Open a class for extension only with a reason.
- `[existing]` `final`/`readonly`/promotion are rare in older code. Use them in new classes; do not retrofit existing ones.
- `[always]` Backed enums for closed value sets. No class-constant pseudo-enums in new code.
- `[always]` `match` over `switch` where the result is a value.
- `[always]` Named arguments for boolean or ambiguous parameters: `->flushTags($tags, executeImmediately: true)`.
- `[always]` Guard clauses and early returns over nested conditionals.
- `[always]` `json_encode`/`json_decode` always with `JSON_THROW_ON_ERROR`.
- `[always]` No `@` error suppression. No `mixed` in new signatures.
- `[always]` Errors: catch where you can act (log via the project's logger and return a safe default, or rethrow). Never silently swallow without at least a logged reason. Sentry is wired in every project via `oym/sentry-logger`; unhandled exceptions reach it.
- `[always]` Naming: classes PascalCase, methods camelCase verb-first (`getAllTeams`, `createEntities`), constants UPPER_SNAKE, descriptive variable names (no `$i`, `$a`, `$tmp`).

## PHPStan

- `[new]` PHPStan `level: max` from day one, green before the first deploy. Config: `templates/php/phpstan.neon.dist`.
- `[existing]` Introduce PHPStan at `level: max` with a generated baseline (`templates/php/phpstan-existing.neon.dist`, then `vendor/bin/phpstan analyse --generate-baseline`). The baseline is a ratchet: never add to it, shrink it when touching baselined files. Never lower the level instead of baselining.
- `[always]` Inline suppressions use `@phpstan-ignore` with an identifier and a trailing reason comment. Last resort only.
- Craft note: the Craft/Yii API is magic-heavy. Expect a large baseline on existing sites and targeted, reasoned ignores on new ones. Do not respond by lowering the level.

## PHPUnit

- `[always]` Unit tests only (see `shared.md`): no database, no Craft/Yii/Symfony container boot, no HTTP. Mock adapters and element queries at the service boundary.
- `[new]` `tests/Unit/` mirrors the source namespace; one test class per class under test; config from `templates/php/phpunit.xml.dist`.
- `[existing]` Most PHP projects have no tests. Add a unit test when you add a service, handler, or pure helper with real logic. Create the minimal `phpunit.xml.dist` from templates if none exists. Do not scaffold a retroactive suite unprompted.
- Priority order when introducing tests to a Craft repo: api-simple action collaborators first (pure, modern), then module services, then enums/helpers.

<!-- oym-card:begin stack=php -->
# OYM PHP rules (back-end language layer)

- `declare(strict_types=1);` in every file. Run `vendor/bin/phpcs` before finishing.
- Trailing commas in every multiline call/declaration/closure/array. Arrow functions where possible; `static` closures when `$this` is unused. Sorted `use` statements. Explicit constant visibility.
- Native types on everything: params, returns (incl. `void`/`never`), properties, typed constants where the PHP version allows. Docblocks only for generics, array shapes, and `@property` maps.
- New classes: `final` by default, `readonly` and constructor promotion where the PHP version allows. Do not retrofit old classes.
- Backed enums for closed value sets. `match` over `switch` for values. Named arguments for boolean params.
- Guard clauses over nesting. `JSON_THROW_ON_ERROR` always. No `@` suppression, no new `mixed`.
- Never swallow exceptions silently: log and return a safe default, or rethrow.
- No filler docblocks, no `@author`/`@package`/`@version`. Comments only for non-obvious "why".
- PHPStan level max. New projects: green from day one. Existing projects: max plus generated baseline; the baseline only shrinks, never grows; never lower the level.
- `@phpstan-ignore` needs an identifier and a reason. Last resort.
- Unit tests only (PHPUnit, `tests/Unit/` mirroring src): no DB, no container boot, mock at service boundaries. Existing projects: test new logic you add; do not scaffold suites unprompted.
- Naming: methods camelCase verb-first, constants UPPER_SNAKE, descriptive variables.
<!-- oym-card:end stack=php -->

## Examples in the wild

- Modern style throughout: `oym-api-simple/src/`
- Ruleset source: any current repo's `phpcs.xml` (eredivisie, domestique-cycling, craft-cacheable are identical)
- Templates: `templates/php/phpstan.neon.dist`, `templates/php/phpstan-existing.neon.dist`, `templates/php/phpunit.xml.dist`
