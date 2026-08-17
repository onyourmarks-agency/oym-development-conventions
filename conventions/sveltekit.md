# SvelteKit conventions

Applies to standalone SvelteKit applications. Builds on `svelte-universal.md`. Canonical: `oym-chant-fe`; `oym-paddock-fe` for forms, auth, and CSP patterns.

## Stack

- `[always]` SvelteKit 2 with `adapter-node`, Svelte 5 runes (force runes via `dynamicCompileOptions` in `svelte.config.js`), Vite 7, Tailwind 4, Node 24, pnpm.
- `[always]` UI primitives via shadcn-svelte / bits-ui (`components.json`, `tailwind-variants`, `tailwind-merge`, `@lucide/svelte`). Check `lib/components/ui` for an existing primitive before writing one.
- `[always]` i18n via Paraglide/inlang (`project.inlang/`, `messages/`), locale routing with a `[lang=lang]` matcher where the app is multilingual.
- `[always]` Sentry via `@sentry/sveltekit`, wired in `hooks.server.ts`/`hooks.client.ts`.
- `[always]` Forms: `sveltekit-superforms` + `formsnap` + `zod` (paddock-fe pattern).
- `[always]` Deployed as a Docker image via the `@1.x` CI components (`stage-build-image`, `stage-tag-environment`); `/health` and `/ready` endpoints exist for the platform.

## Structure

- `[always]` Layout under `src/`:

```
lib/
  components/           shadcn ui/ + project components (atomic tree where used)
  server/               server-only code: modules, API clients, secrets
  entities/             domain types and classes
  states/               [name]State.svelte.ts runes state
  modules/              client service classes, BaseModule
  config/  hooks/  actions/  utils/  assets/
routes/                 (groups), [lang=lang], api/ endpoints
hooks.server.ts  hooks.client.ts  app.css
```

- `[always]` Types and interfaces are defined in `entities/`, imported where needed. The one exception: component prop types are inline in the `$props()` destructuring (still importing entity classes).
- `[always]` Domain entities are TypeScript classes (definite assignment or defaults, `static fromObject` constructors, getters), not bare interfaces, when they carry behaviour.
- `[always]` Named exports only; default exports need a specific reason (SvelteKit's own conventions excepted).

## Client/server boundary

- `[always]` Server secrets and `$env/dynamic/private` only under `lib/server/` or `*.server.ts` files. Keep the boundary strict.
- `[always]` Server-side singletons follow the `getInstance()` + lazy `??=` getter pattern (`*.server.ts` modules); client side uses `BaseModule` the same way.
- `[always]` `load` functions in `+page(.server).ts`/`+layout.server.ts` wrap API access via the module layer with try/catch fallbacks returning explicit null data; typed as `PageLoad`/`PageServerLoad`.
- `[always]` REST endpoints as `+server.ts` with `RequestHandler`, `json()`, and stable coded error strings (`chat_error_410` pattern). Preserve existing endpoint contracts (compact JSON, stable field names).
- `[always]` Auth/session checks on API routes are load-bearing; never remove or bypass them in a refactor.

## Behaviour rules

- `[always]` Error contract: services throw, components catch and display. Services never swallow; components never let errors cross the UI boundary unlogged.
- `[always]` Fetch through the project's service classes (`credentials: 'include'`, CSRF handling where present, throw coded errors on `!response.ok`).
- `[always]` Prettier: tabs, no semicolons, single quotes, `trailingComma: none`, printWidth 100, Tailwind plugin with `tailwindStylesheet` pointing at `src/app.css`. Import order via `simple-import-sort` where configured.
- `[always]` The gate is the project's lint/format/check scripts (`pnpm frontend-cs` where defined, otherwise `lint` + `check`); DDEV projects run them as `ddev pnpm ...`.
- `[new]` New apps clone the chant-fe structure; Vitest from day one for entities and states (`templates/svelte/vitest.config.sveltekit.ts`); CSP directives configured in `svelte.config.js` per paddock-fe.
- `[existing]` `nevobo-trainersplatform-fe` predates this standard (adapter-static, no Tailwind, BEM CSS, class-transformer models, Storybook). Inside it, its local style wins.

<!-- oym-card:begin stack=sveltekit -->
# OYM SvelteKit rules

- SvelteKit 2 + adapter-node, Svelte 5 runes forced, Tailwind 4, Node 24, pnpm. Deployed as a Docker image; keep `/health` and `/ready`.
- Layout: `src/lib/{components,server,entities,states,modules}` + `routes/` with groups and `[lang=lang]`.
- Types/interfaces in `entities/` (domain entities as classes with `static fromObject`); only component prop types are inline in `$props()`.
- Named exports only.
- Strict client/server boundary: secrets and `$env/dynamic/private` only in `lib/server` / `*.server.ts`. Server singletons via `getInstance()` + lazy `??=` getters.
- `load` functions typed (`PageLoad`/`PageServerLoad`), API access through the module layer, try/catch with explicit null fallbacks.
- Endpoints as `+server.ts` with `RequestHandler` + `json()` and stable coded error strings; preserve existing response contracts; never weaken auth checks on API routes.
- Services throw, components catch and display; nothing swallowed.
- shadcn-svelte/bits-ui first: check `lib/components/ui` before writing a primitive. Paraglide for all UI strings.
- Prettier: tabs, no semicolons, trailingComma none, width 100, Tailwind plugin. Run the project's lint/format/check scripts before finishing (`ddev pnpm` where DDEV).
- New apps clone `oym-chant-fe` structure, Vitest for entities/states day one, CSP per `oym-paddock-fe`.
<!-- oym-card:end stack=sveltekit -->

## References

- Canonical: `oym-chant-fe` (and its `AGENTS.md`)
- Forms/auth/CSP: `oym-paddock-fe`
- Legacy contrast: `nevobo-trainersplatform-fe`
