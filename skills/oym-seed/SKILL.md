---
name: oym-seed
description: Give an On Your Marks Craft CMS site seed data, or repair its seed after the content model changed. Use when asked to add seed data, demo data, test data or "fixtures" to a Craft site, when `oym-seed check` or the `test:seed-check` CI job reports drift, or when a project needs a working database without a copy of staging or production. Symfony projects are not covered yet.
---

# OYM seed

Seed data is generated content that fills a site from an empty database: no client data, no copy of staging or production. The package `oym/craft-seed` generates most of it from the content model. Your job is the part it cannot know: the project's own field types, its custom tables, data that normally comes from a feed, and a recipe that makes the content fit the site.

Say **seed**, never "fixtures": in several OYM sports sites a fixture is a match.

The package README (`vendor/oym/craft-seed/README.md`) documents every command, option and recipe key. Read it before you start; this skill only says what to do and in which order.

## Rules that always apply

- Invent the content. The one exception is the structure captured with `vendor/bin/oym-seed capture`: singles, and the fixed sections named in the recipe. Never copy news, people, users, comments or submissions, and never paste real content into seed code or the recipe by hand.
- A seeder never calls an external API. Data that normally comes from a feed (Opta, Sportmonks, FOYS, a sync command) is made up in a seeder.
- Never seed the developer's working database. Work in a scratch database (see "Scratch database").
- Change only seed paths: `config/oym-seed.php`, `config/oym-seed.structure.json`, `modules/*/seed/`, `seed.lock`, `.ddev/commands/web/seed`, and the `oym/craft-seed` lines in `composer.json`, `composer.lock`, `config/project/project.yaml` and `.gitignore`. If seeding exposes a bug elsewhere, report it; do not fix it in passing.
- The drift check decides whether something drifted, not you. Do not edit `seed.lock` by hand.
- Edit `.env` on the host, never inside the ddev container: the file sync can put an old version back.

## Which mode

Run `vendor/bin/oym-seed check --json` in the project root.

| Result | Mode |
|---|---|
| The command does not exist | **init** |
| `hasLock` is false | **init**, from step 3 |
| `clean` is false | **fix** |
| `clean` is true | Nothing to do. Say so and stop. |

## Scratch database

OYM sites load `.env` in a way that overrides environment variables, so the database can only be switched in `.env` itself.

1. `ddev mysql -uroot -proot -e "CREATE DATABASE IF NOT EXISTS seed; GRANT ALL ON seed.* TO 'db'@'%';"`
2. In `.env` on the host set `CRAFT_DB_DATABASE='seed'` and give `REDIS_DATABASE` a number the project does not use. Note the old values.
3. When you are done, put the old values back and tell the user. `ddev restart` also restores `.env` in projects that generate it from Infisical.

Installing Craft into the empty scratch database:

```bash
ddev exec php craft install/craft --interactive=0 --email=seed-admin@example.test \
  --username=seed-admin@example.test --password='OnYourMarks01!' \
  --site-name="<name>" --site-url="<primary site url>" --language=<language>
git checkout -- config/project
ddev exec php craft project-config/apply --force
ddev exec php craft migrate/all --interactive=0
```

The `git checkout` matters: the installer rewrites the primary site's name and base URL in `config/project`. The `migrate/all` matters too: the installer does not run the project's own migrations.

If the install fails for a reason that lies in the project or in a plugin, stop and report it.

## init

1. `ddev composer require oym/craft-seed`, add `/public/seed-assets` and `/storage/oym-seed` to `.gitignore` (use the project's real web root).
2. Set up the scratch database, install Craft into it, then `ddev exec php craft plugin/install oym-seed`. Check that the only change in `config/project` is the `oym-seed` entry under `plugins`.
3. Run `vendor/bin/oym-seed check --json` and handle everything under `uncovered` (see "Writing seed code").
4. Write `config/oym-seed.php`: sensible counts for the sections that carry the site (news, events, products), `skip` for sections that make no sense to fill, and one `pinned` entry per section that end-to-end tests are likely to open. Keep it short; do not list sections that are fine with the default.
5. Add `.ddev/commands/web/seed` from `${CLAUDE_PLUGIN_ROOT}/templates/ddev/seed` if the project has none.
6. Make the structure match the live site, when the developer has a local copy of it: list the fixed sections (formats, categories, static pages, overview pages) under `structure.sections` in the recipe, run `vendor/bin/oym-seed capture --database=<name of the local copy>`, and read `config/oym-seed.structure.json` before keeping it. Do not list sections that hold people or user input. Without a local copy, skip this step and say so.
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

**A table or feed seeder** implements `oym\craft\seed\base\SeederInterface` and is registered under `seeders` (and under `tables` for each table it fills). It runs after all entries exist; register it under `seedersBefore` instead when entry fields point at its data. `purge()` must remove exactly what `seed()` created.

Skip a field type only when it cannot be filled before the entry is saved, or when it is display-only. Say why in a comment next to it in the recipe.

Follow the project's PHP conventions (`oym-php` skill): strict types, the OYM phpcs standard, no new dependencies.

## Verify

A seed is only done when all of this holds. Run each check; do not assume.

1. `ddev seed --fresh` ends without warnings and without "Field types without a seeder".
2. The site is in dev mode and an unknown URL returns 404. A site in production mode answers 200 on everything, which would make the next check worthless.
3. Every seeded entry URL returns 200 on every site, following the trailing-slash redirect. Take the URLs from `elements_sites.uri` of the seeded entries.
4. The web log of today (`storage/logs/web-<date>.log`) has no new `[error]` lines.
5. `public/seed-assets` holds the placeholder images. Nothing was uploaded to the real asset storage.
6. A second `ddev seed --fresh` gives the same slugs.
7. `vendor/bin/oym-seed check --strict` exits 0.
8. Where a structure was captured: menus, footer and sidebar show the same items as the live site, and the homepage and one page of each main type have the same blocks as their live counterparts. Compare with the public site; name every difference you leave.

If a page fails, fix the seed (a missing field seeder, a wrong value shape), not the template. A template that breaks on valid content is a bug to report.

## Report back

State what you seeded (entries, sections, sites), what you wrote and why, every obstacle you met with its cause, and anything you skipped. Say plainly when a check above did not pass.
