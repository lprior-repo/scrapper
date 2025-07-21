/**
 * GraphCanvas Component with React 19 Suspense and Error Boundaries
 *
 * This component implements modern React patterns:
 * - React 19 Suspense for data fetching
 * - Error boundaries for error handling
 * - Promise caching to prevent re-fetching
 * - Proper loading states and error recovery
 */

import React, { Suspense } from 'react'

import { CytoscapeGraphComponent } from './CytoscapeGraphComponent'
import { GraphErrorBoundary } from './ErrorBoundary'
import { GraphLoadingSpinner } from './LoadingSpinner'
import { useGraphData } from '../hooks/useGraphData'

interface GraphCanvasProps {
  readonly organization: string
  readonly useTopics: boolean
}

/**
 * GraphRenderer - The component that actually uses the data
 * This component will suspend when data is not available
 */
const GraphRenderer: React.FC<GraphCanvasProps> = ({
  organization,
  useTopics,
}) => {
  // This will suspend the component if data is not available
  const { nodes, edges } = useGraphData(organization, useTopics)

  // Handle empty data case
  if (!nodes || !edges || (nodes.length === 0 && edges.length === 0)) {
    return (
      <div
        className="flex min-h-screen items-center justify-center bg-dark-bg px-4"
        data-testid="empty-graph-state"
      >
        <div className="text-center">
          <div className="mb-4">
            <svg
              className="mx-auto h-12 w-12 text-dark-text-secondary"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M9 9l3-3 3 3"
              />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-dark-text mb-2">
            No Graph Data Available
          </h3>
          <p className="text-sm text-dark-text-secondary mb-4">
            No codeowners or repository data found for{' '}
            <span className="font-medium text-accent-blue">{organization}</span>
            {useTopics && ' with topics enabled'}.
          </p>
          <p className="text-xs text-dark-text-secondary">
            Make sure the organization exists and has public repositories with
            CODEOWNERS files.
          </p>
        </div>
      </div>
    )
  }

  console.log('Rendering graph with data:', {
    organization,
    useTopics,
    nodesCount: nodes.length,
    edgesCount: edges.length,
  })

  return (
    <div data-testid="graph-canvas" className="w-full h-screen bg-dark-bg">
      <CytoscapeGraphComponent nodes={nodes} edges={edges} />
    </div>
  )
}

/**
 * Main GraphCanvas component with Suspense and Error Boundary
 */
export const GraphCanvas: React.FC<GraphCanvasProps> = ({
  organization,
  useTopics,
}) => {
  // Early validation to provide better error messages
  if (
    !organization ||
    typeof organization !== 'string' ||
    organization.trim().length === 0
  ) {
    return (
      <div
        className="flex min-h-screen items-center justify-center bg-dark-bg px-4"
        data-testid="graph-canvas-error"
      >
        <div className="text-center">
          <div className="mb-4">
            <svg
              className="mx-auto h-12 w-12 text-accent-red"
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
          <h3 className="text-lg font-semibold text-accent-red mb-2">
            Invalid Configuration
          </h3>
          <p className="text-sm text-dark-text-secondary">
            Organization parameter is required and must be a non-empty string.
          </p>
        </div>
      </div>
    )
  }

  return (
    <GraphErrorBoundary
      title="Graph Visualization Error"
      showDetails={process.env.NODE_ENV === 'development'}
      resetKeys={[organization, useTopics]} // Reset error boundary when these props change
      onError={(error, errorInfo) => {
        console.error('GraphCanvas Error Boundary caught an error:', {
          error,
          errorInfo,
          organization,
          useTopics,
        })

        // TODO: Send to error monitoring service
        // Example: Sentry.captureException(error, { contexts: { react: errorInfo } })
      }}
    >
      <Suspense
        fallback={
          <GraphLoadingSpinner stage="fetching" organization={organization} />
        }
      >
        <GraphRenderer organization={organization} useTopics={useTopics} />
      </Suspense>
    </GraphErrorBoundary>
  )
}
