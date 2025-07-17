export default {
  packageManager: 'bun',
  reporters: ['html', 'clear-text', 'progress'],
  testRunner: 'jest',
  coverageAnalysis: 'perTest',
  jest: {
    projectType: 'custom',
    configFile: 'jest.config.js'
  },
  mutate: [
    'src/**/*.ts',
    'src/**/*.tsx',
    '!src/**/*.test.ts',
    '!src/**/*.test.tsx',
    '!src/**/*.d.ts',
    '!src/setupTests.ts',
    '!src/main.tsx'
  ],
  thresholds: {
    high: 90,
    low: 80,
    break: 70
  },
  checkers: ['typescript'],
  tsconfigFile: 'tsconfig.json',
  tempDirName: 'stryker-tmp',
  htmlReporter: {
    fileName: 'mutation-report.html'
  }
};