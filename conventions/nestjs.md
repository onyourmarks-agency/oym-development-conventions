# NestJS conventions

Applies to the NodeJS back-end services (the `*-agent` services and future NestJS APIs). Builds on `shared.md`. Canonical project: `unisport-agent`.

## Stack

- `[always]` NestJS 11, TypeScript, pnpm, Node 24, ESM tooling. Prisma 7 with the `@prisma/adapter-mariadb` driver against MariaDB is the standard data layer (`domestique-newsarticle-ai` on Mongoose is the historical outlier, not the pattern).
- `[always]` Sentry via `@sentry/nestjs`: `instrument.ts` is the first import in `main.ts`, `SentryGlobalFilter` registered as global `APP_FILTER`. Never reorder the import.
- `[always]` Swagger mounted at `/docs`; health endpoints via Terminus (`/health` liveness, `/ready` with a DB ping); CLI entry via nest-commander (`src/cli.ts`) where the service has one-off jobs.

## Structure

- `[always]` Feature modules under `src/modules/<feature>/`, composition root `src/modules/app.module.ts`. Controller, service, gateway, `dto/`, providers, and commands co-located per feature. Shared filters/pipes under `src/common/`.
- `[always]` Repository pattern: the Prisma client is touched only inside repositories (`modules/prisma/repository/`). Services depend on repositories, never on `PrismaService` directly.
- `[always]` External integrations live behind an interface (`LlmProviderInterface` pattern) and are injected; add a new provider by implementing the interface and registering it in the module.

## Validation, config, errors

- `[always]` Global `ValidationPipe({ whitelist: true, transform: true })` in `main.ts`. Every HTTP body/query has a class-validator DTO; no untyped `@Body()`. WebSocket payloads validated through the WS validation pipe with a DTO class.
- `[always]` Config exclusively via `@nestjs/config` (`isGlobal: true`) and `ConfigService.getOrThrow` for required values, resolved at boot so the service fails fast. Never `process.env` outside the config layer.
- `[always]` Errors fall through to Sentry. Do not swallow exceptions: log with the Nest `Logger` and rethrow, or return an explicit graceful fallback.
- `[always]` All repository and provider methods return Promises; never call them without `await`. Long-running operations accept an `AbortSignal` and respect `signal.aborted`.

## TypeScript strictness

The current repos ship `noImplicitAny: false` and `no-explicit-any: off`. That is a known defect, not a convention.

- `[always]` New code is `any`-free regardless of what the compiler flags allow. Type Prisma results, provider responses, and DTOs fully.
- `[new]` New services start with full `strict: true` (including `noImplicitAny`) and `@typescript-eslint/no-explicit-any: error`. Deltas: `templates/nestjs/`.
- `[existing]` Do not flip the strictness flags drive-by; that breaks the build repo-wide and is a dedicated migration task. Meanwhile: no new `any`, and remove `any` from files you touch.

## Testing

- `[always]` Jest with ts-jest, `*.spec.ts` next to the source. Unit tests only: mock repositories in service tests, mock services in controller/gateway tests. No test database, no supertest e2e.
- `[existing]` Coverage is uneven (some services have none). Every new or modified service method gets a spec.
- LLM-dependent tests follow the unisport pattern: real-API and judge suites exist but are gated behind explicit env flags (`RUN_LLM_INTEGRATION`), excluded from the default `pnpm test`.

## Documentation

- `[always]` Agent-service repos carry a nested `AGENTS.md` tree (root plus per-module files) following unisport-agent: root documents stack, module map, required env vars (table: var, consumer, notes), and scripts; module files document contracts. Read the nearest `AGENTS.md` before editing; it outranks this document.
- `[always]` Comment policy from `shared.md` applies. The bar for an inline comment is a non-obvious operational reason (the connection-pool sizing rationale in unisport's `prisma.service.ts` is the model).

<!-- oym-card:begin stack=nestjs -->
# OYM NestJS rules

- Reference project: `unisport-agent`. NestJS 11, pnpm, Node 24, Prisma 7 + MariaDB adapter.
- Feature modules under `src/modules/<feature>/` with dto/, providers, commands co-located; composition root `app.module.ts`; shared filters/pipes in `src/common/`.
- Repository pattern: Prisma only inside `modules/prisma/repository/`; services depend on repositories, never `PrismaService`.
- External integrations behind injected interfaces (`LlmProviderInterface` pattern).
- Global `ValidationPipe({ whitelist: true, transform: true })`; every body/query has a class-validator DTO; WS payloads validated via DTO pipe.
- Config only via global `@nestjs/config` with `getOrThrow` at boot; never `process.env` in business code.
- Sentry: `instrument.ts` first import in `main.ts`; `SentryGlobalFilter` as global filter; never swallow exceptions, log and rethrow or return an explicit fallback.
- Always `await` repository/provider calls; long-running ops respect `AbortSignal`.
- `any` is forbidden in new code even where compiler flags allow it (`noImplicitAny: false` in old repos is a defect). Do not flip strictness flags drive-by; remove `any` from files you touch.
- Jest unit tests only: mock repositories in service specs, mock services in controller specs; no test DB, no supertest e2e. New/modified service methods get a spec.
- Read the nearest nested `AGENTS.md` first; it outranks this card.
<!-- oym-card:end stack=nestjs -->

## References

- Canonical: `unisport-agent` (root and nested `AGENTS.md`, `src/main.ts`, `modules/prisma/repository/`)
- Also current: `runningcoach-agent`, `zendesk-support-agent`
- Templates: `templates/nestjs/tsconfig.strict.delta.jsonc`, `templates/nestjs/eslint.any-rules.snippet.mjs`
