import type { ElementDefinition } from 'cytoscape'

import type { GraphNode, GraphEdge } from '../services'

/**
 * Validates node ID
 */
const hasValidNodeId = (node: GraphNode): boolean =>
  Boolean(node?.id) && typeof node.id === 'string' && node.id.trim().length > 0

/**
 * Validates node type
 */
const hasValidNodeType = (node: GraphNode): boolean =>
  Boolean(node.type) && typeof node.type === 'string'

/**
 * Validates that a node has a valid ID and required properties
 */
const isValidNode = (node: GraphNode): boolean =>
  hasValidNodeId(node) && hasValidNodeType(node)

/**
 * Validates edge ID
 */
const hasValidEdgeId = (edge: GraphEdge): boolean =>
  Boolean(edge?.id) && typeof edge.id === 'string' && edge.id.trim().length > 0

/**
 * Validates edge source
 */
const hasValidEdgeSource = (edge: GraphEdge): boolean =>
  Boolean(edge.source) &&
  typeof edge.source === 'string' &&
  edge.source.trim().length > 0

/**
 * Validates edge target
 */
const hasValidEdgeTarget = (edge: GraphEdge): boolean =>
  Boolean(edge.target) &&
  typeof edge.target === 'string' &&
  edge.target.trim().length > 0

/**
 * Validates edge type
 */
const hasValidEdgeType = (edge: GraphEdge): boolean => Boolean(edge.type)

/**
 * Validates that an edge has valid ID, source, and target
 */
const isValidEdge = (edge: GraphEdge): boolean =>
  hasValidEdgeId(edge) &&
  hasValidEdgeSource(edge) &&
  hasValidEdgeTarget(edge) &&
  hasValidEdgeType(edge)

/**
 * Checks if edge references valid nodes
 */
const hasValidNodeReferences = (
  edge: GraphEdge,
  validNodeIds: ReadonlySet<string>
): boolean => {
  const hasValidSource = validNodeIds.has(edge.source)
  const hasValidTarget = validNodeIds.has(edge.target)
  return hasValidSource && hasValidTarget
}

/**
 * Validates that all edges reference existing nodes
 */
const getValidEdges = (
  edges: readonly GraphEdge[],
  validNodeIds: ReadonlySet<string>
): readonly GraphEdge[] =>
  edges.filter((edge) =>
    isValidEdge(edge)
      ? hasValidNodeReferences(edge, validNodeIds)
        ? true
        : (console.warn(
            `Filtering out edge ${edge.id}: references invalid node(s)`
          ),
          false)
      : (console.warn('Filtering out edge with invalid properties:', edge),
        false)
  )

/**
 * Extracts a valid label from node data name field
 */
const extractNameLabel = (node: GraphNode): string | null =>
  node.data?.name &&
  typeof node.data.name === 'string' &&
  node.data.name.trim().length > 0
    ? node.data.name.trim()
    : null

/**
 * Extracts a valid label from node data login field
 */
const extractLoginLabel = (node: GraphNode): string | null =>
  node.data?.login &&
  typeof node.data.login === 'string' &&
  node.data.login.trim().length > 0
    ? node.data.login.trim()
    : null

/**
 * Creates a fallback label using node type and short ID
 */
const createFallbackLabel = (node: GraphNode): string => {
  const shortId = node.id.length > 8 ? `${node.id.substring(0, 8)}...` : node.id
  return `${node.type} ${shortId}`
}

/**
 * Generates a meaningful label for nodes with empty labels
 */
const generateNodeLabel = (node: GraphNode): string =>
  node.label && node.label.trim().length > 0
    ? node.label.trim()
    : (extractNameLabel(node) ??
      extractLoginLabel(node) ??
      createFallbackLabel(node))

/**
 * Creates an element definition from a valid node
 */
const createNodeElement = (node: GraphNode): ElementDefinition => ({
  data: {
    id: node.id.trim(),
    label: generateNodeLabel(node),
    type: node.type,
    ...node.data,
  },
  position: node.position
    ? {
        x: node.position.x,
        y: node.position.y,
      }
    : undefined,
})

/**
 * Filters and logs invalid nodes
 */
const filterValidNodes = (nodes: readonly GraphNode[]): readonly GraphNode[] =>
  nodes.filter((node) =>
    isValidNode(node)
      ? true
      : (console.warn('Filtering out node with invalid properties:', node),
        false)
  )

export const transformNodesToElements = (
  nodes: readonly GraphNode[]
): readonly ElementDefinition[] =>
  filterValidNodes(nodes).map(createNodeElement)

/**
 * Creates an element definition from a valid edge
 */
const createEdgeElement = (edge: GraphEdge): ElementDefinition => ({
  data: {
    id: edge.id.trim(),
    source: edge.source.trim(),
    target: edge.target.trim(),
    label: edge.label,
    type: edge.type,
  },
})

export const transformEdgesToElements = (
  edges: readonly GraphEdge[],
  validNodeIds: ReadonlySet<string>
): readonly ElementDefinition[] =>
  getValidEdges(edges, validNodeIds).map(createEdgeElement)

/**
 * Creates a set of valid node IDs from node elements
 */
const createValidNodeIds = (
  nodeElements: readonly ElementDefinition[]
): ReadonlySet<string> =>
  new Set(nodeElements.map((el) => el.data.id)) as ReadonlySet<string>

/**
 * Logs element creation summary
 */
const logElementCreation = (nodeCount: number, edgeCount: number): void => {
  const totalElements = nodeCount + edgeCount
  console.warn(
    `Created ${totalElements} elements: ${nodeCount} nodes, ${edgeCount} edges`
  )
}

/**
 * Processes valid nodes and edges into cytoscape elements
 */
const processValidElements = (
  nodeElements: readonly ElementDefinition[],
  edges: readonly GraphEdge[] | undefined
): readonly ElementDefinition[] => {
  const validNodeIds = createValidNodeIds(nodeElements)
  const edgeElements = edges
    ? transformEdgesToElements(edges, validNodeIds)
    : []
  logElementCreation(nodeElements.length, edgeElements.length)
  return [...nodeElements, ...edgeElements]
}

export const createCytoscapeElements = (
  nodes: readonly GraphNode[] | undefined,
  edges: readonly GraphEdge[] | undefined
): readonly ElementDefinition[] =>
  !nodes || nodes.length === 0
    ? (console.warn('No nodes provided to createCytoscapeElements'), [])
    : (() => {
        const nodeElements = transformNodesToElements(nodes)
        return nodeElements.length === 0
          ? (console.warn('No valid nodes after transformation'), [])
          : processValidElements(nodeElements, edges)
      })()
