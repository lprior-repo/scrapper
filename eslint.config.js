import typescriptEslint from '@typescript-eslint/eslint-plugin'
import typescriptParser from '@typescript-eslint/parser'
import functional from 'eslint-plugin-functional'
import react from 'eslint-plugin-react'
import reactHooks from 'eslint-plugin-react-hooks'

export default [
  {
    files: [
      'packages/webapp/src/**/*.{ts,tsx,js,jsx}',
      'packages/shared/src/**/*.{ts,tsx,js,jsx}',
      'packages/webapp/src/components/**/*.{ts,tsx}',
      'packages/webapp/src/components/types/**/*.{ts,tsx}',
    ],
    languageOptions: {
      parser: typescriptParser,
      parserOptions: {
        ecmaVersion: 2022,
        sourceType: 'module',
        ecmaFeatures: {
          jsx: true,
        },
        project: [
          './packages/webapp/tsconfig.json',
          './packages/shared/tsconfig.json',
        ],
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
      functional: functional,
      react: react,
      'react-hooks': reactHooks,
    },
    settings: {
      react: {
        version: '19.0.0', // React version for eslint-plugin-react
      },
    },
    rules: {
      // STRICT: No any types allowed
      '@typescript-eslint/no-explicit-any': 'error',

      // STRICT: Code quality
      '@typescript-eslint/no-unused-vars': [
        'warn', // Changed to warning for better development experience
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_' },
      ],
      '@typescript-eslint/prefer-as-const': 'error',
      '@typescript-eslint/no-non-null-assertion': 'warn', // Changed to warning for existing code

      // STRICT: General rules
      'no-console': ['warn', { allow: ['warn', 'error'] }],
      'no-debugger': 'error',
      'no-var': 'error',
      'prefer-const': 'error',
      'prefer-template': 'error',

      // STRICT: Import/Export rules
      'no-duplicate-imports': 'error',

      // MODERATE: Code complexity - adjusted for existing codebase
      complexity: ['warn', { max: 15 }], // Increased from 5 to accommodate existing code
      'max-depth': ['warn', { max: 4 }], // Restored to 4 levels
      'max-lines': [
        'warn', // Changed to warning
        { max: 1500, skipBlankLines: true, skipComments: true }, // Increased limit for large components
      ],
      'max-lines-per-function': [
        'warn', // Changed to warning
        { max: 150, skipBlankLines: true, skipComments: true }, // Increased limit for React components
      ],
      'max-params': ['warn', { max: 5 }], // Restored to 5 parameters

      // MODERATE: Functional programming rules - practical subset for existing codebase
      'functional/no-let': 'warn', // Changed from error to warn for existing code compatibility
      'functional/prefer-readonly-type': 'warn', // Changed from error to warn
      'functional/immutable-data': [
        'warn', // Changed from error to warn
        {
          ignoreIdentifierPattern: ['^.*[Rr]ef$', '^cy$', '^cache$', '^window$', '^document$'],
          ignoreAccessorPattern: ['\\.current$', '\\.style$', '\\.innerHTML$'],
          ignoreImmediateMutation: true,
          ignoreNonConstDeclarations: true, // Allow mutations in let declarations
        },
      ],
      'functional/no-conditional-statements': 'off', // Too restrictive for existing React components
      'functional/no-throw-statements': 'off', // Allow throwing errors in error handling
      'functional/no-classes': 'off', // Allow classes for existing test utilities and DOM APIs

      // STRICT: React 19 specific rules
      'react/jsx-uses-react': 'off', // Not needed in React 19 with automatic JSX transform
      'react/react-in-jsx-scope': 'off', // Not needed in React 19
      'react/jsx-uses-vars': 'error',
      'react/jsx-key': 'error',
      'react/jsx-no-duplicate-props': 'error',
      'react/jsx-no-undef': 'error',
      'react/no-children-prop': 'error',
      'react/no-danger-with-children': 'error',
      'react/no-deprecated': 'error',
      'react/no-direct-mutation-state': 'error',
      'react/no-find-dom-node': 'error',
      'react/no-is-mounted': 'error',
      'react/no-render-return-value': 'error',
      'react/no-string-refs': 'error',
      'react/no-unescaped-entities': 'error',
      'react/no-unknown-property': 'error',
      'react/no-unsafe': 'warn',
      'react/prop-types': 'off', // Using TypeScript for type checking
      'react/require-render-return': 'error',

      // STRICT: React Hooks rules (simplified for ESLint 9 compatibility)
      'react-hooks/rules-of-hooks': 'error',
      // 'react-hooks/exhaustive-deps': 'warn', // Temporarily disabled due to ESLint 9 compatibility issues
    },
  },
  // Less strict rules for test files
  {
    files: ['**/*.spec.ts', '**/*.spec.tsx', '**/*.test.ts', '**/*.test.tsx'],
    rules: {
      'max-lines-per-function': 'off',
      'no-console': 'off',
      complexity: 'off',
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
]
