#!/usr/bin/env node
/**
 * Validates that every conventions/*.md contains exactly one well-formed card block
 * whose stack attribute matches the filename, within the size budget. Run in CI.
 */

import { readdirSync, readFileSync } from 'node:fs';
import { basename, dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const MAX_CARD_LINES = 200;
const conventionsDirectory = resolve(dirname(fileURLToPath(import.meta.url)), '..', 'conventions');

const problems = [];

for (const filename of readdirSync(conventionsDirectory).filter((name) => name.endsWith('.md'))) {
    const stack = basename(filename, '.md');
    const content = readFileSync(join(conventionsDirectory, filename), 'utf8');

    const beginMarker = `<!-- oym-card:begin stack=${stack} -->`;
    const endMarker = `<!-- oym-card:end stack=${stack} -->`;
    const beginCount = content.split(beginMarker).length - 1;
    const endCount = content.split(endMarker).length - 1;

    if (beginCount !== 1 || endCount !== 1) {
        problems.push(`${filename}: expected exactly one card block for stack=${stack} (found ${beginCount} begin, ${endCount} end)`);
        continue;
    }

    const cardBody = content.split(beginMarker)[1].split(endMarker)[0];
    const cardLineCount = cardBody.trim().split('\n').length;
    if (cardLineCount > MAX_CARD_LINES) {
        problems.push(`${filename}: card is ${cardLineCount} lines, budget is ${MAX_CARD_LINES}`);
    }
}

if (problems.length > 0) {
    for (const problem of problems) {
        console.error(problem);
    }
    process.exit(1);
}
console.log('all convention cards valid');
