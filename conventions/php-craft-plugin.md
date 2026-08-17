# Craft plugin conventions

Applies to repositories with composer `"type": "craft-plugin"`. Builds on `php-universal.md`. Reference: `craft-cacheable` (newest, strictest).

## Structure

- `[always]` Namespace `oym\craft\<handle>\`, autoloaded from `src/`. Composer `extra` block carries `name`, `handle`, `class`, `documentationUrl`, `changelogUrl`.
- `[always]` Layout:

```
src/
  <Plugin>.php            extends craft\base\Plugin
  services/               Yii-component services
  models/Settings.php     extends craft\base\Model, typed public properties
  migrations/             Install.php + m<timestamp>_*.php
  controllers/            web + controllers/console/
  jobs/  enums/  records/  fields/  helpers/
  web/twig/  web/assetbundles/
  templates/  resources/
icon.svg
CHANGELOG.md
```

- `[always]` The plugin class follows the same pattern as a site module class: `public static self $instance`, `@plugins/<handle>` alias, console `controllerNamespace` switch, `@property` service map docblock, event registration in `init()`.
- `[new]` Declare services as a typed class constant array and wire them with `setComponents()` in `init()` (craft-cacheable pattern). Guard required environment variables at boot with a thrown `RuntimeException`.
- `[always]` Settings live in `models/Settings.php` with typed public properties. Services read settings from the model, never scattered `getenv`/`$_ENV` reads.

## Migrations

- `[always]` Schema changes only via `migrations/` (`m<YYMMDD>_<HHMMSS>_<description>.php` plus an install migration). A migration is immutable once a version containing it is tagged.

## Release discipline

Plugins are consumed by many sites; treat every release as public.

- `[always]` Keep a Changelog format, strictly: `## [x.y.z] - YYYY-MM-DD` with `### Added/Changed/Fixed` sections. SemVer: behaviour breaks bump major.
- `[always]` Versions are git tags; the private Repman registry serves the package.
- `[always]` End-user documentation lives in the plugin's separate MkDocs docs repo (`cacheable-docs` pattern). The README stays terse: install line, requirements, changelog pointer, docs link.
- `[new]` PHPStan max with no baseline: plugins are small enough to be clean. Unit tests for services from day one.
- `[existing]` Older plugins (`craft-social-poster`, `craft-seofy`) run PSR-12-only phpcs and Craft 4 era idioms (constants over enums, static `config()` component maps, inline closures). Match the plugin's existing style; upgrades to the current standard are their own task.

<!-- oym-card:begin stack=php-craft-plugin -->
# OYM Craft plugin rules

- Namespace `oym\craft\<handle>\` from `src/`; plugin class with `@property` service map, services in `src/services/`, settings in `src/models/Settings.php` (typed public properties, single source of config).
- Composer `extra` carries name/handle/class/documentationUrl/changelogUrl.
- Schema changes only via `src/migrations/`; migrations are immutable once tagged.
- Keep a Changelog strictly (`## [x.y.z] - YYYY-MM-DD`, Added/Changed/Fixed); SemVer with majors for behaviour breaks; released as git tags via Repman.
- End-user docs in the separate MkDocs docs repo; README stays terse.
- New plugins: PHPStan max without baseline, unit tests for services from day one, typed const component maps, boot-time env guards.
- Existing plugins: match the plugin's era (some are Craft 4, constants instead of enums); modernization is its own task.
<!-- oym-card:end stack=php-craft-plugin -->

## References

- Newest pattern: `craft-cacheable`
- Docs repo pattern: `cacheable-docs`, `social-poster-docs`, `seofy-docs` (MkDocs)
