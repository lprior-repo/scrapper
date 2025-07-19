import type { ElementDefinition } from 'cytoscape';

import type { GraphNode, GraphEdge } from '../services';

export const transformNodesToElements = (nodes: readonly GraphNode[]): readonly ElementDefinition[] =>
  nodes.map((node) => ({
    data: {
      id: node.id,
      label: node.label,
      type: node.type,
      ...node.data,
    },
    position: node.position ? {
      x: node.position.x,
      y: node.position.y,
    } : undefined,
  }));

export const transformEdgesToElements = (edges: readonly GraphEdge[]): readonly ElementDefinition[] =>
  edges.map((edge) => ({
    data: {
      id: edge.id,
      source: edge.source,
      target: edge.target,
      label: edge.label,
      type: edge.type,
    },
  }));

export const createCytoscapeElements = (
  nodes: readonly GraphNode[] | undefined,
  edges: readonly GraphEdge[] | undefined,
): readonly ElementDefinition[] => {
  const nodeElements = nodes ? transformNodesToElements(nodes) : [];
  const edgeElements = edges ? transformEdgesToElements(edges) : [];
  
  return [...nodeElements, ...edgeElements];
};