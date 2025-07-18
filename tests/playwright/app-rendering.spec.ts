import { test, expect } from '@playwright/test';

test.describe('App Rendering and Component Display', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should render App component without errors', async ({ page }) => {
    // Check that the app loads without console errors
    const consoleErrors: string[] = [];
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text());
      }
    });

    await page.waitForLoadState('networkidle');
    
    // No console errors should occur
    expect(consoleErrors).toHaveLength(0);
    
    // Main app elements should be visible
    await expect(page.locator('h1')).toHaveText('GitHub Codeowners Visualization');
    await expect(page.locator('input[placeholder="Enter organization name"]')).toBeVisible();
    await expect(page.locator('button:has-text("Load Graph")')).toBeVisible();
  });

  test('should display correct initial loading states', async ({ page }) => {
    // Check initial state
    const initialMessage = page.locator('text=/Enter an organization name/i');
    await expect(initialMessage).toBeVisible();
    
    // Input should be empty
    const orgInput = page.locator('input[placeholder="Enter organization name"]');
    await expect(orgInput).toHaveValue('');
    
    // Button should be disabled
    const loadButton = page.locator('button:has-text("Load Graph")');
    await expect(loadButton).toBeDisabled();
    
    // Checkbox should be unchecked
    const checkbox = page.locator('input[type="checkbox"]');
    await expect(checkbox).not.toBeChecked();
  });

  test('should show loading state when fetching graph data', async ({ page }) => {
    // Enter organization
    await page.fill('input[placeholder="Enter organization name"]', 'github');
    
    // Set up response delay to observe loading state
    await page.route('**/api/graph/**', async (route) => {
      await new Promise(resolve => setTimeout(resolve, 1000));
      await route.fulfill({
        status: 200,
        body: JSON.stringify({
          data: {
            nodes: [],
            edges: []
          }
        })
      });
    });
    
    // Click load button
    await page.click('button:has-text("Load Graph")');
    
    // Check loading state appears
    await expect(page.locator('text=Loading graph data...')).toBeVisible();
    
    // Wait for loading to complete
    await expect(page.locator('[data-testid="graph-canvas"]')).toBeVisible({ timeout: 5000 });
  });

  test('should handle error states gracefully', async ({ page }) => {
    // Enter organization
    await page.fill('input[placeholder="Enter organization name"]', 'error-test-org');
    
    // Mock error response
    await page.route('**/api/graph/**', (route) => {
      route.fulfill({
        status: 500,
        body: 'Internal Server Error'
      });
    });
    
    // Click load button
    await page.click('button:has-text("Load Graph")');
    
    // Check error display
    await expect(page.locator('text=Error loading graph')).toBeVisible();
    await expect(page.locator('text=/HTTP error! status: 500/i')).toBeVisible();
  });

  test('should validate API responses with Zod schemas', async ({ page }) => {
    // Set up a malformed response to test schema validation
    await page.route('**/api/graph/**', (route) => {
      route.fulfill({
        status: 200,
        body: JSON.stringify({
          data: {
            nodes: [
              {
                id: 'node-1',
                // Missing required fields: type, label, data, position
              }
            ],
            edges: []
          }
        })
      });
    });
    
    await page.fill('input[placeholder="Enter organization name"]', 'schema-test');
    await page.click('button:has-text("Load Graph")');
    
    // The component should show an error due to schema validation failure
    await expect(page.locator('text=Error loading graph')).toBeVisible();
  });

  test('should properly display nodes and edges in GraphCanvas', async ({ page }) => {
    const mockGraphData = {
      data: {
        nodes: [
          {
            id: 'org-github',
            type: 'organization',
            label: 'github',
            data: {},
            position: { x: 0, y: 0 }
          },
          {
            id: 'repo-1',
            type: 'repository',
            label: 'awesome-project',
            data: { stars: 1000 },
            position: { x: 100, y: 100 }
          },
          {
            id: 'user-1',
            type: 'user',
            label: 'johndoe',
            data: {},
            position: { x: 200, y: 200 }
          }
        ],
        edges: [
          {
            id: 'edge-1',
            source: 'org-github',
            target: 'repo-1',
            type: 'owns',
            label: 'owns'
          },
          {
            id: 'edge-2',
            source: 'repo-1',
            target: 'user-1',
            type: 'maintained_by',
            label: 'maintained by'
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
    
    await page.fill('input[placeholder="Enter organization name"]', 'github');
    await page.click('button:has-text("Load Graph")');
    
    // Wait for graph canvas
    const graphCanvas = page.locator('[data-testid="graph-canvas"]');
    await expect(graphCanvas).toBeVisible();
    
    // Verify canvas has correct dimensions
    const canvasBox = await graphCanvas.boundingBox();
    expect(canvasBox?.width).toBeGreaterThan(0);
    expect(canvasBox?.height).toBeGreaterThan(0);
    
    // Check that cytoscape container is created
    await expect(graphCanvas).toBeVisible();
  });

  test('should update graph when switching between teams/topics view', async ({ page }) => {
    // Mock different responses for teams vs topics
    await page.route('**/api/graph/**', async (route) => {
      const url = route.request().url();
      const useTopics = url.includes('useTopics=true');
      
      const data = useTopics ? {
        nodes: [
          {
            id: 'topic-1',
            type: 'topic',
            label: 'javascript',
            data: { name: 'javascript', count: 5 },
            position: { x: 0, y: 0 }
          }
        ],
        edges: []
      } : {
        nodes: [
          {
            id: 'team-1',
            type: 'team',
            label: 'engineering',
            data: {},
            position: { x: 0, y: 0 }
          }
        ],
        edges: []
      };
      
      await route.fulfill({
        status: 200,
        body: JSON.stringify({ data })
      });
    });
    
    // Load with teams view (default)
    await page.fill('input[placeholder="Enter organization name"]', 'test-org');
    await page.click('button:has-text("Load Graph")');
    
    // Wait for first graph to load
    await expect(page.locator('[data-testid="graph-canvas"]')).toBeVisible();
    
    // Now switch to topics view
    await page.check('input[type="checkbox"]');
    await page.click('button:has-text("Load Graph")');
    
    // Verify the API was called with useTopics parameter
    const topicsResponse = await page.waitForResponse(response => 
      response.url().includes('useTopics=true')
    );
    expect(topicsResponse.status()).toBe(200);
  });

  test('should maintain UI state during graph transitions', async ({ page }) => {
    // Mock successful response
    await page.route('**/api/graph/**', (route) => {
      route.fulfill({
        status: 200,
        body: JSON.stringify({
          data: { nodes: [], edges: [] }
        })
      });
    });
    
    // Set initial state
    await page.fill('input[placeholder="Enter organization name"]', 'org1');
    await page.check('input[type="checkbox"]');
    
    // Load first graph
    await page.click('button:has-text("Load Graph")');
    await expect(page.locator('[data-testid="graph-canvas"]')).toBeVisible();
    
    // Change organization
    await page.fill('input[placeholder="Enter organization name"]', 'org2');
    
    // Checkbox state should be preserved
    await expect(page.locator('input[type="checkbox"]')).toBeChecked();
    
    // Load new graph
    await page.click('button:has-text("Load Graph")');
    
    // Canvas should update without losing checkbox state
    await expect(page.locator('[data-testid="graph-canvas"]')).toBeVisible();
    await expect(page.locator('input[type="checkbox"]')).toBeChecked();
  });
});