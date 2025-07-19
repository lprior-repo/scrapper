import React, { useMemo, useState, useEffect, useRef } from 'react'
import { Effect } from 'effect'
import { Network } from 'vis-network/standalone'
import { type GraphNode, type GraphEdge } from '../services'
import 'vis-network/styles/vis-network.css'

interface GraphCanvasProps {
  organization: string
  useTopics: boolean
}

// Loading component
const LoadingSpinner: React.FC = () => (
  <div
    style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh',
      fontSize: '1.5rem',
      color: '#58a6ff',
    }}
  >
    <div>Loading graph data...</div>
  </div>
)

// Error display component
const ErrorDisplay: React.FC<{ error: unknown }> = ({ error }) => (
  <div
    style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh',
      color: '#f85149',
      padding: '2rem',
      textAlign: 'center',
    }}
  >
    <div>
      <h2>Error loading graph</h2>
      <p>{error instanceof Error ? error.message : String(error)}</p>
    </div>
  </div>
)

// Vis.js graph component
interface VisGraphComponentProps {
  nodes?: GraphNode[]
  edges?: GraphEdge[]
}

const VisGraphComponent: React.FC<VisGraphComponentProps> = ({
  nodes,
  edges,
}) => {
  const containerRef = useRef<HTMLDivElement>(null)
  const networkRef = useRef<Network | null>(null)

  const visOptions = useMemo(
    () => ({
      nodes: {
        shape: 'dot',
        size: 25,
        font: {
          size: 14,
          color: '#c9d1d9',
          face: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
        },
        borderWidth: 2,
        shadow: {
          enabled: true,
          color: 'rgba(0,0,0,0.5)',
          size: 10,
          x: 3,
          y: 3,
        },
      },
      edges: {
        width: 2,
        color: {
          color: '#30363d',
          highlight: '#58a6ff',
          hover: '#58a6ff',
        },
        arrows: {
          to: {
            enabled: true,
            scaleFactor: 0.8,
          },
        },
        smooth: {
          enabled: true,
          type: 'curvedCW',
          roundness: 0.2,
        },
        font: {
          size: 12,
          color: '#8b949e',
          strokeWidth: 0,
          align: 'middle',
        },
      },
      physics: {
        enabled: true,
        solver: 'forceAtlas2Based',
        forceAtlas2Based: {
          gravitationalConstant: -50,
          centralGravity: 0.01,
          springLength: 150,
          springConstant: 0.08,
          damping: 0.4,
          avoidOverlap: 0.8,
        },
        stabilization: {
          enabled: true,
          iterations: 1000,
          updateInterval: 25,
        },
      },
      interaction: {
        hover: true,
        tooltipDelay: 300,
        zoomView: true,
        dragView: true,
        navigationButtons: true,
        keyboard: {
          enabled: true,
          speed: { x: 10, y: 10, zoom: 0.02 },
        },
      },
      layout: {
        improvedLayout: true,
        randomSeed: 42,
      },
      groups: {
        organization: {
          color: {
            background: '#238636',
            border: '#2ea043',
            highlight: { background: '#2ea043', border: '#3fb950' },
          },
          shape: 'hexagon',
          size: 35,
        },
        repository: {
          color: {
            background: '#1f6feb',
            border: '#388bfd',
            highlight: { background: '#388bfd', border: '#58a6ff' },
          },
          shape: 'box',
          shapeProperties: {
            borderRadius: 6,
          },
        },
        user: {
          color: {
            background: '#8b5cf6',
            border: '#a78bfa',
            highlight: { background: '#a78bfa', border: '#c4b5fd' },
          },
          shape: 'dot',
          size: 20,
        },
        team: {
          color: {
            background: '#f59e0b',
            border: '#fbbf24',
            highlight: { background: '#fbbf24', border: '#fcd34d' },
          },
          shape: 'square',
          size: 25,
        },
        topic: {
          color: {
            background: '#fb7185',
            border: '#f97316',
            highlight: { background: '#f97316', border: '#fb923c' },
          },
          shape: 'diamond',
          size: 28,
        },
      },
    }),
    []
  )

  // Transform nodes to vis format
  const visNodes = useMemo(
    () =>
      nodes?.map((node) => ({
        id: node.id,
        label: node.label,
        group: node.type,
        title: `${node.type}: ${node.label}`,
        x: node.position.x,
        y: node.position.y,
      })) || [],
    [nodes]
  )

  // Transform edges to vis format
  const visEdges = useMemo(
    () =>
      edges?.map((edge) => ({
        id: edge.id,
        from: edge.source,
        to: edge.target,
        label: edge.label,
        title: `${edge.type}: ${edge.label}`,
      })) || [],
    [edges]
  )

  const graphData = useMemo(() => ({
    nodes: visNodes,
    edges: visEdges,
  }), [visNodes, visEdges])

  useEffect(() => {
    if (containerRef.current && graphData.nodes.length > 0) {
      // Destroy existing network if it exists
      if (networkRef.current) {
        networkRef.current.destroy()
      }

      // Create new network
      networkRef.current = new Network(containerRef.current, graphData, visOptions)
      
      console.log('Network instance created:', networkRef.current)
    }

    // Cleanup on unmount
    return () => {
      if (networkRef.current) {
        networkRef.current.destroy()
        networkRef.current = null
      }
    }
  }, [graphData, visOptions])

  return (
    <div 
      ref={containerRef}
      data-testid="graph-canvas"
      style={{ width: '100%', height: '100vh' }}
    />
  )
}

// Main graph canvas component using Effect
export const GraphCanvas: React.FC<GraphCanvasProps> = ({
  organization,
  useTopics,
}) => {
  const [state, setState] = useState<
    | { type: 'loading' }
    | { type: 'error'; error: unknown }
    | { type: 'success'; data: { nodes: GraphNode[]; edges: GraphEdge[] } }
  >({ type: 'loading' })

  React.useEffect(() => {
    const loadData = Effect.gen(function* () {
      const url = `http://localhost:8081/api/graph/${organization}${useTopics ? '?useTopics=true' : ''}`
      const response = yield* Effect.tryPromise(() =>
        fetch(url).then((res) => {
          if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
          return res.json()
        })
      )
      return response.data
    })

    Effect.runPromise(loadData)
      .then((data) => {
        setState({ type: 'success', data })
      })
      .catch((error) => {
        setState({ type: 'error', error })
      })
  }, [organization, useTopics])

  switch (state.type) {
    case 'loading':
      return <LoadingSpinner />
    case 'error':
      return <ErrorDisplay error={state.error} />
    case 'success':
      return (
        <VisGraphComponent nodes={state.data.nodes} edges={state.data.edges} />
      )
  }
}
