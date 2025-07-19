import React, { useRef, useEffect } from 'react';
import cytoscape, { type Core } from 'cytoscape';

import { createCytoscapeElements } from './graph-elements';
import { createCytoscapeStyles } from './graph-styles';
import type { GraphNode, GraphEdge } from '../services';

interface ICytoscapeGraphComponentProps {
  readonly nodes?: readonly GraphNode[];
  readonly edges?: readonly GraphEdge[];
}

const createCytoscapeConfig = (
  container: HTMLDivElement,
  elements: ReturnType<typeof createCytoscapeElements>,
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
  wheelSensitivity: 0.2,
  minZoom: 0.1,
  maxZoom: 3.0,
});

const initializeCytoscape = (
  containerRef: React.RefObject<HTMLDivElement>,
  elements: ReturnType<typeof createCytoscapeElements>,
): Core | null => {
  if (!containerRef.current) {
    return null;
  }

  const config = createCytoscapeConfig(containerRef.current, elements);
  const cyInstance = cytoscape(config);
  
  console.warn('Cytoscape instance created:', cyInstance);
  return cyInstance;
};

export const CytoscapeGraphComponent: React.FC<ICytoscapeGraphComponentProps> = ({
  nodes,
  edges,
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const cyRef = useRef<Core | null>(null);

  useEffect(() => {
    const elements = createCytoscapeElements(nodes, edges);
    
    if (elements.length === 0) {
      return;
    }

    // Destroy existing instance
    if (cyRef.current) {
      cyRef.current.destroy();
    }

    // Create new instance
    // eslint-disable-next-line functional/immutable-data
    cyRef.current = initializeCytoscape(containerRef, elements);

    // Cleanup function
    return () => {
      if (cyRef.current) {
        cyRef.current.destroy();
        // eslint-disable-next-line functional/immutable-data
        cyRef.current = null;
      }
    };
  }, [nodes, edges]);

  return (
    <div 
      ref={containerRef}
      data-testid="graph-canvas"
      style={{ width: '100%', height: '100vh' }}
    />
  );
};