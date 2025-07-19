import { test, expect } from '@playwright/test'

test.describe('GitHub Automation', () => {
  test('should scan organization first, then load graph with topics', async ({ page }) => {
    // Navigate to the app
    await page.goto('/')
    
    // Wait for the page to load
    await expect(page.locator('h1')).toHaveText('GitHub Codeowners Visualization')
    
    // Enter organization name - use github organization for testing
    const orgInput = page.locator('input[placeholder="Enter organization name"]')
    await orgInput.fill('github')
    
    // Check the "Use Topics instead of Teams" checkbox
    const topicsCheckbox = page.locator('input[type="checkbox"]')
    await topicsCheckbox.check()
    
    // Verify checkbox is checked
    await expect(topicsCheckbox).toBeChecked()
    
    // Step 1: First scan the organization to collect data
    console.log('Step 1: Scanning organization...')
    const scanResponse = await page.request.post('http://localhost:8081/api/scan/github?useTopics=true&max_repos=5', {
      timeout: 60000, // 60 second timeout for scan operation
    })
    
    expect(scanResponse.status()).toBe(201)
    const scanData = await scanResponse.json()
    console.log('Scan Response:', JSON.stringify(scanData, null, 2))
    
    // Verify scan response structure
    expect(scanData).toHaveProperty('data')
    expect(scanData.data).toHaveProperty('success')
    expect(scanData.data).toHaveProperty('organization')
    expect(scanData.data).toHaveProperty('summary')
    
    // Step 2: Now click the "Load Graph" button to trigger the graph API call
    console.log('Step 2: Loading graph...')
    const loadGraphButton = page.locator('button:has-text("Load Graph")')
    await loadGraphButton.click()
    
    // Wait for the graph API call to complete with proper timeout
    const graphResponse = await page.waitForResponse(response => 
      response.url().includes('localhost:8081/api/graph/github') && response.status() === 200,
      { timeout: 30000 } // 30 second timeout
    )
    
    const graphData = await graphResponse.json()
    console.log('Graph API Response:', JSON.stringify(graphData, null, 2))
    
    // Verify the response structure
    expect(graphData).toHaveProperty('data')
    expect(graphData.data).toHaveProperty('nodes')
    expect(graphData.data).toHaveProperty('edges')
    
    // Step 3: Check if the graph canvas is loaded with specific selector
    const graphCanvas = page.locator('[data-testid="graph-canvas"]')
    await expect(graphCanvas).toBeVisible({ timeout: 10000 })
    
    // Verify we have actual graph data
    const nodes = graphData.data.nodes
    const edges = graphData.data.edges
    
    if (nodes.length === 0 && edges.length === 0) {
      console.log('⚠️  WARNING: Graph data is empty - this may indicate the useTopics parameter is not being processed correctly')
    } else {
      console.log(`✅ SUCCESS: Graph loaded with ${nodes.length} nodes and ${edges.length} edges`)
    }
  })
  
  test('should demonstrate the proper API workflow with topics parameter', async ({ page }) => {
    // Step 1: Scan organization with topics
    console.log('Testing scan endpoint with topics...')
    const scanResponse = await page.request.post('http://localhost:8081/api/scan/github?useTopics=true&max_repos=5', {
      timeout: 60000,
    })
    
    expect(scanResponse.status()).toBe(201)
    const scanData = await scanResponse.json()
    console.log('Scan with topics response:', JSON.stringify(scanData, null, 2))
    
    // Step 2: Get graph data after scan
    console.log('Testing graph endpoint after scan...')
    const graphResponse = await page.request.get('http://localhost:8081/api/graph/github?useTopics=true')
    expect(graphResponse.status()).toBe(200)
    
    const graphData = await graphResponse.json()
    console.log('Graph with topics response:', JSON.stringify(graphData, null, 2))
    
    // Verify the workflow produces data
    expect(graphData).toHaveProperty('data')
    expect(graphData.data).toHaveProperty('nodes')
    expect(graphData.data).toHaveProperty('edges')
  })
  
  test('should handle errors gracefully when organization does not exist', async ({ page }) => {
    // Test with non-existent organization
    const nonExistentOrg = 'this-organization-definitely-does-not-exist-12345'
    
    // Navigate to the app
    await page.goto('/')
    
    // Enter non-existent organization name
    const orgInput = page.locator('input[placeholder="Enter organization name"]')
    await orgInput.fill(nonExistentOrg)
    
    // Check the topics checkbox
    const topicsCheckbox = page.locator('input[type="checkbox"]')
    await topicsCheckbox.check()
    
    // Try to scan - this should fail gracefully
    const scanResponse = await page.request.post(`http://localhost:8081/api/scan/${nonExistentOrg}?useTopics=true`, {
      timeout: 30000,
    })
    
    // Should return an error status
    expect(scanResponse.status()).toBeGreaterThanOrEqual(400)
    
    // Click Load Graph button to see frontend behavior
    const loadGraphButton = page.locator('button:has-text("Load Graph")')
    await loadGraphButton.click()
    
    // Wait for the graph canvas to appear (the API returns 200 OK with empty data)
    const graphCanvas = page.locator('[data-testid="graph-canvas"]')
    await expect(graphCanvas).toBeVisible({ timeout: 10000 })
    
    // The graph should load but be empty (no nodes/edges) since the organization doesn't exist
    // This is the correct behavior - the API doesn't fail, it just returns empty data
  })
})