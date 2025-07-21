import React, {
  useRef,
  useEffect,
  useState,
  useImperativeHandle,
  forwardRef,
  useCallback,
} from 'react'
import cytoscape, {
  type Core,
  type EventObject,
  type ElementDefinition,
} from 'cytoscape'

import { createCytoscapeElements } from './graph-elements'
import { createCytoscapeStyles } from './graph-styles'
import type { GraphNode, GraphEdge } from '../services'

interface ICytoscapeGraphComponentProps {
  readonly nodes?: readonly GraphNode[]
  readonly edges?: readonly GraphEdge[]
}

interface CytoscapeGraphRef {
  readonly zoomToFit: () => void
  readonly centerGraph: () => void
  readonly resetZoom: () => void
  readonly zoomToNode: (nodeId: string) => void
  readonly getSelectedElements: () => readonly string[]
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
  wheelSensitivity: 0.5,
  zoomingEnabled: true,
  userZoomingEnabled: true,
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
 * Keyboard shortcut configuration
 */
interface KeyboardShortcut {
  readonly key: string
  readonly ctrlKey?: boolean
  readonly description: string
  readonly action: (cy: Core) => void
}

/**
 * Pan distance for arrow key navigation (in pixels)
 */
const PAN_DISTANCE = 50

/**
 * Zoom step for +/- keys
 */
const ZOOM_STEP = 0.2

/**
 * Creates visual feedback for keyboard shortcuts
 */
const createShortcutFeedback = (action: string, key: string): void => {
  // Remove existing feedback
  const existingFeedback = document.getElementById('keyboard-feedback')
  existingFeedback?.remove()

  const feedbackElement = document.createElement('div')
  feedbackElement.id = 'keyboard-feedback'
  feedbackElement.style.cssText = `
    position: fixed;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(0, 0, 0, 0.8);
    color: #f0f6fc;
    padding: 8px 16px;
    border-radius: 4px;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    font-size: 12px;
    z-index: 1000;
    animation: fadeInOut 1.5s ease-out forwards;
    pointer-events: none;
  `

  feedbackElement.textContent = `${action} (${key})`

  // Add CSS animation
  const style = document.createElement('style')
  style.textContent = `
    @keyframes fadeInOut {
      0% { opacity: 0; transform: translateX(-50%) translateY(20px); }
      20% { opacity: 1; transform: translateX(-50%) translateY(0); }
      80% { opacity: 1; transform: translateX(-50%) translateY(0); }
      100% { opacity: 0; transform: translateX(-50%) translateY(-10px); }
    }
  `
  document.head.appendChild(style)

  document.body.appendChild(feedbackElement)

  // Auto-remove after animation
  setTimeout(() => {
    feedbackElement.remove()
    style.remove()
  }, 1500)
}

/**
 * Pan the graph in a specific direction
 */
const panGraph = (
  cy: Core,
  direction: 'up' | 'down' | 'left' | 'right'
): void => {
  const pan = cy.pan()
  const newPan = { ...pan }

  switch (direction) {
    case 'up':
      newPan.y += PAN_DISTANCE
      break
    case 'down':
      newPan.y -= PAN_DISTANCE
      break
    case 'left':
      newPan.x += PAN_DISTANCE
      break
    case 'right':
      newPan.x -= PAN_DISTANCE
      break
  }

  cy.pan(newPan)
}

/**
 * Zoom the graph by a specific factor
 */
const zoomGraph = (cy: Core, factor: number): void => {
  const currentZoom = cy.zoom()
  const newZoom = Math.max(
    cy.minZoom(),
    Math.min(cy.maxZoom(), currentZoom + factor)
  )
  cy.zoom(newZoom)
}

/**
 * Keyboard shortcuts configuration
 */
const createKeyboardShortcuts = (): readonly KeyboardShortcut[] => [
  // Arrow keys for panning
  {
    key: 'ArrowUp',
    description: 'Pan Up',
    action: (cy: Core) => {
      panGraph(cy, 'up')
      createShortcutFeedback('Pan Up', '↑')
    },
  },
  {
    key: 'ArrowDown',
    description: 'Pan Down',
    action: (cy: Core) => {
      panGraph(cy, 'down')
      createShortcutFeedback('Pan Down', '↓')
    },
  },
  {
    key: 'ArrowLeft',
    description: 'Pan Left',
    action: (cy: Core) => {
      panGraph(cy, 'left')
      createShortcutFeedback('Pan Left', '←')
    },
  },
  {
    key: 'ArrowRight',
    description: 'Pan Right',
    action: (cy: Core) => {
      panGraph(cy, 'right')
      createShortcutFeedback('Pan Right', '→')
    },
  },
  // Zoom controls
  {
    key: '+',
    description: 'Zoom In',
    action: (cy: Core) => {
      zoomGraph(cy, ZOOM_STEP)
      createShortcutFeedback('Zoom In', '+')
    },
  },
  {
    key: '=',
    description: 'Zoom In',
    action: (cy: Core) => {
      zoomGraph(cy, ZOOM_STEP)
      createShortcutFeedback('Zoom In', '+')
    },
  },
  {
    key: '-',
    description: 'Zoom Out',
    action: (cy: Core) => {
      zoomGraph(cy, -ZOOM_STEP)
      createShortcutFeedback('Zoom Out', '-')
    },
  },
  {
    key: 'PageUp',
    description: 'Zoom In',
    action: (cy: Core) => {
      zoomGraph(cy, ZOOM_STEP)
      createShortcutFeedback('Zoom In', 'Page Up')
    },
  },
  {
    key: 'PageDown',
    description: 'Zoom Out',
    action: (cy: Core) => {
      zoomGraph(cy, -ZOOM_STEP)
      createShortcutFeedback('Zoom Out', 'Page Down')
    },
  },
  // Reset and fit controls
  {
    key: ' ',
    description: 'Reset Zoom and Center',
    action: (cy: Core) => {
      cy.zoom(1.0)
      cy.center()
      createShortcutFeedback('Reset Zoom & Center', 'Space')
    },
  },
  {
    key: 'Home',
    description: 'Reset Zoom and Center',
    action: (cy: Core) => {
      cy.zoom(1.0)
      cy.center()
      createShortcutFeedback('Reset Zoom & Center', 'Home')
    },
  },
  {
    key: 'f',
    description: 'Zoom to Fit',
    action: (cy: Core) => {
      cy.fit(undefined, 50)
      createShortcutFeedback('Zoom to Fit', 'F')
    },
  },
  {
    key: 'F',
    description: 'Zoom to Fit',
    action: (cy: Core) => {
      cy.fit(undefined, 50)
      createShortcutFeedback('Zoom to Fit', 'F')
    },
  },
  {
    key: '0',
    description: 'Zoom to Fit',
    action: (cy: Core) => {
      cy.fit(undefined, 50)
      createShortcutFeedback('Zoom to Fit', '0')
    },
  },
  // Selection controls
  {
    key: 'Escape',
    description: 'Clear Selections',
    action: (cy: Core) => {
      const selectedCount = cy.elements(':selected').length
      cy.elements().unselect()
      createShortcutFeedback(`Clear ${selectedCount} Selection(s)`, 'Esc')

      // Also clear info displays
      const existingEdgeDisplay = document.getElementById('edge-info-display')
      const existingNodeDisplay = document.getElementById('node-info-display')
      existingEdgeDisplay?.remove()
      existingNodeDisplay?.remove()
    },
  },
  {
    key: 'a',
    ctrlKey: true,
    description: 'Select All Elements',
    action: (cy: Core) => {
      cy.elements().select()
      const totalCount = cy.elements().length
      createShortcutFeedback(`Select All (${totalCount})`, 'Ctrl+A')
    },
  },
]

/**
 * Keyboard shortcut hook for graph navigation
 */
const useKeyboardShortcuts = (
  cyRef: React.MutableRefObject<Core | null>
): void => {
  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      // Don't handle shortcuts if user is typing in an input field
      if (
        event.target instanceof HTMLInputElement ||
        event.target instanceof HTMLTextAreaElement ||
        (event.target as Element)?.isContentEditable
      ) {
        return
      }

      // Don't handle shortcuts if cytoscape instance is not available
      if (!cyRef.current) {
        return
      }

      const shortcuts = createKeyboardShortcuts()
      const matchingShortcut = shortcuts.find(
        (shortcut) =>
          shortcut.key === event.key &&
          Boolean(shortcut.ctrlKey) === (event.ctrlKey || event.metaKey)
      )

      if (matchingShortcut) {
        event.preventDefault()
        event.stopPropagation()
        matchingShortcut.action(cyRef.current)
      }
    },
    [cyRef]
  )

  useEffect(() => {
    document.addEventListener('keydown', handleKeyDown)
    return () => {
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [handleKeyDown])
}

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
  if (!edge) return

  if (!edge.selected()) {
    edge.style(edgeDefaultStyles)
  } else {
    // Keep selection styling when selected
    edge.style(edgeSelectedStyles)
  }
}

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
 * Displays edge relationship details in console and UI feedback
 */
const displayEdgeDetails = (edge: cytoscape.EdgeSingular): void => {
  const edgeData = edge.data()
  const sourceNode = edge.source().data()
  const targetNode = edge.target().data()

  const edgeInfo = {
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

  // Console output for debugging
  console.warn('🔗 Edge Relationship Details:', {
    id: edgeInfo.id,
    relationship: edgeInfo.relationship,
    label: edgeInfo.label,
    source: edgeInfo.source,
    target: edgeInfo.target,
    fullData: edgeData,
  })

  // Create visual feedback element
  createEdgeInfoDisplay(edgeInfo)
}

/**
 * Edge information interface
 */
interface EdgeInfo {
  readonly id: string
  readonly relationship: string
  readonly label: string
  readonly source: {
    readonly id: string
    readonly label: string
    readonly type: string
  }
  readonly target: {
    readonly id: string
    readonly label: string
    readonly type: string
  }
}

/**
 * Creates a temporary UI element to display edge information
 */
const createEdgeInfoDisplay = (edgeInfo: EdgeInfo): void => {
  // Remove existing edge info display
  const existingDisplay = document.getElementById('edge-info-display')
  existingDisplay?.remove()

  const infoElement = document.createElement('div')
  infoElement.id = 'edge-info-display'
  infoElement.style.cssText = `
    position: fixed;
    top: 20px;
    right: 20px;
    background: #1c2128;
    border: 1px solid #30363d;
    border-radius: 6px;
    padding: 16px;
    color: #f0f6fc;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    font-size: 14px;
    max-width: 350px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    z-index: 1000;
    animation: slideIn 0.3s ease-out;
  `

  const formatRelationshipType = (type: string): string => {
    const typeMap: Record<string, string> = {
      owns: '👑 Owns',
      member_of: '👥 Member Of',
      codeowner: '🛡️ Code Owner',
      maintained_by: '🔧 Maintained By',
      has_topic: '🏷️ Has Topic',
      has: '📁 Has',
    }
    return typeMap[type] || `🔗 ${type}`
  }

  infoElement.innerHTML = `
    <div style="margin-bottom: 12px; padding-bottom: 8px; border-bottom: 1px solid #30363d;">
      <strong style="color: #58a6ff;">Relationship Details</strong>
    </div>
    <div style="margin-bottom: 8px;">
      <strong>Type:</strong> ${formatRelationshipType(edgeInfo.relationship)}
    </div>
    ${
      edgeInfo.label !== 'No label'
        ? `
    <div style="margin-bottom: 8px;">
      <strong>Label:</strong> ${edgeInfo.label}
    </div>
    `
        : ''
    }
    <div style="margin-bottom: 8px;">
      <strong>Source:</strong> ${edgeInfo.source.label} (${edgeInfo.source.type})
    </div>
    <div style="margin-bottom: 12px;">
      <strong>Target:</strong> ${edgeInfo.target.label} (${edgeInfo.target.type})
    </div>
    <div style="font-size: 12px; color: #7d8590;">
      Click anywhere to dismiss
    </div>
  `

  // Add CSS animation
  const style = document.createElement('style')
  style.textContent = `
    @keyframes slideIn {
      from {
        transform: translateX(100%);
        opacity: 0;
      }
      to {
        transform: translateX(0);
        opacity: 1;
      }
    }
  `
  document.head.appendChild(style)

  document.body.appendChild(infoElement)

  // Auto-remove after 10 seconds or on click
  const removeDisplay = (): void => {
    infoElement.remove()
    style.remove()
  }

  setTimeout(removeDisplay, 10000)
  infoElement.addEventListener('click', removeDisplay)
  document.addEventListener('click', removeDisplay, { once: true })
}

/**
 * Handles edge-specific click events with enhanced functionality
 */
const handleEdgeClick = (event: EventObject): void => {
  // Prevent event bubbling to background
  event.stopPropagation()

  const edge = event.target
  const isEdgeValid = edge && edge.isEdge()

  isEdgeValid
    ? (() => {
        const isCtrlClick =
          event.originalEvent?.ctrlKey || event.originalEvent?.metaKey || false
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
        const existingNodeDisplay = document.getElementById('node-info-display')
        existingNodeDisplay?.remove()
      })()
    : void 0
}

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
 * Node information interface
 */
interface NodeInfo {
  readonly id: string
  readonly type: string
  readonly label: string
  readonly data?: Record<string, unknown>
}

/**
 * Creates a temporary UI element to display node information
 */
const createNodeInfoDisplay = (nodeInfo: NodeInfo): void => {
  // Remove existing node info display
  const existingDisplay = document.getElementById('node-info-display')
  existingDisplay?.remove()

  const infoElement = document.createElement('div')
  infoElement.id = 'node-info-display'
  infoElement.style.cssText = `
    position: fixed;
    top: 20px;
    left: 20px;
    background: #1c2128;
    border: 1px solid #30363d;
    border-radius: 6px;
    padding: 16px;
    color: #f0f6fc;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    font-size: 14px;
    max-width: 400px;
    max-height: 80vh;
    overflow-y: auto;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    z-index: 1000;
    animation: slideInLeft 0.3s ease-out;
  `

  const formatNodeType = (type: string): string => {
    const typeMap: Record<string, string> = {
      organization: '🏢 Organization',
      repository: '📁 Repository',
      team: '👥 Team',
      user: '👤 User',
      topic: '🏷️ Topic',
    }
    return typeMap[type] || `🔵 ${type}`
  }

  const formatDataValue = (value: unknown): string =>
    typeof value === 'object' ? JSON.stringify(value, null, 2) : String(value)

  const createDataRow = ([key, value]: readonly [string, unknown]): string => `
  <div style="margin-bottom: 6px; padding: 6px; background: #21262d; border-radius: 3px;">
    <strong style="color: #79c0ff;">${key}:</strong>
    <div style="margin-top: 2px; word-break: break-all;">${formatDataValue(value)}</div>
  </div>
`

  const createDataTable = (data: Record<string, unknown>): string =>
    !data || Object.keys(data).length === 0
      ? '<div style="color: #7d8590; font-style: italic;">No additional data</div>'
      : Object.entries(data).map(createDataRow).join('')

  infoElement.innerHTML = `
    <div style="margin-bottom: 12px; padding-bottom: 8px; border-bottom: 1px solid #30363d;">
      <strong style="color: #58a6ff;">Node Details</strong>
    </div>
    <div style="margin-bottom: 8px;">
      <strong>Type:</strong> ${formatNodeType(nodeInfo.type)}
    </div>
    <div style="margin-bottom: 8px;">
      <strong>ID:</strong> <code style="background: #21262d; padding: 2px 4px; border-radius: 3px;">${nodeInfo.id}</code>
    </div>
    <div style="margin-bottom: 12px;">
      <strong>Label:</strong> ${nodeInfo.label}
    </div>
    ${
      Object.keys(nodeInfo.data || {}).length > 0
        ? `
    <div style="margin-bottom: 8px;">
      <strong>Additional Data:</strong>
    </div>
    <div style="margin-bottom: 12px;">
      ${createDataTable(nodeInfo.data)}
    </div>
    `
        : ''
    }
    <div style="font-size: 12px; color: #7d8590;">
      Click anywhere to dismiss
    </div>
  `

  // Add CSS animation for left slide-in
  const style = document.createElement('style')
  style.textContent = `
    @keyframes slideInLeft {
      from {
        transform: translateX(-100%);
        opacity: 0;
      }
      to {
        transform: translateX(0);
        opacity: 1;
      }
    }
  `
  document.head.appendChild(style)

  document.body.appendChild(infoElement)

  // Auto-remove after 10 seconds or on click
  const removeDisplay = (): void => {
    infoElement.remove()
    style.remove()
  }

  setTimeout(removeDisplay, 10000)
  infoElement.addEventListener('click', removeDisplay)
  document.addEventListener('click', removeDisplay, { once: true })
}

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
  createNodeInfoDisplay(nodeData as NodeInfo)
}

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
  if (!isCtrlClick) {
    // Regular click - clear other selections first
    element.cy().elements().unselect()
  }

  // Toggle selection state
  element.selected() ? element.unselect() : element.select()
}

/**
 * Handles node-specific click events with enhanced functionality
 */
const handleNodeClick = (event: EventObject): void => {
  // Prevent event bubbling to background
  event.stopPropagation()

  const node = event.target
  const isNodeValid = node && node.isNode()

  isNodeValid
    ? (() => {
        const isCtrlClick =
          event.originalEvent?.ctrlKey || event.originalEvent?.metaKey || false
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
        const existingEdgeDisplay = document.getElementById('edge-info-display')
        existingEdgeDisplay?.remove()
      })()
    : void 0
}

/**
 * Enhanced edge information interface for double-click details
 */
interface ExtendedEdgeInfo extends EdgeInfo {
  readonly sourceData: Record<string, unknown>
  readonly targetData: Record<string, unknown>
  readonly edgeProperties: Record<string, unknown>
}

/**
 * Creates an expanded modal display for double-clicked edge details
 */
const createExpandedEdgeModal = (edgeInfo: ExtendedEdgeInfo): void => {
  // Remove existing modal
  const existingModal = document.getElementById('expanded-edge-modal')
  existingModal?.remove()

  const modal = document.createElement('div')
  modal.id = 'expanded-edge-modal'
  modal.style.cssText = `
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
    animation: modalFadeIn 0.3s ease-out;
  `

  const modalContent = document.createElement('div')
  modalContent.style.cssText = `
    background: #1c2128;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 24px;
    color: #f0f6fc;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    font-size: 14px;
    max-width: 80vw;
    max-height: 80vh;
    overflow-y: auto;
    box-shadow: 0 16px 32px rgba(0, 0, 0, 0.3);
    position: relative;
  `

  const closeButton = document.createElement('button')
  closeButton.innerHTML = '✕'
  closeButton.style.cssText = `
    position: absolute;
    top: 16px;
    right: 16px;
    background: none;
    border: none;
    color: #7d8590;
    font-size: 20px;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 4px;
  `

  const formatRelationshipType = (type: string): string => {
    const typeMap: Record<string, string> = {
      owns: '👑 Owns',
      member_of: '👥 Member Of',
      codeowner: '🛡️ Code Owner',
      maintained_by: '🔧 Maintained By',
      has_topic: '🏷️ Has Topic',
      has: '📁 Has',
    }
    return typeMap[type] || `🔗 ${type}`
  }

  const formatDataSection = (
    title: string,
    data: Record<string, unknown>
  ): string => {
    if (!data || Object.keys(data).length === 0) {
      return `<div style="color: #7d8590; font-style: italic;">No ${title.toLowerCase()} data</div>`
    }

    return Object.entries(data)
      .map(
        ([key, value]) => `
        <div style="margin-bottom: 8px; padding: 8px; background: #21262d; border-radius: 4px;">
          <strong style="color: #79c0ff;">${key}:</strong>
          <div style="margin-top: 4px; word-break: break-all; font-family: 'SFMono-Regular', monospace; background: #161b22; padding: 8px; border-radius: 3px; font-size: 12px;">
            ${typeof value === 'object' ? JSON.stringify(value, null, 2) : String(value)}
          </div>
        </div>
      `
      )
      .join('')
  }

  modalContent.innerHTML = `
    <div style="margin-bottom: 24px; padding-bottom: 16px; border-bottom: 2px solid #30363d;">
      <h2 style="color: #58a6ff; margin: 0; font-size: 20px;">🔗 Expanded Relationship Details</h2>
    </div>
    
    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 24px; margin-bottom: 24px;">
      <div>
        <h3 style="color: #f85149; margin-bottom: 12px;">📋 Relationship Overview</h3>
        <div style="margin-bottom: 12px;">
          <strong>Type:</strong> ${formatRelationshipType(edgeInfo.relationship)}
        </div>
        <div style="margin-bottom: 12px;">
          <strong>ID:</strong> <code style="background: #21262d; padding: 4px 8px; border-radius: 3px;">${edgeInfo.id}</code>
        </div>
        ${
          edgeInfo.label !== 'No label'
            ? `
        <div style="margin-bottom: 12px;">
          <strong>Label:</strong> ${edgeInfo.label}
        </div>
        `
            : ''
        }
      </div>
      
      <div>
        <h3 style="color: #f85149; margin-bottom: 12px;">🎯 Connection Map</h3>
        <div style="background: #21262d; padding: 12px; border-radius: 6px;">
          <div style="margin-bottom: 8px;">
            <strong style="color: #56d364;">Source:</strong> ${edgeInfo.source.label}
          </div>
          <div style="margin-bottom: 8px; padding-left: 16px; color: #7d8590;">
            Type: ${edgeInfo.source.type} | ID: ${edgeInfo.source.id}
          </div>
          <div style="text-align: center; margin: 12px 0; color: #58a6ff; font-size: 18px;">
            ⬇️ ${formatRelationshipType(edgeInfo.relationship)} ⬇️
          </div>
          <div style="margin-bottom: 8px;">
            <strong style="color: #ff7b72;">Target:</strong> ${edgeInfo.target.label}
          </div>
          <div style="padding-left: 16px; color: #7d8590;">
            Type: ${edgeInfo.target.type} | ID: ${edgeInfo.target.id}
          </div>
        </div>
      </div>
    </div>

    <div style="margin-bottom: 24px;">
      <h3 style="color: #f85149; margin-bottom: 12px;">📊 Source Node Data</h3>
      ${formatDataSection('source', edgeInfo.sourceData)}
    </div>

    <div style="margin-bottom: 24px;">
      <h3 style="color: #f85149; margin-bottom: 12px;">📊 Target Node Data</h3>
      ${formatDataSection('target', edgeInfo.targetData)}
    </div>

    <div style="margin-bottom: 24px;">
      <h3 style="color: #f85149; margin-bottom: 12px;">🔧 Edge Properties</h3>
      ${formatDataSection('edge', edgeInfo.edgeProperties)}
    </div>

    <div style="text-align: center; padding-top: 16px; border-top: 1px solid #30363d; color: #7d8590; font-size: 12px;">
      Double-click for expanded details • Single-click for quick info • Press ESC or click outside to close
    </div>
  `

  modalContent.appendChild(closeButton)
  modal.appendChild(modalContent)

  // Add CSS animation
  const style = document.createElement('style')
  style.textContent = `
    @keyframes modalFadeIn {
      from {
        opacity: 0;
        transform: scale(0.9);
      }
      to {
        opacity: 1;
        transform: scale(1);
      }
    }
  `
  document.head.appendChild(style)

  document.body.appendChild(modal)

  // Close handlers
  const closeModal = (): void => {
    modal.remove()
    style.remove()
  }

  closeButton.addEventListener('click', closeModal)
  modal.addEventListener('click', (e) => {
    if (e.target === modal) closeModal()
  })
  document.addEventListener(
    'keydown',
    (e) => {
      if (e.key === 'Escape') closeModal()
    },
    { once: true }
  )
}

/**
 * Enhanced node information interface for double-click details
 */
interface ExtendedNodeInfo extends NodeInfo {
  readonly connections: {
    readonly incoming: readonly string[]
    readonly outgoing: readonly string[]
  }
  readonly metrics: {
    readonly totalConnections: number
    readonly nodeType: string
  }
}

/**
 * Creates an expanded modal display for double-clicked node details
 */
const createExpandedNodeModal = (nodeInfo: ExtendedNodeInfo): void => {
  // Remove existing modal
  const existingModal = document.getElementById('expanded-node-modal')
  existingModal?.remove()

  const modal = document.createElement('div')
  modal.id = 'expanded-node-modal'
  modal.style.cssText = `
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
    animation: modalFadeIn 0.3s ease-out;
  `

  const modalContent = document.createElement('div')
  modalContent.style.cssText = `
    background: #1c2128;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 24px;
    color: #f0f6fc;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    font-size: 14px;
    max-width: 80vw;
    max-height: 80vh;
    overflow-y: auto;
    box-shadow: 0 16px 32px rgba(0, 0, 0, 0.3);
    position: relative;
  `

  const closeButton = document.createElement('button')
  closeButton.innerHTML = '✕'
  closeButton.style.cssText = `
    position: absolute;
    top: 16px;
    right: 16px;
    background: none;
    border: none;
    color: #7d8590;
    font-size: 20px;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 4px;
  `

  const formatNodeType = (type: string): string => {
    const typeMap: Record<string, string> = {
      organization: '🏢 Organization',
      repository: '📁 Repository',
      team: '👥 Team',
      user: '👤 User',
      topic: '🏷️ Topic',
    }
    return typeMap[type] || `🔵 ${type}`
  }

  const formatDataSection = (
    title: string,
    data: Record<string, unknown>
  ): string => {
    if (!data || Object.keys(data).length === 0) {
      return `<div style="color: #7d8590; font-style: italic;">No ${title.toLowerCase()} data</div>`
    }

    return Object.entries(data)
      .map(
        ([key, value]) => `
        <div style="margin-bottom: 8px; padding: 8px; background: #21262d; border-radius: 4px;">
          <strong style="color: #79c0ff;">${key}:</strong>
          <div style="margin-top: 4px; word-break: break-all; font-family: 'SFMono-Regular', monospace; background: #161b22; padding: 8px; border-radius: 3px; font-size: 12px;">
            ${typeof value === 'object' ? JSON.stringify(value, null, 2) : String(value)}
          </div>
        </div>
      `
      )
      .join('')
  }

  const formatConnectionsList = (
    connections: readonly string[],
    title: string
  ): string => {
    if (connections.length === 0) {
      return `<div style="color: #7d8590; font-style: italic;">No ${title.toLowerCase()} connections</div>`
    }

    return connections
      .map(
        (conn) => `
        <div style="margin-bottom: 4px; padding: 6px 12px; background: #21262d; border-radius: 4px; border-left: 3px solid #58a6ff;">
          ${conn}
        </div>
      `
      )
      .join('')
  }

  modalContent.innerHTML = `
    <div style="margin-bottom: 24px; padding-bottom: 16px; border-bottom: 2px solid #30363d;">
      <h2 style="color: #58a6ff; margin: 0; font-size: 20px;">🔍 Expanded Node Details</h2>
    </div>
    
    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 24px; margin-bottom: 24px;">
      <div>
        <h3 style="color: #f85149; margin-bottom: 12px;">📋 Node Overview</h3>
        <div style="margin-bottom: 12px;">
          <strong>Type:</strong> ${formatNodeType(nodeInfo.type)}
        </div>
        <div style="margin-bottom: 12px;">
          <strong>ID:</strong> <code style="background: #21262d; padding: 4px 8px; border-radius: 3px;">${nodeInfo.id}</code>
        </div>
        <div style="margin-bottom: 12px;">
          <strong>Label:</strong> ${nodeInfo.label}
        </div>
      </div>
      
      <div>
        <h3 style="color: #f85149; margin-bottom: 12px;">📊 Connection Metrics</h3>
        <div style="background: #21262d; padding: 12px; border-radius: 6px;">
          <div style="margin-bottom: 8px;">
            <strong style="color: #56d364;">Total Connections:</strong> ${nodeInfo.metrics.totalConnections}
          </div>
          <div style="margin-bottom: 8px;">
            <strong style="color: #ff7b72;">Incoming:</strong> ${nodeInfo.connections.incoming.length}
          </div>
          <div>
            <strong style="color: #79c0ff;">Outgoing:</strong> ${nodeInfo.connections.outgoing.length}
          </div>
        </div>
      </div>
    </div>

    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 24px; margin-bottom: 24px;">
      <div>
        <h3 style="color: #f85149; margin-bottom: 12px;">📥 Incoming Connections</h3>
        <div style="max-height: 200px; overflow-y: auto;">
          ${formatConnectionsList(nodeInfo.connections.incoming, 'incoming')}
        </div>
      </div>
      
      <div>
        <h3 style="color: #f85149; margin-bottom: 12px;">📤 Outgoing Connections</h3>
        <div style="max-height: 200px; overflow-y: auto;">
          ${formatConnectionsList(nodeInfo.connections.outgoing, 'outgoing')}
        </div>
      </div>
    </div>

    <div style="margin-bottom: 24px;">
      <h3 style="color: #f85149; margin-bottom: 12px;">📊 Node Data</h3>
      ${formatDataSection('node', nodeInfo.data || {})}
    </div>

    <div style="text-align: center; padding-top: 16px; border-top: 1px solid #30363d; color: #7d8590; font-size: 12px;">
      Double-click for expanded details • Single-click for quick info • Press ESC or click outside to close
    </div>
  `

  modalContent.appendChild(closeButton)
  modal.appendChild(modalContent)

  // Add CSS animation
  const style = document.createElement('style')
  style.textContent = `
    @keyframes modalFadeIn {
      from {
        opacity: 0;
        transform: scale(0.9);
      }
      to {
        opacity: 1;
        transform: scale(1);
      }
    }
  `
  document.head.appendChild(style)

  document.body.appendChild(modal)

  // Close handlers
  const closeModal = (): void => {
    modal.remove()
    style.remove()
  }

  closeButton.addEventListener('click', closeModal)
  modal.addEventListener('click', (e) => {
    if (e.target === modal) closeModal()
  })
  document.addEventListener(
    'keydown',
    (e) => {
      if (e.key === 'Escape') closeModal()
    },
    { once: true }
  )
}

/**
 * Handles edge double-click events for expanded details
 */
const handleEdgeDoubleClick = (event: EventObject): void => {
  event.stopPropagation()

  const edge = event.target
  if (!edge || !edge.isEdge()) return

  const edgeData = edge.data()
  const sourceNode = edge.source()
  const targetNode = edge.target()

  // Apply visual feedback - zoom to edge and highlight
  const cy = edge.cy()
  cy.fit([edge, sourceNode, targetNode], 100)

  // Enhanced highlight effect for double-click
  const doubleClickHighlight = {
    'line-color': '#f85149',
    'target-arrow-color': '#f85149',
    width: 6,
    'overlay-color': '#f85149',
    'overlay-opacity': 0.4,
    'z-index': 999,
  }

  const originalStyles = edge.style()
  edge.style(doubleClickHighlight)

  setTimeout(() => {
    edge.style(originalStyles)
  }, 1500)

  // Create extended edge info
  const extendedEdgeInfo: ExtendedEdgeInfo = {
    id: edgeData.id,
    relationship: edgeData.type,
    label: edgeData.label || 'No label',
    source: {
      id: sourceNode.data('id'),
      label: sourceNode.data('label'),
      type: sourceNode.data('type'),
    },
    target: {
      id: targetNode.data('id'),
      label: targetNode.data('label'),
      type: targetNode.data('type'),
    },
    sourceData: sourceNode.data('data') || {},
    targetData: targetNode.data('data') || {},
    edgeProperties: { ...edgeData },
  }

  createExpandedEdgeModal(extendedEdgeInfo)

  console.warn('🔗🔗 Edge double-clicked - Expanded details:', extendedEdgeInfo)
}

/**
 * Handles node double-click events for expanded details
 */
const handleNodeDoubleClick = (event: EventObject): void => {
  event.stopPropagation()

  const node = event.target
  if (!node || !node.isNode()) return

  const nodeData = node.data()
  const cy = node.cy()

  // Apply visual feedback - zoom to node and highlight
  cy.fit(node, 150)

  // Enhanced highlight effect for double-click
  const doubleClickHighlight = {
    'border-width': 8,
    'border-color': '#f85149',
    'overlay-color': '#f85149',
    'overlay-opacity': 0.4,
    'z-index': 999,
  }

  const originalStyles = node.style()
  node.style(doubleClickHighlight)

  setTimeout(() => {
    node.style(originalStyles)
  }, 1500)

  // Collect connection information
  const connectedEdges = node.connectedEdges()
  const incomingConnections = connectedEdges
    .filter((edge) => edge.target().id() === node.id())
    .map((edge) => `${edge.source().data('label')} → ${edge.data('type')}`)

  const outgoingConnections = connectedEdges
    .filter((edge) => edge.source().id() === node.id())
    .map((edge) => `${edge.data('type')} → ${edge.target().data('label')}`)

  // Create extended node info
  const extendedNodeInfo: ExtendedNodeInfo = {
    id: nodeData.id,
    type: nodeData.type,
    label: nodeData.label,
    data: nodeData.data || {},
    connections: {
      incoming: incomingConnections,
      outgoing: outgoingConnections,
    },
    metrics: {
      totalConnections: connectedEdges.length,
      nodeType: nodeData.type,
    },
  }

  createExpandedNodeModal(extendedNodeInfo)

  console.warn('🔍🔍 Node double-clicked - Expanded details:', extendedNodeInfo)
}

/**
 * Handles element tap events for selection (legacy handler)
 */
const handleElementTap = (event: EventObject): void => {
  const element = event.target
  if (!element) return

  if (element.isEdge()) {
    handleEdgeClick(event)
  } else if (element.isNode()) {
    handleNodeClick(event)
  }
}

/**
 * Handles element double-tap events for expanded details
 */
const handleElementDoubleTap = (event: EventObject): void => {
  const element = event.target
  if (!element) return

  if (element.isEdge()) {
    handleEdgeDoubleClick(event)
  } else if (element.isNode()) {
    handleNodeDoubleClick(event)
  }
}

/**
 * Enhanced background tap handler with detailed logging
 */
const handleBackgroundTap =
  (cy: Core) =>
  (event: EventObject): void => {
    if (event.target === cy && cy) {
      const selectedElements = cy.elements(':selected')
      const selectedCount = selectedElements.length

      if (selectedCount > 0) {
        console.info(
          `🎯 Background clicked - Deselecting ${selectedCount} element(s)`
        )
        selectedElements.unselect()
      } else {
        console.info('🎯 Background clicked - No elements to deselect')
      }

      // Clear both info displays when clicking background
      const existingEdgeDisplay = document.getElementById('edge-info-display')
      const existingNodeDisplay = document.getElementById('node-info-display')
      existingEdgeDisplay?.remove()
      existingNodeDisplay?.remove()
    }
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
  cy.on('dbltap', 'node, edge', handleElementDoubleTap)
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
  outline: 'none', // Remove focus outline
} as const

/**
 * Renders error display component
 */
const renderErrorDisplay = (error: CytoscapeError): React.ReactElement => (
  <div data-testid="cytoscape-error" style={errorContainerStyle}>
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
  <div
    ref={containerRef}
    data-testid="cytoscape-container"
    style={canvasStyle}
    tabIndex={0}
    role="application"
    aria-label="Interactive graph visualization. Use arrow keys to pan, +/- to zoom, F to fit, Space to reset, Escape to clear selection, Ctrl+A to select all."
  />
)

/**
 * Destroys existing cytoscape instance
 */
const destroyInstance = (cyRef: React.MutableRefObject<Core | null>): void => {
  cyRef.current
    ? (() => {
        try {
          cyRef.current?.destroy()
        } catch (error) {
          console.warn('Error destroying cytoscape instance:', error)
        }
      })()
    : void 0
}

/**
 * Creates cleanup function for cytoscape instance
 */
const createCleanup =
  (cyRef: React.MutableRefObject<Core | null>) => (): void => {
    destroyInstance(cyRef)
    // Note: This mutation is necessary for React ref pattern
     
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
   
  cyRef.current = initializeCytoscape(containerRef, elements, handleError)
  return createCleanup(cyRef)
}

/**
 * Utility functions for graph navigation
 */
const createGraphNavigation = (cyRef: React.MutableRefObject<Core | null>) => ({
  /**
   * Fits entire graph in viewport
   */
  zoomToFit: (): void => {
    if (!cyRef.current) return
    cyRef.current.fit(undefined, 50) // 50px padding
  },

  /**
   * Centers the graph in viewport
   */
  centerGraph: (): void => {
    if (!cyRef.current) return
    cyRef.current.center()
  },

  /**
   * Resets zoom to 1.0
   */
  resetZoom: (): void => {
    if (!cyRef.current) return
    cyRef.current.zoom(1.0)
    cyRef.current.center()
  },

  /**
   * Zooms to specific node
   */
  zoomToNode: (nodeId: string): void => {
    if (!cyRef.current) return
    const node = cyRef.current.$(`#${nodeId}`)
    if (node.length > 0) {
      cyRef.current.fit(node, 100) // 100px padding around node
      node.select() // Also select the node for visual feedback
    }
  },

  /**
   * Returns currently selected nodes and edges
   */
  getSelectedElements: (): readonly string[] => {
    if (!cyRef.current) return []
    return cyRef.current.elements(':selected').map((element) => element.id())
  },
})

export const CytoscapeGraphComponent = forwardRef<
  CytoscapeGraphRef,
  ICytoscapeGraphComponentProps
>(({ nodes, edges }, ref) => {
  const containerRef = useRef<HTMLDivElement>(null)
  const cyRef = useRef<Core | null>(null)
  const [error, setError] = useState<CytoscapeError | null>(null)

  const handleError = (cytoscapeError: CytoscapeError): void => {
    console.error('Cytoscape error:', cytoscapeError.message)
    setError(cytoscapeError)
  }

  // Create navigation functions
  const navigation = createGraphNavigation(cyRef)

  // Set up keyboard shortcuts
  useKeyboardShortcuts(cyRef)

  // Expose navigation functions via imperative handle
  useImperativeHandle(ref, () => navigation, [])

  useEffect(() => {
    setError(null)
    const elements = createCytoscapeElements(nodes, edges)
    return elements.length === 0
      ? undefined
      : initializeWithElements(elements, containerRef, cyRef, handleError)
  }, [nodes, edges])

  return error ? renderErrorDisplay(error) : renderCanvas(containerRef)
})

// Export the ref type for external usage
export type { CytoscapeGraphRef }
