import { test, expect } from '@playwright/test';

test.describe('Visual Regression Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    // Set a consistent viewport for visual tests
    await page.setViewportSize({ width: 1280, height: 720 });
  });

  test('should match initial app state screenshot', async ({ page }) => {
    // Wait for app to fully load
    await page.waitForLoadState('networkidle');
    
    // Take screenshot of initial state
    await expect(page).toHaveScreenshot('initial-app-state.png', {
      fullPage: true,
      animations: 'disabled'
    });
  });

  test('should match header component appearance', async ({ page }) => {
    const header = page.locator('header');
    await expect(header).toBeVisible();
    
    await expect(header).toHaveScreenshot('header-component.png', {
      animations: 'disabled'
    });
  });

  test('should match input states visually', async ({ page }) => {
    const inputContainer = page.locator('div:has(input[placeholder="Enter organization name"])');
    
    // Empty state
    await expect(inputContainer).toHaveScreenshot('input-empty-state.png');
    
    // Filled state
    await page.fill('input[placeholder="Enter organization name"]', 'github');
    await expect(inputContainer).toHaveScreenshot('input-filled-state.png');
    
    // With checkbox checked
    await page.check('input[type="checkbox"]');
    await expect(inputContainer).toHaveScreenshot('input-with-checkbox.png');
  });

  test('should match button states visually', async ({ page }) => {
    const button = page.locator('button:has-text("Load Graph")');
    
    // Disabled state
    await expect(button).toHaveScreenshot('button-disabled.png');
    
    // Enabled state
    await page.fill('input[placeholder="Enter organization name"]', 'test');
    await expect(button).toHaveScreenshot('button-enabled.png');
    
    // Hover state
    await button.hover();
    await expect(button).toHaveScreenshot('button-hover.png');
  });

  test('should match loading state appearance', async ({ page }) => {
    // Delay the response to capture loading state
    await page.route('**/api/graph/**', async (route) => {
      await new Promise(resolve => setTimeout(resolve, 2000));
      await route.fulfill({
        status: 200,
        body: JSON.stringify({ data: { nodes: [], edges: [] } })
      });
    });
    
    await page.fill('input[placeholder="Enter organization name"]', 'github');
    await page.click('button:has-text("Load Graph")');
    
    // Capture loading state
    await page.waitForSelector('text=Loading graph data...');
    await expect(page).toHaveScreenshot('loading-state.png', {
      fullPage: true,
      animations: 'disabled'
    });
  });

  test('should match error state appearance', async ({ page }) => {
    // Mock error response
    await page.route('**/api/graph/**', (route) => {
      route.fulfill({
        status: 404,
        body: 'Not Found'
      });
    });
    
    await page.fill('input[placeholder="Enter organization name"]', 'error-test');
    await page.click('button:has-text("Load Graph")');
    
    // Wait for error state
    await page.waitForSelector('text=Error loading graph');
    await expect(page).toHaveScreenshot('error-state.png', {
      fullPage: true,
      animations: 'disabled'
    });
  });

  test('should match graph canvas with different node types', async ({ page }) => {
    const mockGraphData = {
      data: {
        nodes: [
          {
            id: 'org-1',
            type: 'organization',
            label: 'Test Org',
            data: {},
            position: { x: 400, y: 200 }
          },
          {
            id: 'repo-1',
            type: 'repository',
            label: 'repo-one',
            data: {},
            position: { x: 200, y: 300 }
          },
          {
            id: 'repo-2',
            type: 'repository',
            label: 'repo-two',
            data: {},
            position: { x: 600, y: 300 }
          },
          {
            id: 'team-1',
            type: 'team',
            label: 'frontend-team',
            data: {},
            position: { x: 200, y: 400 }
          },
          {
            id: 'user-1',
            type: 'user',
            label: 'alice',
            data: {},
            position: { x: 100, y: 500 }
          },
          {
            id: 'user-2',
            type: 'user',
            label: 'bob',
            data: {},
            position: { x: 300, y: 500 }
          }
        ],
        edges: [
          {
            id: 'e1',
            source: 'org-1',
            target: 'repo-1',
            type: 'owns',
            label: 'owns'
          },
          {
            id: 'e2',
            source: 'org-1',
            target: 'repo-2',
            type: 'owns',
            label: 'owns'
          },
          {
            id: 'e3',
            source: 'team-1',
            target: 'repo-1',
            type: 'maintains',
            label: 'maintains'
          },
          {
            id: 'e4',
            source: 'user-1',
            target: 'team-1',
            type: 'member_of',
            label: 'member of'
          },
          {
            id: 'e5',
            source: 'user-2',
            target: 'team-1',
            type: 'member_of',
            label: 'member of'
          }
        ]
      }
    };
    
    await page.route('**/api/graph/**', (route) => {
      route.fulfill({
        status: 200,
        body: JSON.stringify(mockGraphData)
      });
    });
    
    await page.fill('input[placeholder="Enter organization name"]', 'visual-test');
    await page.click('button:has-text("Load Graph")');
    
    // Wait for graph to render and stabilize
    await page.waitForSelector('[data-testid="graph-canvas"]');
    await page.waitForTimeout(2000); // Allow physics to stabilize
    
    await expect(page).toHaveScreenshot('graph-with-teams.png', {
      fullPage: true,
      animations: 'disabled'
    });
  });

  test('should match graph with topics view', async ({ page }) => {
    const mockTopicsData = {
      data: {
        nodes: [
          {
            id: 'org-1',
            type: 'organization',
            label: 'Test Org',
            data: {},
            position: { x: 400, y: 200 }
          },
          {
            id: 'topic-1',
            type: 'topic',
            label: 'javascript',
            data: { name: 'javascript', count: 10 },
            position: { x: 200, y: 350 }
          },
          {
            id: 'topic-2',
            type: 'topic',
            label: 'typescript',
            data: { name: 'typescript', count: 8 },
            position: { x: 600, y: 350 }
          },
          {
            id: 'topic-3',
            type: 'topic',
            label: 'react',
            data: { name: 'react', count: 5 },
            position: { x: 400, y: 450 }
          }
        ],
        edges: [
          {
            id: 'e1',
            source: 'org-1',
            target: 'topic-1',
            type: 'uses',
            label: 'uses'
          },
          {
            id: 'e2',
            source: 'org-1',
            target: 'topic-2',
            type: 'uses',
            label: 'uses'
          },
          {
            id: 'e3',
            source: 'org-1',
            target: 'topic-3',
            type: 'uses',
            label: 'uses'
          }
        ]
      }
    };
    
    await page.route('**/api/graph/**', (route) => {
      route.fulfill({
        status: 200,
        body: JSON.stringify(mockTopicsData)
      });
    });
    
    await page.fill('input[placeholder="Enter organization name"]', 'visual-topics');
    await page.check('input[type="checkbox"]');
    await page.click('button:has-text("Load Graph")');
    
    // Wait for graph to render and stabilize
    await page.waitForSelector('[data-testid="graph-canvas"]');
    await page.waitForTimeout(2000); // Allow physics to stabilize
    
    await expect(page).toHaveScreenshot('graph-with-topics.png', {
      fullPage: true,
      animations: 'disabled'
    });
  });

  test('should match empty graph state', async ({ page }) => {
    await page.route('**/api/graph/**', (route) => {
      route.fulfill({
        status: 200,
        body: JSON.stringify({ data: { nodes: [], edges: [] } })
      });
    });
    
    await page.fill('input[placeholder="Enter organization name"]', 'empty-org');
    await page.click('button:has-text("Load Graph")');
    
    await page.waitForSelector('[data-testid="graph-canvas"]');
    await page.waitForTimeout(1000);
    
    await expect(page).toHaveScreenshot('empty-graph.png', {
      fullPage: true,
      animations: 'disabled'
    });
  });
});