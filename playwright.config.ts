import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests/playwright',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: [
    {
      command: 'task api-start',
      port: 8081,
      reuseExistingServer: !process.env.CI,
    },
    {
      command: 'cd packages/webapp && bun run dev',
      port: 3000,
      reuseExistingServer: !process.env.CI,
    },
  ],
})
