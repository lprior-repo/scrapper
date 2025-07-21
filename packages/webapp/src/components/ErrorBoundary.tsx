/**
 * Error Boundary Component using react-error-boundary
 *
 * Provides comprehensive error handling with retry functionality
 * and user-friendly error messages styled with GitHub dark theme
 */

import React from 'react'
import {
  ErrorBoundary as ReactErrorBoundary,
  FallbackProps,
} from 'react-error-boundary'

interface ErrorFallbackProps extends FallbackProps {
  readonly title?: string
  readonly showDetails?: boolean
}

/**
 * Error fallback component with GitHub dark theme styling
 */
const ErrorFallback: React.FC<ErrorFallbackProps> = ({
  error,
  resetErrorBoundary,
  title = 'Something went wrong',
  showDetails = true,
}) => {
  const errorMessage = error instanceof Error ? error.message : String(error)
  const errorStack = error instanceof Error ? error.stack : undefined

  const handleCopyError = () => {
    const errorInfo = `Error: ${errorMessage}\n\nStack: ${errorStack || 'No stack trace available'}\n\nTimestamp: ${new Date().toISOString()}`
    navigator.clipboard.writeText(errorInfo).catch(() => {
      // Silently fail if clipboard is not available
    })
  }

  return (
    <div
      role="alert"
      className="flex min-h-screen items-center justify-center bg-dark-bg px-4 py-8"
      data-testid="error-boundary"
    >
      <div className="w-full max-w-md rounded-lg border border-dark-border bg-dark-bg p-6 shadow-lg">
        {/* Error Icon */}
        <div className="mb-4 flex justify-center">
          <div className="rounded-full bg-accent-red/10 p-3">
            <svg
              className="h-8 w-8 text-accent-red"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
        </div>

        {/* Error Title */}
        <h1 className="mb-3 text-center text-xl font-semibold text-dark-text">
          {title}
        </h1>

        {/* Error Message */}
        <div className="mb-4 rounded-md bg-accent-red/10 p-3">
          <p className="text-sm text-accent-red font-medium">{errorMessage}</p>
        </div>

        {/* Error Details (collapsible) */}
        {showDetails && errorStack && (
          <details className="mb-4">
            <summary className="mb-2 cursor-pointer text-sm text-dark-text-secondary hover:text-dark-text">
              View technical details
            </summary>
            <div className="rounded-md bg-dark-bg border border-dark-border p-3">
              <pre className="text-xs text-dark-text-secondary overflow-auto whitespace-pre-wrap break-words">
                {errorStack}
              </pre>
            </div>
          </details>
        )}

        {/* Action Buttons */}
        <div className="space-y-2">
          {/* Retry Button */}
          <button
            onClick={resetErrorBoundary}
            className="w-full rounded-md bg-accent-blue px-4 py-2 text-sm font-medium text-white hover:bg-accent-blue/90 focus:outline-none focus:ring-2 focus:ring-accent-blue focus:ring-offset-2 focus:ring-offset-dark-bg transition-colors"
            data-testid="retry-button"
          >
            Try Again
          </button>

          {/* Copy Error Button */}
          <button
            onClick={handleCopyError}
            className="w-full rounded-md border border-dark-border bg-transparent px-4 py-2 text-sm font-medium text-dark-text hover:bg-dark-border/20 focus:outline-none focus:ring-2 focus:ring-dark-border focus:ring-offset-2 focus:ring-offset-dark-bg transition-colors"
            data-testid="copy-error-button"
          >
            Copy Error Details
          </button>

          {/* Reload Page Button */}
          <button
            onClick={() => window.location.reload()}
            className="w-full rounded-md border border-dark-border bg-transparent px-4 py-2 text-sm font-medium text-dark-text-secondary hover:text-dark-text hover:bg-dark-border/20 focus:outline-none focus:ring-2 focus:ring-dark-border focus:ring-offset-2 focus:ring-offset-dark-bg transition-colors"
            data-testid="reload-button"
          >
            Reload Page
          </button>
        </div>

        {/* Help Text */}
        <p className="mt-4 text-center text-xs text-dark-text-secondary">
          If this problem persists, please contact support with the error
          details above.
        </p>
      </div>
    </div>
  )
}

interface GraphErrorBoundaryProps {
  readonly children: React.ReactNode
  readonly title?: string
  readonly showDetails?: boolean
  readonly onError?: (
    error: Error,
    errorInfo: { readonly componentStack: string }
  ) => void
  readonly resetKeys?: readonly unknown[]
  readonly resetOnPropsChange?: boolean
  readonly fallback?: React.ComponentType<FallbackProps>
}

/**
 * Graph-specific Error Boundary with customized error handling
 */
export const GraphErrorBoundary: React.FC<GraphErrorBoundaryProps> = ({
  children,
  title,
  showDetails = true,
  onError,
  resetKeys,
  resetOnPropsChange = true,
  fallback,
}) => {
  const handleError = (error: Error, errorInfo: { readonly componentStack: string }) => {
    // Log error for debugging (in development)
    if (process.env.NODE_ENV === 'development') {
      console.error('ErrorBoundary caught an error:', error)
      console.error('Component stack:', errorInfo.componentStack)
    }

    // Call custom error handler if provided
    onError?.(error, errorInfo)

    // TODO: Send to error reporting service in production
    // Example: Sentry.captureException(error, { contexts: { react: errorInfo } })
  }

  const FallbackComponent =
    fallback ||
    ((props: FallbackProps) => (
      <ErrorFallback {...props} title={title} showDetails={showDetails} />
    ))

  return (
    <ReactErrorBoundary
      FallbackComponent={FallbackComponent}
      onError={handleError}
      resetKeys={resetKeys}
      resetOnPropsChange={resetOnPropsChange}
    >
      {children}
    </ReactErrorBoundary>
  )
}

/**
 * Simple Error Boundary for general use cases
 */
export const ErrorBoundary: React.FC<{
  readonly children: React.ReactNode
  readonly fallback?: React.ComponentType<FallbackProps>
  readonly onError?: (
    error: Error,
    errorInfo: { readonly componentStack: string }
  ) => void
}> = ({ children, fallback, onError }) => {
  return (
    <GraphErrorBoundary onError={onError} fallback={fallback}>
      {children}
    </GraphErrorBoundary>
  )
}

/**
 * Hook to manually trigger error boundary reset
 */
export const useErrorBoundaryReset = () => {
  const [resetKey, setResetKey] = React.useState(0)

  const reset = React.useCallback(() => {
    setResetKey((prev) => prev + 1)
  }, [])

  return { resetKey, reset }
}

/**
 * Higher-order component that wraps a component with error boundary
 */
export const withErrorBoundary = <P extends object>(
  Component: React.ComponentType<P>,
  errorBoundaryProps?: Omit<GraphErrorBoundaryProps, 'children'>
) => {
  const WrappedComponent: React.FC<P> = (props) => (
    <GraphErrorBoundary {...errorBoundaryProps}>
      <Component {...props} />
    </GraphErrorBoundary>
  )

  WrappedComponent.displayName = `withErrorBoundary(${Component.displayName || Component.name})`

  return WrappedComponent
}

export default ErrorBoundary
