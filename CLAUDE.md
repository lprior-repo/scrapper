# GitHub Codeowners Visualization - Bun Runtime Project

## Important: Bun-Only Environment

This project **exclusively uses Bun** as the JavaScript runtime and package manager. There is **NO Node.js** anywhere in this stack.

### Runtime Environment
- ✅ **Bun**: Primary JavaScript runtime and package manager
- ❌ **Node.js**: Not used, not installed, not supported
- ❌ **npm**: Not used, replaced with Bun
- ❌ **yarn**: Not used, replaced with Bun
- ❌ **pnpm**: Not used, replaced with Bun

### Key Bun Features Used
1. **Native Test Runner**: Using `bun test` instead of Jest
2. **DOM Environment**: Using `happy-dom` for browser simulation
3. **TypeScript Support**: Built-in TypeScript compilation
4. **Hot Reloading**: Native development server with `--hot`
5. **Bundle Building**: Using `bun build` for production builds

### Project Architecture

#### Frontend (packages/webapp/)
- **Runtime**: Bun
- **Framework**: React 18
- **Graph Library**: Cytoscape.js (migrated from vis-network)
- **Testing**: Bun native test runner + Playwright
- **Build Tool**: Bun build
- **Dev Server**: Bun with hot reload

#### Backend
- **Language**: Go (separate from JavaScript stack)
- **Database**: Neo4j
- **API**: REST endpoints

### Development Commands

All commands use Bun exclusively:

```bash
# Development
bun run dev                 # Start development server
bun run dev:webapp         # Start webapp only
bun run build              # Build for production
bun run build:webapp       # Build webapp only

# Testing
bun run test               # Run all tests (unit + E2E)
bun run test:unit          # Run Bun native unit tests
bun run test:unit:watch    # Watch mode for unit tests
bun run test:unit:coverage # Unit tests with coverage
bun run test:e2e           # Run Playwright E2E tests
bun run test:e2e:ui        # E2E tests with UI
bun run test:e2e:debug     # Debug E2E tests
bun run test:visual        # Visual regression tests
bun run test:integration   # Backend integration tests

# Code Quality
bun run lint               # Lint all packages
bun run lint:fix           # Fix linting issues
bun run format             # Format code with Prettier
bun run type-check         # TypeScript type checking
```

### Test Stack (Bun Native)

#### Unit Tests
- **Framework**: Bun native test runner
- **DOM**: happy-dom (lighter than jsdom)
- **Mocking**: Bun's built-in `mock()` function
- **Location**: `packages/webapp/src/**/*.bun.test.{ts,tsx}`

#### Integration Tests
- **Framework**: Playwright
- **Browsers**: Chrome, Firefox, Safari, Mobile
- **Location**: `tests/playwright/*.spec.ts`

#### Visual Regression
- **Framework**: Playwright screenshots
- **Comparison**: Automated pixel-perfect comparisons
- **Threshold**: 0.3 tolerance for minor rendering differences

### Configuration Files

#### Bun Configuration (`bunfig.toml`)
```toml
[test]
preload = ["./src/test/setup.ts"]
timeout = 5000

[test.coverage]
include = ["src/**/*.{ts,tsx}"]
exclude = ["src/index.tsx", "src/types.d.ts", "src/**/*.test.{ts,tsx}"]
```

#### Playwright Configuration
- Uses `bun run dev` for dev server
- Supports multiple browsers and devices
- Visual regression testing enabled

### Migration from Node.js Stack

This project was fully migrated from Node.js dependencies:

#### Removed Node.js Dependencies
- ❌ `jest` → ✅ `bun test`
- ❌ `@testing-library/jest-dom` → ✅ `happy-dom`
- ❌ `ts-jest` → ✅ Native Bun TypeScript
- ❌ `jest-environment-jsdom` → ✅ `happy-dom`
- ❌ `identity-obj-proxy` → ✅ Not needed with Bun

#### Replaced with Bun-Native
- ✅ `bun-types`: Bun TypeScript definitions
- ✅ `happy-dom`: Lightweight DOM for testing
- ✅ Native test runner with mocking
- ✅ Built-in TypeScript compilation
- ✅ Native module resolution

### Graph Visualization Migration

Successfully migrated from vis-network to cytoscape.js:

#### Before (vis-network)
```typescript
import { Network } from 'vis-network/standalone'
import 'vis-network/styles/vis-network.css'

const network = new Network(container, data, options)
```

#### After (cytoscape.js)
```typescript
import cytoscape from 'cytoscape'

const cy = cytoscape({
  container,
  elements,
  style,
  layout: { name: 'cose' }
})
```

### Performance Benefits of Bun

1. **Faster Test Execution**: ~3x faster than Jest
2. **Faster Package Installation**: ~5x faster than npm
3. **Instant TypeScript**: No compilation step needed
4. **Smaller Bundle Size**: More efficient bundling
5. **Hot Reload**: Near-instant updates in development

### Compatibility Notes

#### What Works with Bun
- ✅ React and React DOM
- ✅ TypeScript (native support)
- ✅ ESLint and Prettier
- ✅ Playwright (Node.js-based, runs separately)
- ✅ Effect-TS ecosystem
- ✅ Cytoscape.js
- ✅ Modern ES modules

#### What Doesn't Work
- ❌ Node.js specific APIs (not needed)
- ❌ Jest and Jest ecosystem
- ❌ CommonJS modules (use ESM instead)
- ❌ Node.js test runners

### Development Environment Setup

1. **Install Bun**:
   ```bash
   curl -fsSL https://bun.sh/install | bash
   ```

2. **Install Dependencies**:
   ```bash
   bun install
   ```

3. **Start Development**:
   ```bash
   bun run dev
   ```

4. **Run Tests**:
   ```bash
   bun run test
   ```

### CI/CD Configuration

Ensure CI environments use Bun:

```yaml
# GitHub Actions example
- name: Setup Bun
  uses: oven-sh/setup-bun@v1
  with:
    bun-version: latest

- name: Install dependencies
  run: bun install

- name: Run tests
  run: bun run test
```

### Troubleshooting

#### Common Issues
1. **"npm not found"**: Good! Use `bun` instead
2. **"jest command not found"**: Use `bun test` instead
3. **Module resolution errors**: Ensure using ESM imports
4. **TypeScript errors**: Bun has strict TypeScript checking

#### Solutions
- Always use `bun` commands
- Update imports to use ESM syntax
- Check `bunfig.toml` for configuration
- Use `bun run type-check` for TypeScript validation

### Project Goals Achieved

✅ **Full Bun Migration**: No Node.js dependencies remain  
✅ **Modern Graph Library**: Cytoscape.js provides better performance  
✅ **Comprehensive Testing**: Unit, integration, and visual tests  
✅ **Type Safety**: Full TypeScript coverage with Effect-TS  
✅ **Developer Experience**: Fast dev server, instant tests  
✅ **Production Ready**: Optimized builds with Bun bundler  

---

**Remember**: This is a Bun-only project. Any references to Node.js, npm, Jest, or other Node.js ecosystem tools should be replaced with their Bun equivalents.