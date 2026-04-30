import { defineConfig, devices } from '@playwright/test';

/**
 * Playwright config for Media Collector E2E tests.
 *
 * Tests run against the Vite dev server. The Wails Go bindings won't be
 * available in pure-browser mode, so we mock them in the test fixtures.
 * For full integration tests, run `wails dev` and point baseURL at the
 * Wails dev server (http://localhost:34115).
 */
export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 1,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['html', { open: 'never' }],
    ['list'],
  ],
  timeout: 30_000,
  expect: {
    timeout: 10_000,
  },

  use: {
    /* Tests hit the Vite dev server by default */
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        launchOptions: {
          // Use locally installed Playwright Chromium
          executablePath: undefined,
        },
      },
    },
  ],

  /* Start Vite dev server before tests */
  webServer: {
    command: 'npm run dev -- --port 5173',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
    timeout: 30_000,
  },
});
