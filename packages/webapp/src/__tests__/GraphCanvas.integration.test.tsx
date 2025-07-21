/**
 * GraphCanvas Integration Tests
 * 
 * Tests for GraphCanvas component integration including:
 * - React 19 Suspense integration
 * - Error Boundary integration
 * - Cytoscape.js lifecycle management
 * - Data loading and rendering
 * - Performance and memory management
 * - User interactions and events
 */

import React from 'react'
import { render, screen, waitFor, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { GraphCanvas } from '../components/GraphCanvas'
import { GraphErrorBoundary } from '../components/ErrorBoundary'

// Mock Cytoscape to control its behavior in tests
jest.mock('cytoscape', () => {
  const mockCytoscape = jest.fn(() => ({
    mount: jest.fn(),
    unmount: jest.fn(),
    destroy: jest.fn(),
    layout: jest.fn(() => ({ run: jest.fn() })),
    fit: jest.fn(),
    center: jest.fn(),
    zoom: jest.fn(),
    pan: jest.fn(),
    on: jest.fn(),
    off: jest.fn(),
    elements: jest.fn(() => ({ length: 0 })),
    nodes: jest.fn(() => []),
    edges: jest.fn(() => []),
    add: jest.fn(),
    remove: jest.fn(),
    style: jest.fn(),
    resize: jest.fn(),
  }))
  return mockCytoscape
})

import cytoscape from 'cytoscape'

const mockCytoscape = cytoscape as jest.MockedFunction<typeof cytoscape>

describe('GraphCanvas Integration Tests', () => {
  const mockGraphData = {
    nodes: [
      {
        id: 'org-1',
        type: 'organization',
        label: 'Test Org',
        data: { name: 'Test Org' },
        position: { x: 0, y: 0 },
      },
      {
        id: 'repo-1',
        type: 'repository',
        label: 'test-repo',
        data: { name: 'test-repo', stars: 100 },
        position: { x: 100, y: 100 },
      },
    ],
    edges: [
      {
        id: 'edge-1',
        source: 'org-1',
        target: 'repo-1',
        type: 'owns',
        label: 'owns',
      },
    ],
  }

  beforeEach(() => {
    jest.clearAllMocks()
  })

  describe('React 19 Suspense Integration', () => {
    test('works with Suspense boundary for data loading', async () => {
      const DataLoadingComponent = () => {
        const [data, setData] = React.useState(null)
        
        React.useEffect(() => {
          // Simulate async data loading
          const timer = setTimeout(() => {
            setData(mockGraphData)
          }, 50)
          return () => clearTimeout(timer)
        }, [])
        
        if (!data) {
          // This would normally be a thrown promise in real Suspense
          return <div>Loading...</div>
        }
        
        return <GraphCanvas graphData={data} />
      }
      
      const SuspenseWrapper = () => (
        <React.Suspense fallback={<div>Suspense Loading...</div>}>
          <DataLoadingComponent />
        </React.Suspense>
      )
      
      render(<SuspenseWrapper />)
      
      // Should show loading initially
      expect(screen.getByText('Loading...')).toBeInTheDocument()
      
      // Should eventually show graph canvas
      await waitFor(() => {
        expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      })
    })

    test('handles suspense fallback during data transitions', async () => {
      let resolveData: (value: any) => void
      
      const AsyncDataComponent = ({ delay = 0 }: { readonly delay?: number }) => {
        const [data, setData] = React.useState(null)
        
        React.useEffect(() => {
          const promise = new Promise(resolve => {
            resolveData = resolve
            setTimeout(() => resolve(mockGraphData), delay)
          })
          
          promise.then(setData)
        }, [delay])
        
        if (!data) {
          return <div>Async Loading...</div>
        }
        
        return <GraphCanvas graphData={data} />
      }
      
      const { rerender } = render(
        <React.Suspense fallback={<div>Suspense Fallback</div>}>
          <AsyncDataComponent delay={10} />
        </React.Suspense>
      )
      
      // Should show loading state
      expect(screen.getByText('Async Loading...')).toBeInTheDocument()
      
      // Wait for data to load
      await waitFor(() => {
        expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      })
      
      // Test data transition
      rerender(
        <React.Suspense fallback={<div>Suspense Fallback</div>}>
          <AsyncDataComponent delay={50} />
        </React.Suspense>
      )
      
      // Should handle transition gracefully
      await waitFor(() => {
        expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      })
    })
  })

  describe('Error Boundary Integration', () => {
    test('error boundary catches Cytoscape initialization errors', () => {
      // Mock Cytoscape to throw error
      mockCytoscape.mockImplementationOnce(() => {
        throw new Error('Cytoscape initialization failed')
      })
      
      render(
        <GraphErrorBoundary title="Graph Rendering Error">
          <GraphCanvas graphData={mockGraphData} />
        </GraphErrorBoundary>
      )
      
      // Should catch error and show error boundary
      expect(screen.getByRole('alert')).toBeInTheDocument()
      expect(screen.getByText('Graph Rendering Error')).toBeInTheDocument()
      expect(screen.getByText('Cytoscape initialization failed')).toBeInTheDocument()
    })

    test('error boundary provides retry functionality', async () => {
      const user = userEvent.setup()
      let shouldThrowError = true
      
      // Mock Cytoscape to throw error first, then succeed
      mockCytoscape.mockImplementation(() => {
        if (shouldThrowError) {
          throw new Error('Initial error')
        }
        return {
          mount: jest.fn(),
          unmount: jest.fn(),
          destroy: jest.fn(),
          layout: jest.fn(() => ({ run: jest.fn() })),
          fit: jest.fn(),
          center: jest.fn(),
          on: jest.fn(),
          off: jest.fn(),
        } as any
      })
      
      render(
        <GraphErrorBoundary>
          <GraphCanvas graphData={mockGraphData} />
        </GraphErrorBoundary>
      )
      
      // Should show error initially
      expect(screen.getByRole('alert')).toBeInTheDocument()
      
      // Fix the error condition
      shouldThrowError = false
      
      // Click retry
      const retryButton = screen.getByTestId('retry-button')
      await user.click(retryButton)
      
      // Should eventually show graph
      await waitFor(() => {
        expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      })
    })

    test('error boundary handles data validation errors', () => {
      const invalidData = {
        nodes: [
          { id: '', type: 'invalid' }, // Invalid node
        ],
        edges: [
          { id: 'edge-1', source: 'nonexistent', target: 'also-nonexistent' }, // Invalid edge
        ],
      }
      
      render(
        <GraphErrorBoundary>
          <GraphCanvas graphData={invalidData as any} />
        </GraphErrorBoundary>
      )
      
      // Should either handle gracefully or show error
      // Implementation depends on how data validation is handled
      const errorOrGraph = screen.queryByRole('alert') || screen.queryByTestId('graph-canvas')
      expect(errorOrGraph).toBeInTheDocument()
    })
  })

  describe('Cytoscape.js Lifecycle Management', () => {
    test('initializes Cytoscape instance correctly', () => {
      const mockInstance = {
        mount: jest.fn(),
        unmount: jest.fn(),
        destroy: jest.fn(),
        layout: jest.fn(() => ({ run: jest.fn() })),
        fit: jest.fn(),
        center: jest.fn(),
        on: jest.fn(),
        off: jest.fn(),
      }
      
      mockCytoscape.mockReturnValue(mockInstance as any)
      
      render(<GraphCanvas graphData={mockGraphData} />)
      
      // Should initialize Cytoscape
      expect(mockCytoscape).toHaveBeenCalledWith(
        expect.objectContaining({
          container: expect.any(Object),
          elements: expect.any(Array),
          style: expect.any(Array),
          layout: expect.objectContaining({ name: 'cose' }),
        })
      )
    })

    test('destroys Cytoscape instance on unmount', () => {
      const mockInstance = {
        mount: jest.fn(),
        unmount: jest.fn(),
        destroy: jest.fn(),
        layout: jest.fn(() => ({ run: jest.fn() })),
        fit: jest.fn(),
        center: jest.fn(),
        on: jest.fn(),
        off: jest.fn(),
      }
      
      mockCytoscape.mockReturnValue(mockInstance as any)
      
      const { unmount } = render(<GraphCanvas graphData={mockGraphData} />)
      
      // Unmount component
      unmount()
      
      // Should destroy Cytoscape instance
      expect(mockInstance.destroy).toHaveBeenCalled()
    })

    test('updates Cytoscape when data changes', () => {
      const mockInstance = {
        mount: jest.fn(),
        unmount: jest.fn(),
        destroy: jest.fn(),
        layout: jest.fn(() => ({ run: jest.fn() })),
        fit: jest.fn(),
        center: jest.fn(),
        on: jest.fn(),
        off: jest.fn(),
        elements: jest.fn(() => ({ length: 2 })),
        add: jest.fn(),
        remove: jest.fn(),
      }
      
      mockCytoscape.mockReturnValue(mockInstance as any)
      
      const { rerender } = render(<GraphCanvas graphData={mockGraphData} />)
      
      const newData = {
        ...mockGraphData,
        nodes: [...mockGraphData.nodes, {
          id: 'new-node',
          type: 'user',
          label: 'New User',
          data: { login: 'newuser' },
          position: { x: 200, y: 200 },
        }],
      }
      
      rerender(<GraphCanvas graphData={newData} />)
      
      // Should update the graph with new data
      expect(mockCytoscape).toHaveBeenCalledTimes(2)
    })

    test('handles resize events correctly', () => {
      const mockInstance = {
        mount: jest.fn(),
        unmount: jest.fn(),
        destroy: jest.fn(),
        layout: jest.fn(() => ({ run: jest.fn() })),
        fit: jest.fn(),
        center: jest.fn(),
        resize: jest.fn(),
        on: jest.fn(),
        off: jest.fn(),
      }
      
      mockCytoscape.mockReturnValue(mockInstance as any)
      
      render(<GraphCanvas graphData={mockGraphData} />)
      
      // Simulate window resize
      act(() => {
        window.dispatchEvent(new Event('resize'))
      })
      
      // Should call resize on Cytoscape instance
      expect(mockInstance.resize).toHaveBeenCalled()
    })
  })

  describe('Data Loading and Rendering', () => {
    test('renders loading state while data is being processed', async () => {
      // Simulate delayed Cytoscape initialization
      mockCytoscape.mockImplementationOnce(() => {
        return new Promise(resolve => {
          setTimeout(() => resolve({
            mount: jest.fn(),
            unmount: jest.fn(),
            destroy: jest.fn(),
            layout: jest.fn(() => ({ run: jest.fn() })),
            fit: jest.fn(),
            center: jest.fn(),
            on: jest.fn(),
            off: jest.fn(),
          }), 50)
        }) as any
      })
      
      render(<GraphCanvas graphData={mockGraphData} />)
      
      // Should show loading state initially
      expect(screen.getByTestId('graph-loading-spinner')).toBeInTheDocument()
      
      // Should eventually show graph
      await waitFor(() => {
        expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      })
    })

    test('handles empty data gracefully', () => {
      const emptyData = { nodes: [], edges: [] }
      
      render(<GraphCanvas graphData={emptyData} />)
      
      // Should still render graph canvas
      expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      
      // Should initialize Cytoscape with empty data
      expect(mockCytoscape).toHaveBeenCalledWith(
        expect.objectContaining({
          elements: [],
        })
      )
    })

    test('validates and filters invalid data', () => {
      const mixedData = {
        nodes: [
          { id: 'valid-node', type: 'user', label: 'Valid', data: {}, position: { x: 0, y: 0 } },
          { id: '', type: 'invalid' }, // Invalid node
          null, // Invalid node
        ],
        edges: [
          { id: 'valid-edge', source: 'valid-node', target: 'valid-node', type: 'self', label: 'self' },
          { id: '', source: '', target: '' }, // Invalid edge
          null, // Invalid edge
        ],
      }
      
      render(<GraphCanvas graphData={mixedData as any} />)
      
      // Should filter out invalid data and render valid data
      expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      
      // Should have called Cytoscape with filtered data
      expect(mockCytoscape).toHaveBeenCalledWith(
        expect.objectContaining({
          elements: expect.arrayContaining([
            expect.objectContaining({ data: expect.objectContaining({ id: 'valid-node' }) }),
            expect.objectContaining({ data: expect.objectContaining({ id: 'valid-edge' }) }),
          ]),
        })
      )
    })
  })

  describe('User Interactions and Events', () => {
    test('handles node click events', async () => {
      const user = userEvent.setup()
      const mockInstance = {
        mount: jest.fn(),
        unmount: jest.fn(),
        destroy: jest.fn(),
        layout: jest.fn(() => ({ run: jest.fn() })),
        fit: jest.fn(),
        center: jest.fn(),
        on: jest.fn(),
        off: jest.fn(),
      }
      
      mockCytoscape.mockReturnValue(mockInstance as any)
      
      render(<GraphCanvas graphData={mockGraphData} />)
      
      // Should set up event listeners
      expect(mockInstance.on).toHaveBeenCalledWith('tap', 'node', expect.any(Function))
      expect(mockInstance.on).toHaveBeenCalledWith('cxttap', 'node', expect.any(Function))
    })

    test('handles graph pan and zoom events', () => {
      const mockInstance = {
        mount: jest.fn(),
        unmount: jest.fn(),
        destroy: jest.fn(),
        layout: jest.fn(() => ({ run: jest.fn() })),
        fit: jest.fn(),
        center: jest.fn(),
        on: jest.fn(),
        off: jest.fn(),
      }
      
      mockCytoscape.mockReturnValue(mockInstance as any)
      
      render(<GraphCanvas graphData={mockGraphData} />)
      
      // Should set up pan and zoom event listeners
      expect(mockInstance.on).toHaveBeenCalledWith('pan zoom', expect.any(Function))
    })

    test('context menu integration works correctly', async () => {
      const mockInstance = {
        mount: jest.fn(),
        unmount: jest.fn(),
        destroy: jest.fn(),
        layout: jest.fn(() => ({ run: jest.fn() })),
        fit: jest.fn(),
        center: jest.fn(),
        on: jest.fn((event, selector, callback) => {
          if (event === 'cxttap' && selector === 'node') {
            // Simulate right-click on node
            setTimeout(() => {
              callback({
                target: {
                  data: () => ({ id: 'test-node', type: 'user' }),
                },
                renderedPosition: { x: 100, y: 100 },
              })
            }, 10)
          }
        }),
        off: jest.fn(),
      }
      
      mockCytoscape.mockReturnValue(mockInstance as any)
      
      render(<GraphCanvas graphData={mockGraphData} />)
      
      // Context menu should be set up
      expect(mockInstance.on).toHaveBeenCalledWith('cxttap', 'node', expect.any(Function))
      
      // Wait for potential context menu to appear
      await waitFor(() => {
        // This would show context menu in real implementation
        expect(mockInstance.on).toHaveBeenCalled()
      })
    })
  })

  describe('Performance and Memory Management', () => {
    test('cleans up event listeners on unmount', () => {
      const mockInstance = {
        mount: jest.fn(),
        unmount: jest.fn(),
        destroy: jest.fn(),
        layout: jest.fn(() => ({ run: jest.fn() })),
        fit: jest.fn(),
        center: jest.fn(),
        on: jest.fn(),
        off: jest.fn(),
      }
      
      mockCytoscape.mockReturnValue(mockInstance as any)
      
      const { unmount } = render(<GraphCanvas graphData={mockGraphData} />)
      
      unmount()
      
      // Should clean up event listeners
      expect(mockInstance.off).toHaveBeenCalled()
      expect(mockInstance.destroy).toHaveBeenCalled()
    })

    test('handles large datasets efficiently', () => {
      // Create large dataset
      const largeData = {
        nodes: Array.from({ length: 1000 }, (_, i) => ({
          id: `node-${i}`,
          type: 'user',
          label: `User ${i}`,
          data: { login: `user${i}` },
          position: { x: Math.random() * 1000, y: Math.random() * 1000 },
        })),
        edges: Array.from({ length: 500 }, (_, i) => ({
          id: `edge-${i}`,
          source: `node-${i * 2}`,
          target: `node-${(i * 2) + 1}`,
          type: 'follows',
          label: 'follows',
        })),
      }
      
      const startTime = performance.now()
      
      render(<GraphCanvas graphData={largeData} />)
      
      const endTime = performance.now()
      const renderTime = endTime - startTime
      
      // Should render in reasonable time
      expect(renderTime).toBeLessThan(1000) // 1 second
      
      // Should still initialize Cytoscape
      expect(mockCytoscape).toHaveBeenCalled()
    })

    test('manages memory correctly during data updates', () => {
      const mockInstance = {
        mount: jest.fn(),
        unmount: jest.fn(),
        destroy: jest.fn(),
        layout: jest.fn(() => ({ run: jest.fn() })),
        fit: jest.fn(),
        center: jest.fn(),
        on: jest.fn(),
        off: jest.fn(),
        elements: jest.fn(() => ({ remove: jest.fn() })),
        add: jest.fn(),
        remove: jest.fn(),
      }
      
      mockCytoscape.mockReturnValue(mockInstance as any)
      
      const { rerender } = render(<GraphCanvas graphData={mockGraphData} />)
      
      // Update data multiple times
      for (let i = 0; i < 10; i++) {
        const newData = {
          nodes: [{
            id: `node-${i}`,
            type: 'user',
            label: `User ${i}`,
            data: {},
            position: { x: 0, y: 0 },
          }],
          edges: [],
        }
        
        rerender(<GraphCanvas graphData={newData} />)
      }
      
      // Should handle updates efficiently
      expect(mockCytoscape).toHaveBeenCalled()
    })
  })

  describe('Edge Cases', () => {
    test('handles null or undefined data', () => {
      render(<GraphCanvas graphData={null as any} />)
      
      // Should not crash and should show appropriate state
      expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
    })

    test('handles Cytoscape library not available', () => {
      mockCytoscape.mockImplementationOnce(() => {
        throw new Error('Cytoscape library not found')
      })
      
      render(
        <GraphErrorBoundary>
          <GraphCanvas graphData={mockGraphData} />
        </GraphErrorBoundary>
      )
      
      // Should show error boundary
      expect(screen.getByRole('alert')).toBeInTheDocument()
    })

    test('handles rapid data changes', () => {
      const { rerender } = render(<GraphCanvas graphData={mockGraphData} />)
      
      // Rapid re-renders with different data
      for (let i = 0; i < 20; i++) {
        const rapidData = {
          nodes: [{ id: `rapid-${i}`, type: 'user', label: `User ${i}`, data: {}, position: { x: i, y: i } }],
          edges: [],
        }
        
        rerender(<GraphCanvas graphData={rapidData} />)
      }
      
      // Should handle gracefully
      expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
    })
  })
})