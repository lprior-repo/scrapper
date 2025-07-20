import { z } from 'zod'
import { Effect, Data } from 'effect'

/**
 * Health check response schema
 * Returns the health status of the API and its dependencies
 */
export const HealthResponseSchema = z.object({
  data: z.object({
    status: z.string().describe('Health status of the system'),
    database: z.string().describe('Database connection status'),
    version: z.string().describe('API version'),
    timestamp: z
      .string()
      .datetime()
      .describe('Current timestamp in ISO 8601 format'),
  }),
})

/**
 * Error response schema
 * Standard error response format for all API endpoints
 */
export const ErrorResponseSchema = z.object({
  error: z.object({
    message: z.string().describe('Human-readable error message'),
    code: z.string().optional().describe('Machine-readable error code'),
  }),
})

/**
 * Scan summary schema
 * Summary statistics from a GitHub organization scan
 */
export const ScanSummarySchema = z.object({
  total_repos: z
    .number()
    .int()
    .min(0)
    .describe('Total number of repositories scanned'),
  repos_with_codeowners: z
    .number()
    .int()
    .min(0)
    .describe('Number of repositories with CODEOWNERS files'),
  total_teams: z.number().int().min(0).describe('Total number of teams found'),
  unique_owners: z
    .array(z.string())
    .describe('List of unique codeowners (users and teams)'),
  api_calls_used: z
    .number()
    .int()
    .min(0)
    .describe('Number of GitHub API calls used'),
  processing_time_ms: z
    .number()
    .int()
    .min(0)
    .describe('Processing time in milliseconds'),
})

/**
 * Scan response schema
 * Response from scanning a GitHub organization
 */
export const ScanResponseSchema = z.object({
  data: z.object({
    success: z.boolean().describe('Whether the scan was successful'),
    organization: z.string().describe('Name of the scanned organization'),
    summary: ScanSummarySchema,
    errors: z
      .array(z.string())
      .describe('List of errors encountered during scanning'),
    data: z
      .object({})
      .passthrough()
      .optional()
      .describe(
        'Raw scan data including organization, repositories, teams, and codeowners'
      ),
  }),
})

/**
 * Graph position schema
 * Coordinates for positioning nodes in the graph visualization
 */
export const GraphPositionSchema = z.object({
  x: z.number().describe('X coordinate'),
  y: z.number().describe('Y coordinate'),
})

/**
 * Graph node type enum
 * Types of nodes that can appear in the graph
 */
export const GraphNodeTypeSchema = z.enum([
  'organization',
  'repository',
  'team',
  'user',
])

/**
 * Graph edge type enum
 * Types of relationships between nodes
 */
export const GraphEdgeTypeSchema = z.enum(['owns', 'member_of', 'codeowner'])

/**
 * Graph node schema
 * Individual node in the organization graph
 */
export const GraphNodeSchema = z.object({
  id: z.string().describe('Unique node identifier'),
  type: GraphNodeTypeSchema.describe('Node type'),
  label: z.string().describe('Display label for the node'),
  data: z.object({}).passthrough().optional().describe('Additional node data'),
  position: GraphPositionSchema.optional(),
})

/**
 * Graph edge schema
 * Connection between two nodes in the graph
 */
export const GraphEdgeSchema = z.object({
  id: z.string().describe('Unique edge identifier'),
  source: z.string().describe('Source node ID'),
  target: z.string().describe('Target node ID'),
  type: GraphEdgeTypeSchema.describe('Edge type'),
  label: z.string().optional().describe('Display label for the edge'),
})

/**
 * Graph response schema
 * Graph visualization data for an organization
 */
export const GraphResponseSchema = z.object({
  data: z.object({
    nodes: z.array(GraphNodeSchema).describe('List of nodes in the graph'),
    edges: z.array(GraphEdgeSchema).describe('List of edges connecting nodes'),
  }),
})

/**
 * Stats response schema
 * Statistical summary of a scanned organization
 */
export const StatsResponseSchema = z.object({
  data: z.object({
    organization: z.string().describe('Organization name'),
    total_repositories: z
      .number()
      .int()
      .min(0)
      .describe('Total number of repositories'),
    total_teams: z.number().int().min(0).describe('Total number of teams'),
    total_users: z.number().int().min(0).describe('Total number of users'),
    total_codeowners: z
      .number()
      .int()
      .min(0)
      .describe('Total number of codeowners entries'),
    codeowner_coverage: z
      .string()
      .describe('Percentage of repositories with codeowners'),
    last_scan_time: z
      .string()
      .datetime()
      .describe('Timestamp of last scan in ISO 8601 format'),
  }),
})

/**
 * Scan request parameters schema
 * Query parameters for scanning an organization
 */
export const ScanRequestParamsSchema = z.object({
  org: z.string().min(1).describe('GitHub organization name'),
  max_repos: z
    .number()
    .int()
    .min(1)
    .max(1000)
    .default(100)
    .optional()
    .describe('Maximum number of repositories to scan'),
  max_teams: z
    .number()
    .int()
    .min(1)
    .max(500)
    .default(50)
    .optional()
    .describe('Maximum number of teams to scan'),
})

/**
 * Graph request parameters schema
 * Query parameters for retrieving graph data
 */
export const GraphRequestParamsSchema = z.object({
  org: z.string().min(1).describe('GitHub organization name'),
  useTopics: z
    .boolean()
    .default(false)
    .optional()
    .describe('Use repository topics instead of teams for graph visualization'),
})

/**
 * Stats request parameters schema
 * Path parameters for retrieving statistics
 */
export const StatsRequestParamsSchema = z.object({
  org: z.string().min(1).describe('GitHub organization name'),
})

// Type inference helpers
export type HealthResponse = z.infer<typeof HealthResponseSchema>
export type ErrorResponse = z.infer<typeof ErrorResponseSchema>
export type ScanSummary = z.infer<typeof ScanSummarySchema>
export type ScanResponse = z.infer<typeof ScanResponseSchema>
export type GraphPosition = z.infer<typeof GraphPositionSchema>
export type GraphNodeType = z.infer<typeof GraphNodeTypeSchema>
export type GraphEdgeType = z.infer<typeof GraphEdgeTypeSchema>
export type GraphNode = z.infer<typeof GraphNodeSchema>
export type GraphEdge = z.infer<typeof GraphEdgeSchema>
export type GraphResponse = z.infer<typeof GraphResponseSchema>
export type StatsResponse = z.infer<typeof StatsResponseSchema>
export type ScanRequestParams = z.infer<typeof ScanRequestParamsSchema>
export type GraphRequestParams = z.infer<typeof GraphRequestParamsSchema>
export type StatsRequestParams = z.infer<typeof StatsRequestParamsSchema>

// Validation Error using Effect.ts Data.TaggedError (functional approach)
export const ValidationError = Data.TaggedError('ValidationError')<{
  readonly message: string
  readonly context?: string
  readonly zodError: z.ZodError
}>()

// Effect-based API response validation
export const validateApiResponse = <T>(
  schema: z.ZodSchema<T>,
  data: unknown,
  context?: string
): Effect.Effect<T, ValidationError> =>
  Effect.try({
    try: () => schema.parse(data),
    catch: (error) =>
      error instanceof z.ZodError
        ? ValidationError({
            message: context
              ? `Validation error in ${context}: ${error.message}`
              : `Validation error: ${error.message}`,
            context,
            zodError: error,
          })
        : ValidationError({
            message: `Unknown validation error: ${String(error)}`,
            context,
            zodError: new z.ZodError([]),
          }),
  })

// Synchronous version for backwards compatibility (uses Effect.runSync)
export const validateApiResponseSync = <T>(
  schema: z.ZodSchema<T>,
  data: unknown,
  context?: string
): T => {
  const result = schema.safeParse(data)

  return result.success
    ? result.data
    : Effect.runSync(
        Effect.fail(
          ValidationError({
            message: context
              ? `Validation error in ${context}: ${result.error.message}`
              : `Validation error: ${result.error.message}`,
            context,
            zodError: result.error,
          })
        )
      )
}

// Re-export commonly used schemas for convenience
export const schemas = {
  health: HealthResponseSchema,
  error: ErrorResponseSchema,
  scan: {
    params: ScanRequestParamsSchema,
    response: ScanResponseSchema,
    summary: ScanSummarySchema,
  },
  graph: {
    params: GraphRequestParamsSchema,
    response: GraphResponseSchema,
    node: GraphNodeSchema,
    edge: GraphEdgeSchema,
    position: GraphPositionSchema,
  },
  stats: {
    params: StatsRequestParamsSchema,
    response: StatsResponseSchema,
  },
} as const
