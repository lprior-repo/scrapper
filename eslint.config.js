import typescriptEslint from '@typescript-eslint/eslint-plugin';
import typescriptParser from '@typescript-eslint/parser';
import functional from 'eslint-plugin-functional';

export default [
  {
    files: ['packages/webapp/src/**/*.{ts,tsx}', 'packages/shared/src/**/*.{ts,tsx}'],
    languageOptions: {
      parser: typescriptParser,
      parserOptions: {
        ecmaVersion: 2022,
        sourceType: 'module',
        ecmaFeatures: {
          jsx: true,
        },
        project: ['./packages/webapp/tsconfig.json', './packages/shared/tsconfig.json'],
      },
      globals: {
        console: 'readonly',
        process: 'readonly',
        Buffer: 'readonly',
        global: 'readonly',
        window: 'readonly',
        document: 'readonly',
        HTMLElement: 'readonly',
        Element: 'readonly',
        Node: 'readonly',
        Event: 'readonly',
        MouseEvent: 'readonly',
        KeyboardEvent: 'readonly',
        FocusEvent: 'readonly',
        InputEvent: 'readonly',
        CustomEvent: 'readonly',
        IntersectionObserver: 'readonly',
        ResizeObserver: 'readonly',
        requestAnimationFrame: 'readonly',
        cancelAnimationFrame: 'readonly',
        fetch: 'readonly',
        URL: 'readonly',
        URLSearchParams: 'readonly',
      },
    },
    plugins: {
      '@typescript-eslint': typescriptEslint,
      'functional': functional,
    },
    rules: {
      // STRICT: No any types allowed
      '@typescript-eslint/no-explicit-any': 'error',
      
      // STRICT: Code quality
      '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
      '@typescript-eslint/prefer-as-const': 'error',
      '@typescript-eslint/no-non-null-assertion': 'error',
      
      // STRICT: General rules
      'no-console': ['warn', { allow: ['warn', 'error'] }],
      'no-debugger': 'error',
      'no-var': 'error',
      'prefer-const': 'error',
      'prefer-template': 'error',
      
      // STRICT: Import/Export rules
      'no-duplicate-imports': 'error',
      
      // STRICT: Code complexity
      'complexity': ['error', { max: 5 }],
      'max-depth': ['error', { max: 4 }],
      'max-lines': ['error', { max: 500, skipBlankLines: true, skipComments: true }],
      'max-lines-per-function': ['error', { max: 50, skipBlankLines: true, skipComments: true }],
      'max-params': ['error', { max: 5 }],
      
      // STRICT: Functional programming rules - practical subset
      'functional/no-let': 'error',
      'functional/prefer-readonly-type': 'error',
      'functional/immutable-data': ['error', {
        ignoreIdentifierPattern: ['^.*[Rr]ef$'],
        ignoreAccessorPattern: ['\\.current$'],
        ignoreImmediateMutation: true,
      }],
      'functional/no-conditional-statements': 'warn',
      'functional/no-throw-statements': 'warn',
      'functional/no-classes': 'warn',
    },
  },
  // Less strict rules for test files
  {
    files: ['**/*.spec.ts', '**/*.spec.tsx', '**/*.test.ts', '**/*.test.tsx'],
    rules: {
      'max-lines-per-function': 'off',
      'no-console': 'off',
      'complexity': 'off',
    },
  },
  // Less strict rules for config files
  {
    files: ['*.config.ts', '*.config.js', 'eslint.config.js'],
    rules: {
      'max-lines': 'off',
      '@typescript-eslint/no-explicit-any': 'off',
    },
  },
];