import { test, expect } from '@playwright/test'

test.describe('Comprehensive API Integration Tests (Covering Hurl Test Suite)', () => {
  const baseUrl = 'http://localhost:8081'
  const testOrgSmall = 'golang'
  const testOrgInvalid = 'thisorgdoesnotexist12345'
  const userAgent = 'Playwright-Test-Suite/1.0'

  test.describe('Health Endpoint Coverage', () => {
    test('should pass all health endpoint scenarios from hurl tests', async ({
      page,
    }) => {
      console.log('🏥 Testing Health API endpoints...')

      // Test 1: Basic health check
      const healthResponse = await page.request.get(`${baseUrl}/api/health`, {
        headers: {
          'User-Agent': userAgent,
          Accept: 'application/json',
        },
      })

      expect(healthResponse.status()).toBe(200)
      expect(healthResponse.headers()['content-type']).toContain(
        'application/json'
      )

      const healthData = await healthResponse.json()
      expect(healthData.data.status).toBe('healthy')
      expect(healthData.data.database).toBe('connected')
      expect(healthData.data.version).toBeTruthy()
      expect(healthData.data.timestamp).toBeTruthy()
      expect(healthData.data.version).toBe('1.0.0')

      // Validate timestamp format (ISO 8601)
      expect(healthData.data.timestamp).toMatch(
        /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.*$/
      )

      // Test 2: Health check with different Accept headers
      const healthResponse2 = await page.request.get(`${baseUrl}/api/health`, {
        headers: {
          'User-Agent': userAgent,
          Accept: 'application/json, text/plain, */*',
        },
      })

      expect(healthResponse2.status()).toBe(200)
      const healthData2 = await healthResponse2.json()
      expect(healthData2.data.status).toBe('healthy')

      // Test 3: Health check with no Accept header
      const healthResponse3 = await page.request.get(`${baseUrl}/api/health`, {
        headers: {
          'User-Agent': userAgent,
        },
      })

      expect(healthResponse3.status()).toBe(200)
      const healthData3 = await healthResponse3.json()
      expect(healthData3.data.status).toBe('healthy')

      // Test 4: Health check response timing (should be under 1 second)
      const startTime = Date.now()
      const healthResponse4 = await page.request.get(`${baseUrl}/api/health`)
      const endTime = Date.now()

      expect(healthResponse4.status()).toBe(200)
      expect(endTime - startTime).toBeLessThan(1000)

      console.log('✅ All health endpoint tests passed')
    })

    test('should validate health check consistency across multiple calls', async ({
      page,
    }) => {
      // Test consistency across multiple calls (hurl test 14)
      const promises = Array.from({ length: 5 }, () =>
        page.request.get(`${baseUrl}/api/health`, {
          headers: { 'User-Agent': userAgent },
        })
      )

      const responses = await Promise.all(promises)

      for (const response of responses) {
        expect(response.status()).toBe(200)
        const data = await response.json()
        expect(data.data.status).toBe('healthy')
        expect(data.data.database).toBe('connected')
      }

      console.log('✅ Health check consistency test passed')
    })
  })

  test.describe('Graph Endpoint Coverage', () => {
    test('should handle valid and invalid organization graph requests', async ({
      page,
    }) => {
      console.log('📊 Testing Graph API endpoints...')

      // Test 1: Basic graph retrieval for valid organization
      const graphResponse = await page.request.get(
        `${baseUrl}/api/graph/${testOrgSmall}`,
        {
          headers: {
            'User-Agent': userAgent,
            Accept: 'application/json',
          },
        }
      )

      expect(graphResponse.status()).toBe(200)
      expect(graphResponse.headers()['content-type']).toContain(
        'application/json'
      )

      const graphData = await graphResponse.json()
      expect(graphData.data.nodes).toBeDefined()
      expect(graphData.data.edges).toBeDefined()
      expect(Array.isArray(graphData.data.nodes)).toBe(true)
      expect(Array.isArray(graphData.data.edges)).toBe(true)

      // Test 2: Graph with invalid organization returns empty graph
      const invalidGraphResponse = await page.request.get(
        `${baseUrl}/api/graph/${testOrgInvalid}`,
        {
          headers: {
            'User-Agent': userAgent,
            Accept: 'application/json',
          },
        }
      )

      expect(invalidGraphResponse.status()).toBe(200)
      const invalidGraphData = await invalidGraphResponse.json()
      expect(invalidGraphData.data.nodes).toBeDefined()
      expect(invalidGraphData.data.edges).toBeDefined()
      expect(invalidGraphData.data.nodes).toHaveLength(0)
      expect(invalidGraphData.data.edges).toHaveLength(0)

      console.log('✅ Graph endpoint basic tests passed')
    })

    test('should handle topics parameter correctly', async ({ page }) => {
      // Test graph with useTopics parameter
      const topicsResponse = await page.request.get(
        `${baseUrl}/api/graph/${testOrgSmall}?useTopics=true`,
        {
          headers: {
            'User-Agent': userAgent,
            Accept: 'application/json',
          },
        }
      )

      expect(topicsResponse.status()).toBe(200)
      const topicsData = await topicsResponse.json()
      expect(topicsData.data.nodes).toBeDefined()
      expect(topicsData.data.edges).toBeDefined()

      // Compare with teams view
      const teamsResponse = await page.request.get(
        `${baseUrl}/api/graph/${testOrgSmall}`,
        {
          headers: {
            'User-Agent': userAgent,
            Accept: 'application/json',
          },
        }
      )

      expect(teamsResponse.status()).toBe(200)
      const teamsData = await teamsResponse.json()

      // Both should be valid responses
      expect(Array.isArray(topicsData.data.nodes)).toBe(true)
      expect(Array.isArray(teamsData.data.nodes)).toBe(true)

      console.log('✅ Topics parameter test passed')
    })
  })

  test.describe('Scan Endpoint Coverage', () => {
    test('should handle organization scanning scenarios', async ({ page }) => {
      console.log('🔍 Testing Scan API endpoints...')

      // Test 1: Basic organization scan with valid org
      const scanResponse = await page.request.post(
        `${baseUrl}/api/scan/${testOrgSmall}?max_repos=2&max_teams=2`,
        {
          headers: {
            'User-Agent': userAgent,
            Accept: 'application/json',
          },
          timeout: 60000, // Scanning can take time
        }
      )

      expect(scanResponse.status()).toBe(201)
      expect(scanResponse.headers()['content-type']).toContain(
        'application/json'
      )

      const scanData = await scanResponse.json()
      expect(scanData.data.success).toBe(true)
      expect(scanData.data.organization).toBe(testOrgSmall)
      expect(scanData.data.summary).toBeDefined()
      expect(scanData.data.summary.total_repos).toBeDefined()
      expect(scanData.data.data).toBeDefined()

      // Test 2: Scan with invalid organization name
      const invalidScanResponse = await page.request.post(
        `${baseUrl}/api/scan/${testOrgInvalid}`,
        {
          headers: {
            'User-Agent': userAgent,
            Accept: 'application/json',
          },
        }
      )

      expect(invalidScanResponse.status()).toBe(404)
      const invalidScanData = await invalidScanResponse.json()
      expect(invalidScanData.error).toBeDefined()

      // Test 3: Scan with minimal parameters
      const minimalScanResponse = await page.request.post(
        `${baseUrl}/api/scan/${testOrgSmall}?max_repos=1&max_teams=1`,
        {
          headers: {
            'User-Agent': userAgent,
            Accept: 'application/json',
          },
          timeout: 60000,
        }
      )

      expect(minimalScanResponse.status()).toBe(201)
      const minimalScanData = await minimalScanResponse.json()
      expect(minimalScanData.data.success).toBe(true)
      expect(minimalScanData.data.summary.total_repos).toBeLessThanOrEqual(1)

      console.log('✅ Scan endpoint tests passed')
    })
  })

  test.describe('Stats Endpoint Coverage', () => {
    test('should handle stats retrieval scenarios', async ({ page }) => {
      console.log('📈 Testing Stats API endpoints...')

      // First ensure we have data by scanning
      await page.request
        .post(`${baseUrl}/api/scan/${testOrgSmall}?max_repos=2`, {
          timeout: 30000,
        })
        .catch(() => {
          // Ignore scan failures - data might already exist
        })

      // Test 1: Basic stats retrieval for valid organization
      const statsResponse = await page.request.get(
        `${baseUrl}/api/stats/${testOrgSmall}`,
        {
          headers: {
            'User-Agent': userAgent,
            Accept: 'application/json',
          },
        }
      )

      expect(statsResponse.status()).toBe(200)
      expect(statsResponse.headers()['content-type']).toContain(
        'application/json'
      )

      const statsData = await statsResponse.json()
      expect(statsData.data.total_repositories).toBeDefined()
      expect(statsData.data.total_codeowners).toBeDefined()

      // Test 2: Stats with invalid organization
      const invalidStatsResponse = await page.request.get(
        `${baseUrl}/api/stats/${testOrgInvalid}`,
        {
          headers: {
            'User-Agent': userAgent,
            Accept: 'application/json',
          },
        }
      )

      expect(invalidStatsResponse.status()).toBe(404)
      const invalidStatsData = await invalidStatsResponse.json()
      expect(invalidStatsData.error).toBeDefined()

      console.log('✅ Stats endpoint tests passed')
    })
  })

  test.describe('API Schema Validation and Error Handling', () => {
    test('should validate response schemas match shared types', async ({
      page,
    }) => {
      console.log('🔍 Testing API schema validation...')

      // Test graph response schema
      const graphResponse = await page.request.get(
        `${baseUrl}/api/graph/${testOrgSmall}`
      )

      if (graphResponse.status() === 200) {
        const graphData = await graphResponse.json()

        // Validate top-level structure
        expect(graphData).toHaveProperty('data')
        expect(graphData.data).toHaveProperty('nodes')
        expect(graphData.data).toHaveProperty('edges')

        // Validate node structure if nodes exist
        if (graphData.data.nodes.length > 0) {
          const node = graphData.data.nodes[0]
          expect(node).toHaveProperty('id')
          expect(node).toHaveProperty('type')
          expect(node).toHaveProperty('label')
          expect(node).toHaveProperty('data')

          // Validate node types are from expected enum
          const validNodeTypes = [
            'organization',
            'repository',
            'user',
            'team',
            'topic',
          ]
          expect(validNodeTypes).toContain(node.type)

          if (node.position) {
            expect(node.position).toHaveProperty('x')
            expect(node.position).toHaveProperty('y')
            expect(typeof node.position.x).toBe('number')
            expect(typeof node.position.y).toBe('number')
          }
        }

        // Validate edge structure if edges exist
        if (graphData.data.edges.length > 0) {
          const edge = graphData.data.edges[0]
          expect(edge).toHaveProperty('id')
          expect(edge).toHaveProperty('source')
          expect(edge).toHaveProperty('target')
          expect(edge).toHaveProperty('type')

          // Validate edge types are from expected enum
          const validEdgeTypes = ['owns', 'member_of', 'codeowner']
          expect(validEdgeTypes).toContain(edge.type)
        }
      }

      console.log('✅ Schema validation tests passed')
    })

    test('should handle malformed requests and edge cases', async ({
      page,
    }) => {
      console.log('⚠️ Testing error handling and edge cases...')

      // Test 1: Invalid endpoints return 404
      const invalidEndpointResponse = await page.request.get(
        `${baseUrl}/api/invalid-endpoint`
      )
      expect(invalidEndpointResponse.status()).toBe(404)

      // Test 2: Invalid HTTP methods
      const invalidMethodResponse = await page.request.delete(
        `${baseUrl}/api/health`
      )
      expect([405, 404]).toContain(invalidMethodResponse.status())

      // Test 3: Large organization names
      const longOrgName = 'a'.repeat(100)
      const longOrgResponse = await page.request.get(
        `${baseUrl}/api/graph/${longOrgName}`
      )
      expect([200, 400, 404]).toContain(longOrgResponse.status())

      // Test 4: Special characters in organization names
      const specialOrgName = 'org-with_special.chars123'
      const specialOrgResponse = await page.request.get(
        `${baseUrl}/api/graph/${specialOrgName}`
      )
      expect([200, 400, 404]).toContain(specialOrgResponse.status())

      // Test 5: Empty organization name
      const emptyOrgResponse = await page.request.get(`${baseUrl}/api/graph/`)
      expect([400, 404]).toContain(emptyOrgResponse.status())

      console.log('✅ Error handling tests passed')
    })
  })

  test.describe('Performance and Reliability', () => {
    test('should meet performance requirements', async ({ page }) => {
      console.log('⚡ Testing API performance...')

      // Test 1: Health check should respond quickly
      const healthStartTime = Date.now()
      const healthResponse = await page.request.get(`${baseUrl}/api/health`)
      const healthEndTime = Date.now()

      expect(healthResponse.status()).toBe(200)
      expect(healthEndTime - healthStartTime).toBeLessThan(1000) // Under 1 second

      // Test 2: Concurrent requests handling
      const concurrentRequests = Array.from({ length: 5 }, () =>
        page.request.get(`${baseUrl}/api/health`)
      )

      const concurrentResponses = await Promise.all(concurrentRequests)

      for (const response of concurrentResponses) {
        expect(response.status()).toBe(200)
      }

      // Test 3: Graph endpoint reasonable response time
      const graphStartTime = Date.now()
      const graphResponse = await page.request.get(
        `${baseUrl}/api/graph/${testOrgSmall}`
      )
      const graphEndTime = Date.now()

      expect(graphResponse.status()).toBe(200)
      expect(graphEndTime - graphStartTime).toBeLessThan(5000) // Under 5 seconds

      console.log('✅ Performance tests passed')
    })

    test('should handle rate limiting gracefully', async ({ page }) => {
      console.log('🚦 Testing rate limiting behavior...')

      // Make multiple rapid requests to test rate limiting
      const rapidRequests = Array.from({ length: 20 }, (_, i) =>
        page.request.get(`${baseUrl}/api/health`, {
          headers: { 'User-Agent': `Test-${i}` },
        })
      )

      const rapidResponses = await Promise.all(rapidRequests)

      // Should handle all requests gracefully (either success or proper rate limit response)
      for (const response of rapidResponses) {
        expect([200, 429, 503]).toContain(response.status())
      }

      console.log('✅ Rate limiting tests passed')
    })
  })

  test.describe('Data Integrity and Referential Integrity', () => {
    test('should maintain data consistency', async ({ page }) => {
      console.log('🔗 Testing data integrity...')

      // Scan organization first to ensure data exists
      await page.request
        .post(`${baseUrl}/api/scan/${testOrgSmall}?max_repos=5`, {
          timeout: 60000,
        })
        .catch(() => {
          // Ignore if already scanned
        })

      // Get graph data
      const graphResponse = await page.request.get(
        `${baseUrl}/api/graph/${testOrgSmall}`
      )

      if (graphResponse.status() === 200) {
        const graphData = await graphResponse.json()

        // Test referential integrity between nodes and edges
        const nodeIds = new Set(
          graphData.data.nodes.map((node: any) => node.id)
        )

        for (const edge of graphData.data.edges) {
          expect(nodeIds.has(edge.source)).toBe(true)
          expect(nodeIds.has(edge.target)).toBe(true)
        }

        // Test unique node IDs
        const nodeIdsList = graphData.data.nodes.map((node: any) => node.id)
        const uniqueNodeIds = new Set(nodeIdsList)
        expect(nodeIdsList.length).toBe(uniqueNodeIds.size)

        // Test unique edge IDs
        const edgeIdsList = graphData.data.edges.map((edge: any) => edge.id)
        const uniqueEdgeIds = new Set(edgeIdsList)
        expect(edgeIdsList.length).toBe(uniqueEdgeIds.size)

        console.log(
          `✅ Data integrity verified: ${graphData.data.nodes.length} nodes, ${graphData.data.edges.length} edges`
        )
      }
    })
  })
})
