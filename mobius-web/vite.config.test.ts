import { defineConfig } from 'vitest/config';
import { sveltekit } from '@sveltejs/kit/vite';

export default defineConfig({
  plugins: [sveltekit()],
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}'],
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/setupTests.ts'],
    // Ensure tests run in browser-like environment
    environmentOptions: {
      jsdom: {
        resources: 'usable'
      }
    },
    pool: 'forks'
  }
});
