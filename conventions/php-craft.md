# Craft CMS site conventions

Applies to Craft CMS 5 site repositories (composer requires `craftcms/cms`, not a plugin). Builds on `php-universal.md`. Canonical starting point: `oym-craft5-boilerplate`.

Files under a top-level `api-simple/` directory follow `php-api-simple.md`, not this file.

## Project shape

- `[always]` Site code lives in Yii modules under `modules/`, autoloaded as `"modules\\": "modules/"`. Craft config in `config/`, Twig in `templates/`, frontend in `resources/` (see `svelte-vite-craft.md`).
- `[always]` PHP 8.3 to 8.5 depending on project age; the OYM plugin suite (`oym/craft-cacheable`, `oym/craft-environment`, `oym/craft-utilities`, `oym/sentry-logger`, usually `oym/craft-seofy` and `oym/api-simple`) comes from the private Repman registry.
- `[always]` Module directory layout:

```
modules/<name>module/
  <Name>module.php        module class (entrypoint)
  services/               business logic, Yii components
  services/adapter/       external API adapters behind an interface
  console/controllers/    console commands
  controllers/            web/CP controllers
  jobs/                   craft\queue\BaseJob subclasses
  events/handlers/        invokable event handlers (newer projects)
  enums/                  backed enums
  variables/              Twig variables
  fieldtypes/             custom fields
  records/                ActiveRecord
  web/twig/               Twig extensions
  templates/              module template roots
  utils/                  stateless helpers
```

## Module class

- `[always]` `class Foomodule extends yii\base\Module` with the established boilerplate: `Craft::setAlias('@modules/foomodule', ...)` and the `controllerNamespace` console switch in `__construct`, `static::setInstance($this)`, event registration in `init()`.
- `[always]` Services are Yii components registered on the module (`setComponents()` or a `components` config array) and accessed via magic getters. The class-level docblock is a complete, alphabetized `@property FooService $fooService` map. This map is the one sanctioned "redundant" docblock: it is what makes `$module->fooService` resolve in the IDE and PHPStan. Keep it in sync with the registered components.
- `[always]` Keep `__construct` minimal; registration work belongs in `init()`.

## Events

- `[new]` One invokable handler class per event concern under `events/handlers/`: `Event::on(Cp::class, Cp::EVENT_REGISTER_CP_NAV_ITEMS, new RegisterCpNavItemsHandler())` with `public function __invoke(RegisterCpNavItemsEvent $event): void`.
- `[existing]` Older projects (eredivisie era) register events with inline static closures in the module class. Match whichever style the module already uses. Migrating inline closures to handler classes is a task on its own, not a drive-by.

## Services, jobs, console

- `[always]` Services hold the business logic; controllers and console commands stay thin.
- `[always]` External APIs are wrapped in an adapter under `services/adapter/` behind an interface, so the service is testable without the API.
- `[always]` Long-running or bulk work goes in a queue job (`extends craft\queue\BaseJob`, `execute($queue): void`, `defaultDescription(): string`), never in a request handler.
- `[always]` Console commands are console controllers: `actionFoo(): int` returning a `yii\console\ExitCode` constant, flags via `options()`/`optionAliases()`, shared logic in an abstract base controller when several commands need it.
- `[always]` Element queries are chained and explicit: `Entry::find()->section(...)->site(...)->status(null)->all()`. Never query inside a loop when one query can do it.
- `[always]` Cache invalidation goes through `oym/craft-cacheable` tag flushing where the project uses it; check for a `flushers/` directory.

## Mode guidance

- `[new]` Scaffold from `oym-craft5-boilerplate`. Wire PHPStan max and PHPUnit (`templates/php/`) before feature work. Use handler classes, backed enums, typed constants from the start.
- `[existing]` Craft sites span a decade of idioms. Read the module you are editing before writing: its event style, its service registration, its naming (module class casing is historically inconsistent; match the file you are in). Local idiom wins over this document.

<!-- oym-card:begin stack=php-craft -->
# OYM Craft CMS site rules

- Site code in Yii modules under `modules/`; services/, console/controllers/, jobs/, events/handlers/, enums/, variables/ inside each module.
- Services are Yii components on the module, exposed via a complete alphabetized `@property FooService $fooService` docblock map on the module class. Keep the map in sync; it drives IDE completion.
- Controllers and console commands stay thin; logic lives in services; external APIs behind an adapter interface under `services/adapter/`.
- Long-running work in queue jobs (`craft\queue\BaseJob`), never in request handlers.
- Console actions return `yii\console\ExitCode` constants.
- Newer projects: one invokable handler class per event under `events/handlers/`. Older projects use inline static closures in the module class: match what the module already does.
- Backed enums for statuses and types; no loose strings across service boundaries.
- Element queries explicit and chained; never query in a loop when one query works.
- Cache invalidation via `oym/craft-cacheable` tag flushing where present.
- Files under a top-level `api-simple/` directory follow the api-simple rules, not these.
- New sites scaffold from `oym-craft5-boilerplate` with PHPStan max and PHPUnit wired first.
- Existing sites: the module you are editing defines the style; mirror it exactly.
<!-- oym-card:end stack=php-craft -->

## References

- Boilerplate: `oym-craft5-boilerplate`
- Newest module style: `domestique-cycling/modules/cyclingdatamodule/`
- Older module style (for contrast): `eredivisie/modules/optamodule/`
