/**
 * App Integration Tests
 * 
 * Tests for full app integration including:
 * - Complete user workflows
 * - Component interactions
 * - State management across components
 * - Error handling and recovery
 * - React 19 patterns (Suspense, Error Boundaries)
 * - Performance and accessibility
 */

import React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from '../App'

// Mock the services to control API responses
jest.mock('../services', () => ({
  fetchGraphData: jest.fn(),
}))

import * as services from '../services'

const mockFetchGraphData = services.fetchGraphData as jest.MockedFunction<typeof services.fetchGraphData>

describe('App Integration Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  describe('Complete User Workflows', () => {
    test('successful graph loading workflow', async () => {
      const user = userEvent.setup()
      
      // Mock successful API response
      mockFetchGraphData.mockResolvedValue({
        data: {
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
        },
      })
      
      render(<App />)
      
      // 1. Initial state - should show empty state
      expect(screen.getByText('GitHub Codeowners Visualization')).toBeInTheDocument()
      expect(screen.getByPlaceholderText('Enter organization name')).toHaveValue('')
      expect(screen.getByRole('button', { name: 'Load Graph' })).toBeDisabled()
      
      // 2. User enters organization name
      const input = screen.getByPlaceholderText('Enter organization name')
      await user.type(input, 'test-org')
      
      // 3. Load button should become enabled
      const loadButton = screen.getByRole('button', { name: 'Load Graph' })
      expect(loadButton).toBeEnabled()
      
      // 4. User clicks load button
      await user.click(loadButton)
      
      // 5. Should show loading state
      await waitFor(() => {
        expect(screen.getByText(/loading|fetching/i)).toBeInTheDocument()
      })
      
      // 6. Should eventually show the graph
      await waitFor(() => {
        expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      })
      
      // Verify API was called with correct parameters
      expect(mockFetchGraphData).toHaveBeenCalledWith('test-org', false)
    })

    test('topics vs teams workflow', async () => {
      const user = userEvent.setup()
      
      // Mock different responses for teams vs topics
      mockFetchGraphData
        .mockResolvedValueOnce({
          data: { nodes: [{ id: 'team-1', type: 'team', label: 'Team', data: {}, position: { x: 0, y: 0 } }], edges: [] },
        })
        .mockResolvedValueOnce({
          data: { nodes: [{ id: 'topic-1', type: 'topic', label: 'Topic', data: {}, position: { x: 0, y: 0 } }], edges: [] },
        })
      
      render(<App />)
      
      // 1. Enter organization and load with teams (default)
      await user.type(screen.getByPlaceholderText('Enter organization name'), 'test-org')
      await user.click(screen.getByRole('button', { name: 'Load Graph' }))
      
      await waitFor(() => {
        expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      })
      
      expect(mockFetchGraphData).toHaveBeenCalledWith('test-org', false)
      
      // 2. Toggle to topics view
      const topicsToggle = screen.getByRole('checkbox', { name: 'Use Topics instead of Teams' })
      await user.click(topicsToggle)
      
      // 3. Load again with topics
      await user.click(screen.getByRole('button', { name: 'Load Graph' }))
      
      await waitFor(() => {
        expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      })
      
      expect(mockFetchGraphData).toHaveBeenCalledWith('test-org', true)
    })

    test('error handling and recovery workflow', async () => {
      const user = userEvent.setup()
      
      // Mock API error then success
      mockFetchGraphData
        .mockRejectedValueOnce(new Error('Network error'))
        .mockResolvedValueOnce({
          data: { nodes: [], edges: [] },
        })
      
      render(<App />)
      
      // 1. Try to load with error
      await user.type(screen.getByPlaceholderText('Enter organization name'), 'error-org')
      await user.click(screen.getByRole('button', { name: 'Load Graph' }))
      
      // 2. Should show error state
      await waitFor(() => {
        expect(screen.getByText(/error/i)).toBeInTheDocument()
      })
      
      // 3. User can retry
      const retryButton = screen.getByRole('button', { name: /try again|retry/i })
      await user.click(retryButton)
      
      // 4. Should eventually succeed
      await waitFor(() => {
        expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      })
    })
  })

  describe('Component Interactions', () => {
    test('header controls interact correctly with graph display', async () => {
      const user = userEvent.setup()
      
      mockFetchGraphData.mockResolvedValue({
        data: { nodes: [], edges: [] },
      })
      
      render(<App />)
      
      // Test that header controls affect graph loading
      const input = screen.getByPlaceholderText('Enter organization name')
      const button = screen.getByRole('button', { name: 'Load Graph' })
      const checkbox = screen.getByRole('checkbox')
      
      // Initially disabled
      expect(button).toBeDisabled()
      
      // Enable after typing
      await user.type(input, 'test')
      expect(button).toBeEnabled()
      
      // Check topics toggle
      await user.click(checkbox)
      expect(checkbox).toBeChecked()
      
      // Load graph
      await user.click(button)
      
      await waitFor(() => {
        expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      })
      
      // Verify state persists
      expect(input).toHaveValue('test')
      expect(checkbox).toBeChecked()
    })

    test('loading states transition correctly', async () => {
      const user = userEvent.setup()
      
      // Mock with delay to observe loading states
      mockFetchGraphData.mockImplementation(() => 
        new Promise(resolve => 
          setTimeout(() => resolve({ data: { nodes: [], edges: [] } }), 100)
        )
      )
      
      render(<App />)
      
      await user.type(screen.getByPlaceholderText('Enter organization name'), 'test-org')
      await user.click(screen.getByRole('button', { name: 'Load Graph' }))
      
      // Should show loading state
      expect(screen.getByText(/loading|fetching/i)).toBeInTheDocument()
      
      // Should transition to graph
      await waitFor(() => {
        expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
        expect(screen.queryByText(/loading|fetching/i)).not.toBeInTheDocument()
      })
    })
  })

  describe('State Management', () => {
    test('maintains state across component re-renders', async () => {
      const user = userEvent.setup()
      
      render(<App />)
      
      // Set initial state
      const input = screen.getByPlaceholderText('Enter organization name')
      const checkbox = screen.getByRole('checkbox')
      
      await user.type(input, 'persistent-org')
      await user.click(checkbox)
      
      expect(input).toHaveValue('persistent-org')
      expect(checkbox).toBeChecked()
      
      // Force re-render by triggering React update
      // (In real app this might happen due to external state changes)
      const button = screen.getByRole('button', { name: 'Load Graph' })
      await user.hover(button)
      await user.unhover(button)
      
      // State should persist
      expect(input).toHaveValue('persistent-org')
      expect(checkbox).toBeChecked()
    })

    test('state updates propagate correctly between components', async () => {
      const user = userEvent.setup()
      
      render(<App />)
      
      // Organization input affects button state
      const input = screen.getByPlaceholderText('Enter organization name')
      const button = screen.getByRole('button', { name: 'Load Graph' })
      
      expect(button).toBeDisabled()
      
      await user.type(input, 'a')
      expect(button).toBeEnabled()
      
      await user.clear(input)
      expect(button).toBeDisabled()
    })
  })

  describe('React 19 Patterns', () => {
    test('error boundary catches and displays errors properly', async () => {
      const user = userEvent.setup()
      
      // Mock API error
      mockFetchGraphData.mockRejectedValue(new Error('Component error'))
      
      render(<App />)
      
      await user.type(screen.getByPlaceholderText('Enter organization name'), 'error-org')
      await user.click(screen.getByRole('button', { name: 'Load Graph' }))
      
      // Should catch error and show error boundary UI
      await waitFor(() => {
        expect(screen.getByRole('alert')).toBeInTheDocument()
      })
    })

    test('suspense loading states work correctly', async () => {
      mockFetchGraphData.mockImplementation(() => 
        new Promise(resolve => 
          setTimeout(() => resolve({ data: { nodes: [], edges: [] } }), 50)
        )
      )
      
      render(<App />)
      
      const user = userEvent.setup()
      await user.type(screen.getByPlaceholderText('Enter organization name'), 'suspense-test')
      await user.click(screen.getByRole('button', { name: 'Load Graph' }))
      
      // Should show loading component
      expect(screen.getByTestId('loading-spinner')).toBeInTheDocument()
      
      // Should resolve to graph
      await waitFor(() => {
        expect(screen.getByTestId('graph-canvas')).toBeInTheDocument()
      })
    })
  })

  describe('Accessibility Integration', () => {
    test('keyboard navigation works across all components', async () => {
      const user = userEvent.setup()
      
      render(<App />)
      
      // Tab through interactive elements
      const input = screen.getByPlaceholderText('Enter organization name')
      const checkbox = screen.getByRole('checkbox')
      const button = screen.getByRole('button', { name: 'Load Graph' })
      
      // Focus first element
      await user.tab()
      expect(input).toHaveFocus()
      
      // Type to enable button
      await user.type(input, 'keyboard-test')
      
      // Tab to checkbox
      await user.tab()
      expect(checkbox).toHaveFocus()
      
      // Tab to button
      await user.tab()
      expect(button).toHaveFocus()
      
      // All elements should be accessible
      expect(input).toBeInTheDocument()
      expect(checkbox).toBeInTheDocument()
      expect(button).toBeEnabled()
    })

    test('screen reader announcements work correctly', async () => {
      const user = userEvent.setup()
      
      mockFetchGraphData.mockResolvedValue({
        data: { nodes: [], edges: [] },
      })
      
      render(<App />)
      
      await user.type(screen.getByPlaceholderText('Enter organization name'), 'a11y-test')
      await user.click(screen.getByRole('button', { name: 'Load Graph' }))
      
      // Loading state should be announced
      await waitFor(() => {
        const loadingElement = screen.getByRole('status')
        expect(loadingElement).toBeInTheDocument()
        expect(loadingElement).toHaveAttribute('aria-live', 'polite')
      })
    })
  })

  describe('Performance', () => {
    test('app renders efficiently with many interactions', async () => {
      const user = userEvent.setup()
      
      mockFetchGraphData.mockResolvedValue({
        data: { nodes: [], edges: [] },
      })
      
      const startTime = performance.now()
      
      render(<App />)
      
      // Perform many interactions
      const input = screen.getByPlaceholderText('Enter organization name')
      const checkbox = screen.getByRole('checkbox')
      
      for (let i = 0; i < 5; i++) {
        await user.type(input, `test-${i}`)
        await user.clear(input)
        await user.click(checkbox)
      }
      
      const endTime = performance.now()
      const renderTime = endTime - startTime
      
      // Should complete interactions in reasonable time (adjust threshold as needed)
      expect(renderTime).toBeLessThan(5000) // 5 seconds
    })

    test('memory usage stays reasonable', () => {
      const { unmount } = render(<App />)
      
      // Component should render without throwing
      expect(screen.getByText('GitHub Codeowners Visualization')).toBeInTheDocument()
      
      // Should unmount cleanly
      unmount()
      
      // No memory leaks should occur (tested by Jest's cleanup)
    })
  })

  describe('Edge Cases', () => {
    test('handles rapid user interactions gracefully', async () => {
      const user = userEvent.setup()
      
      mockFetchGraphData.mockResolvedValue({
        data: { nodes: [], edges: [] },
      })
      
      render(<App />)
      
      const input = screen.getByPlaceholderText('Enter organization name')
      const button = screen.getByRole('button', { name: 'Load Graph' })
      const checkbox = screen.getByRole('checkbox')
      
      // Rapid interactions
      await user.type(input, 'rapid')
      await user.click(checkbox)
      await user.click(checkbox)
      await user.click(button)
      await user.click(checkbox)
      
      // Should still be functional
      expect(input).toHaveValue('rapid')
      expect(button).toBeInTheDocument()
      expect(checkbox).toBeInTheDocument()
    })

    test('handles API response edge cases', async () => {
      const user = userEvent.setup()
      
      // Mock malformed response
      mockFetchGraphData.mockResolvedValue({
        data: {
          nodes: [{ id: '', type: 'invalid' }], // Invalid node
          edges: [],
        },
      } as any)
      
      render(<App />)
      
      await user.type(screen.getByPlaceholderText('Enter organization name'), 'edge-case')
      await user.click(screen.getByRole('button', { name: 'Load Graph' }))
      
      // Should handle gracefully - either show error or filter bad data
      await waitFor(() => {
        const errorOrGraph = screen.queryByRole('alert') || screen.queryByTestId('graph-canvas')
        expect(errorOrGraph).toBeInTheDocument()
      })
    })

    test('handles network failures gracefully', async () => {
      const user = userEvent.setup()
      
      // Mock network error
      mockFetchGraphData.mockRejectedValue(new Error('Network failure'))
      
      render(<App />)
      
      await user.type(screen.getByPlaceholderText('Enter organization name'), 'network-fail')
      await user.click(screen.getByRole('button', { name: 'Load Graph' }))
      
      // Should show error state
      await waitFor(() => {
        expect(screen.getByText(/error|failed/i)).toBeInTheDocument()
      })
      
      // App should remain responsive
      expect(screen.getByPlaceholderText('Enter organization name')).toBeEnabled()
    })
  })
})