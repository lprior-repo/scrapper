/**
 * Loading Spinner Component for React 19 Suspense
 *
 * Accessible loading spinner with GitHub dark theme styling
 * Supports different sizes and loading states
 */

import React from 'react'

interface LoadingSpinnerProps {
  readonly size?: 'sm' | 'md' | 'lg' | 'xl'
  readonly message?: string
  readonly showMessage?: boolean
  readonly className?: string
  readonly fullScreen?: boolean
  readonly color?: 'blue' | 'green' | 'white'
}

const sizeClasses = {
  sm: 'h-4 w-4',
  md: 'h-6 w-6',
  lg: 'h-8 w-8',
  xl: 'h-12 w-12',
} as const

const colorClasses = {
  blue: 'text-accent-blue',
  green: 'text-accent-green',
  white: 'text-white',
} as const

/**
 * Animated spinner SVG component
 */
const SpinnerIcon: React.FC<{
  readonly size: LoadingSpinnerProps['size']
  readonly color: LoadingSpinnerProps['color']
  readonly className?: string
}> = ({ size = 'md', color = 'blue', className = '' }) => (
  <svg
    className={`animate-spin ${sizeClasses[size]} ${colorClasses[color]} ${className}`}
    fill="none"
    viewBox="0 0 24 24"
    role="img"
    aria-label="Loading spinner"
  >
    <circle
      className="opacity-25"
      cx="12"
      cy="12"
      r="10"
      stroke="currentColor"
      strokeWidth="4"
    />
    <path
      className="opacity-75"
      fill="currentColor"
      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
    />
  </svg>
)

/**
 * Pulsing dots animation as an alternative to spinner
 */
const PulsingDots: React.FC<{
  readonly color: LoadingSpinnerProps['color']
  readonly className?: string
}> = ({ color = 'blue', className = '' }) => (
  <div
    className={`flex space-x-1 ${className}`}
    role="img"
    aria-label="Loading animation"
  >
    {[0, 1, 2].map((index) => (
      <div
        key={index}
        className={`h-2 w-2 rounded-full animate-pulse ${colorClasses[color]}`}
        style={{
          backgroundColor: 'currentColor',
          animationDelay: `${index * 0.15}s`,
          animationDuration: '1s',
        }}
      />
    ))}
  </div>
)

/**
 * Main LoadingSpinner component
 */
export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'md',
  message = 'Loading...',
  showMessage = true,
  className = '',
  fullScreen = false,
  color = 'blue',
}) => {
  const containerClasses = fullScreen
    ? 'fixed inset-0 flex flex-col items-center justify-center bg-dark-bg z-50'
    : 'flex flex-col items-center justify-center p-4'

  return (
    <div
      className={`${containerClasses} ${className}`}
      role="status"
      aria-live="polite"
      aria-label={message}
      data-testid="loading-spinner"
    >
      <SpinnerIcon size={size} color={color} />

      {showMessage && (
        <p
          className="mt-3 text-sm text-dark-text-secondary font-medium"
          aria-live="polite"
        >
          {message}
        </p>
      )}

      {/* Screen reader only text */}
      <span className="sr-only">Loading content, please wait...</span>
    </div>
  )
}

/**
 * Skeleton loading component for content placeholders
 */
export const SkeletonLoader: React.FC<{
  readonly className?: string
  readonly lines?: number
  readonly showAvatar?: boolean
}> = ({ className = '', lines = 3, showAvatar = false }) => (
  <div
    className={`animate-pulse ${className}`}
    role="status"
    aria-label="Loading content"
    data-testid="skeleton-loader"
  >
    <div className="flex items-start space-x-4">
      {showAvatar && (
        <div className="h-10 w-10 rounded-full bg-dark-border"></div>
      )}
      <div className="flex-1 space-y-2">
        {Array.from({ length: lines }, (_, index) => (
          <div
            key={index}
            className="h-4 rounded bg-dark-border"
            style={{
              width: index === lines - 1 ? '75%' : '100%',
            }}
          ></div>
        ))}
      </div>
    </div>
  </div>
)

/**
 * Graph-specific loading component with progress indication
 */
export const GraphLoadingSpinner: React.FC<{
  readonly stage?: 'fetching' | 'processing' | 'rendering'
  readonly organization?: string
}> = ({ stage = 'fetching', organization }) => {
  const getStageMessage = () => {
    switch (stage) {
      case 'fetching':
        return organization
          ? `Fetching data for ${organization}...`
          : 'Fetching graph data...'
      case 'processing':
        return 'Processing graph structure...'
      case 'rendering':
        return 'Rendering visualization...'
      default:
        return 'Loading graph...'
    }
  }

  const getProgress = () => {
    switch (stage) {
      case 'fetching':
        return 33
      case 'processing':
        return 66
      case 'rendering':
        return 90
      default:
        return 0
    }
  }

  return (
    <div
      className="flex min-h-screen flex-col items-center justify-center bg-dark-bg px-4"
      role="status"
      aria-live="polite"
      data-testid="graph-loading-spinner"
    >
      {/* Main spinner */}
      <SpinnerIcon size="xl" color="blue" className="mb-6" />

      {/* Loading message */}
      <h2 className="mb-2 text-lg font-semibold text-dark-text">
        {getStageMessage()}
      </h2>

      {/* Progress bar */}
      <div className="w-64 mb-4">
        <div className="h-2 bg-dark-border rounded-full overflow-hidden">
          <div
            className="h-full bg-accent-blue transition-all duration-500 ease-out"
            style={{ width: `${getProgress()}%` }}
            role="progressbar"
            aria-valuenow={getProgress()}
            aria-valuemin={0}
            aria-valuemax={100}
          />
        </div>
        <p className="text-xs text-dark-text-secondary mt-1 text-center">
          {getProgress()}% complete
        </p>
      </div>

      {/* Stage indicators */}
      <div className="flex space-x-6 text-xs">
        {(['fetching', 'processing', 'rendering'] as const).map(
          (stageItem, index) => {
            const isActive = stageItem === stage
            const isCompleted = getProgress() > (index + 1) * 33

            return (
              <div
                key={stageItem}
                className={`flex items-center space-x-1 ${
                  isActive
                    ? 'text-accent-blue'
                    : isCompleted
                      ? 'text-accent-green'
                      : 'text-dark-text-secondary'
                }`}
              >
                <div
                  className={`h-2 w-2 rounded-full ${
                    isActive
                      ? 'bg-accent-blue'
                      : isCompleted
                        ? 'bg-accent-green'
                        : 'bg-dark-border'
                  }`}
                />
                <span className="capitalize">{stageItem}</span>
              </div>
            )
          }
        )}
      </div>

      {/* Accessibility text */}
      <span className="sr-only">
        Loading graph visualization. Current stage: {stage}.{getProgress()}%
        complete.
      </span>
    </div>
  )
}

/**
 * Inline loading component for smaller UI elements
 */
export const InlineLoader: React.FC<{
  readonly size?: 'sm' | 'md'
  readonly message?: string
  readonly dots?: boolean
}> = ({ size = 'sm', message, dots = false }) => (
  <div
    className="flex items-center space-x-2"
    role="status"
    aria-live="polite"
    data-testid="inline-loader"
  >
    {dots ? (
      <PulsingDots color="blue" />
    ) : (
      <SpinnerIcon size={size} color="blue" />
    )}
    {message && (
      <span className="text-sm text-dark-text-secondary">{message}</span>
    )}
  </div>
)

/**
 * Loading overlay for existing content
 */
export const LoadingOverlay: React.FC<{
  readonly children: React.ReactNode
  readonly isLoading: boolean
  readonly message?: string
}> = ({ children, isLoading, message = 'Loading...' }) => (
  <div className="relative">
    {children}
    {isLoading && (
      <div className="absolute inset-0 flex items-center justify-center bg-dark-bg/75 backdrop-blur-sm">
        <LoadingSpinner message={message} color="white" />
      </div>
    )}
  </div>
)

export default LoadingSpinner
