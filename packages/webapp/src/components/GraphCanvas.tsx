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

const createApiUrl = (organization: string, useTopics: boolean): string =>
  !organization ||
  typeof organization !== 'string' ||
  organization.trim().length === 0
    ? Effect.runSync(
        Effect.fail(
          new Error(
            'Organization parameter is required and must be a non-empty string'
          )
        )
      )
    : (() => {
        const cleanOrg = encodeURIComponent(organization.trim())
        return `http://localhost:8081/api/graph/${cleanOrg}${useTopics ? '?useTopics=true' : ''}`
      })()

const validateGraphData = (
  data: unknown
): {
  readonly nodes: readonly GraphNode[]
  readonly edges: readonly GraphEdge[]
} =>
  !data || typeof data !== 'object'
    ? Effect.runSync(
        Effect.fail(new Error('Invalid graph data: expected object'))
      )
    : (() => {
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

const fetchGraphData = (url: string) =>
  Effect.gen(function* () {
    const response = yield* Effect.tryPromise(() =>
      fetch(url).then((res) =>
        res.ok
          ? res.json()
          : Effect.runSync(
              Effect.fail(new Error(`HTTP error! status: ${res.status}`))
            )
      )
    )

    return !response || typeof response !== 'object'
      ? Effect.runSync(Effect.fail(new Error('Invalid API response format')))
      : validateGraphData(response.data)
  })

const renderGraphState = (state: GraphState): React.ReactElement =>
  state.type === 'loading' ? (
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
  )

export const GraphCanvas: React.FC<IGraphCanvasProps> = ({
  organization,
  useTopics,
}) => {
  const [state, setState] = useState<GraphState>({ type: 'loading' })

  useEffect(() => {
    return !organization ||
      typeof organization !== 'string' ||
      organization.trim().length === 0
      ? setState({
          type: 'error',
          error: new Error('Organization name is required'),
        })
      : (() => {
          setState({ type: 'loading' })
          try {
            const url = createApiUrl(organization, useTopics)
            const loadData = fetchGraphData(url)
            Effect.runPromise(loadData)
              .then((data) => {
                console.warn('Graph data loaded successfully:', data)
                setState({ type: 'success', data })
              })
              .catch((error) => {
                console.error('Failed to load graph data:', error)
                setState({ type: 'error', error })
              })
          } catch (error) {
            console.error('Error creating API request:', error)
            setState({ type: 'error', error })
          }
        })()
  }, [organization, useTopics])

  return renderGraphState(state)
}
