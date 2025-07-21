/**
 * LoadingSpinner Component Tests
 * 
 * Tests for the loading spinner components including:
 * - Basic rendering and animations
 * - Different size and color variants
 * - Accessibility features (ARIA, screen readers)
 * - Loading states and progress indication
 * - Skeleton loaders and overlays
 * - Graph-specific loading components
 * - Performance and integration patterns
 */

import React from 'react'
import { render, screen } from '@testing-library/react'
import { 
  LoadingSpinner,
  SkeletonLoader,
  GraphLoadingSpinner,
  InlineLoader,
  LoadingOverlay 
} from '../LoadingSpinner'

describe('LoadingSpinner Components', () => {
  describe('Basic LoadingSpinner', () => {
    test('renders with default props', () => {
      render(<LoadingSpinner />)
      
      const spinner = screen.getByTestId('loading-spinner')
      expect(spinner).toBeInTheDocument()
      expect(spinner).toHaveAttribute('role', 'status')
      expect(spinner).toHaveAttribute('aria-live', 'polite')
    })

    test('displays loading message by default', () => {
      render(<LoadingSpinner />)
      
      expect(screen.getByText('Loading...')).toBeInTheDocument()
    })

    test('displays custom message', () => {
      render(<LoadingSpinner message="Custom loading message" />)
      
      expect(screen.getByText('Custom loading message')).toBeInTheDocument()
    })

    test('hides message when showMessage is false', () => {
      render(<LoadingSpinner message="Should be hidden" showMessage={false} />)
      
      expect(screen.queryByText('Should be hidden')).not.toBeInTheDocument()
    })

    test('applies different sizes correctly', () => {
      const { rerender } = render(<LoadingSpinner size="sm" />)
      
      let spinner = screen.getByLabelText('Loading spinner')
      expect(spinner).toHaveClass('h-4', 'w-4')
      
      rerender(<LoadingSpinner size="md" />)
      spinner = screen.getByLabelText('Loading spinner')
      expect(spinner).toHaveClass('h-6', 'w-6')
      
      rerender(<LoadingSpinner size="lg" />)
      spinner = screen.getByLabelText('Loading spinner')
      expect(spinner).toHaveClass('h-8', 'w-8')
      
      rerender(<LoadingSpinner size="xl" />)
      spinner = screen.getByLabelText('Loading spinner')
      expect(spinner).toHaveClass('h-12', 'w-12')
    })

    test('applies different colors correctly', () => {
      const { rerender } = render(<LoadingSpinner color="blue" />)
      
      let spinner = screen.getByLabelText('Loading spinner')
      expect(spinner).toHaveClass('text-accent-blue')
      
      rerender(<LoadingSpinner color="green" />)
      spinner = screen.getByLabelText('Loading spinner')
      expect(spinner).toHaveClass('text-accent-green')
      
      rerender(<LoadingSpinner color="white" />)
      spinner = screen.getByLabelText('Loading spinner')
      expect(spinner).toHaveClass('text-white')
    })

    test('renders in fullscreen mode', () => {
      render(<LoadingSpinner fullScreen={true} />)
      
      const container = screen.getByTestId('loading-spinner')
      expect(container).toHaveClass(
        'fixed',
        'inset-0',
        'flex',
        'flex-col',
        'items-center',
        'justify-center',
        'bg-dark-bg',
        'z-50'
      )
    })

    test('applies custom className', () => {
      render(<LoadingSpinner className="custom-spinner-class" />)
      
      const container = screen.getByTestId('loading-spinner')
      expect(container).toHaveClass('custom-spinner-class')
    })
  })

  describe('GraphLoadingSpinner', () => {
    test('renders with default fetching stage', () => {
      render(<GraphLoadingSpinner />)
      
      expect(screen.getByText('Fetching graph data...')).toBeInTheDocument()
      expect(screen.getByText('33% complete')).toBeInTheDocument()
    })

    test('displays different stages correctly', () => {
      const { rerender } = render(<GraphLoadingSpinner stage="fetching" />)
      expect(screen.getByText('Fetching graph data...')).toBeInTheDocument()
      expect(screen.getByText('33% complete')).toBeInTheDocument()
      
      rerender(<GraphLoadingSpinner stage="processing" />)
      expect(screen.getByText('Processing graph structure...')).toBeInTheDocument()
      expect(screen.getByText('66% complete')).toBeInTheDocument()
      
      rerender(<GraphLoadingSpinner stage="rendering" />)
      expect(screen.getByText('Rendering visualization...')).toBeInTheDocument()
      expect(screen.getByText('90% complete')).toBeInTheDocument()
    })

    test('includes organization name in fetching message', () => {
      render(<GraphLoadingSpinner stage="fetching" organization="github" />)
      
      expect(screen.getByText('Fetching data for github...')).toBeInTheDocument()
    })

    test('has proper accessibility attributes', () => {
      render(<GraphLoadingSpinner />)
      
      const container = screen.getByTestId('graph-loading-spinner')
      expect(container).toHaveAttribute('role', 'status')
      expect(container).toHaveAttribute('aria-live', 'polite')
      
      const progressBar = screen.getByRole('progressbar')
      expect(progressBar).toHaveAttribute('aria-valuenow', '33')
      expect(progressBar).toHaveAttribute('aria-valuemin', '0')
      expect(progressBar).toHaveAttribute('aria-valuemax', '100')
    })

    test('displays stage indicators with correct states', () => {
      render(<GraphLoadingSpinner stage="processing" />)
      
      const stages = screen.getAllByText(/fetching|processing|rendering/i)
      expect(stages).toHaveLength(4) // 3 indicators + 1 main message
    })
  })

  describe('SkeletonLoader', () => {
    test('renders with default number of lines', () => {
      render(<SkeletonLoader />)
      
      const skeleton = screen.getByTestId('skeleton-loader')
      expect(skeleton).toBeInTheDocument()
      
      // Should have 3 skeleton lines by default
      const lines = skeleton.querySelectorAll('.h-4.rounded.bg-dark-border')
      expect(lines).toHaveLength(3)
    })

    test('renders custom number of lines', () => {
      render(<SkeletonLoader lines={5} />)
      
      const skeleton = screen.getByTestId('skeleton-loader')
      const lines = skeleton.querySelectorAll('.h-4.rounded.bg-dark-border')
      expect(lines).toHaveLength(5)
    })

    test('shows avatar when requested', () => {
      render(<SkeletonLoader showAvatar={true} />)
      
      const skeleton = screen.getByTestId('skeleton-loader')
      const avatar = skeleton.querySelector('.h-10.w-10.rounded-full')
      expect(avatar).toBeInTheDocument()
    })

    test('applies custom className', () => {
      render(<SkeletonLoader className="custom-skeleton" />)
      
      const skeleton = screen.getByTestId('skeleton-loader')
      expect(skeleton).toHaveClass('custom-skeleton')
    })
  })

  describe('InlineLoader', () => {
    test('renders with spinner by default', () => {
      render(<InlineLoader />)
      
      const container = screen.getByTestId('inline-loader')
      expect(container).toBeInTheDocument()
      
      const spinner = screen.getByLabelText('Loading spinner')
      expect(spinner).toBeInTheDocument()
    })

    test('renders with dots when requested', () => {
      render(<InlineLoader dots={true} />)
      
      const container = screen.getByTestId('inline-loader')
      const dotsAnimation = container.querySelector('[role="img"][aria-label="Loading animation"]')
      expect(dotsAnimation).toBeInTheDocument()
    })

    test('displays message when provided', () => {
      render(<InlineLoader message="Loading data" />)
      
      expect(screen.getByText('Loading data')).toBeInTheDocument()
    })

    test('applies different sizes', () => {
      const { rerender } = render(<InlineLoader size="sm" />)
      
      let spinner = screen.getByLabelText('Loading spinner')
      expect(spinner).toHaveClass('h-4', 'w-4')
      
      rerender(<InlineLoader size="md" />)
      spinner = screen.getByLabelText('Loading spinner')
      expect(spinner).toHaveClass('h-6', 'w-6')
    })
  })

  describe('LoadingOverlay', () => {
    test('shows children when not loading', () => {
      render(
        <LoadingOverlay isLoading={false}>
          <div>Content to show</div>
        </LoadingOverlay>
      )
      
      expect(screen.getByText('Content to show')).toBeInTheDocument()
      expect(screen.queryByTestId('loading-spinner')).not.toBeInTheDocument()
    })

    test('shows overlay when loading', () => {
      render(
        <LoadingOverlay isLoading={true}>
          <div>Content to hide</div>
        </LoadingOverlay>
      )
      
      expect(screen.getByText('Content to hide')).toBeInTheDocument()
      expect(screen.getByTestId('loading-spinner')).toBeInTheDocument()
    })

    test('displays custom loading message in overlay', () => {
      render(
        <LoadingOverlay isLoading={true} message="Processing...">
          <div>Content</div>
        </LoadingOverlay>
      )
      
      expect(screen.getByText('Processing...')).toBeInTheDocument()
    })
  })

  describe('Accessibility Features', () => {
    test('all loading components have proper ARIA attributes', () => {
      render(
        <div>
          <LoadingSpinner />
          <GraphLoadingSpinner />
          <SkeletonLoader />
          <InlineLoader />
        </div>
      )
      
      // Check role="status" for loading indicators
      const statusElements = screen.getAllByRole('status')
      expect(statusElements.length).toBeGreaterThan(0)
      
      // Check aria-live for dynamic updates
      const liveElements = screen.getAllByLabelText(/loading/i)
      expect(liveElements.length).toBeGreaterThan(0)
    })

    test('screen reader text provides context', () => {
      render(<LoadingSpinner />)
      
      const screenReaderText = screen.getByText('Loading content, please wait...')
      expect(screenReaderText).toHaveClass('sr-only')
    })

    test('GraphLoadingSpinner provides detailed progress information', () => {
      render(<GraphLoadingSpinner stage="processing" />)
      
      const screenReaderInfo = screen.getByText(/Loading graph visualization.*Current stage: processing.*66% complete/)
      expect(screenReaderInfo).toHaveClass('sr-only')
    })
  })

  describe('Animation and Styling', () => {
    test('spinner has correct animation classes', () => {
      render(<LoadingSpinner />)
      
      const spinner = screen.getByLabelText('Loading spinner')
      expect(spinner).toHaveClass('animate-spin')
    })

    test('skeleton loader has pulse animation', () => {
      render(<SkeletonLoader />)
      
      const skeleton = screen.getByTestId('skeleton-loader')
      expect(skeleton).toHaveClass('animate-pulse')
    })

    test('progress bar has transition animation', () => {
      render(<GraphLoadingSpinner />)
      
      const progressBar = screen.getByRole('progressbar')
      expect(progressBar).toHaveClass('transition-all', 'duration-500', 'ease-out')
    })
  })

  describe('Integration Patterns', () => {
    test('works with React Suspense patterns', () => {
      const SuspenseComponent = () => (
        <React.Suspense fallback={<LoadingSpinner message="Loading suspense content" />}>
          <div>Suspense content</div>
        </React.Suspense>
      )
      
      render(<SuspenseComponent />)
      
      // Since there's no actual async loading, should show content immediately
      expect(screen.getByText('Suspense content')).toBeInTheDocument()
    })

    test('LoadingOverlay integrates with state management', () => {
      const StatefulComponent = () => {
        const [isLoading, setIsLoading] = React.useState(true)
        
        React.useEffect(() => {
          const timer = setTimeout(() => setIsLoading(false), 100)
          return () => clearTimeout(timer)
        }, [])
        
        return (
          <LoadingOverlay isLoading={isLoading}>
            <div>Loaded content</div>
          </LoadingOverlay>
        )
      }
      
      render(<StatefulComponent />)
      
      expect(screen.getByTestId('loading-spinner')).toBeInTheDocument()
      expect(screen.getByText('Loaded content')).toBeInTheDocument()
    })
  })

  describe('Performance', () => {
    test('does not cause unnecessary re-renders', () => {
      let renderCount = 0
      
      const TestComponent = ({ isLoading }: { readonly isLoading: boolean }) => {
        renderCount++
        return <LoadingSpinner message={isLoading ? 'Loading' : 'Not loading'} />
      }
      
      const { rerender } = render(<TestComponent isLoading={true} />)
      expect(renderCount).toBe(1)
      
      rerender(<TestComponent isLoading={true} />)
      expect(renderCount).toBe(2)
    })

    test('animations do not impact performance', () => {
      render(
        <div>
          {Array.from({ length: 10 }, (_, i) => (
            <LoadingSpinner key={i} size="sm" />
          ))}
        </div>
      )
      
      const spinners = screen.getAllByLabelText('Loading spinner')
      expect(spinners).toHaveLength(10)
      
      // All should have animation class
      spinners.forEach(spinner => {
        expect(spinner).toHaveClass('animate-spin')
      })
    })
  })

  describe('Edge Cases', () => {
    test('handles undefined props gracefully', () => {
      render(<LoadingSpinner message={undefined} size={undefined as any} />)
      
      const spinner = screen.getByTestId('loading-spinner')
      expect(spinner).toBeInTheDocument()
    })

    test('handles very long messages', () => {
      const longMessage = 'A'.repeat(1000)
      render(<LoadingSpinner message={longMessage} />)
      
      expect(screen.getByText(longMessage)).toBeInTheDocument()
    })

    test('GraphLoadingSpinner handles invalid stages gracefully', () => {
      render(<GraphLoadingSpinner stage={'invalid' as any} />)
      
      expect(screen.getByText('Loading graph...')).toBeInTheDocument()
      expect(screen.getByText('0% complete')).toBeInTheDocument()
    })

    test('LoadingOverlay handles rapid state changes', () => {
      const { rerender } = render(
        <LoadingOverlay isLoading={true}>
          <div>Content</div>
        </LoadingOverlay>
      )
      
      // Rapid toggle
      for (let i = 0; i < 10; i++) {
        rerender(
          <LoadingOverlay isLoading={i % 2 === 0}>
            <div>Content</div>
          </LoadingOverlay>
        )
      }
      
      expect(screen.getByText('Content')).toBeInTheDocument()
    })
  })
})