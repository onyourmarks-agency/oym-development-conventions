# Symfony conventions (framework layer)

Applies to repositories requiring `symfony/framework-bundle`. Builds on `php.md`. OYM Symfony services are APIs built on API Platform with Doctrine ORM.

## Structure

- `[always]` Standard Symfony skeleton: `src/` with `ApiResource/` (or attribute-mapped `Entity/`), `Entity/`, `Repository/`, `Service/`, `Controller/` (thin, only where API Platform does not cover the route), `config/` with `services.yaml` autowiring, `migrations/`.
- `[always]` Constructor injection with promoted, typed, readonly properties. Autowiring via `config/services.yaml`; no service locators or container pulls in business code.
- `[always]` Controllers stay thin; logic in services; external systems behind an interface.

## API Platform

- `[always]` Resources and DTOs explicitly typed; state providers/processors over ad-hoc controllers for custom operations.
- `[always]` Validation via constraint attributes on the resource/DTO, never manual checks in controllers.
- `[always]` Serialization groups explicit on every resource; no accidental exposure by default serialization.
- `[always]` Auth via the configured JWT/security layer; never weaken security attributes on operations in a refactor.

## Doctrine

- `[always]` Attribute mapping on entities; every schema change ships a migration; repositories own the queries; no DQL strings scattered through services.
- `[always]` Entities carry typed properties and behaviour-light logic; heavy logic lives in services.

## Quality gates

- `[new]` PHPStan max from day one with the `phpstan/phpstan-symfony` and `phpstan/phpstan-doctrine` extensions; Symfony tolerates max well, no baseline expected.
- `[always]` Unit tests for services with mocked repositories; no `WebTestCase`/`KernelTestCase` integration tests.
- `[existing]` Follow the repo's existing directory and resource layout exactly; extend the patterns that are there.

<!-- oym-card:begin stack=framework-symfony -->
# OYM Symfony rules

- API Platform + Doctrine ORM. Standard skeleton: `src/{Entity,Repository,Service,Controller}`, autowired via `config/services.yaml`, migrations for every schema change.
- Constructor injection with promoted readonly typed properties; no service locators.
- API Platform resources/DTOs explicitly typed; custom operations via state providers/processors; validation via constraint attributes only; explicit serialization groups; never weaken security attributes.
- Repositories own queries; no scattered DQL; entities typed, heavy logic in services; thin controllers.
- PHPStan max with the symfony and doctrine extensions; no baseline expected.
- Unit tests for services with mocked repositories; no WebTestCase/KernelTestCase integration tests.
- All PHP language rules apply (strict types, phpcs OYM, final/readonly, comment policy).
<!-- oym-card:end stack=framework-symfony -->

## Templates

- `templates/php/phpcs.xml` (change scan target to `src/`), `templates/php/phpstan.neon.dist` (add the symfony/doctrine extensions), `templates/php/phpunit.xml.dist`
