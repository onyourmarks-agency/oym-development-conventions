# SvelteKit conventions (tool layer)

Applies to standalone SvelteKit applications. Builds on `javascript.md` + `framework-svelte.md`.

## Stack

- `[always]` SvelteKit with `adapter-node`, runes forced via `dynamicCompileOptions` in `svelte.config.js`, CSS-first Tailwind.
- `[always]` UI primitives via shadcn-svelte / bits-ui where the project uses them (`components.json` present): check `lib/components/ui` for an existing primitive before writing one.
- `[always]` i18n via Paraglide/inlang (`project.inlang/`, `messages/`), locale routing with a `[lang=lang]` matcher where the app is multilingual.
- `[always]` Sentry via the SvelteKit SDK, wired in `hooks.server.ts`/`hooks.client.ts`.
- `[always]` Forms: superforms + formsnap + zod where the app has real forms.
- `[always]` Deployed as a Docker image via the docker CI components (`stage-build-image`, `stage-tag-environment`); `/health` and `/ready` endpoints exist for the platform.

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
hooks.server.ts  hooks.client.ts  app.css (or style/main.css)
```

- `[always]` Types and interfaces are defined in `entities/` (or the repo's established types location), imported where needed. The one exception: component prop types are inline in the `$props()` destructuring.
- `[always]` Domain entities are TypeScript classes (definite assignment or defaults, `static fromObject` constructors, getters) when they carry behaviour, not bare interfaces.

## Client/server boundary

- `[always]` Server secrets and `$env/dynamic/private` only under `lib/server/` or `*.server.ts` files. Keep the boundary strict.
- `[always]` Server-side singletons follow the `getInstance()` + lazy `??=` getter pattern (`*.server.ts` modules); client side uses `BaseModule` the same way.
- `[always]` `load` functions in `+page(.server).ts`/`+layout.server.ts` wrap API access via the module layer with try/catch fallbacks returning explicit null data; typed as `PageLoad`/`PageServerLoad`.
- `[always]` REST endpoints as `+server.ts` with `RequestHandler`, `json()`, and stable coded error strings (`chat_error_410` pattern). Preserve existing endpoint contracts (compact JSON, stable field names).
- `[always]` Auth/session checks on API routes are load-bearing; never remove or bypass them in a refactor.
- `[always]` Fetch through the project's service classes (`credentials: 'include'`, CSRF handling where present, throw coded errors on `!response.ok`).

## Quality gates

- `[always]` This family's Prettier config: tabs, single quotes, `trailingComma: 'none'`, printWidth 100, Tailwind plugin with `tailwindStylesheet` pointing at the main CSS file. Semicolons stay on (Prettier default; the config does not disable them). The repo config is authoritative (see `javascript.md`).
- `[always]` Import order via `simple-import-sort` where configured.
- `[always]` The gate is the project's scripts: `pnpm frontend-cs` where defined, otherwise `lint` + `check`; `ddev pnpm ...` in DDEV projects.
- `[new]` New apps use the structure above, Vitest from day one for entities and states (`templates/svelte/vitest.config.sveltekit.ts`), and CSP directives configured in `svelte.config.js`.
- `[existing]` One older SvelteKit app predates this standard (adapter-static, no Tailwind, BEM CSS, class-transformer models, Storybook). Inside it, its local style wins.

<!-- oym-card:begin stack=tool-sveltekit -->
# OYM SvelteKit rules

- `adapter-node`, runes forced, CSS-first Tailwind. Deployed as a Docker image; keep `/health` and `/ready`.
- Layout: `src/lib/{components,server,entities,states,modules}` + `routes/` with groups and `[lang=lang]`.
- Types/interfaces in `entities/` (domain entities as classes with `static fromObject` when they carry behaviour); only component prop types are inline in `$props()`.
- Strict client/server boundary: secrets and `$env/dynamic/private` only in `lib/server` / `*.server.ts`. Server singletons via `getInstance()` + lazy `??=` getters.
- `load` functions typed (`PageLoad`/`PageServerLoad`), API access through the module layer, try/catch with explicit null fallbacks.
- Endpoints as `+server.ts` with `RequestHandler` + `json()` and stable coded error strings; preserve existing response contracts; never weaken auth checks on API routes.
- Fetch via the project's service classes (`credentials: 'include'`, throw coded errors on `!response.ok`).
- shadcn-svelte/bits-ui first where present: check `lib/components/ui` before writing a primitive. Paraglide for all UI strings.
- Family formatting (repo config is authoritative): tabs, single quotes, trailingComma none, width 100, Tailwind plugin; semicolons stay on.
- Run the project's gate before finishing (`frontend-cs` or lint + check; `ddev pnpm` where DDEV).
- New apps: this structure, Vitest for entities/states day one, CSP in `svelte.config.js`.
<!-- oym-card:end stack=tool-sveltekit -->

## Examples in the wild

- Newest: `eredivisie-top1000-fe`; richest: `oym-chant-fe`; forms/auth/CSP: `oym-paddock-fe`
- Legacy contrast: `nevobo-trainersplatform-fe`
