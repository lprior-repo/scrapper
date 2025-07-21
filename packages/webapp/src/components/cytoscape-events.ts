import type { Core, EventObject } from 'cytoscape'

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
 * Edge selection styles
 */
const edgeSelectedStyles = {
  'line-color': '#f85149',
  'target-arrow-color': '#f85149',
  width: 4,
  'overlay-color': '#f85149',
  'overlay-opacity': 0.3,
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
  node && !node.selected()
    ? (node.style(nodeDefaultStyles), node.removestyle('border-color'))
    : void 0
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
  edge
    ? !edge.selected()
      ? edge.style(edgeDefaultStyles)
      : edge.style(edgeSelectedStyles)
    : void 0
}

/**
 * Enhanced background tap handler with detailed logging
 */
const createBackgroundTapHandler =
  (cy: Core) =>
  (event: EventObject): void => {
    event.target === cy && cy
      ? (() => {
          const selectedElements = cy.elements(':selected')
          const selectedCount = selectedElements.length

          selectedCount > 0
            ? (console.warn(
                `🎯 Background clicked - Deselecting ${selectedCount} element(s)`
              ),
              selectedElements.unselect())
            : console.warn('🎯 Background clicked - No elements to deselect')

          // Clear both info displays when clicking background
          const existingEdgeDisplay =
            document.getElementById('edge-info-display')
          const existingNodeDisplay =
            document.getElementById('node-info-display')
          existingEdgeDisplay?.remove()
          existingNodeDisplay?.remove()
        })()
      : void 0
  }

/**
 * Sets up hover effects for nodes and edges since CSS :hover doesn't work in Cytoscape
 */
export const setupInteractiveEvents = (
  cy: Core,
  onNodeClick: (event: EventObject) => void,
  onEdgeClick: (event: EventObject) => void
): void => {
  cy.on('mouseover', 'node', handleNodeMouseover)
  cy.on('mouseout', 'node', handleNodeMouseout)
  cy.on('mouseover', 'edge', handleEdgeMouseover)
  cy.on('mouseout', 'edge', handleEdgeMouseout)
  cy.on('tap', 'node', onNodeClick)
  cy.on('tap', 'edge', onEdgeClick)
  cy.on('tap', createBackgroundTapHandler(cy))
}
