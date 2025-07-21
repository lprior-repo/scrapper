/**
 * Graph Data Hook using React 19's `use` hook
 *
 * This hook integrates with the promise cache to provide
 * efficient data fetching for React Suspense boundaries
 */

import { use } from 'react'
import { Effect } from 'effect'
import { getCachedPromise, createCacheKey } from '../utils/promiseCache'
import type { GraphNode, GraphEdge } from '../services'

interface GraphData {
  readonly nodes: readonly GraphNode[]
  readonly edges: readonly GraphEdge[]
}

interface UseGraphDataOptions {
  readonly enabled?: boolean
  readonly staleTime?: number
  readonly retryOnError?: boolean
}

interface GraphApiResponse {
  readonly data?: {
    readonly nodes: readonly GraphNode[]
    readonly edges: readonly GraphEdge[]
  }
  readonly success?: boolean
  readonly message?: string
}

/**
 * Validate graph API response structure
 */
const validateGraphResponse = (response: unknown): GraphData => {
  if (!response || typeof response !== 'object') {
    throw new Error('Invalid API response: Expected object')
  }

  const apiResponse = response as GraphApiResponse

  if (!apiResponse.data || typeof apiResponse.data !== 'object') {
    throw new Error('Invalid API response: Missing data field')
  }

  const { nodes, edges } = apiResponse.data

  if (!Array.isArray(nodes)) {
    throw new Error('Invalid API response: nodes must be an array')
  }

  if (!Array.isArray(edges)) {
    throw new Error('Invalid API response: edges must be an array')
  }

  return {
    nodes: nodes as readonly GraphNode[],
    edges: edges as readonly GraphEdge[],
  }
}

/**
 * Create API URL with proper encoding and validation
 */
const createApiUrl = (organization: string, useTopics: boolean): string => {
  if (
    !organization ||
    typeof organization !== 'string' ||
    organization.trim().length === 0
  ) {
    throw new Error(
      'Organization parameter is required and must be a non-empty string'
    )
  }

  const cleanOrg = encodeURIComponent(organization.trim())
  const baseUrl = 'http://localhost:8081'
  return `${baseUrl}/api/graph/${cleanOrg}${useTopics ? '?useTopics=true' : ''}`
}

/**
 * Fetch graph data with comprehensive error handling
 */
const fetchGraphData = async (
  organization: string,
  useTopics: boolean
): Promise<GraphData> => {
  const url = createApiUrl(organization, useTopics)

  // Create Effect pipeline for data fetching
  const pipeline = Effect.gen(function* () {
    // Make HTTP request
    const httpResponse = yield* Effect.tryPromise({
      try: () =>
        fetch(url, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            Accept: 'application/json',
            'X-Request-Source': 'webapp-frontend',
            'X-Request-Timestamp': Date.now().toString(),
          },
        }),
      catch: (error) => {
        const errorMessage = `Network error: ${error instanceof Error ? error.message : String(error)}`
        console.error('Network request failed:', errorMessage, {
          url,
          organization,
          useTopics,
        })
        return new Error(errorMessage)
      },
    })

    // Check response status
    if (!httpResponse.ok) {
      const errorMessage = `HTTP error! status: ${httpResponse.status} ${httpResponse.statusText}`
      console.error('HTTP request failed:', errorMessage, {
        url,
        status: httpResponse.status,
        statusText: httpResponse.statusText,
        organization,
        useTopics,
      })

      // Try to get error details from response body
      try {
        const errorBody = yield* Effect.tryPromise({
          try: () => httpResponse.text(),
          catch: () => new Error('Could not read error response'),
        })

        throw new Error(`${errorMessage}\nResponse: ${errorBody}`)
      } catch {
        throw new Error(errorMessage)
      }
    }

    // Parse JSON response
    const apiJsonData = yield* Effect.tryPromise({
      try: () => httpResponse.json(),
      catch: (error) => {
        const errorMessage = `JSON parsing error: ${error instanceof Error ? error.message : String(error)}`
        console.error('JSON parsing failed:', errorMessage, {
          url,
          organization,
          useTopics,
        })
        return new Error(errorMessage)
      },
    })

    // Validate and return data
    const validatedData = validateGraphResponse(apiJsonData)

    console.log('Graph data fetched successfully:', {
      organization,
      useTopics,
      nodesCount: validatedData.nodes.length,
      edgesCount: validatedData.edges.length,
    })

    return validatedData
  })

  // Execute the Effect pipeline as a Promise
  try {
    return await Effect.runPromise(pipeline)
  } catch (error) {
    // Re-throw with additional context
    const contextualError =
      error instanceof Error
        ? new Error(
            `Failed to fetch graph data for "${organization}": ${error.message}`
          )
        : new Error(
            `Failed to fetch graph data for "${organization}": ${String(error)}`
          )

    console.error('Graph data fetch failed:', contextualError, {
      organization,
      useTopics,
      url,
      originalError: error,
    })

    throw contextualError
  }
}

/**
 * Hook to fetch graph data using React 19's use hook with Suspense
 */
export const useGraphData = (
  organization: string,
  useTopics: boolean = false,
  options: UseGraphDataOptions = {}
): GraphData => {
  const {
    enabled = true,
    staleTime = 5 * 60 * 1000,
    retryOnError = true,
  } = options

  // Early validation
  if (
    !organization ||
    typeof organization !== 'string' ||
    organization.trim().length === 0
  ) {
    throw new Error(
      'Organization parameter is required and must be a non-empty string'
    )
  }

  if (!enabled) {
    // Return empty data when disabled
    return { nodes: [], edges: [] }
  }

  // Create cache key
  const cacheKey = createCacheKey('graph-data', organization, useTopics)

  // Get cached promise (this will trigger Suspense if not resolved)
  const promise = getCachedPromise(cacheKey, () =>
    fetchGraphData(organization, useTopics)
  )

  // Use React 19's use hook to suspend until promise resolves
  return use(promise)
}

/**
 * Hook to preload graph data without suspending
 */
export const usePreloadGraphData = (
  organization: string,
  useTopics: boolean = false
): void => {
  if (
    !organization ||
    typeof organization !== 'string' ||
    organization.trim().length === 0
  ) {
    return
  }

  const cacheKey = createCacheKey('graph-data', organization, useTopics)

  // Preload data into cache without suspending
  getCachedPromise(cacheKey, () => fetchGraphData(organization, useTopics))
}

/**
 * Hook to get graph data loading state without suspending
 */
export const useGraphDataState = (
  organization: string,
  useTopics: boolean = false
): {
  readonly isLoading: boolean
  readonly isError: boolean
  readonly isSuccess: boolean
} => {
  if (
    !organization ||
    typeof organization !== 'string' ||
    organization.trim().length === 0
  ) {
    return { isLoading: false, isError: true, isSuccess: false }
  }

  const cacheKey = createCacheKey('graph-data', organization, useTopics)

  // Access cache state without triggering fetch
  const { getDefaultCache } = getCachedPromise as any
  const cache = getDefaultCache()
  const status = cache.getStatus(cacheKey)

  return {
    isLoading: status === 'pending',
    isError: status === 'rejected',
    isSuccess: status === 'resolved',
  }
}

/**
 * Hook to refresh graph data by evicting cache
 */
export const useRefreshGraphData = () => {
  const refresh = (organization: string, useTopics: boolean = false) => {
    if (
      !organization ||
      typeof organization !== 'string' ||
      organization.trim().length === 0
    ) {
      return
    }

    const cacheKey = createCacheKey('graph-data', organization, useTopics)

    // Access default cache and evict the entry
    const { getDefaultCache } = getCachedPromise as any
    const cache = getDefaultCache()
    cache.evict(cacheKey)
  }

  return refresh
}

/**
 * Hook to get cache statistics for debugging
 */
export const useGraphDataCacheStats = () => {
  const { getDefaultCache } = getCachedPromise as any
  const cache = getDefaultCache()
  return cache.getStats()
}

export type { GraphData, UseGraphDataOptions }
