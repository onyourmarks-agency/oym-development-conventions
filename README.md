# oym-development-conventions

Development conventions for On Your Marks, written for both humans and AI coding tools (Claude Code, Devin, Codex). One source of truth, three delivery mechanisms.

## The layered model

Conventions are composable layers, not per-project silos. A project reads `shared` + its language layer + its framework + (for frontends) its tool:

| Project type | Layers |
|---|---|
| Craft CMS site | `shared` + `php` + `framework-craft` (+ `framework-api-simple` per directory) |
| Craft plugin | `shared` + `php` + `framework-craft-plugin` |
| Symfony API | `shared` + `php` + `framework-symfony` |
| NestJS service | `shared` + `javascript` + `framework-nestjs` |
| SvelteKit app | `shared` + `javascript` + `framework-svelte` + `tool-sveltekit` |
| Vite islands in Craft | `shared` + `javascript` + `framework-svelte` + `tool-vite` |
| Webpack islands (legacy) | `shared` + `javascript` + `framework-svelte` + `tool-webpack` |

NestJS deliberately shares the `javascript` layer with the front-ends: typing rules, brace style, naming, and comment policy are identical across Svelte and Nest code.

## Layout

| Path | What it is |
|---|---|
| `conventions/` | The rules, one file per layer, tool-agnostic. Each file embeds a condensed "card" between `oym-card` markers; the card is what gets vendored into product repos. |
| `VERSIONS.md` | The only place version numbers live. Rules never hardcode versions; product repos pin their own in their manifests. |
| `skills/` | Claude Code skills (thin routers into `conventions/`). |
| `templates/` | Drop-in configs for product repos: PHPStan, PHPUnit, Vitest, NestJS strictness deltas, `AGENTS.md`/`CLAUDE.md` starters. |
| `scripts/sync-conventions.mjs` | Vendors convention cards into a product repo's `AGENTS.md`. |
| `.claude-plugin/` | Claude Code plugin and marketplace manifests. |

## The mode model

Every rule is tagged `[always]`, `[new]`, or `[existing]`.

- **New projects** get the full standard from day one: PHPStan max, unit tests wired, current patterns.
- **Existing projects** bind the AI to local context: the surrounding module outranks this repo, no drive-by retrofits, tooling introduced with baselines and ratchets.

The mode is declared per block in each repo's `AGENTS.md` (`mode=existing|new`); tools fall back to heuristics and default to `existing`.

## Using it per tool

### Claude Code

Install once, works in every repo:

```
/plugin marketplace add onyourmarks-agency/oym-development-conventions
/plugin install oym-conventions@oym
```

The skills (`oym-php`, `oym-nestjs`, `oym-frontend`, `oym-new-project`) trigger automatically on matching work, detect the framework/tool from the repo, and read the right layer files. Product repos keep a one-line `CLAUDE.md` (`@AGENTS.md`) so Claude also sees the vendored cards and project-specific rules.

### Devin

Devin reads the product repo's `AGENTS.md`, which carries the vendored cards. No per-repo setup beyond running the sync script once. Optionally add an org-level Knowledge entry per layer pointing at this repo.

### Codex

`AGENTS.md` is Codex's native mechanism, including nested files. The vendored cards work unmodified.

## Wiring a product repo

```bash
# once, in the product repo
cp <this-repo>/templates/AGENTS.template.md AGENTS.md   # fill in the project sections
cp <this-repo>/templates/CLAUDE.template.md CLAUDE.md
node <this-repo>/scripts/sync-conventions.mjs --stacks shared,php,framework-craft --target .

# refresh after conventions change
node <this-repo>/scripts/sync-conventions.mjs --target .

# CI drift gate
node <this-repo>/scripts/sync-conventions.mjs --check --target .
```

The script only rewrites content between the `oym-conventions` markers, preserves the `mode=` attribute, and stamps the source revision and sync date. Pick the layers from the table above; a Craft site with a Vite frontend carries `shared,php,framework-craft,framework-api-simple,javascript,framework-svelte,tool-vite`.

Repos with an existing rich `AGENTS.md` (the NestJS agent services, the newest front-ends) keep their project rules on top and append the vendored blocks under an "Organization conventions" heading; project rules override the blocks.

## Changing the conventions

Rules change via PR to this repo. Keep the card in each file a summary that fits the 200-line budget (CI enforces this and runs the sync self-test). Version bumps to `VERSIONS.md` need no other edits. Tag releases with SemVer; product repos pick up changes on their next sync run, or pin with `--ref`.
