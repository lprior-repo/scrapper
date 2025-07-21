import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests/playwright',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: 'html',
  outputDir: 'test-results/',
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    // Additional browser configurations for comprehensive testing
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    // Mobile configurations
    {
      name: 'mobile-chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'mobile-safari',
      use: { ...devices['iPhone 12'] },
    },
  ],
  webServer: [
    {
      command: 'task api-start',
      port: 8081,
      reuseExistingServer: !process.env.CI,
      timeout: 120 * 1000,
    },
    {
      command: 'cd packages/webapp && bun --hot --port 3000 index.html',
      port: 3000,
      reuseExistingServer: !process.env.CI,
      timeout: 120 * 1000,
    },
  ],
  // Test match patterns to support new modular component structure
  testMatch: [
    '**/tests/playwright/**/*.spec.ts',
    '**/tests/playwright/**/*.test.ts',
  ],
  // Expectations for visual regression testing
  expect: {
    // Threshold for visual comparisons to handle minor rendering differences
    toHaveScreenshot: { threshold: 0.3 },
    toMatchSnapshot: { threshold: 0.3 },
  },
})
