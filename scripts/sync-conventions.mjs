#!/usr/bin/env node
/**
 * Vendors OYM convention cards into a product repo's AGENTS.md.
 *
 * Source of truth: conventions/<stack>.md in oym-development-conventions, where each
 * file carries a condensed card between oym-card markers. This script copies that card
 * into the matching oym-conventions marker block in the target repo's AGENTS.md.
 *
 * Usage:
 *   node sync-conventions.mjs                              sync all blocks present in ./AGENTS.md
 *   node sync-conventions.mjs --target ../eredivisie       sync another repo
 *   node sync-conventions.mjs --stacks php-universal,php-craft   append missing blocks, then sync
 *   node sync-conventions.mjs --check                      exit 1 when cards drifted, write nothing
 *   node sync-conventions.mjs --source /path/to/checkout   use a local conventions checkout
 *   node sync-conventions.mjs --ref v1.0.0                 pin the cloned source to a ref
 *
 * Without --source the script uses its own repo checkout when run from one, and
 * otherwise clones the conventions repo into ~/.cache/oym-conventions.
 */

import { execFileSync } from 'node:child_process';
import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { homedir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { parseArgs } from 'node:util';
import { fileURLToPath } from 'node:url';

const CONVENTIONS_REPO_URL = 'https://github.com/onyourmarks-agency/oym-development-conventions.git';
const CACHE_DIRECTORY = join(homedir(), '.cache', 'oym-conventions');
const AGENTS_FILENAME = 'AGENTS.md';
const DEFAULT_MODE = 'existing';
const EXIT_DRIFT = 1;
const EXIT_USAGE = 2;

const { values: options } = parseArgs({
    options: {
        source: { type: 'string' },
        target: { type: 'string', default: '.' },
        stacks: { type: 'string' },
        ref: { type: 'string' },
        check: { type: 'boolean', default: false },
        help: { type: 'boolean', default: false },
    },
});

if (options.help) {
    console.log('See the header of this file for usage.');
    process.exit(0);
}

const fail = (message) => {
    console.error(`sync-conventions: ${message}`);
    process.exit(EXIT_USAGE);
};

const runGit = (args, cwd) =>
    execFileSync('git', args, { cwd, encoding: 'utf8', stdio: ['ignore', 'pipe', 'pipe'] }).trim();

const resolveSourceDirectory = () => {
    if (options.source) {
        const sourceDirectory = resolve(options.source);
        if (!existsSync(join(sourceDirectory, 'conventions'))) {
            fail(`--source ${sourceDirectory} has no conventions/ directory`);
        }
        return sourceDirectory;
    }

    const scriptRepoRoot = resolve(dirname(fileURLToPath(import.meta.url)), '..');
    if (existsSync(join(scriptRepoRoot, 'conventions'))) {
        return scriptRepoRoot;
    }

    mkdirSync(dirname(CACHE_DIRECTORY), { recursive: true });
    if (!existsSync(join(CACHE_DIRECTORY, '.git'))) {
        runGit(['clone', '--depth', '1', CONVENTIONS_REPO_URL, CACHE_DIRECTORY]);
    } else {
        runGit(['-C', CACHE_DIRECTORY, 'fetch', '--depth', '1', 'origin']);
        runGit(['-C', CACHE_DIRECTORY, 'reset', '--hard', 'origin/HEAD']);
    }
    if (options.ref) {
        runGit(['-C', CACHE_DIRECTORY, 'fetch', '--depth', '1', 'origin', options.ref]);
        runGit(['-C', CACHE_DIRECTORY, 'checkout', options.ref]);
    }
    return CACHE_DIRECTORY;
};

const sourceRevision = (sourceDirectory) => {
    try {
        return runGit(['-C', sourceDirectory, 'rev-parse', '--short', 'HEAD']);
    } catch {
        return 'unknown';
    }
};

const extractCard = (sourceDirectory, stack) => {
    const conventionPath = join(sourceDirectory, 'conventions', `${stack}.md`);
    if (!existsSync(conventionPath)) {
        fail(`unknown stack "${stack}": ${conventionPath} does not exist`);
    }
    const conventionContent = readFileSync(conventionPath, 'utf8');
    const cardPattern = new RegExp(
        `<!-- oym-card:begin stack=${stack} -->\\n([\\s\\S]*?)<!-- oym-card:end stack=${stack} -->`,
    );
    const cardMatch = conventionContent.match(cardPattern);
    if (!cardMatch) {
        fail(`conventions/${stack}.md has no oym-card block`);
    }
    return cardMatch[1].trimEnd() + '\n';
};

const BLOCK_PATTERN =
    /<!-- oym-conventions:begin stack=([a-z0-9-]+)((?:\s+[a-z]+=[^\s>]+)*)\s*-->\n?([\s\S]*?)<!-- oym-conventions:end stack=\1 -->/g;

const parseAttributes = (attributeText) => {
    const attributes = {};
    for (const attributeMatch of attributeText.matchAll(/([a-z]+)=([^\s>]+)/g)) {
        attributes[attributeMatch[1]] = attributeMatch[2];
    }
    return attributes;
};

const main = () => {
    const sourceDirectory = resolveSourceDirectory();
    const revision = sourceRevision(sourceDirectory);
    const today = new Date().toISOString().slice(0, 10);
    const targetDirectory = resolve(options.target);
    const agentsPath = join(targetDirectory, AGENTS_FILENAME);

    if (!existsSync(agentsPath)) {
        fail(`${agentsPath} does not exist; create it from templates/AGENTS.template.md first`);
    }

    let agentsContent = readFileSync(agentsPath, 'utf8');

    const requestedStacks = options.stacks ? options.stacks.split(',').map((stack) => stack.trim()).filter(Boolean) : [];
    const presentStacks = [...agentsContent.matchAll(BLOCK_PATTERN)].map((blockMatch) => blockMatch[1]);
    const missingStacks = requestedStacks.filter((stack) => !presentStacks.includes(stack));

    if (missingStacks.length > 0 && options.check) {
        console.error(`drift: missing blocks for ${missingStacks.join(', ')}`);
        process.exit(EXIT_DRIFT);
    }
    for (const stack of missingStacks) {
        agentsContent =
            agentsContent.trimEnd() +
            `\n\n<!-- oym-conventions:begin stack=${stack} mode=${DEFAULT_MODE} -->\n` +
            `<!-- oym-conventions:end stack=${stack} -->\n`;
    }

    const driftedStacks = [];
    const updatedContent = agentsContent.replace(
        BLOCK_PATTERN,
        (fullBlock, stack, attributeText, currentBody) => {
            const card = extractCard(sourceDirectory, stack);
            const attributes = parseAttributes(attributeText);
            const mode = attributes.mode ?? DEFAULT_MODE;
            if (currentBody.trimEnd() + '\n' === card) {
                return fullBlock;
            }
            driftedStacks.push(stack);
            const header = `<!-- oym-conventions:begin stack=${stack} mode=${mode} rev=${revision} synced=${today} -->`;
            return `${header}\n${card}<!-- oym-conventions:end stack=${stack} -->`;
        },
    );

    if (options.check) {
        if (driftedStacks.length > 0) {
            console.error(`drift: stale blocks for ${driftedStacks.join(', ')} (run sync-conventions.mjs)`);
            process.exit(EXIT_DRIFT);
        }
        console.log('cards up to date');
        return;
    }

    if (updatedContent === readFileSync(agentsPath, 'utf8')) {
        console.log('cards up to date, nothing written');
        return;
    }

    writeFileSync(agentsPath, updatedContent);
    const syncedStacks = [...new Set([...driftedStacks, ...missingStacks])];
    console.log(`synced ${syncedStacks.join(', ')} into ${agentsPath} (source ${revision})`);
};

main();
