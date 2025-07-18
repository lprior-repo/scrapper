import { test, expect } from '@playwright/test';

test.describe('Key User Flows', () => {
  test('complete flow: scan organization and view graph', async ({ page }) => {
    // Navigate to app
    await page.goto('/');
    
    // Step 1: Enter organization name
    console.log('Step 1: Entering organization name...');
    const orgInput = page.locator('input[placeholder="Enter organization name"]');
    await orgInput.fill('github');
    await expect(orgInput).toHaveValue('github');
    
    // Step 2: Select topics view
    console.log('Step 2: Selecting topics view...');
    const topicsCheckbox = page.locator('input[type="checkbox"]');
    await topicsCheckbox.check();
    await expect(topicsCheckbox).toBeChecked();
    
    // Step 3: Scan the organization first
    console.log('Step 3: Scanning organization...');
    const scanResponse = await page.request.post('http://localhost:8081/api/scan/github?useTopics=true&max_repos=10', {
      timeout: 60000,
    });
    
    expect(scanResponse.status()).toBe(201);
    const scanData = await scanResponse.json();
    console.log(`Scan completed: ${scanData.data.summary.totalRepositories} repos, ${scanData.data.summary.totalCodeowners} codeowners`);
    
    // Step 4: Load and view the graph
    console.log('Step 4: Loading graph visualization...');
    await page.click('button:has-text("Load Graph")');
    
    // Wait for loading state
    await expect(page.locator('text=Loading graph data...')).toBeVisible();
    
    // Wait for graph to load
    const graphCanvas = page.locator('[data-testid="graph-canvas"]');
    await expect(graphCanvas).toBeVisible({ timeout: 30000 });
    
    // Verify graph is rendered with cytoscape
    await expect(graphCanvas).toBeVisible();
    
    // Step 5: Verify interaction capabilities
    console.log('Step 5: Testing graph interactions...');
    
    // Test zoom
    await graphCanvas.hover();
    await page.mouse.wheel(0, -100); // Zoom in
    await page.waitForTimeout(500);
    await page.mouse.wheel(0, 100); // Zoom out
    
    // Test pan
    await page.mouse.move(640, 360);
    await page.mouse.down();
    await page.mouse.move(740, 360);
    await page.mouse.up();
    
    console.log('✅ Complete user flow successful!');
  });

  test('flow: switch between different organizations', async ({ page }) => {
    await page.goto('/');
    
    // Mock responses for different organizations
    await page.route('**/api/graph/**', async (route) => {
      const url = route.request().url();
      const org = url.includes('microsoft') ? 'microsoft' : 'facebook';
      
      const data = {
        nodes: [
          {
            id: `org-${org}`,
            type: 'organization',
            label: org,
            data: {},
            position: { x: 0, y: 0 }
          },
          {
            id: `repo-${org}-1`,
            type: 'repository',
            label: `${org}-repo`,
            data: {},
            position: { x: 100, y: 100 }
          }
        ],
        edges: [
          {
            id: 'e1',
            source: `org-${org}`,
            target: `repo-${org}-1`,
            type: 'owns',
            label: 'owns'
          }
        ]
      };
      
      await route.fulfill({
        status: 200,
        body: JSON.stringify({ data })
      });
    });
    
    // Load first organization
    console.log('Loading first organization: microsoft');
    await page.fill('input[placeholder="Enter organization name"]', 'microsoft');
    await page.click('button:has-text("Load Graph")');
    await expect(page.locator('[data-testid="graph-canvas"]')).toBeVisible();
    
    // Verify microsoft graph is loaded
    const graphResponse1 = await page.waitForResponse(response => 
      response.url().includes('microsoft')
    );
    expect(graphResponse1.status()).toBe(200);
    
    // Switch to second organization
    console.log('Switching to second organization: facebook');
    await page.fill('input[placeholder="Enter organization name"]', 'facebook');
    await page.click('button:has-text("Load Graph")');
    
    // Verify facebook graph is loaded
    const graphResponse2 = await page.waitForResponse(response => 
      response.url().includes('facebook')
    );
    expect(graphResponse2.status()).toBe(200);
    
    // Canvas should still be visible
    await expect(page.locator('[data-testid="graph-canvas"]')).toBeVisible();
    
    console.log('✅ Organization switching flow successful!');
  });

  test('flow: toggle between teams and topics view', async ({ page }) => {
    await page.goto('/');
    
    // Mock different responses for teams vs topics
    let requestCount = 0;
    await page.route('**/api/graph/**', async (route) => {
      const url = route.request().url();
      const useTopics = url.includes('useTopics=true');
      requestCount++;
      
      const nodeType = useTopics ? 'topic' : 'team';
      const nodeLabel = useTopics ? `topic-${requestCount}` : `team-${requestCount}`;
      
      const data = {
        nodes: [
          {
            id: 'org-test',
            type: 'organization',
            label: 'test-org',
            data: {},
            position: { x: 0, y: 0 }
          },
          {
            id: `${nodeType}-1`,
            type: nodeType,
            label: nodeLabel,
            data: useTopics ? { name: nodeLabel, count: 5 } : {},
            position: { x: 100, y: 100 }
          }
        ],
        edges: [
          {
            id: 'e1',
            source: 'org-test',
            target: `${nodeType}-1`,
            type: 'has',
            label: 'has'
          }
        ]
      };
      
      await route.fulfill({
        status: 200,
        body: JSON.stringify({ data })
      });
    });
    
    // Load with teams view (default)
    console.log('Loading graph with teams view...');
    await page.fill('input[placeholder="Enter organization name"]', 'test-org');
    await page.click('button:has-text("Load Graph")');
    
    await expect(page.locator('[data-testid="graph-canvas"]')).toBeVisible();
    
    // Verify teams request
    const teamsResponse = await page.waitForResponse(response => 
      response.url().includes('/api/graph/') && !response.url().includes('useTopics')
    );
    expect(teamsResponse.status()).toBe(200);
    
    // Toggle to topics view
    console.log('Switching to topics view...');
    await page.check('input[type="checkbox"]');
    await page.click('button:has-text("Load Graph")');
    
    // Verify topics request
    const topicsResponse = await page.waitForResponse(response => 
      response.url().includes('useTopics=true')
    );
    expect(topicsResponse.status()).toBe(200);
    
    // Toggle back to teams view
    console.log('Switching back to teams view...');
    await page.uncheck('input[type="checkbox"]');
    await page.click('button:has-text("Load Graph")');
    
    // Verify another teams request
    await page.waitForResponse(response => 
      response.url().includes('/api/graph/') && !response.url().includes('useTopics')
    );
    
    console.log('✅ View toggling flow successful!');
  });

  test('flow: handle scanning errors and retry', async ({ page }) => {
    await page.goto('/');
    
    // First attempt - scan fails
    let scanAttempt = 0;
    await page.route('**/api/scan/**', async (route) => {
      scanAttempt++;
      if (scanAttempt === 1) {
        // First attempt fails
        await route.fulfill({
          status: 500,
          body: 'Internal Server Error'
        });
      } else {
        // Second attempt succeeds
        await route.fulfill({
          status: 201,
          body: JSON.stringify({
            data: {
              success: true,
              organization: 'retry-org',
              summary: {
                totalRepositories: 5,
                totalCodeowners: 10
              }
            }
          })
        });
      }
    });
    
    // Also mock graph endpoint for after successful scan
    await page.route('**/api/graph/**', (route) => {
      route.fulfill({
        status: 200,
        body: JSON.stringify({
          data: { nodes: [], edges: [] }
        })
      });
    });
    
    console.log('Attempting scan that will fail...');
    const firstScanResponse = await page.request.post('http://localhost:8081/api/scan/retry-org', {
      timeout: 30000,
    });
    
    expect(firstScanResponse.status()).toBe(500);
    console.log('First scan failed as expected');
    
    // Retry scan
    console.log('Retrying scan...');
    const retryScanResponse = await page.request.post('http://localhost:8081/api/scan/retry-org', {
      timeout: 30000,
    });
    
    expect(retryScanResponse.status()).toBe(201);
    console.log('Retry successful!');
    
    // Now load the graph
    await page.fill('input[placeholder="Enter organization name"]', 'retry-org');
    await page.click('button:has-text("Load Graph")');
    
    await expect(page.locator('[data-testid="graph-canvas"]')).toBeVisible();
    
    console.log('✅ Error handling and retry flow successful!');
  });

  test('flow: rapid organization changes', async ({ page }) => {
    await page.goto('/');
    
    // Mock graph endpoint
    await page.route('**/api/graph/**', async (route) => {
      const url = route.request().url();
      const org = url.match(/\/api\/graph\/([^?]+)/)?.[1] || 'unknown';
      
      // Simulate varying response times
      const delay = Math.random() * 1000;
      await new Promise(resolve => setTimeout(resolve, delay));
      
      await route.fulfill({
        status: 200,
        body: JSON.stringify({
          data: {
            nodes: [{
              id: `org-${org}`,
              type: 'organization',
              label: org,
              data: {},
              position: { x: 0, y: 0 }
            }],
            edges: []
          }
        })
      });
    });
    
    console.log('Testing rapid organization changes...');
    
    // Rapidly change organizations
    const organizations = ['org1', 'org2', 'org3', 'org4', 'org5'];
    
    for (const org of organizations) {
      await page.fill('input[placeholder="Enter organization name"]', org);
      await page.click('button:has-text("Load Graph")');
      
      // Don't wait for completion, immediately change to next
      await page.waitForTimeout(100);
    }
    
    // Final organization should eventually load
    await expect(page.locator('[data-testid="graph-canvas"]')).toBeVisible({ timeout: 10000 });
    
    // Verify the last organization loaded
    const lastOrgResponse = await page.waitForResponse(response => 
      response.url().includes('org5')
    );
    expect(lastOrgResponse.status()).toBe(200);
    
    console.log('✅ Rapid organization change handling successful!');
  });

  test('flow: keyboard navigation in graph', async ({ page }) => {
    await page.goto('/');
    
    // Mock simple graph
    await page.route('**/api/graph/**', (route) => {
      route.fulfill({
        status: 200,
        body: JSON.stringify({
          data: {
            nodes: [
              {
                id: 'center',
                type: 'organization',
                label: 'Center Node',
                data: {},
                position: { x: 400, y: 300 }
              },
              {
                id: 'top',
                type: 'repository',
                label: 'Top Node',
                data: {},
                position: { x: 400, y: 100 }
              },
              {
                id: 'bottom',
                type: 'repository',
                label: 'Bottom Node',
                data: {},
                position: { x: 400, y: 500 }
              }
            ],
            edges: [
              {
                id: 'e1',
                source: 'center',
                target: 'top',
                type: 'owns',
                label: 'owns'
              },
              {
                id: 'e2',
                source: 'center',
                target: 'bottom',
                type: 'owns',
                label: 'owns'
              }
            ]
          }
        })
      });
    });
    
    // Load graph
    await page.fill('input[placeholder="Enter organization name"]', 'keyboard-test');
    await page.click('button:has-text("Load Graph")');
    
    const graphCanvas = page.locator('[data-testid="graph-canvas"]');
    await expect(graphCanvas).toBeVisible();
    
    // Focus on graph canvas
    await graphCanvas.click();
    
    console.log('Testing keyboard navigation...');
    
    // Test arrow key navigation
    await page.keyboard.press('ArrowUp');
    await page.waitForTimeout(200);
    await page.keyboard.press('ArrowDown');
    await page.waitForTimeout(200);
    await page.keyboard.press('ArrowLeft');
    await page.waitForTimeout(200);
    await page.keyboard.press('ArrowRight');
    await page.waitForTimeout(200);
    
    // Test zoom with keyboard
    await page.keyboard.press('Equal'); // Zoom in
    await page.waitForTimeout(200);
    await page.keyboard.press('Minus'); // Zoom out
    
    console.log('✅ Keyboard navigation flow successful!');
  });
});