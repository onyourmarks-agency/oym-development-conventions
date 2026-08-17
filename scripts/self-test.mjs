#!/usr/bin/env node
/**
 * Self-test for sync-conventions.mjs: init, idempotency, and drift detection
 * against a scratch AGENTS.md built from the template.
 */

import { execFileSync } from 'node:child_process';
import { copyFileSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const repoRoot = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const syncScript = join(repoRoot, 'scripts', 'sync-conventions.mjs');
const scratchDirectory = mkdtempSync(join(tmpdir(), 'oym-conventions-test-'));

const runSync = (args, { expectFailure = false } = {}) => {
    try {
        return execFileSync('node', [syncScript, '--source', repoRoot, '--target', scratchDirectory, ...args], {
            encoding: 'utf8',
        });
    } catch (error) {
        if (expectFailure) {
            return error;
        }
        throw error;
    }
};

const assert = (condition, message) => {
    if (!condition) {
        console.error(`FAIL: ${message}`);
        rmSync(scratchDirectory, { recursive: true, force: true });
        process.exit(1);
    }
    console.log(`ok: ${message}`);
};

try {
    copyFileSync(join(repoRoot, 'templates', 'AGENTS.template.md'), join(scratchDirectory, 'AGENTS.md'));

    runSync(['--stacks', 'php-universal,php-craft']);
    const firstPass = readFileSync(join(scratchDirectory, 'AGENTS.md'), 'utf8');
    assert(firstPass.includes('# OYM shared rules'), 'shared card vendored from template block');
    assert(firstPass.includes('# OYM PHP rules'), 'php-universal block appended and filled');
    assert(firstPass.includes('# OYM Craft CMS site rules'), 'php-craft block appended and filled');
    assert(firstPass.includes('mode=existing'), 'mode attribute preserved');

    runSync([]);
    const secondPass = readFileSync(join(scratchDirectory, 'AGENTS.md'), 'utf8');
    assert(firstPass === secondPass, 'second run is idempotent');

    const checkClean = runSync(['--check']);
    assert(checkClean.includes('up to date'), 'check passes when cards are current');

    writeFileSync(join(scratchDirectory, 'AGENTS.md'), secondPass.replace('# OYM PHP rules', '# tampered'));
    const checkDrift = runSync(['--check'], { expectFailure: true });
    assert(checkDrift.status === 1, 'check exits 1 on drift');

    runSync([]);
    const repaired = readFileSync(join(scratchDirectory, 'AGENTS.md'), 'utf8');
    assert(repaired.includes('# OYM PHP rules'), 'sync repairs drifted block');

    console.log('self-test passed');
} finally {
    rmSync(scratchDirectory, { recursive: true, force: true });
}
