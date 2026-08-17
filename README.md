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
| `cmd/`, `internal/`, `main.go` | The `oym-conventions` Go CLI: wires projects and keeps the vendored cards current. Rules and templates are embedded in the binary at release time. |
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

Devin reads the product repo's `AGENTS.md`, which carries the vendored cards. No per-repo setup beyond running `oym-conventions init` once. Optionally add an org-level Knowledge entry per layer pointing at this repo.

### Codex

`AGENTS.md` is Codex's native mechanism, including nested files. The vendored cards work unmodified.

## The `oym-conventions` CLI

Install once (also used in CI):

```bash
curl -fsSL https://raw.githubusercontent.com/onyourmarks-agency/oym-development-conventions/main/install.sh | bash
```

Alternatives: `go install github.com/onyourmarks-agency/oym-development-conventions@latest` (binary lands as `oym-development-conventions`; rename it), or download an archive from the releases page.

The binary embeds the conventions and templates of the release it was built from; a release tag is a versioned snapshot of the rules.

```bash
oym-conventions init     # in a product repo: detects stacks, interactive confirm,
                         # writes/extends AGENTS.md, creates CLAUDE.md, offers
                         # quality-gate templates. Flags: --stacks, --mode, --yes.
oym-conventions sync     # refresh the vendored blocks (idempotent)
oym-conventions check    # exit 1 on stale/missing blocks (CI gate)
oym-conventions update   # self-update to the latest release
oym-conventions version
```

`init` and `sync` only rewrite content between the `oym-conventions` markers, preserve the `mode=` attribute, and stamp the tool version and sync date. Detection covers all project types in the table above; a Craft site with an `api-simple/` directory and a Vite frontend gets `shared,php,framework-craft,framework-api-simple,javascript,framework-svelte,tool-vite`.

Repos that already have a rich `AGENTS.md` keep their project rules on top; `init` appends the vendored blocks under an "Organization conventions" heading, and project rules override the blocks.

## Staying current

- On request: run `oym-conventions update && oym-conventions sync` in the repo.
- Monthly: create a scheduled GitLab pipeline with the jobs in `ci/gitlab-conventions-sync.yml`; the check variant warns on drift, the opt-in autosync variant opens a merge request with the refreshed cards.

## Changing the conventions

Rules change via PR to this repo. Keep the card in each file a summary that fits the 200-line budget (CI enforces this and runs the sync lifecycle tests). Version bumps to `VERSIONS.md` need no other edits. Tag releases with SemVer; product repos pick up changes with `oym-conventions update && oym-conventions sync`, or pin by installing a specific release's binary from the releases page.
