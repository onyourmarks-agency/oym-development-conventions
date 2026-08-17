# AGENTS.md

<!--
  Starter for OYM product repos.
  1. Fill in the project summary and project-specific rules below.
  2. Set mode=existing or mode=new on each convention block.
  3. Keep the marker blocks; their content is machine-managed by
     scripts/sync-conventions.mjs from oym-development-conventions.
     Init/refresh: node scripts/sync-conventions.mjs --stacks <stack,...> --target .
-->

## What this project is

<!-- One paragraph: what the repo does, entry points, how to run it (ddev/pnpm commands). -->

## Project-specific rules

<!-- Only true deltas: contracts to preserve, files not to touch, local quirks.
     Rules here override the organization conventions below. -->

## Organization conventions

Rules below are vendored from `oym-development-conventions`. Project-specific rules above override them. Do not edit inside the marker blocks.

<!-- oym-conventions:begin stack=shared mode=existing -->
<!-- oym-conventions:end stack=shared -->

<!-- Add one empty block pair per applicable stack (same syntax as the shared block
     above, with stack=php-universal, stack=php-craft, ...), or let the sync script
     append them via its stacks flag: sync-conventions.mjs with stacks=php-universal,php-craft -->
