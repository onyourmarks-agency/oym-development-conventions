# NestJS conventions (framework layer)

Applies to NestJS back-end services. Builds on `javascript.md` (NestJS shares the JavaScript style layer with the frontend: same typing rules, brace rule, naming, comment policy).

## Stack

- `[always]` NestJS with pnpm; Prisma with the MariaDB adapter is the standard data layer (one legacy service uses Mongoose; that is history, not the pattern).
- `[always]` Sentry: `instrument.ts` is the first import in `main.ts`, `SentryGlobalFilter` registered as global `APP_FILTER`. Never reorder the import.
- `[always]` Swagger mounted at `/docs`; health endpoints via Terminus (`/health` liveness, `/ready` with a DB ping); CLI entry via nest-commander (`src/cli.ts`) where the service has one-off jobs.

## Structure

- `[always]` Feature modules under `src/modules/<feature>/`, composition root `src/modules/app.module.ts`. Controller, service, gateway, `dto/`, providers, and commands co-located per feature. Shared filters/pipes under `src/common/`.

```
src/
  main.ts            bootstrap: global ValidationPipe, Swagger, CORS
  cli.ts             nest-commander entry (when needed)
  instrument.ts      Sentry init, imported first
  common/            shared filters, pipes, guards
  modules/
    app.module.ts    composition root
    <feature>/       controller, service, gateway, dto/, provider/, command/
    prisma/          PrismaService + repository/
    health/          Terminus health checks
```

- `[always]` Repository pattern: the Prisma client is touched only inside repositories (`modules/prisma/repository/`). Services depend on repositories, never on `PrismaService` directly.
- `[always]` External integrations live behind an interface (an `FooProviderInterface`) and are injected; add a new provider by implementing the interface and registering it in the module.

## Validation, config, errors

- `[always]` Global `ValidationPipe({ whitelist: true, transform: true })` in `main.ts`. Every HTTP body/query has a class-validator DTO; no untyped `@Body()`. WebSocket payloads validated through a WS validation pipe with a DTO class.

  ```ts
  export class ChatMessageDto {
    @IsUUID(4)
    conversationId: string;

    @IsString()
    @IsNotEmpty()
    @MaxLength(4000)
    content: string;
  }
  ```
- `[always]` Config exclusively via `@nestjs/config` (`isGlobal: true`) and `ConfigService.getOrThrow` for required values, resolved at boot so the service fails fast. Never `process.env` outside the config layer.
- `[always]` Errors fall through to Sentry. Do not swallow exceptions: log with the Nest `Logger` and rethrow, or return an explicit graceful fallback.
- `[always]` All repository and provider methods return Promises; never call them without `await`. Long-running operations accept an `AbortSignal` and respect `signal.aborted`.

## TypeScript strictness

Some existing services ship `noImplicitAny: false` and `no-explicit-any: off`. That is a known defect, not a convention; `javascript.md`'s no-`any` rule applies in full.

- `[new]` New services start with full `strict: true` and `@typescript-eslint/no-explicit-any: error`. Deltas: `templates/nestjs/`.
- `[existing]` Do not flip the strictness flags drive-by; that breaks the build repo-wide and is a dedicated migration task. Meanwhile: no new `any`, and remove `any` from files you touch.

## Testing

- `[always]` Jest with ts-jest, `*.spec.ts` next to the source. Unit tests only: mock repositories in service tests, mock services in controller/gateway tests. No test database, no supertest e2e.
- `[existing]` Coverage is uneven. Every new or modified service method gets a spec.
- Tests that hit real external APIs (LLM integration/judge suites) are gated behind explicit env flags and excluded from the default `pnpm test`.

## Documentation

- `[always]` Service repos carry a nested `AGENTS.md` tree (root plus per-module files): root documents what the service is, the module map, required env vars (table: var, consumer, notes), and scripts; module files document contracts. Read the nearest `AGENTS.md` before editing; it outranks this document.

<!-- oym-card:begin stack=framework-nestjs -->
# OYM NestJS rules

- All JavaScript layer rules apply (no `any`, arrow consts, braces, naming, comment policy); NestJS shares that layer with the frontend.
- Feature modules under `src/modules/<feature>/` with dto/, providers, commands co-located; composition root `app.module.ts`; shared filters/pipes in `src/common/`.
- Repository pattern: Prisma only inside `modules/prisma/repository/`; services depend on repositories, never `PrismaService`.
- External integrations behind injected interfaces.
- Global `ValidationPipe({ whitelist: true, transform: true })`; every body/query has a class-validator DTO; WS payloads validated via DTO pipe.
- Config only via global `@nestjs/config` with `getOrThrow` at boot; never `process.env` in business code.
- Sentry: `instrument.ts` first import in `main.ts`; `SentryGlobalFilter` as global filter; never swallow exceptions, log and rethrow or return an explicit fallback.
- Always `await` repository/provider calls; long-running ops respect `AbortSignal`.
- Swagger at `/docs`; Terminus `/health` + `/ready`; nest-commander for one-off jobs.
- Legacy `noImplicitAny: false` flags are a defect: no new `any`, remove `any` from touched files, flag flips are a dedicated migration.
- Jest unit tests only: mock repositories in service specs, mock services in controller specs; no test DB, no supertest e2e. New/modified service methods get a spec. Real-API test suites stay env-gated.
- Read the nearest nested `AGENTS.md` first; it outranks this card.
<!-- oym-card:end stack=framework-nestjs -->

## Templates

- `templates/nestjs/tsconfig.strict.delta.jsonc`, `templates/nestjs/eslint.any-rules.snippet.mjs`
