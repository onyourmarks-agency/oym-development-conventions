# Svelte conventions (framework layer)

Applies to every Svelte component in OYM projects, regardless of tool (Vite islands, webpack islands, SvelteKit apps). Builds on `javascript.md`. The tool file (`tool-vite.md`, `tool-webpack.md`, `tool-sveltekit.md`) defines project structure, build, and deploy.

## Components

- `[always]` Runes exclusively: `$state`, `$derived`, `$effect`, `$props`. No legacy `$:` reactive statements, no `export let` props, no `svelte/store` (`writable`/`readable`/`derived`) in new code.
- `[always]` Event attributes, not directives: `onclick`, not `on:click`.
- `[always]` Every component is `<script lang="ts">`.
- `[always]` Props are declared with `$props()` destructuring and typed inline, never in a separate interface. Formatting: each prop on its own line; optional props come after required ones, use `?:`, and always have a default:

  ```svelte
  <script lang="ts">
  	let {
  		text,
  		size = 'md',
  		onclick
  	}: {
  		text: string;
  		size?: 'sm' | 'md' | 'lg';
  		onclick?: () => void;
  	} = $props();
  </script>
  ```

  Domain types used inside the prop type are still imported from the project's entities/types location.
- `[always]` Component files are PascalCase. Check the existing component tree (a shadcn `ui/` directory or the atomic `atoms/` tree, per the tool file) before writing a new primitive.
- `[always]` Shared reactive state lives in `[name]State.svelte.ts` files using runes: `export const FooState = $state<...>({...})` or a runes class with private fields and getters.
- `[always]` Error contract: services throw, components catch and display. Components never let an error cross the UI boundary unlogged.
- `[always]` All UI strings go through the project's i18n (i18next or Paraglide, per the tool file); no hardcoded copy where translations exist.

## Styling (Tailwind projects)

Projects on CSS-first Tailwind define their design tokens in the `@theme` block of the main CSS file. That file is the authoritative token list.

- `[always]` Tailwind utility classes inline are the primary styling method. `<style>` blocks only when no utility exists or a large identical class group repeats across many elements; `@apply` only for genuinely repeated groups.
- `[always]` No arbitrary bracket values in final output: no `gap-[24px]`, `text-[14px]`, `leading-[1.2]`, `bg-[#...]`, `text-[color:var(...)]`. Every value maps to a token class.
- `[always]` Spacing derives from the px-divided-by-4 scale: `24px` becomes `gap-6`/`p-6`; half-steps for values not divisible by 4. Combine identical axes (`p-6`, not `px-6 py-6`).
- `[always]` Typography and colors use only the tokens defined in `@theme` (`text-32`, `leading-120`, `tracking-10`, `bg-<token>`); no default Tailwind sizes like `text-xl` in token-based projects.
- `[always]` A value with no matching token: add a properly named token to `@theme` if it is intentional/reused, or map to the nearest token with an uppercase `<!-- TODO: ... -->` comment if it is a one-off. Never hardcode the raw value.
- `[always]` Design-tool output (Figma MCP and similar) is never copied verbatim; run every emitted class through the token derivation above.
- `[always]` Hover styles via `hover:`/`group-hover:` variants, never JavaScript mouse events or `class:` toggling driven by `$state`.
- `[existing]` Pre-Tailwind projects use component-scoped CSS with BEM class names and CSS custom properties; follow that instead.

## Testing

- `[always]` Vitest, unit tests only: state files, entities, pure utils, service logic with fetch mocked. Never component-render or browser tests (see `javascript.md`).

<!-- oym-card:begin stack=framework-svelte -->
# OYM Svelte rules (all tools)

- Runes exclusively (`$state`/`$derived`/`$effect`/`$props`); no `$:`, no `export let`, no `svelte/store`. Events as `onclick`, not `on:click`.
- `<script lang="ts">` always.
- Props via `$props()` destructuring, typed inline (never a separate interface): each prop on its own line; optional props after required, `?:`, always with a default. Domain types imported from entities/types.
- PascalCase component files; check the existing primitives (`ui/` or `atoms/`) before writing a new one.
- Shared state in `[name]State.svelte.ts` runes files.
- Services throw, components catch and display; nothing unlogged.
- All UI strings through the project's i18n; no hardcoded copy.
- Tailwind (token projects): utilities inline; `@theme` in the main CSS is the authoritative token list; NO arbitrary bracket values (`gap-[24px]` becomes `gap-6`, `text-[32px]` becomes `text-32`); spacing = px/4 scale, combine identical axes; missing value: add a token or map to nearest + uppercase TODO comment; never copy design-tool output verbatim; hover via `hover:`/`group-hover:`, never JS mouse events.
- Pre-Tailwind projects: scoped CSS + BEM + custom properties; follow local style.
- Vitest unit tests for states/entities/utils only; no component-render tests.
<!-- oym-card:end stack=framework-svelte -->

