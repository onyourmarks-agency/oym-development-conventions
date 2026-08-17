# Shared conventions (all code)

Applies to every On Your Marks repository. The conventions are layered; a project composes:

- **Back-end (PHP)**: this file + `php.md` + one framework file (`framework-craft.md`, `framework-craft-plugin.md`, `framework-api-simple.md`, `framework-symfony.md`)
- **Front-end (JavaScript)**: this file + `javascript.md` + `framework-svelte.md` + one tool file (`tool-vite.md`, `tool-webpack.md`, `tool-sveltekit.md`)
- **NestJS**: this file + `javascript.md` + `framework-nestjs.md` (NestJS shares the JavaScript style layer with the frontend)

## Precedence

From strongest to weakest:

1. The task itself. Do what was asked, nothing more.
2. The project's own `AGENTS.md` rules above the vendored convention blocks.
3. Tool file, then framework file, then language file (`php.md`/`javascript.md`), then this file.

If local code contradicts an `[always]` rule in an existing project: follow the local code, mention the conflict in your summary or PR description, and do not mass-refactor to resolve it.

## Project mode: existing vs new

Every rule in every layer is tagged:

- `[always]` applies in every project.
- `[new]` applies when starting a project from scratch. New projects get the full standard from day one.
- `[existing]` applies when working in a project that already has substantial code. In existing projects the surrounding code is the strongest convention signal, stronger than this document. Mirror local idioms exactly. Do not retrofit newer patterns, tooling, or formatting unless that migration is the task.

How to determine the mode:

1. If the `AGENTS.md` convention block header declares `mode=existing` or `mode=new`, that wins.
2. Otherwise: a `phpstan-baseline.neon` file, or a repository with substantial committed feature code, means existing.
3. When unsure, treat the project as existing. Mirroring local context never breaks anything; retrofitting does.

## Versions

Version numbers do not appear in convention rules. Current org-wide targets live in `VERSIONS.md` at the root of the conventions repo; each product repo pins its own versions in `composer.json`, `package.json` (`engines`, `packageManager`), and `.nvmrc`, and those manifests are the runtime truth. When a rule depends on a language feature, it says so ("where the PHP version allows") instead of naming a version.

## Commits

- `[always]` Conventional Commits with optional scope: `feat(optamodule): add live standings consumer`. Types in use: `feat`, `fix`, `chore`, `docs`, `style`, `refactor`, `perf`, `ci`, `build`, `test`.
- `[always]` Scope is a module or feature name, not a filename.
- `[always]` Subject in lower case imperative: `feat: add x`, not `feat: Added x`.
- `[always]` Breaking changes: `feat!:` or a `BREAKING CHANGE:` footer. Version bumps are derived from these patterns (gsemver: `feat|chore|build|ci|refactor|perf` bump minor, `!`/`BREAKING CHANGE` bump major).
- `[always]` The team commit tool is the OYM commitizen fork, invoked as `git cz`. Hand-written messages must follow the same format.
- Code, comments, identifiers, and commit messages are written in English.

## Comments

AI agents consistently write too many comments. The OYM standard is the opposite of the AI default. When in doubt, omit the comment.

- `[always]` A comment explains a non-obvious "why", never a "what". If the code states it, do not repeat it.
- `[always]` No filler docblocks: no `/** Constructor. */`, no `Class Foo`, no parameter docblocks that repeat native types. PHP enforces this via phpcs (Slevomat `UselessFunctionDocComment`, `ForbiddenComments`); the same rule applies to TypeScript and Svelte where no linter enforces it.
- `[always]` Banned annotations: `@author`, `@created`, `@copyright`, `@license`, `@package`, `@version`, `@since`.
- `[always]` A docblock is justified only when it carries information the signature cannot express: generics, array shapes, magic `@property` maps, an external contract (`@link` to an endpoint).
- `[always]` `// TODO:` in uppercase, with enough context that someone else can pick it up.
- `[always]` Do not add or remove comments in code you are not otherwise changing.

## Typing

The codebase is the IDE's autocompletion database. Types are not optional documentation, they are the interface.

- `[always]` Everything is typed: parameters, returns, properties, class constants where the language supports it.
- `[always]` No escape hatches in new code: no `any` (TypeScript), no `mixed` in new signatures (PHP), no error suppression to silence the type checker.
- `[always]` Closed value sets are backed enums (PHP) or union literal types (TypeScript), never loose strings crossing a boundary.
- `[always]` Explicit return types everywhere, including `void` and `never`.

## Testing

- `[always]` Unit tests only. Never integration tests: no database, no HTTP, no filesystem, no framework container boot. Mock at the boundary.
- `[always]` Test targets are logic: services, entities, handlers, pure functions. Do not test framework glue, getters, or configuration.
- `[always]` Runners: PHPUnit (PHP), Jest (NestJS), Vitest (Svelte).
- `[new]` New projects wire the test runner before feature work starts and every logic-bearing class gets a test alongside it.
- `[existing]` Most existing projects have no test suite. Do not scaffold one unprompted. When you add a new service or handler with real logic, add a unit test alongside, creating the minimal config from `templates/` if none exists.

## Tooling and infrastructure

- `[always]` JavaScript: pnpm, ESM (`"type": "module"`), Node version pinned per repo (`engines` + `.nvmrc`).
- `[always]` PHP: Composer with the private Repman registry (`onyourmarks.repo.repman.io`); `oym/sentry-logger` is a standard dependency.
- `[always]` CI is GitLab with the `oym-gitlab-deploy-templates` component catalog: the PHP deploy flow (build-vendors, test-phpcs, deploy via Deployer) for PHP/Craft projects, the docker flow (`stage-build-image`, `stage-tag-environment`) for containerized Node apps. Component versions: `VERSIONS.md`.
- `[always]` Pre-commit hooks (husky + lint-staged) are part of the contract. Never bypass with `--no-verify`.
- `[always]` Anything published as a package (Craft plugins, npm packages, this repo) follows Keep a Changelog and SemVer.
- `[new]` New repositories start from the canonical boilerplate for their stack where one exists. Never hand-roll project structure.

<!-- oym-card:begin stack=shared -->
# OYM shared rules

Mode: `existing` projects mirror local code exactly; the surrounding module outranks this document. `new` projects apply the full standard from day one. Rules below marked (new) only apply to new projects.

Layering: back-end = shared + PHP + framework. Front-end = shared + JavaScript + Svelte + tool. NestJS = shared + JavaScript + NestJS. More specific layers win.

Versions are never hardcoded in rules: the repo's own manifests (composer.json, package.json engines, .nvmrc) are the truth; org targets live in the conventions repo's VERSIONS.md.

## Commits
- Conventional Commits, optional scope, lower case imperative subject: `feat(optamodule): add live standings consumer`.
- Breaking changes via `!` or `BREAKING CHANGE:` footer.

## Comments (AI: your default verbosity is wrong here)
- Comments explain non-obvious "why" only. Never restate what code does.
- No filler docblocks. No `@author`/`@package`/`@version`/`@since`.
- Docblocks only for what signatures cannot express: generics, array shapes, `@property` maps.
- Do not touch comments in code you are not otherwise changing.

## Typing
- Type everything: params, returns, properties, constants.
- No `any` (TS), no new `mixed` (PHP), no suppressions.
- Backed enums / union literals over magic strings.

## Testing
- Unit tests only. No DB, no HTTP, no container boot. Mock at the boundary.
- Test logic (services, entities, handlers), not framework glue.
- Existing projects: add tests alongside new logic only, do not scaffold suites unprompted. (new) Wire the runner day one.

## Tooling
- pnpm + ESM for JS; Composer + Repman for PHP.
- GitLab CI via `oym-gitlab-deploy-templates` components.
- Never bypass pre-commit hooks.
- English everywhere.
<!-- oym-card:end stack=shared -->

## References

- Commit tool: `github.com/onyourmarks-agency/commitizen` (`git cz`)
- CI catalog: `oym-gitlab-deploy-templates` (GitLab component catalog)
- Deploy image: `oym-docker-deploy-image` (Alpine + PHP + Deployer)
