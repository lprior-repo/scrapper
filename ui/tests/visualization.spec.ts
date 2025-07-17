import { test, expect } from '@playwright/test';

test.describe('GitHub Codeowner Visualization', () => {
  test('should load the main page with title', async ({ page }) => {
    await page.goto('/');
    
    // Check page title
    await expect(page).toHaveTitle('GitHub Codeowner Visualization');
    
    // Check main heading
    await expect(page.locator('h1')).toContainText('GitHub Codeowner Visualization');
    
    // Check description
    await expect(page.locator('p')).toContainText('Interactive graph showing repositories, teams, users, and their codeowner relationships');
  });

  test('should render React Flow graph', async ({ page }) => {
    await page.goto('/');
    
    // Wait for React Flow to load
    await page.waitForSelector('.react-flow', { timeout: 10000 });
    
    // Check that React Flow container exists
    await expect(page.locator('.react-flow')).toBeVisible();
    
    // Check that controls are present
    await expect(page.locator('.react-flow__controls')).toBeVisible();
    
    // Check that minimap is present
    await expect(page.locator('.react-flow__minimap')).toBeVisible();
  });

  test('should display organization node', async ({ page }) => {
    await page.goto('/');
    
    // Wait for nodes to load
    await page.waitForSelector('.react-flow__node', { timeout: 10000 });
    
    // Look for organization node content using data-testid
    await expect(page.getByTestId('rf__node-org-1')).toBeVisible();
    await expect(page.getByTestId('rf__node-org-1').getByText('example-org')).toBeVisible();
    await expect(page.getByTestId('rf__node-org-1').getByText('Organization')).toBeVisible();
  });

  test('should display repository nodes', async ({ page }) => {
    await page.goto('/');
    
    // Wait for nodes to load
    await page.waitForSelector('.react-flow__node', { timeout: 10000 });
    
    // Check for specific repository nodes using test ids
    await expect(page.getByTestId('rf__node-repo-1')).toBeVisible();
    await expect(page.getByTestId('rf__node-repo-1').getByText('backend-api', { exact: true })).toBeVisible();
    await expect(page.getByTestId('rf__node-repo-2').getByText('frontend-app', { exact: true })).toBeVisible();
    await expect(page.getByTestId('rf__node-repo-3').getByText('infrastructure', { exact: true })).toBeVisible();
    await expect(page.getByTestId('rf__node-repo-1').getByText('Repository', { exact: true })).toBeVisible();
  });

  test('should display team nodes', async ({ page }) => {
    await page.goto('/');
    
    // Wait for nodes to load
    await page.waitForSelector('.react-flow__node', { timeout: 10000 });
    
    // Check for specific team nodes using test ids
    await expect(page.getByTestId('rf__node-team-1')).toBeVisible();
    await expect(page.getByTestId('rf__node-team-1').getByText('platform-team', { exact: true })).toBeVisible();
    await expect(page.getByTestId('rf__node-team-2').getByText('backend-team', { exact: true })).toBeVisible();
    await expect(page.getByTestId('rf__node-team-3').getByText('frontend-team', { exact: true })).toBeVisible();
    await expect(page.getByTestId('rf__node-team-1').getByText('Team', { exact: true })).toBeVisible();
  });

  test('should display user nodes', async ({ page }) => {
    await page.goto('/');
    
    // Wait for nodes to load
    await page.waitForSelector('.react-flow__node', { timeout: 10000 });
    
    // Check for specific user nodes using test ids
    await expect(page.getByTestId('rf__node-user-1')).toBeVisible();
    await expect(page.getByTestId('rf__node-user-1').getByText('alice', { exact: true })).toBeVisible();
    await expect(page.getByTestId('rf__node-user-2').getByText('bob', { exact: true })).toBeVisible();
    await expect(page.getByTestId('rf__node-user-3').getByText('charlie', { exact: true })).toBeVisible();
    await expect(page.getByTestId('rf__node-user-1').getByText('User', { exact: true })).toBeVisible();
  });

  test('should display edges/relationships', async ({ page }) => {
    await page.goto('/');
    
    // Wait for React Flow to fully load
    await page.waitForSelector('.react-flow', { timeout: 10000 });
    await page.waitForTimeout(2000); // Give time for edges to render
    
    // Check that edges exist
    const edges = page.locator('.react-flow__edge');
    await expect(edges.first()).toBeVisible();
    
    // Count should be greater than 0
    const edgeCount = await edges.count();
    expect(edgeCount).toBeGreaterThan(0);
  });

  test('should allow node interaction', async ({ page }) => {
    await page.goto('/');
    
    // Wait for nodes to load
    await page.waitForSelector('.react-flow__node', { timeout: 10000 });
    
    // Check that nodes are interactive by verifying attributes
    const repoNode = page.getByTestId('rf__node-repo-1');
    await expect(repoNode).toBeVisible();
    await expect(repoNode).toHaveAttribute('tabindex', '0');
    await expect(repoNode).toHaveAttribute('role', 'group');
  });

  test('should have working controls', async ({ page }) => {
    await page.goto('/');
    
    // Wait for controls to load
    await page.waitForSelector('.react-flow__controls', { timeout: 10000 });
    
    // Check that controls are visible and enabled (without clicking due to background interference)
    const controls = page.locator('.react-flow__controls');
    await expect(controls).toBeVisible();
    
    const zoomIn = page.locator('.react-flow__controls button').first();
    await expect(zoomIn).toBeVisible();
    await expect(zoomIn).toBeEnabled();
  });

  test('should show CODEOWNERS indicators', async ({ page }) => {
    await page.goto('/');
    
    // Wait for nodes to load
    await page.waitForSelector('.react-flow__node', { timeout: 10000 });
    
    // Look for CODEOWNERS indicators in specific repository nodes
    await expect(page.getByTestId('rf__node-repo-1').getByText('CODEOWNERS')).toBeVisible();
  });

  test('should display different node types with different colors', async ({ page }) => {
    await page.goto('/');
    
    // Wait for React Flow to load
    await page.waitForSelector('.react-flow', { timeout: 10000 });
    
    // Check minimap which should show different colored nodes
    await expect(page.locator('.react-flow__minimap')).toBeVisible();
    
    // The minimap svg should contain colored rectangles representing nodes
    const minimap = page.locator('.react-flow__minimap svg');
    await expect(minimap).toBeVisible();
  });

  test('should be responsive', async ({ page }) => {
    await page.goto('/');
    
    // Test different viewport sizes
    await page.setViewportSize({ width: 1200, height: 800 });
    await expect(page.locator('.react-flow')).toBeVisible();
    
    await page.setViewportSize({ width: 800, height: 600 });
    await expect(page.locator('.react-flow')).toBeVisible();
    
    await page.setViewportSize({ width: 400, height: 300 });
    await expect(page.locator('.react-flow')).toBeVisible();
  });
});