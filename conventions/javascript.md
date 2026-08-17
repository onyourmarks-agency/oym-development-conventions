# JavaScript/TypeScript conventions (language layer)

Applies to all JavaScript and TypeScript in OYM repositories: Svelte front-ends AND NestJS back-ends share this layer. Framework files (`framework-svelte.md`, `framework-nestjs.md`) and tool files (`tool-vite.md`, `tool-webpack.md`, `tool-sveltekit.md`) build on it.

## TypeScript

- `[always]` Everything is TypeScript. `strict: true` stays on; never weaken compiler flags to make an error go away.
- `[always]` No `any`: use `unknown` plus narrowing, generics, or proper types. This holds even in repos whose legacy compiler flags would allow `any` (that is a defect, not permission). Remove `any` from files you touch.
- `[always]` Explicit return types on exported functions and methods, including `void` and `never`.
- `[always]` Closed value sets are union literal types (`'sm' | 'md' | 'lg'`), never loose strings.
- `[always]` Non-null assertions (`!`) need a reason comment or a guard instead.
- `[always]` Prefer `??` over `||` for defaults; respect the repo's lint rules on nullish handling.

## Style

- `[always]` Functions are arrow expressions assigned to `const`: `const fnName = () => {}`, never `function fnName()`.
- `[always]` Statement bodies always get braces, even single-line:

  ```javascript
  if (variable) {
  	return;
  }
  ```

  A concise arrow body is exempt and stays concise, also when it wraps an assignment: `const method = () => (a = b);`.
- `[always]` Named exports by default; a default export needs a specific reason (framework files list their exceptions).
- `[always]` Full descriptive names. No `o`, `a`, `n`, `tmp`, `idx`; write `index`, `option`, `response`.
- `[always]` Imports sorted per the repo's lint setup (`simple-import-sort` where configured); respect existing group order.
- `[always]` Errors: never swallow silently. Throw with a meaningful message or coded error string, log where you can act. Framework files define where throwing vs catching happens.
- `[always]` Comment policy from `shared.md` applies with full force to TS/Svelte, where no linter enforces it: no filler JSDoc, no comments restating code. JSDoc only when it adds contract information the types cannot carry.

## Formatting

- `[always]` The repo's own Prettier config is authoritative. Run the project's format script; never restyle by hand and never carry a style between repos. Tool files document what each family's config looks like so you recognize it; the config file itself always wins.
- `[always]` ESLint is flat config. Fix lint findings, do not disable rules; a `// eslint-disable` needs a reason comment and a narrow scope.

## Runtime and packaging

- `[always]` pnpm (`packageManager` pinned per repo), ESM (`"type": "module"`), Node version per the repo's `engines`/`.nvmrc` (org target: `VERSIONS.md`).
- `[always]` Husky + lint-staged pre-commit hooks are part of the contract; never `--no-verify`.
- `[always]` The aggregate check script (`frontend-cs` in most repos, or the repo's `lint` + `check`/`test` scripts) runs before finishing. In DDEV projects, prefix with `ddev`: `ddev pnpm frontend-cs`.

## Testing

- `[always]` Unit tests only (see `shared.md`). Runner per framework: Vitest for Svelte projects, Jest for NestJS. Test entities, state, services, pure utils; mock at boundaries (fetch, repositories); no component-render, browser, or e2e tests.
- `[existing]` Most JS projects have no suite. Add specs alongside new logic-bearing code; configs in `templates/`.

<!-- oym-card:begin stack=javascript -->
# OYM JavaScript/TypeScript rules (shared by Svelte and NestJS)

- TypeScript everywhere, `strict: true`, never weaken flags.
- No `any` (use `unknown`/generics), even where legacy flags allow it; remove `any` from files you touch. Explicit return types on exports.
- Union literal types for closed value sets. Non-null assertions need a reason.
- Functions: `const fnName = () => {}`, never `function fnName()`.
- Braces on every statement body, single-line included; concise arrow bodies exempt (`const method = () => (a = b);`).
- Named exports by default; default exports need a specific reason.
- Full descriptive names: `index` not `i`, `response` not `r`.
- Never swallow errors: throw meaningful/coded errors, log where you act.
- Comments: non-obvious "why" only; no filler JSDoc; do not touch comments in unrelated code.
- The repo's Prettier config is authoritative; run the format script, never restyle by hand or import another repo's style. ESLint flat config; disables need a reason and narrow scope.
- pnpm + ESM; Node per the repo's `engines`/`.nvmrc`.
- Run the aggregate check before finishing (`frontend-cs` or lint + check; `ddev pnpm ...` in DDEV repos). Never `--no-verify`.
- Unit tests only: Vitest (Svelte) / Jest (NestJS); mock at boundaries; no render/browser/e2e tests. Add specs with new logic.
<!-- oym-card:end stack=javascript -->

