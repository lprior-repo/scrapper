import type { Core } from 'cytoscape'
import type {
  ContextMenuConfig,
  ContextMenuItem,
  MenuPosition,
  ElementData,
  HiddenElementsState,
} from './types/context-menu'

/**
 * Context menu action handlers interface
 */
export interface ContextMenuActions {
  readonly onViewDetails: (elementData: ElementData) => void
  readonly onZoomToNode: (nodeId: string) => void
  readonly onSelectConnected: (nodeId: string) => void
  readonly onHideNode: (nodeId: string) => void
  readonly onViewEdgeDetails: (edgeId: string) => void
  readonly onHideEdge: (edgeId: string) => void
  readonly onSelectSourceTarget: (edgeId: string) => void
  readonly onZoomToFit: () => void
  readonly onCenterGraph: () => void
  readonly onClearSelections: () => void
  readonly onShowAll: () => void
}

/**
 * Creates context menu items for node interactions
 */
const createNodeMenuItems = (
  nodeData: ElementData,
  actions: ContextMenuActions
): readonly ContextMenuItem[] => [
  {
    id: 'view-details',
    label: 'View Details',
    icon: '🔍',
    action: () => actions.onViewDetails(nodeData),
  },
  {
    id: 'separator-1',
    label: '',
    separator: true,
    action: () => {},
  },
  {
    id: 'zoom-to-node',
    label: 'Zoom to Node',
    icon: '🎯',
    action: () => actions.onZoomToNode(nodeData.id),
  },
  {
    id: 'select-connected',
    label: 'Select Connected',
    icon: '🔗',
    action: () => actions.onSelectConnected(nodeData.id),
  },
  {
    id: 'separator-2',
    label: '',
    separator: true,
    action: () => {},
  },
  {
    id: 'hide-node',
    label: 'Hide Node',
    icon: '👁️‍🗨️',
    action: () => actions.onHideNode(nodeData.id),
  },
]

/**
 * Creates context menu items for edge interactions
 */
const createEdgeMenuItems = (
  edgeData: ElementData,
  actions: ContextMenuActions
): readonly ContextMenuItem[] => [
  {
    id: 'view-edge-details',
    label: 'View Edge Details',
    icon: '🔍',
    action: () => actions.onViewEdgeDetails(edgeData.id),
  },
  {
    id: 'separator-1',
    label: '',
    separator: true,
    action: () => {},
  },
  {
    id: 'select-source-target',
    label: 'Select Source/Target',
    icon: '🎯',
    action: () => actions.onSelectSourceTarget(edgeData.id),
  },
  {
    id: 'separator-2',
    label: '',
    separator: true,
    action: () => {},
  },
  {
    id: 'hide-edge',
    label: 'Hide Edge',
    icon: '👁️‍🗨️',
    action: () => actions.onHideEdge(edgeData.id),
  },
]

/**
 * Creates context menu items for background interactions
 */
const createBackgroundMenuItems = (
  actions: ContextMenuActions,
  hiddenElements: HiddenElementsState
): readonly ContextMenuItem[] => [
  {
    id: 'zoom-to-fit',
    label: 'Zoom to Fit',
    icon: '🔍',
    action: actions.onZoomToFit,
  },
  {
    id: 'center-graph',
    label: 'Center Graph',
    icon: '🎯',
    action: actions.onCenterGraph,
  },
  {
    id: 'separator-1',
    label: '',
    separator: true,
    action: () => {},
  },
  {
    id: 'clear-selections',
    label: 'Clear Selections',
    icon: '✖️',
    action: actions.onClearSelections,
  },
  {
    id: 'separator-2',
    label: '',
    separator: true,
    action: () => {},
  },
  {
    id: 'show-all',
    label: 'Show All',
    icon: '👁️',
    disabled:
      hiddenElements.nodes.length === 0 && hiddenElements.edges.length === 0,
    action: actions.onShowAll,
  },
]

/**
 * Extracts element data from cytoscape element
 */
const extractElementData = (
  element: cytoscape.NodeSingular | cytoscape.EdgeSingular
): ElementData => {
  const data = element.data()
  return {
    id: data.id,
    label: data.label || data.id,
    type: data.type || 'unknown',
    data: data.data,
  }
}

/**
 * Creates node context menu configuration
 */
export const createNodeContextMenu = (
  node: cytoscape.NodeSingular,
  position: MenuPosition,
  actions: ContextMenuActions
): ContextMenuConfig => {
  const nodeData = extractElementData(node)
  return {
    type: 'node',
    targetId: nodeData.id,
    position,
    items: createNodeMenuItems(nodeData, actions),
  }
}

/**
 * Creates edge context menu configuration
 */
export const createEdgeContextMenu = (
  edge: cytoscape.EdgeSingular,
  position: MenuPosition,
  actions: ContextMenuActions
): ContextMenuConfig => {
  const edgeData = extractElementData(edge)
  return {
    type: 'edge',
    targetId: edgeData.id,
    position,
    items: createEdgeMenuItems(edgeData, actions),
  }
}

/**
 * Creates background context menu configuration
 */
export const createBackgroundContextMenu = (
  position: MenuPosition,
  actions: ContextMenuActions,
  hiddenElements: HiddenElementsState
): ContextMenuConfig => ({
  type: 'background',
  position,
  items: createBackgroundMenuItems(actions, hiddenElements),
})

/**
 * Determines if event position is within valid bounds for context menu
 */
export const isValidMenuPosition = (position: MenuPosition): boolean =>
  position.x >= 0 &&
  position.y >= 0 &&
  position.x <= window.innerWidth &&
  position.y <= window.innerHeight

/**
 * Extracts menu position from mouse event
 */
export const extractMenuPosition = (
  event: MouseEvent | cytoscape.EventObject
): MenuPosition => ({
  x: 'clientX' in event ? event.clientX : (event.originalEvent?.clientX ?? 0),
  y: 'clientY' in event ? event.clientY : (event.originalEvent?.clientY ?? 0),
})

/**
 * Selects all connected nodes to a given node
 */
export const selectConnectedNodes = (cy: Core, nodeId: string): void => {
  const node = cy.$(`#${nodeId}`)
  if (node.length === 0) return

  // Get connected edges and their connected nodes
  const connectedEdges = node.connectedEdges()
  const connectedNodes = connectedEdges.connectedNodes()

  // Select the original node, connected edges, and connected nodes
  cy.elements().unselect()
  node.select()
  connectedEdges.select()
  connectedNodes.select()
}

/**
 * Selects source and target nodes of an edge
 */
export const selectEdgeSourceTarget = (cy: Core, edgeId: string): void => {
  const edge = cy.$(`#${edgeId}`)
  if (edge.length === 0) return

  const source = edge.source()
  const target = edge.target()

  cy.elements().unselect()
  edge.select()
  source.select()
  target.select()
}

/**
 * Hides a node and its connected edges temporarily
 */
export const hideNodeAndEdges = (
  cy: Core,
  nodeId: string
): readonly string[] => {
  const node = cy.$(`#${nodeId}`)
  if (node.length === 0) return []

  const connectedEdges = node.connectedEdges()
  const elementsToHide = node.union(connectedEdges)

  elementsToHide.style('display', 'none')

  return connectedEdges.map((edge) => edge.id())
}

/**
 * Hides an edge temporarily
 */
export const hideEdge = (cy: Core, edgeId: string): void => {
  const edge = cy.$(`#${edgeId}`)
  if (edge.length === 0) return

  edge.style('display', 'none')
}

/**
 * Shows all hidden elements
 */
export const showAllElements = (cy: Core): void => {
  cy.elements().style('display', 'element')
}
