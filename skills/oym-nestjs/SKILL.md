---
name: oym-nestjs
description: OYM NestJS conventions. Use when working in On Your Marks NestJS services - writing modules, controllers, services, repositories, DTOs, gateways, or configuring the app. Covers structure, config and Sentry bootstrap rules, typing strictness, and Jest unit-test policy.
---

# OYM NestJS

Layered: read `${CLAUDE_PLUGIN_ROOT}/conventions/shared.md` + `conventions/javascript.md` + `conventions/framework-nestjs.md`. NestJS shares the JavaScript style layer with the Svelte front-ends: same typing rules, brace rule, naming, and comment policy apply here.

## First: read the nearest AGENTS.md

Service repos carry a nested `AGENTS.md` tree (root + per-module). The file closest to what you are editing outranks this skill. Read it before writing.

## Mode

`AGENTS.md` block header `mode=` wins; otherwise a repo with substantial code is `existing`: mirror its local structure (some put features directly under `src/`, the standard is `src/modules/`), do not restructure.

## Non-negotiables

- Feature modules under `src/modules/<feature>/`; controller, service, `dto/`, providers co-located; composition root `app.module.ts`; shared filters/pipes in `src/common/`.
- Repository pattern: Prisma client only inside repositories; services depend on repositories, never `PrismaService` directly.
- External integrations behind injected interfaces; add providers by implementing the interface and registering in the module.
- Every HTTP body/query has a class-validator DTO (global `ValidationPipe({ whitelist: true, transform: true })`); WS payloads validated with a DTO pipe. No untyped `@Body()`.
- Config only via global `@nestjs/config` and `getOrThrow` at boot. Never `process.env` in business code.
- Sentry: `instrument.ts` stays the first import of `main.ts`. Never swallow exceptions: `Logger` + rethrow, or explicit fallback.
- Always `await` repository/provider calls; long-running operations respect `AbortSignal`.
- `any` is forbidden in new code even where legacy compiler flags allow it (that is a defect, not a convention). Remove `any` from files you touch. Do NOT flip the compiler flags drive-by; that is a dedicated migration task.
- Jest unit tests only: mock repositories in service specs, mock services in controller/gateway specs; no test DB, no supertest e2e. Every new or modified service method gets a spec.
- Comments: non-obvious operational "why" only; keep the default verbosity down.
- Versions come from the repo's package.json; org targets in `${CLAUDE_PLUGIN_ROOT}/VERSIONS.md`.

## New services

Follow the structure in `framework-nestjs.md` with full `strict: true` and `no-explicit-any: error` (`templates/nestjs/`), Swagger at `/docs`, Terminus `/health` + `/ready`, and a root `AGENTS.md` documenting the module map and required env vars. Full checklist: the `oym-new-project` skill.
