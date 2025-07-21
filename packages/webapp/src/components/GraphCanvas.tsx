import React, { useState, useEffect } from 'react'
import { Effect } from 'effect'

import { CytoscapeGraphComponent } from './CytoscapeGraphComponent'
import { GraphErrorDisplay } from './GraphErrorDisplay'
import { GraphLoadingSpinner } from './GraphLoadingSpinner'
import type { GraphNode, GraphEdge } from '../services'

interface IGraphCanvasProps {
  readonly organization: string
  readonly useTopics: boolean
}

type GraphState =
  | { readonly type: 'loading' }
  | { readonly type: 'error'; readonly error: unknown }
  | {
      readonly type: 'success'
      readonly data: {
        readonly nodes: readonly GraphNode[]
        readonly edges: readonly GraphEdge[]
      }
    }

const createApiUrl = (
  organization: string,
  useTopics: boolean
): Effect.Effect<string, Error> =>
  !organization ||
  typeof organization !== 'string' ||
  organization.trim().length === 0
    ? Effect.fail(
        new Error(
          'Organization parameter is required and must be a non-empty string'
        )
      )
    : Effect.succeed(
        (() => {
          const cleanOrg = encodeURIComponent(organization.trim())
          return `http://localhost:8081/api/graph/${cleanOrg}${useTopics ? '?useTopics=true' : ''}`
        })()
      )

const validateGraphData = (
  data: unknown
): Effect.Effect<
  {
    readonly nodes: readonly GraphNode[]
    readonly edges: readonly GraphEdge[]
  },
  Error
> =>
  !data || typeof data !== 'object'
    ? Effect.fail(new Error('Invalid graph data: expected object'))
    : Effect.succeed(
        (() => {
          const graphData = data as {
            readonly nodes?: unknown
            readonly edges?: unknown
          }
          const nodes = Array.isArray(graphData.nodes) ? graphData.nodes : []
          const edges = Array.isArray(graphData.edges) ? graphData.edges : []
          console.warn(
            `Received ${nodes.length} nodes and ${edges.length} edges from API`
          )
          return {
            nodes: nodes as readonly GraphNode[],
            edges: edges as readonly GraphEdge[],
          }
        })()
      )

const fetchGraphData = (url: string) =>
  Effect.gen(function* () {
    const fetchResponse = yield* Effect.tryPromise({
      try: () => fetch(url),
      catch: (error) => new Error(`Network error: ${error instanceof Error ? error.message : String(error)}`)
    })

    if (!fetchResponse.ok) {
      yield* Effect.fail(new Error(`HTTP error! status: ${fetchResponse.status}`))
    }

    const graphApiJson = yield* Effect.tryPromise({
      try: () => fetchResponse.json(),
      catch: (error) => new Error(`JSON parsing error: ${error instanceof Error ? error.message : String(error)}`)
    })

    if (!graphApiJson || typeof graphApiJson !== 'object') {
      yield* Effect.fail(new Error('Invalid API response format'))
    }

    return yield* validateGraphData(graphApiJson.data)
  })

const renderGraphState = (state: GraphState): React.ReactElement => (
  <div data-testid="graph-canvas" style={{ width: '100%', height: '100vh' }}>
    {state.type === 'loading' ? (
      <GraphLoadingSpinner />
    ) : state.type === 'error' ? (
      <GraphErrorDisplay error={state.error} />
    ) : state.type === 'success' ? (
      (() => {
        const { nodes, edges } = state.data
        return !nodes && !edges ? (
          <GraphErrorDisplay error={new Error('No graph data received')} />
        ) : (
          <CytoscapeGraphComponent nodes={nodes} edges={edges} />
        )
      })()
    ) : (
      <GraphErrorDisplay error={new Error('Unknown graph state')} />
    )}
  </div>
)

export const GraphCanvas: React.FC<IGraphCanvasProps> = ({
  organization,
  useTopics,
}) => {
  const [state, setState] = useState<GraphState>({ type: 'loading' })

  useEffect(() => {
    if (
      !organization ||
      typeof organization !== 'string' ||
      organization.trim().length === 0
    ) {
      setState({
        type: 'error',
        error: new Error('Organization name is required'),
      })
      return
    }

    setState({ type: 'loading' })

    const pipeline = Effect.gen(function* () {
      const url = yield* createApiUrl(organization, useTopics)
      return yield* fetchGraphData(url)
    })

    Effect.runPromise(pipeline)
      .then((graphResult) => {
        console.warn('Graph data loaded successfully:', graphResult)
        setState({ type: 'success', data: graphResult })
      })
      .catch((error) => {
        console.error('Failed to load graph data:', error)
        setState({ type: 'error', error })
      })
  }, [organization, useTopics])

  return renderGraphState(state)
}
