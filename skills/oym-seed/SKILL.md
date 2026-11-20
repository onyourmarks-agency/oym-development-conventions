---
name: oym-seed
description: Give an On Your Marks Craft CMS site seed data, or repair its seed after the content model changed. Use when asked to add seed data, demo data, test data or "fixtures" to a Craft site, when `oym-seed check` or the `test:seed-check` CI job reports drift, or when a project needs a working database without a copy of staging or production. Symfony projects are not covered yet.
---

# OYM seed

Seed data is generated content that fills a site from an empty database: no client data, no copy of staging or production. The package `oym/craft-seed` generates most of it from the content model. Your job is the part it cannot know: the project's own field types, its custom tables, data that normally comes from a feed, and a recipe that makes the content fit the site.

Say **seed**, never "fixtures": in several OYM sports sites a fixture is a match.

The package README (`vendor/oym/craft-seed/README.md`) documents every command, option and recipe key. Read it before you write anything (in init mode it only exists after step 1); this skill only says what to do and in which order.

Run PHP, Composer, Craft and `vendor/bin/oym-seed` inside ddev (`ddev exec …`, `ddev composer …`), never on the host: the host has neither the PHP extensions nor the database hostname.

## Rules that always apply

- Invent the content. The one exception is the structure captured with `vendor/bin/oym-seed capture`: singles, and the fixed sections named in the recipe. Never copy news, people, users, comments or submissions, and never paste real content into seed code or the recipe by hand.
- A seeder never calls an external API. Data that normally comes from a feed (Opta, Sportmonks, FOYS, a sync command) is made up in a seeder.
- Never seed the developer's working database. Work in a scratch database (see "Scratch database").
- Change only seed paths: `config/oym-seed.php`, `config/oym-seed.structure.json`, `modules/*/seed/`, `seed.lock`, `.ddev/commands/web/seed`, and the `oym/craft-seed` lines in `composer.json`, `composer.lock`, `.gitignore` and the project config (the `plugins` key of `config/project/project.yaml`, or `config/project/plugins/plugins.yaml` where the config is split). If seeding exposes a bug elsewhere, report it; do not fix it in passing (see "When the project is in the way").
- Seed code never repairs the project: no seeder that creates a missing table or works around a broken template.
- The drift check decides whether something drifted, not you. Do not edit `seed.lock` by hand.
- Edit `.env` on the host, never inside the ddev container: the file sync can put an old version back. Change single lines with `sed`; do not open the file, it holds secrets.
- ddev syncs files between host and container with a delay. After a command in the container that writes files (install, plugin install, capture, check `--update`), run `ddev mutagen sync` before any `git` command or file check on the host.

## Before you start

All three must hold, otherwise stop and report; they are not seed problems.

1. `.env` exists in the project root. If it does not, look for an Infisical hook in `.ddev/config.yaml`; without one, create it with `infisical export --env=dev > .env` on the host. A missing `.env` shows as a fatal error in `EnvironmentLoader`, served with status 200.
2. `ddev composer install` succeeds.
3. `ddev exec php craft install/check` answers, and `git status --short config/project` is clean. If it is not clean, those changes are someone's work: do not lose them in the steps below.

## Which mode

Run `ddev exec vendor/bin/oym-seed check --json` in the project root.

| Result | Mode |
|---|---|
| The command does not exist | **init** |
| `hasLock` is false | **init**, from step 3 |
| `clean` is false | **fix** |
| `clean` is true | Nothing to do. Say so and stop. |

## Scratch database

OYM sites load `.env` in a way that overrides environment variables, so the database can only be switched in `.env` itself.

1. `ddev mysql -uroot -proot -e "CREATE DATABASE IF NOT EXISTS seed; GRANT ALL ON seed.* TO 'db'@'%';"`
2. In `.env` on the host set `CRAFT_DB_DATABASE='seed'` and give `REDIS_DATABASE` another number than the one it has (not 15: `capture` uses that one). Note the old values. Keep the table prefix (`CRAFT_DB_TABLE_PREFIX`) as it is; every table name in this skill then carries that prefix.
3. When you are done, or when you stop early, put the old values back and tell the user, unless the user asked to look at the seeded site. `ddev restart` only restores `.env` in projects whose `.ddev/config.yaml` has the Infisical hook.

Installing Craft into the empty scratch database:

```bash
ddev exec php craft install/craft --interactive=0 --email=seed-admin@example.test \
  --username=seed-admin@example.test --password='OnYourMarks01!' \
  --site-name="<name>" --site-url="<primary site url>" --language=<language>
ddev mutagen sync
git checkout -- config/project && git clean -fd config/project
ddev exec php craft project-config/apply --force
ddev exec php craft migrate/all --interactive=0
```

Take the site name and language from the primary site in `config/project/sites/` (or `project.yaml`), and the URL from the ddev project.

The `git checkout` and `git clean` matter: the installer rewrites the primary site's name and base URL in `config/project` and adds files of its own (a public GraphQL schema). Both commands are only safe because the project config was clean before you started. The `migrate/all` matters too: the installer does not run the project's own migrations.

If `project-config/apply --force` fails on a single with required fields ("Couldn't save single entry … due to validation errors"), restore `config/project` once more and run `project-config/apply` without `--force`. `ddev exec php craft project-config/diff` must then show no pending changes.

If the install fails for a reason that lies in the project or in a plugin, stop and report it.

## init

1. `ddev composer require oym/craft-seed` (under `require`, not `require-dev`: the plugin is in the project config, so it must be installed everywhere; it refuses to run on production). Add `/public/seed-assets` and `/storage/oym-seed` to `.gitignore` (use the project's real web root). Then run `ddev exec vendor/bin/oym-seed check --json`: a plugin under `stale` that Composer does not install will make the install in step 2 fail, so report it now.
2. Set up the scratch database, install Craft into it, then `ddev exec php craft plugin/install oym-seed`. Check that the only changes in `config/project` are the `oym-seed` entry under `plugins` and `dateModified`.
3. Run `vendor/bin/oym-seed check --json` and handle everything under `uncovered` (see "Writing seed code").
4. Write `config/oym-seed.php`: sensible counts for the sections that carry the site (news, events, products), `skip` for sections that make no sense to fill, and one `pinned` entry per section that end-to-end tests are likely to open. Keep it short; do not list sections that are fine with the default.
5. Add `.ddev/commands/web/seed` from `${CLAUDE_PLUGIN_ROOT}/templates/ddev/seed` if the project has none.
6. Make the structure match the live site, when the developer has a local copy of it: list the fixed sections (formats, categories, static pages, overview pages) under `structure.sections` in the recipe, run `ddev exec vendor/bin/oym-seed capture --database=<name of the local copy>`, and read `config/oym-seed.structure.json` before keeping it. Without a local copy, skip this step and say so.
   - Do not list sections that hold people, user input, or other organisations (clients, partners, sponsors): those are generated.
   - The copy must be on the current content model. If `capture` fails or the file misses fields, ask the developer to bring the copy up to date; do not write to it yourself.
   - Fields of plugins the package does not know are not captured. When menus, footer or buttons use such a field (a link field, for example), write a capture for it (`fieldCaptures` in the README), otherwise those links and labels are generated.
   - Name in your report every person the captured text mentions.
7. Run the seed and verify (see "Verify").
8. `ddev exec php craft oym-seed/check --update --strict` must exit 0. Commit `seed.lock` with the seed code.

## fix

1. Read the report. Open only the files it names: the entry types, fields or tables under `drift`, `uncovered` and `stale`.
2. Decide per item:

   | Report | What it needs |
   |---|---|
   | `drift.added` / `changed`, field type already covered | Usually nothing: the generator follows the model. Check whether the recipe should mention the new section. |
   | `drift.removed` | Remove what the recipe or a seeder still says about it. |
   | `uncovered.fieldTypes` | A field seeder, or `skipFieldTypes` with a one-line reason. |
   | `uncovered.tables` | A seeder under `tables`, or `false` when the table should stay empty (logs, queues, caches). |
   | `stale` | Remove or correct that recipe entry. |

3. Run the seed and verify.
4. `ddev exec php craft oym-seed/check --update --strict` must exit 0.

## Writing seed code

Seed code lives in the module that owns the thing it seeds: `modules/<name>module/seed/`, namespace `modules\<name>module\seed`.

**A field seeder** implements `oym\craft\seed\base\FieldSeederInterface` and is registered under `fieldSeeders`. Before writing one, read the field type class to learn what it stores: look at `normalizeValue()`, `serializeValue()` and the input template. Return a value in exactly that shape, or `null` to leave the field empty. Use `$context->faker` for anything random, so the output stays the same for the same seed.

**A field capture** implements `oym\craft\seed\base\FieldCaptureInterface` and is registered under `fieldCaptures`; the same class may also be the field's seeder. It stores entries by section and slug, never by id, and no personal data.

**A table or feed seeder** implements `oym\craft\seed\base\SeederInterface` and is registered under `seeders` (and under `tables` for each table it fills). It runs after all entries exist; register it under `seedersBefore` instead when entry fields point at its data. `purge()` must remove exactly what `seed()` created.

Skip a field type only when it cannot be filled before the entry is saved, or when it is display-only. Say why in a comment next to it in the recipe.

Follow the project's PHP conventions (`oym-php` skill): strict types, the OYM phpcs standard, no new dependencies.

## Verify

A seed is only done when all of this holds. Run each check; do not assume.

1. `ddev seed --fresh` ends without warnings and without "Field types without a seeder".
2. The site is in dev mode and an unknown URL returns 404. A site in production mode answers 200 on everything, which would make the next check worthless.
3. Every seeded entry URL returns 200 on every site, following the trailing-slash redirect. Take the URIs from `elements_sites.uri` (with the table prefix) of enabled entries that are no draft, no revision and not nested (`entries.primaryOwnerId IS NULL`), singles included. Build each URL from the base URL of its site, on that site's own hostname and language path; `__home__` is the base URL itself.
4. The web log of today (`storage/logs/web-<date>.log`) has no new error lines. Note the size of the file before the URL check and read only what was added; Craft 5 writes them as `[web.ERROR]`, older setups as `[error]`.
5. `public/seed-assets` holds the placeholder images. Nothing was uploaded to the real asset storage.
6. A second `ddev seed --fresh` gives the same slugs.
7. `vendor/bin/oym-seed check --strict` exits 0.
8. Where a structure was captured: menus, footer and sidebar show the same items as the live site, and the homepage and one page of each main type have the same blocks as their live counterparts. Compare with the public site; name every difference you leave.

If a page fails, fix the seed (a missing field seeder, a wrong value shape), not the template.

## When the project is in the way

A fresh install and generated content find bugs that a copy of production hides: a plugin in the project config that Composer no longer installs, a plugin without an install migration, a template that breaks on content the CMS allows.

- Do not fix them and do not hide them in seed code or the recipe. The one exception is `skip` for a section that cannot render at all (its template does not exist), with the reason in a comment.
- If the bug stops the install or the seed run, stop there. Otherwise finish the seed, leave the pages failing, and list every failing URL with its cause (file and line).
- A seed with failing checks is not done. Write `seed.lock` so the drift check works, but say in your report that the seed is not finished and what has to be fixed first.

## Report back

State what you seeded (entries, sections, sites), what you wrote and why, every obstacle you met with its cause, and anything you skipped. Say plainly when a check above did not pass.
