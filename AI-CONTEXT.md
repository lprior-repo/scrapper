# AI-CONTEXT.md

This file helps AI assistants understand your project structure, conventions, and patterns.
It's automatically loaded by AI tools to provide consistent context across sessions.

ACTION_REQUIRED: UPDATE_ALL_PLACEHOLDERS
IMMEDIATE: true

## Project Overview

**Project**: @overseer/monorepo
**Framework**: React
**Language**: TypeScript
**Last Updated**: 2025-07-19

### Description

[REPLACE: Brief description of what this project does]

### Key Features

- [REPLACE: Feature 1]
- [REPLACE: Feature 2]
- [REPLACE: Feature 3]

## Tech Stack

### Core Technologies

- **Framework**: React
- **Language**: TypeScript
- **Package Manager**: [npm/yarn/pnpm]
- **Build Tool**: [vite/webpack/etc]

### Key Dependencies

- [REPLACE: List main dependencies and their purpose]

## Project Structure

```
src/
├── components/     # [UI components]
├── pages/         # [Page components]
├── utils/         # [Utility functions]
├── styles/        # [Global styles]
└── types/         # [TypeScript types]
```

## Development Guidelines

### Code Style

- **Naming**: [DETECT: camelCase/PascalCase conventions from actual code]
- **Components**: [DETECT: Functional/Class components from src/]
- **State Management**: [DETECT: useState/Redux/Context from imports]
- **Styling**: [DETECT: CSS Modules/Styled Components/Tailwind from files]

### Best Practices

1. [EXTRACT: Practice 1 from codebase]
2. [EXTRACT: Practice 2 from codebase]
3. [EXTRACT: Practice 3 from codebase]

### Common Patterns

- **Component Structure**: [ANALYZE: How components are organized]
- **Data Flow**: [ANALYZE: How data flows through the app]
- **Error Handling**: [ANALYZE: How errors are handled]

## API Patterns

### Endpoints

- **Base URL**: [API base URL]
- **Authentication**: [How auth works]
- **Error Format**: [How API errors are structured]

### Common Operations

```typescript
// Example API call pattern
[Show typical API call]
```

## Testing Strategy

### Test Types

- **Unit Tests**: [Jest/Vitest]
- **Integration Tests**: [Testing Library]
- **E2E Tests**: [Cypress/Playwright]

### Coverage Goals

- Minimum coverage: [X%]
- Critical paths must have tests

## Performance Considerations

### Optimization Strategies

- [Strategy 1]
- [Strategy 2]

### Bundle Size Targets

- Initial load: <[X]kb
- Lazy loaded chunks: <[X]kb

## Security Guidelines

### Authentication

- [How auth is handled]

### Data Protection

- [How sensitive data is protected]

### API Security

- [CORS, rate limiting, etc.]

## Deployment

### Environments

- **Development**: [URL]
- **Staging**: [URL]
- **Production**: [URL]

### CI/CD Pipeline

- [Brief description of deployment process]

## Quick Commands

```bash
# Development
npm run dev

# Testing
npm test

# Build
npm run build

# Lint
npm run lint
```

## AI Assistant Notes

### When generating code:

1. Follow the patterns in AI-PATTERNS.md
2. Use existing utilities before creating new ones
3. Maintain consistent naming conventions
4. Add appropriate error handling
5. Include TypeScript types

### Common tasks:

- To add a new component: Follow patterns in src/components
- To add an API endpoint: See API patterns above
- To add a test: Match existing test structure

---

## ACTION REQUIRED: COMPLETE THIS FILE

EXECUTE_IMMEDIATELY:

1. READ all project files
2. REPLACE all [REPLACE:], [DETECT:], [EXTRACT:], and [ANALYZE:] markers
3. FILL all bracketed placeholders with actual project data
4. VALIDATE no placeholders remain

VALIDATION:
{
"action": "CHECK_PLACEHOLDERS",
"fail_if_found": ["[REPLACE:", "[DETECT:", "[EXTRACT:", "[ANALYZE:", "[", "]"],
"required": "zero_placeholders"
}

_This file is part of the AI Context system. Update it when project structure changes significantly._
