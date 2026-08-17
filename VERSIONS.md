# Current version targets

The only place in this repo where versions live. Conventions and skills refer here; product repos pin their own versions in `composer.json`, `package.json` (`engines`, `packageManager`), and `.nvmrc`, which are always the runtime truth for that repo.

Update this table when the org moves; nothing else in this repo needs touching.

| Target | Current | Notes |
|---|---|---|
| PHP | 8.5 | New projects; existing sites range 8.0 to 8.5 |
| Craft CMS | 5 | Boilerplate: `oym-craft5-boilerplate` |
| oym/api-simple | 2.x | Composer package via Repman |
| Symfony | 7 | With API Platform 4 |
| PHPStan | latest, level max | See `conventions/php.md` for the baseline policy |
| Node | 24 | `engines.node` + `.nvmrc` in every JS repo |
| pnpm | 10/11 | Pinned per repo via `packageManager` |
| TypeScript | 5/6 | Newest projects on 6 |
| Svelte | 5 (runes) | |
| SvelteKit | 2 | `adapter-node` |
| Vite | 7/8 | |
| Tailwind | 4 (CSS-first) | `@theme` in CSS, no tailwind.config.js |
| NestJS | 11 | |
| Prisma | 7 | With the MariaDB adapter |
| ESLint | 9/10 | Flat config only |
| GitLab CI, PHP deploy flow | `@4.0.8` | `oym-gitlab-deploy-templates` components (build-vendors, test-phpcs, deploy) |
| GitLab CI, docker flow | `@1.0.8` | Components `stage-build-image`, `stage-tag-environment` |
