# Vite islands conventions (tool layer)

Applies to Svelte islands embedded in Craft CMS projects via plain Vite (no SvelteKit). Builds on `javascript.md` + `framework-svelte.md`. New projects of this type scaffold from the `oym-vite-boilerplate` repo.

## Structure

- `[always]` Frontend lives in `resources/` at the Craft project root: `resources/js/`, `resources/css/`. Path aliases `$src`, `$js`, `$css` (tsconfig + vite.config.ts).
- `[always]` `resources/js/` layout:

```
app.ts                    entry: new Bootstrap(), module boot
app/base/                 BaseModule and core services
components/
  _base/                  primitives (Svg, ...)
  atoms/  molecules/  organisms/    atomic design
modules/                  feature module classes
states/                   [Name]State.svelte.ts runes state
_types/                   shared type definitions
i18n/<locale>/            i18next resources
isolated/                 standalone islands
```

## Architecture

- `[always]` `BaseModule` singleton with lazy getters (`this._utils ??= new UtilService()`) orchestrates services (UtilService, TranslationService, FocusTrapService, KeyPressService, WindowResizeService). Follow it exactly; do not introduce a DI container or a second service pattern.
- `[always]` Islands architecture: Svelte mounts into Craft Twig templates. Entry points stay thin; heavy features load via dynamic `import()` gated on DOM presence or window flags, with `.catch()` routed to the error service.
- `[always]` i18n via i18next with resources under `resources/js/i18n/`.
- `[always]` Long Tailwind class lists may be extracted to a `const fooClassNames = [...]` array; tokens per `framework-svelte.md`.
- `[always]` SVGs go through the sprite build (`pnpm svg-sprite`, spritemap plugin), rendered via the base `Svg` component.

## Quality gates and build

- `[always]` This family's Prettier config: semicolons on, 2-space indent, single quotes, `trailingComma: 'all'`, printWidth 100, Twig formatted via the prettier Twig plugin. The repo config is authoritative (see `javascript.md`).
- `[always]` ESLint: flat config with the type-checked strict presets; `prefer-nullish-coalescing` and `no-unnecessary-condition` are errors; `no-console` warns except `warn`/`error`.
- `[always]` Stylelint on CSS (standard config plus BEM pattern plugin, Tailwind at-rules ignored).
- `[always]` The aggregate gate is `pnpm frontend-cs` (format + eslint + stylelint + ts-check + svelte-check). Run it before finishing; `ddev pnpm frontend-cs` in DDEV projects.
- `[always]` Build output goes to `public/` with a manifest consumed by Craft Twig (`renderManifestTag`); the dev server is DDEV-aware. Do not restructure the vite config casually.
- `[new]` Scaffold from `oym-vite-boilerplate`; add `templates/svelte/vitest.config.vite-craft.ts` for unit tests.

<!-- oym-card:begin stack=tool-vite -->
# OYM Vite islands rules (Svelte in Craft)

- Frontend in `resources/js` with aliases `$src`/`$js`/`$css`; components in atomic design tree (`_base`, `atoms`, `molecules`, `organisms`); feature classes in `modules/`; runes state in `states/[Name]State.svelte.ts`.
- `BaseModule` singleton with lazy `??=` getters is the service layer; follow it, never add a DI container or second pattern.
- Islands: Svelte mounts into Craft Twig; thin entries; heavy features via dynamic `import()` gated on DOM/window flags with errors routed to the error service.
- i18next for all UI strings (`resources/js/i18n/`).
- SVGs via the sprite build and the base `Svg` component.
- Family formatting (repo config is authoritative): semicolons, 2-space, single quotes, trailingComma all, width 100. ESLint type-checked strict presets. Stylelint on CSS.
- The gate is `pnpm frontend-cs` (format + eslint + stylelint + ts-check + svelte-check); run before finishing (`ddev pnpm frontend-cs` in DDEV).
- Build outputs to `public/` with a manifest consumed by Twig; vite config is DDEV-aware, do not restructure casually.
- New projects scaffold from `oym-vite-boilerplate`.
<!-- oym-card:end stack=tool-vite -->

## Examples in the wild

- Boilerplate: `oym-vite-boilerplate`
- Production example: `domestique-cycling` (`resources/js/`)
