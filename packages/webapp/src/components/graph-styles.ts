import type { Stylesheet } from 'cytoscape';

export const createBaseNodeStyle = (): Stylesheet => ({
  selector: 'node',
  style: {
    'background-color': '#666',
    'label': 'data(label)',
    'text-valign': 'center',
    'text-halign': 'center',
    'font-size': '14px',
    'font-family': '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    'color': '#c9d1d9',
    'text-outline-width': 2,
    'text-outline-color': '#0d1117',
    'border-width': 2,
    'width': 50,
    'height': 50,
  },
});

const createOrganizationNodeStyle = (): Stylesheet => ({
  selector: 'node[type="organization"]',
  style: {
    'background-color': '#238636',
    'border-color': '#2ea043',
    'shape': 'hexagon',
    'width': 70,
    'height': 70,
  },
});

const createRepositoryNodeStyle = (): Stylesheet => ({
  selector: 'node[type="repository"]',
  style: {
    'background-color': '#1f6feb',
    'border-color': '#388bfd',
    'shape': 'roundrectangle',
    'width': 60,
    'height': 40,
  },
});

const createUserNodeStyle = (): Stylesheet => ({
  selector: 'node[type="user"]',
  style: {
    'background-color': '#8b5cf6',
    'border-color': '#a78bfa',
    'shape': 'ellipse',
    'width': 40,
    'height': 40,
  },
});

const createTeamNodeStyle = (): Stylesheet => ({
  selector: 'node[type="team"]',
  style: {
    'background-color': '#f59e0b',
    'border-color': '#fbbf24',
    'shape': 'rectangle',
    'width': 50,
    'height': 50,
  },
});

const createTopicNodeStyle = (): Stylesheet => ({
  selector: 'node[type="topic"]',
  style: {
    'background-color': '#fb7185',
    'border-color': '#f97316',
    'shape': 'diamond',
    'width': 56,
    'height': 56,
  },
});

const createHoverNodeStyle = (): Stylesheet => ({
  selector: 'node:hover',
  style: {
    'border-width': 4,
    'border-color': '#58a6ff',
  },
});

export const createNodeTypeStyles = (): readonly Stylesheet[] => [
  createOrganizationNodeStyle(),
  createRepositoryNodeStyle(),
  createUserNodeStyle(),
  createTeamNodeStyle(),
  createTopicNodeStyle(),
  createHoverNodeStyle(),
];

export const createEdgeStyles = (): readonly Stylesheet[] => [
  {
    selector: 'edge',
    style: {
      'width': 2,
      'line-color': '#30363d',
      'target-arrow-color': '#30363d',
      'target-arrow-shape': 'triangle',
      'curve-style': 'bezier',
      'label': 'data(label)',
      'font-size': '12px',
      'color': '#8b949e',
      'text-background-color': '#0d1117',
      'text-background-opacity': 0.8,
      'text-background-padding': '4px',
    },
  },
  {
    selector: 'edge:hover',
    style: {
      'line-color': '#58a6ff',
      'target-arrow-color': '#58a6ff',
      'width': 3,
    },
  },
];

export const createCytoscapeStyles = (): readonly Stylesheet[] => [
  createBaseNodeStyle(),
  ...createNodeTypeStyles(),
  ...createEdgeStyles(),
];