# Test Suite Documentation

This directory contains comprehensive tests for the GitHub Codeowners Visualization application. The test suite is designed to ensure the app renders correctly, components display expected data, and all user flows work as intended.

## Test Structure

### Unit Tests (`packages/webapp/src/__tests__/`)

- **Framework**: Bun native test runner with happy-dom
- **Coverage**: React components, services, and utility functions
- **Focus**: Component behavior, props handling, state management
- **Files**: `*.bun.test.{ts,tsx}` format

#### Test Files:

- `App.bun.test.tsx` - Main App component functionality
- `components/GraphCanvas.bun.test.tsx` - Graph visualization component
- `services.bun.test.ts` - API client and Zod schema validation

### End-to-End Tests (`tests/playwright/`)

- **Framework**: Playwright
- **Coverage**: Full user workflows and browser interactions
- **Focus**: Real user scenarios and integration testing

#### Test Files:

- `github-automation.spec.ts` - Original automation tests
- `app-rendering.spec.ts` - App rendering and component display
- `user-flows.spec.ts` - Complete user workflows
- `visual-regression.spec.ts` - Visual appearance consistency
- `api-integration.spec.ts` - Frontend-backend integration

## Test Categories

### 1. Component Rendering Tests ✅

- App component renders without errors
- GraphCanvas displays correctly with different data
- Loading states are shown properly
- Error states are handled gracefully

### 2. Data Validation Tests ✅

- API responses validated with Zod schemas
- GraphNode and GraphEdge structure validation
- Error handling for malformed responses
- Type safety verification

### 3. User Flow Tests ✅

- **Scanning workflow**: Enter org → Scan → View graph
- **View switching**: Toggle between teams and topics
- **Organization switching**: Change between different orgs
- **Error recovery**: Handle API failures gracefully
- **Rapid interactions**: Handle fast user input changes

### 4. Visual Regression Tests ✅

- Initial app state appearance
- Loading state UI
- Error state display
- Graph visualization with different node types
- Button and input states

### 5. Integration Tests ✅

- Real API endpoint communication
- Schema validation with actual responses
- Frontend-backend data flow
- Network error handling

## Running Tests

### All Tests

```bash
# Run complete test suite
bun run test

# Watch mode for development
bun run test:unit:watch
```

### Unit Tests Only

```bash
# Run all unit tests
bun run test:unit

# Run with coverage
bun run test:unit:coverage

# Watch mode
bun run test:unit:watch
```

### E2E Tests Only

```bash
# Run all E2E tests
bun run test:e2e

# Run with UI mode (browser visible)
bun run test:e2e:ui

# Debug mode (step through)
bun run test:e2e:debug
```

### Specific Test Categories

```bash
# Visual regression tests only
bun run test:visual

# Integration tests only
bun run test:integration
```

## Test Environment Setup

### Prerequisites

1. **Bun Runtime**: Ensure Bun is installed (`curl -fsSL https://bun.sh/install | bash`)
2. **Backend Running**: The Go backend must be running on `localhost:8081`
3. **Frontend Running**: The React app must be running on `localhost:3000`
4. **Dependencies**: All Bun dependencies installed (`bun install`)

### Environment Variables

Tests automatically detect if the backend is available and adjust expectations accordingly.

### Mock Data

- Unit tests use Bun's native `mock()` function for API responses
- E2E tests can use both mocked and real API calls via Playwright
- Visual tests use controlled mock data for consistency
- happy-dom provides lightweight DOM simulation

## Test Data Coverage

### Node Types Tested

- ✅ Organization nodes
- ✅ Repository nodes
- ✅ User nodes
- ✅ Team nodes
- ✅ Topic nodes

### Edge Types Tested

- ✅ Ownership relationships
- ✅ Maintenance relationships
- ✅ Team membership
- ✅ Topic associations

### API Endpoints Tested

- ✅ `GET /api/graph/{org}` - Graph data retrieval
- ✅ `GET /api/graph/{org}?useTopics=true` - Topics view
- ✅ `POST /api/scan/{org}` - Organization scanning
- ✅ Error responses (404, 500, etc.)

## Key Test Scenarios

### 1. Happy Path 🎯

```
User enters "github" → Checks topics → Scans org → Loads graph → Views visualization
```

### 2. Error Handling 🚨

```
Network fails → Shows error → User retries → Success
Invalid org → Empty graph → User switches org → Success
```

### 3. State Management 📊

```
Teams view → Switch to topics → Data updates → Switch back → Consistent state
```

### 4. Performance 🚀

```
Rapid org changes → Cancels old requests → Shows latest data
Large datasets → Renders without freezing → Interactive controls work
```

## Browser Support

E2E tests run on:

- ✅ Chrome (Desktop)
- ✅ Firefox (Desktop)
- ✅ Safari (Desktop)
- ✅ Chrome Mobile (Pixel 5)
- ✅ Safari Mobile (iPhone 12)

## Continuous Integration

### GitHub Actions Integration

Tests are designed to run in CI environments:

- Headless browser mode
- Reduced timeouts for faster execution
- Screenshot comparison with threshold tolerance
- Retry logic for flaky network tests

### Test Reports

- Unit test coverage reports (Jest)
- E2E test HTML reports (Playwright)
- Visual regression comparison images
- Performance metrics

## Writing New Tests

### Unit Test Guidelines

1. Test component behavior, not implementation
2. Use Bun's native `mock()` function for external dependencies
3. Use descriptive test names
4. Group related tests with `describe` blocks
5. Use `*.bun.test.{ts,tsx}` naming convention

### E2E Test Guidelines

1. Test real user workflows
2. Use data-testid attributes for reliable selectors
3. Handle async operations properly
4. Clean up state between tests

### Visual Test Guidelines

1. Use consistent viewport sizes
2. Disable animations for stable screenshots
3. Use mock data for predictable visuals
4. Set appropriate comparison thresholds

## Troubleshooting

### Common Issues

**Tests failing locally but passing in CI:**

- Check backend is running (`localhost:8081`)
- Verify frontend is accessible (`localhost:3000`)
- Ensure test dependencies are installed

**Visual tests showing differences:**

- Check screen resolution/scaling
- Verify browser version consistency
- Look for timing issues with animations

**API integration tests failing:**

- Confirm backend is running and healthy
- Check network connectivity
- Verify API endpoints are responsive

### Debug Commands

```bash
# Run specific test file
bun run test:e2e app-rendering.spec.ts

# Run with verbose output
bun run test:unit --verbose

# Generate coverage report
bun run test:unit:coverage
```

## Future Enhancements

### Planned Additions

- [ ] Accessibility testing (a11y)
- [ ] Performance benchmarking
- [ ] Cross-browser compatibility matrix
- [ ] API contract testing
- [ ] Load testing for large datasets

### Test Infrastructure

- [ ] Parallel test execution optimization
- [ ] Test result caching
- [ ] Flaky test detection
- [ ] Performance regression detection

---

This test suite provides comprehensive coverage of the application's functionality, ensuring reliability and maintainability as the codebase evolves.
