// OYM NestJS eslint delta: drop into the flat config's rules object.
// Replaces the legacy relaxations (`no-explicit-any: off`, `no-floating-promises: warn`)
// that older agent services still carry. New services start with this; existing
// services adopt it together with the tsconfig strictness migration.
export const oymNestjsRules = {
  '@typescript-eslint/no-explicit-any': 'error',
  '@typescript-eslint/no-floating-promises': 'error',
  '@typescript-eslint/no-unsafe-argument': 'error',
  '@typescript-eslint/explicit-function-return-type': [
    'error',
    { allowExpressions: true },
  ],
};
