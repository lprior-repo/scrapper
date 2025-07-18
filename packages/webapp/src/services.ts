import { Context, Effect, Layer } from 'effect'
import { Schema, ParseResult } from '@effect/schema'
import { UnknownException } from 'effect/Cause'

// Graph Data Types
export const GraphNode = Schema.Struct({
  id: Schema.String,
  type: Schema.String,
  label: Schema.String,
  data: Schema.Record(Schema.String, Schema.Unknown),
  position: Schema.Struct({
    x: Schema.Number,
    y: Schema.Number,
  }),
})

export type GraphNode = Schema.Schema.Type<typeof GraphNode>

export const GraphEdge = Schema.Struct({
  id: Schema.String,
  source: Schema.String,
  target: Schema.String,
  type: Schema.String,
  label: Schema.String,
})

export type GraphEdge = Schema.Schema.Type<typeof GraphEdge>

export const GraphTopic = Schema.Struct({
  name: Schema.String,
  count: Schema.Number,
})

export type GraphTopic = Schema.Schema.Type<typeof GraphTopic>

export const GraphResponse = Schema.Struct({
  nodes: Schema.Array(GraphNode),
  edges: Schema.Array(GraphEdge),
})

export type GraphResponse = Schema.Schema.Type<typeof GraphResponse>

// Utility functions
export const extractTopicsFromNodes = (nodes: GraphNode[]): GraphTopic[] => {
  return nodes
    .filter((node) => node.type === 'topic')
    .map((node) => ({
      name: node.data.name as string,
      count: node.data.count as number,
    }))
}

export const extractNodesByType = (
  nodes: GraphNode[],
  type: string
): GraphNode[] => {
  return nodes.filter((node) => node.type === type)
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
        const url = `${baseUrl}/api/graph/${org}${useTopics ? '?use_topics=true' : ''}`

        const response = yield* Effect.tryPromise(() =>
          fetch(url, {
            method: 'GET',
            headers: {
              Accept: 'application/json',
            },
          })
        )

        const json = yield* Effect.tryPromise(() => response.json())

        return yield* Schema.decodeUnknown(GraphResponse)(json)
      }),
  })
)
