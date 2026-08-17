# oym-development-conventions

Development conventions for On Your Marks, written for both humans and AI coding tools (Claude Code, Devin, Codex). One source of truth, three delivery mechanisms.

## Layout

| Path | What it is |
|---|---|
| `conventions/` | The rules. One markdown file per stack, tool-agnostic. Each file embeds a condensed "card" between `oym-card` markers; the card is what gets vendored into product repos. |
| `skills/` | Claude Code skills (thin routers into `conventions/`). |
| `templates/` | Drop-in configs for product repos: `phpstan.neon.dist`, `phpunit.xml.dist`, Vitest configs, NestJS strictness deltas, `AGENTS.md`/`CLAUDE.md` starters. |
| `scripts/sync-conventions.mjs` | Vendors convention cards into a product repo's `AGENTS.md`. |
| `.claude-plugin/` | Claude Code plugin and marketplace manifests. |

Stacks: `shared`, `php-universal`, `php-craft`, `php-craft-plugin`, `php-api-simple`, `php-symfony`, `nestjs`, `svelte-universal`, `svelte-vite-craft`, `sveltekit`.

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

The skills (`oym-php`, `oym-nestjs`, `oym-svelte`, `oym-new-project`) trigger automatically on matching work and read the full convention files. Product repos keep a one-line `CLAUDE.md` (`@AGENTS.md`) so Claude also sees the vendored cards and project-specific rules.

### Devin

Devin reads the product repo's `AGENTS.md`, which carries the vendored cards. No per-repo setup beyond running the sync script once. Optionally add an org-level Knowledge entry per stack pointing at this repo.

### Codex

`AGENTS.md` is Codex's native mechanism, including nested files. The vendored cards work unmodified.

## Wiring a product repo

```bash
# once, in the product repo
cp <this-repo>/templates/AGENTS.template.md AGENTS.md   # fill in the project sections
cp <this-repo>/templates/CLAUDE.template.md CLAUDE.md
node <this-repo>/scripts/sync-conventions.mjs --stacks shared,php-universal,php-craft --target .

# refresh after conventions change
node <this-repo>/scripts/sync-conventions.mjs --target .

# CI drift gate
node <this-repo>/scripts/sync-conventions.mjs --check --target .
```

The script only rewrites content between the `oym-conventions` markers, preserves the `mode=` attribute, and stamps the source revision and sync date. Pick the stacks that apply; a Craft site with a Svelte frontend carries `shared,php-universal,php-craft,php-api-simple,svelte-universal,svelte-vite-craft`.

Repos with an existing rich `AGENTS.md` (unisport-agent, oym-chant-fe) keep their project rules on top and append the vendored blocks under an "Organization conventions" heading; project rules override the blocks.

## Changing the conventions

Rules change via PR to this repo. Keep the card in each file a summary that fits the 200-line budget (CI enforces this and runs the sync self-test). Tag releases with SemVer; product repos pick up changes on their next sync run, or pin with `--ref`.
