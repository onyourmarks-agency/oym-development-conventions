// OYM Vitest config, Vite-in-Craft family.
// Unit tests only: states, module/service classes, pure utils under resources/js.
// No component-render tests, no browser environment.
import { defineConfig } from 'vitest/config';
import path from 'node:path';

export default defineConfig({
  resolve: {
    alias: {
      $src: path.resolve(__dirname, 'resources'),
      $js: path.resolve(__dirname, 'resources/js'),
      $css: path.resolve(__dirname, 'resources/css'),
    },
  },
  test: {
    include: ['resources/js/**/*.spec.ts'],
    environment: 'node',
    globals: false,
  },
});
