# api-simple conventions

Applies to code under a top-level `api-simple/` directory in a project, and to the `oym/api-simple` package itself. Builds on `php-universal.md`.

api-simple is the in-house micro-router for high-load endpoints that must not boot Craft/Yii: near-zero dependencies, pure PHP, optional FrankenPHP worker mode. Package docs: the `api-simple-docs` MkDocs repo (routes, caching, rate limiting, database, worker mode, responses, configuration).

## Action anatomy

- `[always]` Entrypoint is `public/api-simple/index.php` (verbatim from the boilerplate): defines the path constants, loads dotenv, then `(new \oym\api\simple\Bootstrap())->init(new CraftContext())`.
- `[always]` Actions live in the top-level `api-simple/` directory, namespace `api\simple\<module>`, grouped by feature (`api-simple/content/Statistics.php` handles `?module=content&action=statistics`). The Bootstrap autoloader maps the namespace to the directory; composer PSR-4 does not apply here.
- `[always]` An action extends `\oym\api\simple\route\ApiSimple` and implements `handle(): array`. The returned array is the JSON response.
- `[always]` The class docblock carries the endpoint contract: `@link /{language}/api-simple/?module=content&action=statistics`.
- `[always]` Declare the response shape with an array-shape docblock on `handle()` so PHPStan and the IDE know the payload.
- `[always]` Read request input via the base-class request helpers where available; validate and cast everything from `$_GET` before use.
- `[always]` Use the framework traits for cross-cutting concerns: cache headers (`Cache-Control`, `Surrogate-Key`), Redis rate limiting, database access. Do not hand-roll these per action.
- `[always]` Consume the composer package `oym/api-simple` (^2.x). Do not copy the framework into the repo (only the oldest project does this).

## Style ceiling

This is the most modern PHP in the org; the full php-universal standard applies with no age excuses.

- `[always]` `final` classes, `readonly` value objects, constructor property promotion, `: never` for terminating paths, named arguments.
- `[always]` Keep `handle()` thin: parse input, delegate to small pure collaborators, shape the response. Pure collaborators are the top unit-test target in any Craft repo.
- `[always]` Performance posture: these endpoints exist because they are hot. No heavyweight dependencies in the hot path, no framework bootstrapping, prefer pre-generated JSON cache files on disk over live queries where the project already does so.
- `[always]` PHPStan max applies with essentially zero baseline tolerance. A baseline entry pointing into `api-simple/` is a smell.

<!-- oym-card:begin stack=php-api-simple -->
# OYM api-simple rules

- Actions in top-level `api-simple/<module>/`, namespace `api\simple\<module>`, extend `ApiSimple`, implement `handle(): array`. Returned array is the JSON response.
- Class docblock documents the endpoint: `@link /api-simple/?module=x&action=y`. `handle()` gets an array-shape return docblock.
- Validate and cast all request input; use the framework traits for cache headers, rate limiting, and database access instead of hand-rolling.
- Consume the `oym/api-simple` composer package; never vendor a copy.
- Full modern style, no age excuses: `final`, `readonly`, promoted constructors, `: never`, named arguments.
- `handle()` stays thin; logic in small pure collaborators. Those collaborators are the first things to unit test in any Craft repo.
- Hot path discipline: no heavy deps, no framework boot, prefer pre-generated JSON cache files where the project uses them.
- PHPStan max with zero baseline tolerance here.
<!-- oym-card:end stack=php-api-simple -->

## References

- Framework: `oym-api-simple` (`src/Bootstrap.php`, `src/route/ApiSimple.php`, `src/route/traits/`)
- Modern consumer: `domestique-cycling/api-simple/`
- Docs: `api-simple-docs` repo
