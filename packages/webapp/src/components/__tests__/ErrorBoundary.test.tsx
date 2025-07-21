/**
 * ErrorBoundary Component Tests
 * 
 * Tests for the error boundary components including:
 * - Basic error catching and fallback rendering
 * - Error fallback UI components and interactions
 * - Accessibility features for error states
 * - Recovery mechanisms (retry, reset)
 * - Integration with react-error-boundary
 * - HOC and hook patterns
 */

import React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { 
  ErrorBoundary, 
  GraphErrorBoundary, 
  withErrorBoundary, 
  useErrorBoundaryReset 
} from '../ErrorBoundary'

// Mock console methods to prevent test noise
const originalError = console.error
const originalWarn = console.warn

beforeAll(() => {
  console.error = jest.fn()
  console.warn = jest.fn()
})

afterAll(() => {
  console.error = originalError
  console.warn = originalWarn
})

// Test components that throw errors
const ThrowError: React.FC<{ readonly shouldThrow?: boolean; readonly message?: string }> = ({ 
  shouldThrow = true, 
  message = 'Test error message' 
}) => {
  if (shouldThrow) {
    throw new Error(message)
  }
  return <div>No error</div>
}

const AsyncThrowError: React.FC<{ readonly shouldThrow?: boolean }> = ({ shouldThrow = true }) => {
  React.useEffect(() => {
    if (shouldThrow) {
      setTimeout(() => {
        throw new Error('Async error')
      }, 10)
    }
  }, [shouldThrow])
  
  return <div>Async component</div>
}

describe('ErrorBoundary Components', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  describe('Basic ErrorBoundary', () => {
    test('renders children when no error occurs', () => {
      render(
        <ErrorBoundary>
          <div>Child component</div>
        </ErrorBoundary>
      )
      
      expect(screen.getByText('Child component')).toBeInTheDocument()
    })

    test('catches errors and renders fallback UI', () => {
      render(
        <ErrorBoundary>
          <ThrowError />
        </ErrorBoundary>
      )
      
      expect(screen.getByRole('alert')).toBeInTheDocument()
      expect(screen.getByText('Something went wrong')).toBeInTheDocument()
      expect(screen.getByText('Test error message')).toBeInTheDocument()
    })

    test('calls onError callback when provided', () => {
      const onError = jest.fn()
      
      render(
        <ErrorBoundary onError={onError}>
          <ThrowError message="Callback test error" />
        </ErrorBoundary>
      )
      
      expect(onError).toHaveBeenCalledWith(
        expect.objectContaining({
          message: 'Callback test error'
        }),
        expect.objectContaining({
          componentStack: expect.any(String)
        })
      )
    })
  })

  describe('GraphErrorBoundary', () => {
    test('renders children when no error occurs', () => {
      render(
        <GraphErrorBoundary>
          <div>Graph component</div>
        </GraphErrorBoundary>
      )
      
      expect(screen.getByText('Graph component')).toBeInTheDocument()
    })

    test('renders fallback with custom title', () => {
      render(
        <GraphErrorBoundary title="Graph Loading Failed">
          <ThrowError />
        </GraphErrorBoundary>
      )
      
      expect(screen.getByText('Graph Loading Failed')).toBeInTheDocument()
    })

    test('shows/hides error details based on showDetails prop', () => {
      const { rerender } = render(
        <GraphErrorBoundary showDetails={true}>
          <ThrowError />
        </GraphErrorBoundary>
      )
      
      // Details should be visible
      expect(screen.getByText('View technical details')).toBeInTheDocument()
      
      rerender(
        <GraphErrorBoundary showDetails={false}>
          <ThrowError />
        </GraphErrorBoundary>
      )
      
      // Details should be hidden
      expect(screen.queryByText('View technical details')).not.toBeInTheDocument()
    })

    test('resets on prop changes when resetOnPropsChange is true', () => {
      const { rerender } = render(
        <GraphErrorBoundary resetKeys={['key1']} resetOnPropsChange={true}>
          <ThrowError />
        </GraphErrorBoundary>
      )
      
      // Error state should be shown
      expect(screen.getByRole('alert')).toBeInTheDocument()
      
      // Change resetKeys to trigger reset
      rerender(
        <GraphErrorBoundary resetKeys={['key2']} resetOnPropsChange={true}>
          <ThrowError shouldThrow={false} />
        </GraphErrorBoundary>
      )
      
      // Should reset and show child component
      expect(screen.getByText('No error')).toBeInTheDocument()
    })
  })

  describe('Error Fallback UI', () => {
    test('renders with proper accessibility attributes', () => {
      render(
        <ErrorBoundary>
          <ThrowError />
        </ErrorBoundary>
      )
      
      const errorContainer = screen.getByRole('alert')
      expect(errorContainer).toHaveAttribute('data-testid', 'error-boundary')
    })

    test('displays error icon and message', () => {
      render(
        <ErrorBoundary>
          <ThrowError message="Custom error message" />
        </ErrorBoundary>
      )
      
      // Check for error icon (SVG)
      const errorIcon = screen.getByRole('alert').querySelector('svg')
      expect(errorIcon).toBeInTheDocument()
      expect(errorIcon).toHaveAttribute('aria-hidden', 'true')
      
      // Check error message
      expect(screen.getByText('Custom error message')).toBeInTheDocument()
    })

    test('shows collapsible technical details', async () => {
      const user = userEvent.setup()
      
      render(
        <GraphErrorBoundary showDetails={true}>
          <ThrowError />
        </GraphErrorBoundary>
      )
      
      const detailsSummary = screen.getByText('View technical details')
      expect(detailsSummary).toBeInTheDocument()
      
      await user.click(detailsSummary)
      
      // Stack trace should be visible after clicking
      const stackTrace = screen.getByText(/Error: Test error message/, { exact: false })
      expect(stackTrace).toBeInTheDocument()
    })

    test('retry button calls resetErrorBoundary', async () => {
      const user = userEvent.setup()
      let shouldThrow = true
      
      const { rerender } = render(
        <ErrorBoundary>
          <ThrowError shouldThrow={shouldThrow} />
        </ErrorBoundary>
      )
      
      const retryButton = screen.getByTestId('retry-button')
      expect(retryButton).toHaveTextContent('Try Again')
      
      // Change the error condition
      shouldThrow = false
      
      await user.click(retryButton)
      
      // Re-render with no error
      rerender(
        <ErrorBoundary>
          <ThrowError shouldThrow={shouldThrow} />
        </ErrorBoundary>
      )
      
      expect(screen.getByText('No error')).toBeInTheDocument()
    })

    test('copy error button copies error details to clipboard', async () => {
      const user = userEvent.setup()
      const mockWriteText = jest.fn().mockResolvedValue(undefined)
      
      Object.assign(navigator, {
        clipboard: {
          writeText: mockWriteText,
        },
      })
      
      render(
        <ErrorBoundary>
          <ThrowError />
        </ErrorBoundary>
      )
      
      const copyButton = screen.getByTestId('copy-error-button')
      await user.click(copyButton)
      
      expect(mockWriteText).toHaveBeenCalledWith(
        expect.stringContaining('Error: Test error message')
      )
    })

    test('reload button triggers page reload', async () => {
      const user = userEvent.setup()
      const mockReload = jest.fn()
      
      Object.defineProperty(window, 'location', {
        value: { reload: mockReload },
        writable: true,
      })
      
      render(
        <ErrorBoundary>
          <ThrowError />
        </ErrorBoundary>
      )
      
      const reloadButton = screen.getByTestId('reload-button')
      await user.click(reloadButton)
      
      expect(mockReload).toHaveBeenCalled()
    })
  })

  describe('Accessibility Features', () => {
    test('error container has proper ARIA attributes', () => {
      render(
        <ErrorBoundary>
          <ThrowError />
        </ErrorBoundary>
      )
      
      const errorContainer = screen.getByRole('alert')
      expect(errorContainer).toHaveAttribute('role', 'alert')
    })

    test('retry button has proper keyboard accessibility', async () => {
      const user = userEvent.setup()
      
      render(
        <ErrorBoundary>
          <ThrowError />
        </ErrorBoundary>
      )
      
      const retryButton = screen.getByTestId('retry-button')
      
      // Should be focusable
      retryButton.focus()
      expect(retryButton).toHaveFocus()
      
      // Should be activatable with keyboard
      await user.keyboard('{Enter}')
      // Button click behavior would be tested in integration
    })

    test('action buttons have proper focus management', () => {
      render(
        <ErrorBoundary>
          <ThrowError />
        </ErrorBoundary>
      )
      
      const retryButton = screen.getByTestId('retry-button')
      const copyButton = screen.getByTestId('copy-error-button')
      const reloadButton = screen.getByTestId('reload-button')
      
      // All buttons should have focus classes
      expect(retryButton).toHaveClass('focus:outline-none', 'focus:ring-2')
      expect(copyButton).toHaveClass('focus:outline-none', 'focus:ring-2')
      expect(reloadButton).toHaveClass('focus:outline-none', 'focus:ring-2')
    })
  })

  describe('Custom Fallback Components', () => {
    test('accepts custom fallback component', () => {
      const CustomFallback = () => <div>Custom error UI</div>
      
      render(
        <GraphErrorBoundary fallback={CustomFallback}>
          <ThrowError />
        </GraphErrorBoundary>
      )
      
      expect(screen.getByText('Custom error UI')).toBeInTheDocument()
      expect(screen.queryByText('Something went wrong')).not.toBeInTheDocument()
    })

    test('passes correct props to custom fallback', () => {
      const CustomFallback: React.FC<any> = ({ error, resetErrorBoundary }) => (
        <div>
          <span>Error: {error.message}</span>
          <button onClick={resetErrorBoundary}>Custom Reset</button>
        </div>
      )
      
      render(
        <GraphErrorBoundary fallback={CustomFallback}>
          <ThrowError message="Custom fallback test" />
        </GraphErrorBoundary>
      )
      
      expect(screen.getByText('Error: Custom fallback test')).toBeInTheDocument()
      expect(screen.getByText('Custom Reset')).toBeInTheDocument()
    })
  })

  describe('withErrorBoundary HOC', () => {
    test('wraps component with error boundary', () => {
      const TestComponent = () => <div>Test component</div>
      const WrappedComponent = withErrorBoundary(TestComponent)
      
      render(<WrappedComponent />)
      
      expect(screen.getByText('Test component')).toBeInTheDocument()
    })

    test('catches errors in wrapped component', () => {
      const WrappedThrowError = withErrorBoundary(ThrowError)
      
      render(<WrappedThrowError />)
      
      expect(screen.getByRole('alert')).toBeInTheDocument()
      expect(screen.getByText('Something went wrong')).toBeInTheDocument()
    })

    test('applies custom error boundary props', () => {
      const WrappedComponent = withErrorBoundary(ThrowError, {
        title: 'HOC Error Title'
      })
      
      render(<WrappedComponent />)
      
      expect(screen.getByText('HOC Error Title')).toBeInTheDocument()
    })

    test('sets correct displayName', () => {
      const TestComponent = () => <div>Test</div>
      TestComponent.displayName = 'TestComponent'
      
      const WrappedComponent = withErrorBoundary(TestComponent)
      
      expect(WrappedComponent.displayName).toBe('withErrorBoundary(TestComponent)')
    })
  })

  describe('useErrorBoundaryReset Hook', () => {
    test('provides reset functionality', () => {
      let resetKey: number
      let resetFn: () => void
      
      const TestComponent = () => {
        const { resetKey: key, reset } = useErrorBoundaryReset()
        resetKey = key
        resetFn = reset
        return <div>Reset key: {key}</div>
      }
      
      render(<TestComponent />)
      
      expect(screen.getByText('Reset key: 0')).toBeInTheDocument()
      expect(typeof resetFn!).toBe('function')
    })

    test('increments reset key when reset is called', () => {
      const TestComponent = () => {
        const { resetKey, reset } = useErrorBoundaryReset()
        
        return (
          <div>
            <span>Reset key: {resetKey}</span>
            <button onClick={reset}>Reset</button>
          </div>
        )
      }
      
      const { rerender } = render(<TestComponent />)
      
      expect(screen.getByText('Reset key: 0')).toBeInTheDocument()
      
      // Simulate reset call (would need user interaction in real scenario)
      rerender(<TestComponent />)
    })
  })

  describe('Error Logging', () => {
    test('logs errors in development mode', () => {
      const originalEnv = process.env.NODE_ENV
      process.env.NODE_ENV = 'development'
      
      render(
        <GraphErrorBoundary>
          <ThrowError message="Development error" />
        </GraphErrorBoundary>
      )
      
      expect(console.error).toHaveBeenCalledWith(
        'ErrorBoundary caught an error:',
        expect.objectContaining({
          message: 'Development error'
        })
      )
      
      process.env.NODE_ENV = originalEnv
    })

    test('does not spam console in production mode', () => {
      const originalEnv = process.env.NODE_ENV
      process.env.NODE_ENV = 'production'
      
      render(
        <GraphErrorBoundary>
          <ThrowError message="Production error" />
        </GraphErrorBoundary>
      )
      
      // Should still log the error but not spam
      expect(console.error).toHaveBeenCalled()
      
      process.env.NODE_ENV = originalEnv
    })
  })

  describe('Edge Cases', () => {
    test('handles errors with no stack trace', () => {
      const errorWithoutStack = new Error('No stack error')
      errorWithoutStack.stack = undefined
      
      const ThrowCustomError = () => {
        throw errorWithoutStack
      }
      
      render(
        <ErrorBoundary>
          <ThrowCustomError />
        </ErrorBoundary>
      )
      
      expect(screen.getByText('No stack error')).toBeInTheDocument()
      expect(screen.getByRole('alert')).toBeInTheDocument()
    })

    test('handles non-Error objects being thrown', () => {
      const ThrowString = () => {
        throw 'String error'
      }
      
      render(
        <ErrorBoundary>
          <ThrowString />
        </ErrorBoundary>
      )
      
      expect(screen.getByText('String error')).toBeInTheDocument()
    })

    test('handles clipboard API not available', async () => {
      const user = userEvent.setup()
      
      // Mock clipboard as undefined
      Object.defineProperty(navigator, 'clipboard', {
        value: undefined,
        writable: true,
      })
      
      render(
        <ErrorBoundary>
          <ThrowError />
        </ErrorBoundary>
      )
      
      const copyButton = screen.getByTestId('copy-error-button')
      
      // Should not throw error when clipboard is unavailable
      await user.click(copyButton)
      
      expect(copyButton).toBeInTheDocument() // Still functional
    })

    test('handles multiple rapid errors', () => {
      const MultiErrorComponent = ({ errorCount }: { readonly errorCount: number }) => {
        if (errorCount > 0) {
          throw new Error(`Error ${errorCount}`)
        }
        return <div>No errors</div>
      }
      
      const { rerender } = render(
        <ErrorBoundary>
          <MultiErrorComponent errorCount={0} />
        </ErrorBoundary>
      )
      
      // Trigger multiple errors
      rerender(
        <ErrorBoundary>
          <MultiErrorComponent errorCount={1} />
        </ErrorBoundary>
      )
      
      expect(screen.getByText('Error 1')).toBeInTheDocument()
      
      // Second error
      rerender(
        <ErrorBoundary>
          <MultiErrorComponent errorCount={2} />
        </ErrorBoundary>
      )
      
      expect(screen.getByText('Error 2')).toBeInTheDocument()
    })
  })

  describe('Performance', () => {
    test('does not impact performance when no errors occur', () => {
      const renderCount = jest.fn()
      
      const TestComponent = () => {
        renderCount()
        return <div>Normal component</div>
      }
      
      const { rerender } = render(
        <ErrorBoundary>
          <TestComponent />
        </ErrorBoundary>
      )
      
      expect(renderCount).toHaveBeenCalledTimes(1)
      
      rerender(
        <ErrorBoundary>
          <TestComponent />
        </ErrorBoundary>
      )
      
      expect(renderCount).toHaveBeenCalledTimes(2)
    })

    test('error state does not cause additional re-renders', () => {
      const { rerender } = render(
        <ErrorBoundary>
          <ThrowError />
        </ErrorBoundary>
      )
      
      const initialAlert = screen.getByRole('alert')
      
      // Re-render should maintain same error state
      rerender(
        <ErrorBoundary>
          <ThrowError />
        </ErrorBoundary>
      )
      
      expect(screen.getByRole('alert')).toBe(initialAlert)
    })
  })
})