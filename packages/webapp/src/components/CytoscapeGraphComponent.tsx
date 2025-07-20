import React, { useRef, useEffect, useState } from 'react'
import cytoscape, { type Core, type EventObject } from 'cytoscape'

import { createCytoscapeElements } from './graph-elements'
import { createCytoscapeStyles } from './graph-styles'
import type { GraphNode, GraphEdge } from '../services'

interface ICytoscapeGraphComponentProps {
  readonly nodes?: readonly GraphNode[]
  readonly edges?: readonly GraphEdge[]
}

interface CytoscapeError {
  readonly message: string
  readonly timestamp: number
}

const createCytoscapeConfig = (
  container: HTMLDivElement,
  elements: ReturnType<typeof createCytoscapeElements>
) => ({
  container,
  elements,
  style: createCytoscapeStyles(),
  layout: {
    name: 'cose',
    animate: true,
    animationDuration: 1000,
    nodeRepulsion: 400000,
    nodeOverlap: 10,
    idealEdgeLength: 100,
    edgeElasticity: 100,
    nestingFactor: 5,
    gravity: 80,
    numIter: 1000,
    randomize: false,
  },
  // wheelSensitivity: 0.2, // Removed to prevent zoom issues with different mice
  minZoom: 0.1,
  maxZoom: 3.0,
})

/**
 * Hover styles for nodes
 */
const nodeHoverStyles = {
  'border-width': 4,
  'border-color': '#58a6ff',
  'overlay-color': '#58a6ff',
  'overlay-opacity': 0.2,
} as const

/**
 * Default styles for nodes when not hovering
 */
const nodeDefaultStyles = {
  'border-width': 2,
  'overlay-opacity': 0,
} as const

/**
 * Hover styles for edges
 */
const edgeHoverStyles = {
  'line-color': '#58a6ff',
  'target-arrow-color': '#58a6ff',
  width: 3,
  'overlay-color': '#58a6ff',
  'overlay-opacity': 0.2,
} as const

/**
 * Default styles for edges when not hovering
 */
const edgeDefaultStyles = {
  'line-color': '#30363d',
  'target-arrow-color': '#30363d',
  width: 2,
  'overlay-opacity': 0,
} as const

/**
 * Handles node mouseover events
 */
const handleNodeMouseover = (event: EventObject): void => {
  event.target.style(nodeHoverStyles)
}

/**
 * Handles node mouseout events
 */
const handleNodeMouseout = (event: EventObject): void => {
  const node = event.target
  !node.selected() &&
    (node.style(nodeDefaultStyles), node.removestyle('border-color'))
}

/**
 * Handles edge mouseover events
 */
const handleEdgeMouseover = (event: EventObject): void => {
  event.target.style(edgeHoverStyles)
}

/**
 * Handles edge mouseout events
 */
const handleEdgeMouseout = (event: EventObject): void => {
  const edge = event.target
  !edge.selected() && edge.style(edgeDefaultStyles)
}

/**
 * Handles element tap events for selection
 */
const handleElementTap = (event: EventObject): void => {
  const element = event.target
  element.selected() ? element.unselect() : element.select()
}

/**
 * Handles background tap events to deselect all
 */
const handleBackgroundTap =
  (cy: Core) =>
  (event: EventObject): void => {
    event.target === cy && cy.elements().unselect()
  }

/**
 * Sets up hover effects for nodes and edges since CSS :hover doesn't work in Cytoscape
 */
const setupInteractiveEvents = (cy: Core): void => {
  cy.on('mouseover', 'node', handleNodeMouseover)
  cy.on('mouseout', 'node', handleNodeMouseout)
  cy.on('mouseover', 'edge', handleEdgeMouseover)
  cy.on('mouseout', 'edge', handleEdgeMouseout)
  cy.on('tap', 'node, edge', handleElementTap)
  cy.on('tap', handleBackgroundTap(cy))
}

const createCytoscapeInstance = (
  container: HTMLDivElement,
  elements: ReturnType<typeof createCytoscapeElements>,
  onError: (error: CytoscapeError) => void
): Core => {
  const config = createCytoscapeConfig(container, elements)
  const cyInstance = cytoscape(config)

  // Add error handling for cytoscape events
  cyInstance.on('error', (event) => {
    const error = {
      message: `Cytoscape error: ${event.message || 'Unknown error'}`,
      timestamp: Date.now(),
    }
    onError(error)
  })

  // Set up interactive events
  setupInteractiveEvents(cyInstance)

  console.warn(
    `Cytoscape instance created successfully with ${elements.length} elements`
  )
  return cyInstance
}

const initializeCytoscape = (
  containerRef: React.RefObject<HTMLDivElement>,
  elements: ReturnType<typeof createCytoscapeElements>,
  onError: (error: CytoscapeError) => void
): Core | null =>
  !containerRef.current
    ? (onError({
        message: 'Container ref is not available',
        timestamp: Date.now(),
      }),
      null)
    : !elements || elements.length === 0
      ? (console.warn('No elements to render in Cytoscape'), null)
      : (() => {
          try {
            return createCytoscapeInstance(
              containerRef.current,
              elements,
              onError
            )
          } catch (error) {
            const cytoscapeError = {
              message: `Failed to initialize Cytoscape: ${error instanceof Error ? error.message : String(error)}`,
              timestamp: Date.now(),
            }
            onError(cytoscapeError)
            return null
          }
        })()

/**
 * Error display styles
 */
const errorContainerStyle = {
  width: '100%',
  height: '100vh',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  backgroundColor: '#f8f9fa',
  border: '1px solid #dee2e6',
  borderRadius: '4px',
} as const

const errorContentStyle = {
  textAlign: 'center',
  padding: '20px',
} as const

const errorTitleStyle = {
  color: '#dc3545',
  marginBottom: '10px',
} as const

const errorMessageStyle = {
  color: '#6c757d',
  fontSize: '14px',
} as const

const errorTimeStyle = {
  color: '#6c757d',
  fontSize: '12px',
} as const

/**
 * Canvas container styles
 */
const canvasStyle = {
  width: '100%',
  height: '100vh',
} as const

/**
 * Renders error display component
 */
const renderErrorDisplay = (error: CytoscapeError): React.ReactElement => (
  <div data-testid="graph-error" style={errorContainerStyle}>
    <div style={errorContentStyle}>
      <h3 style={errorTitleStyle}>Graph Rendering Error</h3>
      <p style={errorMessageStyle}>{error.message}</p>
      <p style={errorTimeStyle}>
        Time: {new Date(error.timestamp).toLocaleTimeString()}
      </p>
    </div>
  </div>
)

/**
 * Renders canvas component
 */
const renderCanvas = (
  containerRef: React.RefObject<HTMLDivElement>
): React.ReactElement => (
  <div ref={containerRef} data-testid="graph-canvas" style={canvasStyle} />
)

/**
 * Destroys existing cytoscape instance
 */
const destroyInstance = (cyRef: React.MutableRefObject<Core | null>): void => {
  cyRef.current?.destroy()
}

/**
 * Creates cleanup function for cytoscape instance
 */
const createCleanup =
  (cyRef: React.MutableRefObject<Core | null>) => (): void => {
    destroyInstance(cyRef)
    // Note: This mutation is necessary for React ref pattern
    // eslint-disable-next-line functional/immutable-data
    cyRef.current = null
  }

/**
 * Initializes cytoscape with elements
 */
const initializeWithElements = (
  elements: readonly ElementDefinition[],
  containerRef: React.RefObject<HTMLDivElement>,
  cyRef: React.MutableRefObject<Core | null>,
  handleError: (error: CytoscapeError) => void
): (() => void) | undefined => {
  destroyInstance(cyRef)
  // Note: This mutation is necessary for React ref pattern
  // eslint-disable-next-line functional/immutable-data
  cyRef.current = initializeCytoscape(containerRef, elements, handleError)
  return createCleanup(cyRef)
}

export const CytoscapeGraphComponent: React.FC<
  ICytoscapeGraphComponentProps
> = ({ nodes, edges }) => {
  const containerRef = useRef<HTMLDivElement>(null)
  const cyRef = useRef<Core | null>(null)
  const [error, setError] = useState<CytoscapeError | null>(null)

  const handleError = (cytoscapeError: CytoscapeError): void => {
    console.error('Cytoscape error:', cytoscapeError.message)
    setError(cytoscapeError)
  }

  useEffect(() => {
    setError(null)
    const elements = createCytoscapeElements(nodes, edges)
    return elements.length === 0
      ? undefined
      : initializeWithElements(elements, containerRef, cyRef, handleError)
  }, [nodes, edges])

  return error ? renderErrorDisplay(error) : renderCanvas(containerRef)
}
