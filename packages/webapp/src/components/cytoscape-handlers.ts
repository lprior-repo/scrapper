import type { EventObject } from 'cytoscape'
import {
  createEdgeInfoDisplay,
  createNodeInfoDisplay,
  type EdgeInfo,
  type NodeInfo,
} from './cytoscape-ui-display'

/**
 * Selection state styles for temporary highlights
 */
const temporaryHighlightStyles = {
  'border-width': 6,
  'border-color': '#f85149',
  'overlay-color': '#f85149',
  'overlay-opacity': 0.3,
  'z-index': 999,
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
 * Default edge styles
 */
const edgeDefaultStyles = {
  'line-color': '#30363d',
  'target-arrow-color': '#30363d',
  width: 2,
  'overlay-opacity': 0,
} as const

/**
 * Provides temporary visual feedback for clicks
 */
const applyTemporaryHighlight = (
  element: cytoscape.NodeSingular | cytoscape.EdgeSingular
): void => {
  const originalStyles = element.style()
  element.style(temporaryHighlightStyles)

  setTimeout(() => {
    element.style(originalStyles)
  }, 800)
}

/**
 * Handles multi-select logic with ctrl+click
 */
const handleMultiSelect = (
  element: cytoscape.NodeSingular | cytoscape.EdgeSingular,
  isCtrlClick: boolean
): void => {
  !isCtrlClick ? element.cy().elements().unselect() : void 0
  element.selected() ? element.unselect() : element.select()
}

/**
 * Creates edge information object from cytoscape data
 */
const createEdgeInfo = (edge: cytoscape.EdgeSingular): EdgeInfo => {
  const edgeData = edge.data()
  const sourceNode = edge.source().data()
  const targetNode = edge.target().data()

  return {
    id: edgeData.id,
    relationship: edgeData.type,
    label: edgeData.label || 'No label',
    source: {
      id: sourceNode.id,
      label: sourceNode.label,
      type: sourceNode.type,
    },
    target: {
      id: targetNode.id,
      label: targetNode.label,
      type: targetNode.type,
    },
  }
}

/**
 * Displays edge relationship details in console and UI feedback
 */
const displayEdgeDetails = (edge: cytoscape.EdgeSingular): void => {
  const edgeInfo = createEdgeInfo(edge)

  // Console output for debugging
  console.warn('🔗 Edge Relationship Details:', {
    id: edgeInfo.id,
    relationship: edgeInfo.relationship,
    label: edgeInfo.label,
    source: edgeInfo.source,
    target: edgeInfo.target,
    fullData: edge.data(),
  })

  // Create visual feedback element
  createEdgeInfoDisplay(edgeInfo)
}

/**
 * Creates node information object from data
 */
const createNodeInfo = (nodeData: Record<string, unknown>): NodeInfo => ({
  id: String(nodeData.id),
  type: String(nodeData.type),
  label: String(nodeData.label),
  data: nodeData.data as Record<string, unknown>,
})

/**
 * Displays detailed node information in console and UI
 */
const displayNodeInfo = (nodeData: Record<string, unknown>): void => {
  const nodeInfo = {
    ID: nodeData.id,
    Type: nodeData.type,
    Label: nodeData.label,
    'Data Count':
      nodeData.data && typeof nodeData.data === 'object'
        ? Object.keys(nodeData.data).length
        : 0,
  }

  console.warn('🔍 Node Details:', nodeInfo)

  nodeData.data && typeof nodeData.data === 'object'
    ? console.warn('📊 Additional Node Data:', nodeData.data)
    : void 0

  // Create UI display
  createNodeInfoDisplay(createNodeInfo(nodeData))
}

/**
 * Removes existing display by ID
 */
const removeExistingDisplay = (displayId: string): void => {
  const existingDisplay = document.getElementById(displayId)
  existingDisplay?.remove()
}

/**
 * Checks if event has ctrl or meta key pressed
 */
const isCtrlOrMetaPressed = (event: EventObject): boolean =>
  event.originalEvent?.ctrlKey || event.originalEvent?.metaKey || false

/**
 * Handles edge-specific click events with enhanced functionality
 */
export const handleEdgeClick = (event: EventObject): void => {
  // Prevent event bubbling to background
  event.stopPropagation()

  const edge = event.target
  const isEdgeValid = edge && edge.isEdge()

  isEdgeValid
    ? (() => {
        const isCtrlClick = isCtrlOrMetaPressed(event)
        const edgeData = edge.data()

        // Apply temporary visual feedback
        applyTemporaryHighlight(edge)

        // Handle selection with multi-select support
        handleMultiSelect(edge, isCtrlClick)

        // Apply selection styling
        edge.selected()
          ? edge.style(edgeSelectedStyles)
          : edge.style(edgeDefaultStyles)

        // Display edge details
        displayEdgeDetails(edge)

        console.warn(
          `🔗 Edge clicked: ${edgeData.id} ${isCtrlClick ? '(Ctrl+Click)' : ''}`
        )

        // Remove any existing node info display when clicking edges
        removeExistingDisplay('node-info-display')
      })()
    : void 0
}

/**
 * Handles node-specific click events with enhanced functionality
 */
export const handleNodeClick = (event: EventObject): void => {
  // Prevent event bubbling to background
  event.stopPropagation()

  const node = event.target
  const isNodeValid = node && node.isNode()

  isNodeValid
    ? (() => {
        const isCtrlClick = isCtrlOrMetaPressed(event)
        const nodeData = node.data()

        // Apply temporary visual feedback
        applyTemporaryHighlight(node)

        // Display detailed node information
        displayNodeInfo(nodeData)

        // Handle selection with multi-select support
        handleMultiSelect(node, isCtrlClick)

        console.warn(
          `✨ Node clicked: ${nodeData.id} ${isCtrlClick ? '(Ctrl+Click)' : ''}`
        )

        // Remove any existing edge info display when clicking nodes
        removeExistingDisplay('edge-info-display')
      })()
    : void 0
}
