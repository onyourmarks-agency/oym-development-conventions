# Symfony conventions

Applies to repositories requiring `symfony/framework-bundle`. Builds on `php-universal.md`.

The org has one Symfony project: `nevobo-trainersplatform-api` (Symfony 7, API Platform 4, Doctrine ORM 3, Lexik JWT, GraphQL). That project is the convention: mirror its structure, DI style, and API Platform resource patterns rather than generic Symfony tutorials.

## Rules

- `[always]` Everything in `php-universal.md`: strict types, phpcs OYM ruleset, final/readonly/promotion defaults, PHPStan policy, unit-test-only policy.
- `[always]` Constructor injection with promoted, typed, readonly properties. Autowiring via `config/services.yaml`; no service locators in business code.
- `[always]` API Platform: resources and DTOs explicitly typed; validation via constraint attributes, never manual checks in controllers; serialization groups explicit.
- `[always]` Doctrine: attributes for mapping, migrations for every schema change, repositories for queries; no DQL strings scattered through services.
- `[always]` Controllers stay thin; logic in services; external systems behind an interface.
- `[new]` A second Symfony project starts by copying the trainersplatform structure. PHPStan max from day one with `phpstan/phpstan-symfony` and `phpstan/phpstan-doctrine` extensions; Symfony tolerates max well, no baseline expected.
- `[existing]` The trainersplatform has a `tests/` bootstrap; extend it with unit tests for services (mock repositories), not WebTestCase integration tests.

<!-- oym-card:begin stack=php-symfony -->
# OYM Symfony rules

- Reference project: `nevobo-trainersplatform-api` (Symfony 7 + API Platform 4). Mirror its structure and patterns; it outranks generic Symfony habits.
- Constructor injection with promoted readonly typed properties; autowiring; no service locators.
- API Platform resources/DTOs explicitly typed; validation via constraint attributes only; explicit serialization groups.
- Doctrine attributes + migrations + repositories; no scattered DQL.
- Thin controllers, logic in services, external systems behind interfaces.
- PHPStan max with symfony/doctrine extensions; no baseline expected on Symfony.
- Unit tests for services with mocked repositories; no WebTestCase integration tests.
- All PHP universal rules apply (strict types, phpcs OYM, final/readonly, comment policy).
<!-- oym-card:end stack=php-symfony -->

## References

- Reference project: `nevobo-trainersplatform-api`
