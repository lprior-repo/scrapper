import { Context, Effect, Layer } from 'effect'
import { ParseResult } from '@effect/schema'
import { UnknownException } from 'effect/Cause'
import {
  GraphResponseSchema,
  validateApiResponseSync,
  type GraphNode,
  type GraphEdge,
  type GraphResponse,
} from '@overseer/shared'

// Re-export types for compatibility
export type { GraphNode, GraphEdge, GraphResponse }

// Utility functions
export const extractTopicsFromNodes = (
  nodes: readonly GraphNode[]
): readonly GraphTopic[] => {
  return nodes
    .filter((node) => node.type === 'topic')
    .map((node) => ({
      name: node.data.name as string,
      count: node.data.count as number,
    }))
}

export const extractNodesByType = (
  nodes: readonly GraphNode[],
  type: string
): readonly GraphNode[] => {
  return nodes.filter((node) => node.type === type)
}

// Graph topic type for utility functions
export interface GraphTopic {
  readonly name: string
  readonly count: number
}

// API Client Service
export interface ApiClient {
  readonly getGraph: (
    org: string,
    useTopics?: boolean
  ) => Effect.Effect<GraphResponse, ParseResult.ParseError | UnknownException>
}

export const ApiClient = Context.GenericTag<ApiClient>('ApiClient')

// Live Implementation
export const ApiClientLive = Layer.succeed(
  ApiClient,
  ApiClient.of({
    getGraph: (org: string, useTopics?: boolean) =>
      Effect.gen(function* () {
        const baseUrl = 'http://localhost:8081' // Your Go backend URL
        const url = `${baseUrl}/api/graph/${org}${useTopics ? '?useTopics=true' : ''}`

        const httpResponse = yield* Effect.tryPromise({
          try: () => fetch(url, {
            method: 'GET',
            headers: {
              Accept: 'application/json',
            },
          }),
          catch: (error) => new Error(`Network error: ${error instanceof Error ? error.message : String(error)}`)
        })

        httpResponse.ok ? void 0 : yield* Effect.fail(new Error(`HTTP error! status: ${httpResponse.status}`))

        const apiJsonData = yield* Effect.tryPromise({
          try: () => httpResponse.json(),
          catch: (error) => new Error(`JSON parsing error: ${error instanceof Error ? error.message : String(error)}`)
        })

        // Use the shared schema validation
        return validateApiResponseSync(
          GraphResponseSchema,
          apiJsonData,
          `Graph API response for ${org}`
        )
      }),
  })
)
